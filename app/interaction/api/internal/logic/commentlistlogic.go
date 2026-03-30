// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"miniblog/app/interaction/api/internal/svc"
	"miniblog/app/interaction/api/internal/types"
	"miniblog/app/interaction/rpc/interaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	rpcResp, err := l.svcCtx.InteractionRpc.CommentList(l.ctx, &interaction.CommentListRequest{
		PostId:   req.PostId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 【核心修复 3】：使用 make 强行初始化为空切片，绝对不给前端返回 null
	list := make([]types.CommentItem, 0)

	if rpcResp.List != nil {
		for _, item := range rpcResp.List {
			list = append(list, types.CommentItem{
				Id:         item.Id,
				PostId:     item.PostId,
				UserId:     item.UserId,
				Content:    item.Content,
				CreateTime: item.CreateTime,
			})
		}
	}

	return &types.CommentListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
