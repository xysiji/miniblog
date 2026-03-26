package logic

import (
	"context"

	"miniblog/app/notice/api/internal/svc"
	"miniblog/app/notice/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 单条设置为已读
func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ReadReq) (resp *types.ReadResp, err error) {
	err = l.svcCtx.NoticeModel.ReadById(l.ctx, req.NoticeId)
	if err != nil {
		l.Logger.Errorf("单条已读执行失败: noticeId=%d, err=%v", req.NoticeId, err)
		return nil, err
	}
	return &types.ReadResp{}, nil
}
