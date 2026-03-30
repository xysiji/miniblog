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
	total, err := l.svcCtx.CommentModel.CountByPostIdShard(l.ctx, in.PostId)
	if err != nil {
		l.Logger.Errorf("查询评论总数失败: %v", err)
		return nil, err
	}

	var list []*interaction.CommentItem
	if total > 0 {
		// 【核心修复 2】：计算真正的数据库偏移量 (Offset)，而不是直接传 Page！
		offset := (in.Page - 1) * in.PageSize

		// 传入计算好的 offset
		comments, err := l.svcCtx.CommentModel.FindPageListByPostIdShard(l.ctx, in.PostId, int(offset), int(in.PageSize))
		if err != nil {
			l.Logger.Errorf("查询评论列表失败: %v", err)
			return nil, err
		}

		for _, c := range comments {
			list = append(list, &interaction.CommentItem{
				Id:         c.Id,
				PostId:     c.PostId,
				UserId:     c.UserId,
				Content:    c.Content,
				CreateTime: c.CreateTime.Unix(),
			})
		}
	}

	return &interaction.CommentListResponse{
		List:  list,
		Total: total,
	}, nil
}
