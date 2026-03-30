package logic

import (
	"context"

	"miniblog/app/interaction/model"
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

func (l *CommentListLogic) CommentList(in *interaction.CommentListRequest) (*interaction.CommentListResponse, error) {
	// 1. 查询该博文下一级主评论 (root_id=0) 的总数
	total, err := l.svcCtx.CommentModel.CountByPostIdRootIdShard(l.ctx, in.PostId, 0)
	if err != nil {
		l.Logger.Errorf("查询主评论总数失败: %v", err)
		return nil, err
	}

	var list []*interaction.CommentItem
	if total > 0 {
		offset := (in.Page - 1) * in.PageSize

		// 2. 查询主评论列表
		mainComments, err := l.svcCtx.CommentModel.FindPageListByPostIdRootIdShard(l.ctx, in.PostId, 0, int(offset), int(in.PageSize))
		if err != nil {
			l.Logger.Errorf("查询主评论列表失败: %v", err)
			return nil, err
		}

		var rootIds []int64
		for _, c := range mainComments {
			rootIds = append(rootIds, c.Id)
			list = append(list, l.modelToProto(c))
		}

		// 3. 批量查询子评论（传入 postId 确保路由到同一个分表）
		subComments, err := l.svcCtx.CommentModel.FindAllByRootIdsShard(l.ctx, in.PostId, rootIds)
		if err != nil {
			l.Logger.Errorf("批量查询子评论失败: %v", err)
		}

		for _, sc := range subComments {
			list = append(list, l.modelToProto(sc))
		}
	}

	return &interaction.CommentListResponse{
		List:  list,
		Total: total,
	}, nil
}

func (l *CommentListLogic) modelToProto(m *model.Comment) *interaction.CommentItem {
	return &interaction.CommentItem{
		Id:            m.Id,
		PostId:        m.PostId,
		UserId:        m.UserId,
		Content:       m.Content,
		CreateTime:    m.CreateTime.Unix(),
		RootId:        m.RootId,
		ParentId:      m.ParentId,
		ReplyToUserId: m.ReplyToUserId,
	}
}
