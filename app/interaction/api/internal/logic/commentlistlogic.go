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

	// ⚠️ 架构师修改：遵循 types.go，使用指针切片 []*types.CommentItem
	roots := make([]*types.CommentItem, 0)
	subsMap := make(map[int64][]*types.CommentItem)

	if rpcResp.List != nil {
		for _, item := range rpcResp.List {
			apiItem := &types.CommentItem{
				Id:            item.Id,
				PostId:        item.PostId,
				UserId:        item.UserId,
				Username:      "待聚合",
				Avatar:        "",
				Content:       item.Content,
				CreateTime:    item.CreateTime,
				RootId:        item.RootId,
				ParentId:      item.ParentId,
				ReplyToUserId: item.ReplyToUserId,
				ReplyToName:   "",
				Children:      make([]*types.CommentItem, 0), // 初始化防止前端 null 报错
			}

			if item.RootId == 0 {
				roots = append(roots, apiItem)
			} else {
				subsMap[item.RootId] = append(subsMap[item.RootId], apiItem)
			}
		}

		// 树形挂载
		for _, root := range roots {
			if children, ok := subsMap[root.Id]; ok {
				root.Children = children
				root.ChildrenCount = int64(len(children))
			}
		}
	}

	return &types.CommentListResp{
		List:  roots,
		Total: rpcResp.Total,
	}, nil
}
