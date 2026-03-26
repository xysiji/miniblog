package logic

import (
	"context"

	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CommentList 获取评论列表
func (l *CommentListLogic) CommentList(in *interaction.CommentListRequest) (*interaction.CommentListResponse, error) {
	// 【核心修复】：调用带有 Shard 后缀的分表查询方法 (走从库读取)
	total, err := l.svcCtx.CommentModel.CountByPostIdShard(l.ctx, in.PostId)
	if err != nil {
		l.Logger.Errorf("查询评论总数失败: %v", err)
		return nil, err
	}

	var list []*interaction.CommentItem
	if total > 0 {
		// 【核心修复】：调用带有 Shard 后缀的分表查询方法 (走从库读取)
		comments, err := l.svcCtx.CommentModel.FindPageListByPostIdShard(l.ctx, in.PostId, int(in.Page), int(in.PageSize))
		if err != nil {
			l.Logger.Errorf("查询评论列表失败: %v", err)
			return nil, err
		}

		for _, c := range comments {
			list = append(list, &interaction.CommentItem{
				Id:      c.Id,
				PostId:  c.PostId,
				UserId:  c.UserId,
				Content: c.Content,
				// 注意：如果你的 c.CreateTime 是 time.Time 类型，请保留 .Unix()。如果是 int64，直接赋值即可。
				CreateTime: c.CreateTime.Unix(),
			})
		}
	}

	return &interaction.CommentListResponse{
		List:  list,
		Total: total,
	}, nil
}
