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
	timelineKey := "biz:post:timeline:global"
	start := (in.Page - 1) * in.PageSize
	stop := start + in.PageSize - 1

	postIdsStrs, err := l.svcCtx.BizRedis.ZrevrangeCtx(l.ctx, timelineKey, start, stop)
	if err != nil {
		l.Logger.Errorf("读取 Timeline 缓存失败，准备降级: %v", err)
	}

	totalCount, _ := l.svcCtx.BizRedis.ZcardCtx(l.ctx, timelineKey)
	var list []*post.PostItem

	if len(postIdsStrs) > 0 {
		l.Logger.Infof("Timeline 命中缓存, 获取到 %d 条 ID", len(postIdsStrs))

		for _, idStr := range postIdsStrs {
			postId, _ := strconv.ParseInt(idStr, 10, 64)
			postInfo, err := l.svcCtx.PostSlaveModel.FindOne(l.ctx, postId)

			if err == nil && postInfo != nil {
				// ✅ 修复：处理 sql.NullString
				imagesStr := ""
				if postInfo.Images.Valid {
					imagesStr = postInfo.Images.String
				}

				list = append(list, &post.PostItem{
					Id:         postInfo.Id,
					UserId:     postInfo.UserId,
					Content:    postInfo.Content,
					Images:     imagesStr,
					CreateTime: postInfo.CreateTime.Unix(),
					// 👇【核心补齐 1】：命中缓存时，将统计数据传给 RPC
					LikeCount:    postInfo.LikeCount,
					CommentCount: postInfo.CommentCount,
				})
			}
		}
	} else {
		l.Logger.Infof("🚨 Timeline 未命中！触发从库降级扫描 (Fallback to MySQL)")

		dbCount, err := l.svcCtx.PostSlaveModel.Count(l.ctx)
		if err == nil {
			totalCount = int(dbCount)
		}

		dbPosts, err := l.svcCtx.PostSlaveModel.FindPageListByPage(l.ctx, in.Page, in.PageSize)
		if err != nil {
			l.Logger.Errorf("降级查询 MySQL 失败: %v", err)
			return nil, fmt.Errorf("获取列表失败，底层服务异常")
		}

		for _, p := range dbPosts {
			// ✅ 修复：处理 sql.NullString
			imagesStr := ""
			if p.Images.Valid {
				imagesStr = p.Images.String
			}

			list = append(list, &post.PostItem{
				Id:         p.Id,
				UserId:     p.UserId,
				Content:    p.Content,
				Images:     imagesStr,
				CreateTime: p.CreateTime.Unix(),
				// 👇【核心补齐 2】：降级扫库时，将统计数据传给 RPC
				LikeCount:    p.LikeCount,
				CommentCount: p.CommentCount,
			})

			_, _ = l.svcCtx.BizRedis.ZaddCtx(l.ctx, timelineKey, p.CreateTime.Unix(), strconv.FormatInt(p.Id, 10))
		}
		l.Logger.Infof("✅ 缓存重建完毕，共恢复 %d 条数据至 Redis", len(dbPosts))
	}

	return &post.ListResponse{
		List:  list,
		Total: int64(totalCount),
	}, nil
}
