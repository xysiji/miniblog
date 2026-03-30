package logic

import (
	"context"

	"miniblog/app/relation/rpc/internal/svc"
	"miniblog/app/relation/rpc/relation"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckIsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIsFollowLogic {
	return &CheckIsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckIsFollowLogic) CheckIsFollow(in *relation.CheckIsFollowReq) (*relation.CheckIsFollowResp, error) {
	// todo: add your logic here and delete this line

	return &relation.CheckIsFollowResp{}, nil
}
