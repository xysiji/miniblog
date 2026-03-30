package logic

import (
	"context"

	"miniblog/app/relation/rpc/internal/svc"
	"miniblog/app/relation/rpc/relation"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationCountLogic {
	return &GetRelationCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationCountLogic) GetRelationCount(in *relation.RelationCountReq) (*relation.RelationCountResp, error) {
	// todo: add your logic here and delete this line

	return &relation.RelationCountResp{}, nil
}
