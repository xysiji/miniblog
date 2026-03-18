package logic

import (
	"context"
	"fmt"
	"strconv"

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

// List 获取微型博客列表 (支持多级缓存与读写分离)
func (l *ListLogic) List(in *post.ListRequest) (*post.ListResponse, error) {
	// 1. 定义 Timeline 的 Redis Key
	timelineKey := "biz:post:timeline:global"

	// 2. 计算 Redis ZSet 的分页参数 (0-based)
	start := (in.Page - 1) * in.PageSize
	stop := start + in.PageSize - 1

	// 3. 从 Redis 极速拉取当前页的博文 ID 列表 (按时间倒序)
	postIdsStrs, err := l.svcCtx.BizRedis.ZrevrangeCtx(l.ctx, timelineKey, start, stop)
	if err != nil {
		l.Logger.Errorf("读取 Timeline 缓存失败，准备降级: %v", err)
	}

	// 4. 获取 Timeline 中的总博文数，用于前端分页
	totalCount, _ := l.svcCtx.BizRedis.ZcardCtx(l.ctx, timelineKey)

	var list []*post.PostItem

	if len(postIdsStrs) > 0 {
		l.Logger.Infof("Timeline 命中缓存, 获取到 %d 条 ID", len(postIdsStrs))

		// 5. 遍历缓存中的 ID，去数据库查询完整信息
		for _, idStr := range postIdsStrs {
			postId, _ := strconv.ParseInt(idStr, 10, 64)

			// 6. 【分布式架构核心】：读操作回源查询时，路由到 Slave 从库
			postInfo, err := l.svcCtx.PostSlaveModel.FindOne(l.ctx, postId)
			if err == nil && postInfo != nil {
				list = append(list, &post.PostItem{
					Id:         postInfo.Id,
					UserId:     postInfo.UserId,
					Content:    postInfo.Content,
					CreateTime: postInfo.CreateTime.Unix(),
				})
			}
		}
	} else {
		l.Logger.Infof("🚨 Timeline 未命中！触发从库降级扫描 (Fallback to MySQL)")

		// 【修复1】：兜底查询：从 MySQL 查总数
		dbCount, err := l.svcCtx.PostSlaveModel.Count(l.ctx)
		if err == nil {
			totalCount = int(dbCount)
		}

		// 【修复2】：兜底查询：从 MySQL 查当前页的列表
		dbPosts, err := l.svcCtx.PostSlaveModel.FindPageListByPage(l.ctx, in.Page, in.PageSize)
		if err != nil {
			l.Logger.Errorf("降级查询 MySQL 失败: %v", err)
			return nil, fmt.Errorf("获取列表失败，底层服务异常")
		}

		// 【修复3】：组装数据，并【重建 Redis 缓存】
		for _, p := range dbPosts {
			list = append(list, &post.PostItem{
				Id:         p.Id,
				UserId:     p.UserId,
				Content:    p.Content,
				CreateTime: p.CreateTime.Unix(),
			})

			// 缓存预热：把从 MySQL 查出的数据，重新塞回 Redis 的 ZSet 中
			// Score 用创建时间戳，Value 用博文 ID
			_, _ = l.svcCtx.BizRedis.ZaddCtx(l.ctx, timelineKey, p.CreateTime.Unix(), strconv.FormatInt(p.Id, 10))
		}
		l.Logger.Infof("✅ 缓存重建完毕，共恢复 %d 条数据至 Redis", len(dbPosts))
	}

	return &post.ListResponse{
		List:  list,
		Total: int64(totalCount),
	}, nil
}
