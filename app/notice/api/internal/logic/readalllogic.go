package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/notice/api/internal/svc"
	"miniblog/app/notice/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadAllLogic {
	return &ReadAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadAllLogic) ReadAll(req *types.ReadAllReq) (resp *types.ReadAllResp, err error) {
	var userId int64
	if uidNum, ok := l.ctx.Value("userId").(json.Number); ok {
		userId, _ = uidNum.Int64()
	}

	if userId == 0 {
		l.Logger.Errorf("获取用户信息失败，Token异常")
		return nil, fmt.Errorf("未授权访问，请重新登录") // 【已修复 logx.Errorf】
	}

	// 你的 servicecontext.go 里确实已经挂载了 NoticeModel，所以直接调用没问题
	err = l.svcCtx.NoticeModel.ReadAllByUserId(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("一键已读执行失败: userId=%d, err=%v", userId, err)
		return nil, fmt.Errorf("操作失败，请稍后再试")
	}

	return &types.ReadAllResp{}, nil
}
