// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"miniblog/app/user/api/internal/svc"
	"miniblog/app/user/api/internal/types"
	"miniblog/app/user/rpc/userclient" // 引入 rpc 客户端的 types

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 1. 调用底层的 RPC 服务
	rpcRes, err := l.svcCtx.UserRpc.Register(l.ctx, &userclient.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	// 2. 组装返回给前端的 HTTP 响应
	return &types.RegisterResp{
		UserId: rpcRes.UserId,
	}, nil
}
