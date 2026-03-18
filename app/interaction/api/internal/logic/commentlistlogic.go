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

// 获取博文的评论列表
func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	// 1. 设置分页默认值 (防御性编程)
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 2. 呼叫底层 RPC 获取数据
	rpcResp, err := l.svcCtx.InteractionRpc.CommentList(l.ctx, &interaction.CommentListRequest{
		PostId:   req.PostId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 3. 将 RPC 响应数据映射为前端 API 所需的 JSON 结构
	var list []types.CommentItem
	for _, item := range rpcResp.List {
		list = append(list, types.CommentItem{
			Id:         item.Id,
			PostId:     item.PostId,
			UserId:     item.UserId,
			Content:    item.Content,
			CreateTime: item.CreateTime,
		})
	}

	// 4. 返回组装好的数据
	return &types.CommentListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
