package logic

import (
	"context"

	"miniblog/app/interaction/rpc/interaction"
	"miniblog/app/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikedListLogic {
	return &LikedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikedListLogic) LikedList(in *interaction.LikedListRequest) (*interaction.LikedListResponse, error) {
	// todo: add your logic here and delete this line

	return &interaction.LikedListResponse{}, nil
}
