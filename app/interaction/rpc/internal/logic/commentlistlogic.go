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

// 获取评论列表
func (l *CommentListLogic) CommentList(in *interaction.CommentListRequest) (*interaction.CommentListResponse, error) {
	// 1. 查询该博文下的评论总数
	total, err := l.svcCtx.CommentModel.CountByPostId(l.ctx, in.PostId)
	if err != nil || total == 0 {
		return &interaction.CommentListResponse{}, nil
	}

	// 2. 分页查询评论列表数据
	dbList, err := l.svcCtx.CommentModel.FindPageListByPostId(l.ctx, in.PostId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	// 3. 将数据库模型映射为 RPC Protobuf 返回模型
	var list []*interaction.CommentItem
	for _, item := range dbList {
		list = append(list, &interaction.CommentItem{
			Id:         item.Id,
			PostId:     item.PostId,
			UserId:     item.UserId,
			Content:    item.Content,
			CreateTime: item.CreateTime.Unix(), // 转换为时间戳传给前端网关
		})
	}

	return &interaction.CommentListResponse{
		List:  list,
		Total: total,
	}, nil
}
