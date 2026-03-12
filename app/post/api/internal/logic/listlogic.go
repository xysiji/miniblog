// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"miniblog/app/post/api/internal/svc"
	"miniblog/app/post/api/internal/types"
	"miniblog/app/post/rpc/postclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListReq) (resp *types.ListResp, err error) {
	// 1. 发起呼叫：调用底层 Post RPC 服务，索要分页数据
	rpcResp, err := l.svcCtx.PostRpc.List(l.ctx, &postclient.ListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 2. 数据清洗与转化 (BFF 模式核心表现)
	// 将底层的 rpcResp 转化为直接给前端看的 types.PostItem 数组
	var list []types.PostItem
	for _, item := range rpcResp.List {
		list = append(list, types.PostItem{
			Id:      item.Id,
			UserId:  item.UserId,
			Content: item.Content,
			// 核心转化：利用 Go 独特的 2006-01-02 15:04:05 诞生时间来格式化
			CreateAt: time.Unix(item.CreateTime, 0).Format("2006-01-02 15:04:05"),
		})
	}

	// 3. 组装最终给前端的 JSON 响应
	return &types.ListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
