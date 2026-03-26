package logic

import (
	"context"

	"miniblog/app/post/rpc/internal/svc"
	"miniblog/app/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增：详情和删除 RPC 方法
func (l *DetailLogic) Detail(in *post.DetailRequest) (*post.DetailResponse, error) {
	// todo: add your logic here and delete this line

	return &post.DetailResponse{}, nil
}
