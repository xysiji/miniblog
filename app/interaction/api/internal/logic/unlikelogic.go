package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/interaction/api/internal/svc"
	"miniblog/app/interaction/api/internal/types"
	"miniblog/app/interaction/rpc/interaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnlikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeLogic {
	return &UnlikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlikeLogic) Unlike(req *types.UnlikeReq) (resp *types.UnlikeResp, err error) {
	// 1. 【终极兼容版】从 JWT Token 解析上下文中获取当前登录用户的 ID
	var userId int64
	if uidVal := l.ctx.Value("userId"); uidVal != nil {
		switch v := uidVal.(type) {
		case json.Number:
			userId, _ = v.Int64()
		case float64:
			userId = int64(v)
		case string:
			fmt.Sscanf(v, "%d", &userId)
		}
	}

	l.Logger.Infof("=> [API网关层] 接收到取消点赞请求: 解析出 UserId=%d, 目标 PostId=%d", userId, req.PostId)

	// 2. 调用底层的 Interaction RPC 服务
	_, err = l.svcCtx.InteractionRpc.Unlike(l.ctx, &interaction.UnlikeRequest{
		UserId: userId,
		PostId: req.PostId,
	})

	if err != nil {
		l.Logger.Errorf("=> [API网关层] 取消点赞RPC返回拦截错误: %v", err)
		return nil, fmt.Errorf("操作被拦截: %v", err)
	}

	return &types.UnlikeResp{}, nil
}
