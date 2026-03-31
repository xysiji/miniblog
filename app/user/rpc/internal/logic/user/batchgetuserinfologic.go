package userlogic

import (
	"context"

	"miniblog/app/user/rpc/internal/svc"
	"miniblog/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUserInfoLogic {
	return &BatchGetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 【核心新增】：批量接口
func (l *BatchGetUserInfoLogic) BatchGetUserInfo(in *user.BatchUserInfoReq) (*user.BatchUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &user.BatchUserInfoResp{}, nil
}
