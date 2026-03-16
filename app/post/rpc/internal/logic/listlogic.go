package logic

import (
	"context"
	"strconv" // 新增：用于将 Redis 里取出的字符串 ID 转换为 int64

	"miniblog/app/post/rpc/internal/svc"
	"miniblog/app/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 内部获取列表的 RPC 方法 (多级缓存改造版)
func (l *ListLogic) List(in *post.ListRequest) (*post.ListResponse, error) {
	// 1. 参数防御与默认值设置
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	// 2. 第一步：先查总条数 (暂时直接查 MySQL)
	total, err := l.svcCtx.PostModel.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	// 3. 如果总数为 0，直接短路返回
	if total == 0 {
		return &post.ListResponse{
			List:  []*post.PostItem{},
			Total: 0,
		}, nil
	}

	// ==========================================
	// 本次新增核心逻辑：4. 多级缓存读取机制
	// ==========================================
	timelineKey := "biz:post:global_timeline"
	start := (in.Page - 1) * in.PageSize
	stop := start + in.PageSize - 1

	var respList []*post.PostItem

	// 第一级缓存：尝试从业务 Redis ZSet 中获取当前页的博文 ID 列表
	redisIds, err := l.svcCtx.BizRedis.ZrevrangeCtx(l.ctx, timelineKey, start, stop)

	if err == nil && len(redisIds) > 0 {
		// 【缓存命中分支】
		l.Logger.Infof("Timeline 命中 Redis 缓存, 获取到 %d 条 ID", len(redisIds))

		for _, idStr := range redisIds {
			id, parseErr := strconv.ParseInt(idStr, 10, 64)
			if parseErr != nil {
				continue
			}

			// 第二级缓存：FindOne 底层会自动走 go-zero 的 Redis 行级缓存 (sqlc)
			// 此时框架会自动拦截，如果行缓存里有这篇博文，直接返回；如果没有，它会查 MySQL 并写回缓存。
			p, findErr := l.svcCtx.PostModel.FindOne(l.ctx, id)
			if findErr == nil && p != nil {
				respList = append(respList, &post.PostItem{
					Id:         p.Id,
					UserId:     p.UserId,
					Content:    p.Content,
					CreateTime: p.CreateTime.Unix(),
				})
			}
		}
	} else {
		// 【缓存穿透/未命中分支（降级回源 MySQL）】
		l.Logger.Infof("Timeline 缓存未命中或为空，降级查询 MySQL")

		posts, findErr := l.svcCtx.PostModel.FindPageListByPage(l.ctx, in.Page, in.PageSize)
		if findErr != nil {
			return nil, findErr
		}

		for _, p := range posts {
			respList = append(respList, &post.PostItem{
				Id:         p.Id,
				UserId:     p.UserId,
				Content:    p.Content,
				CreateTime: p.CreateTime.Unix(),
			})
		}
	}

	// 5. 返回最终拼装好的结果集
	return &post.ListResponse{
		List:  respList,
		Total: total,
	}, nil
}
