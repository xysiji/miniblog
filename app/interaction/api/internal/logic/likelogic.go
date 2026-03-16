package logic

import (
	"context"
	"encoding/json"
	"fmt" // 新增导入

	"miniblog/app/interaction/api/internal/svc"
	"miniblog/app/interaction/api/internal/types"
	"miniblog/app/interaction/rpc/interaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeLogic) Like(req *types.LikeReq) (resp *types.LikeResp, err error) {
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

	// 打印追踪日志，让你在 API 黑窗口里能清楚看到收到了什么数据
	l.Logger.Infof("=> [API网关层] 接收到点赞请求: 解析出 UserId=%d, 目标 PostId=%d", userId, req.PostId)

	// 2. 组装参数，调用后端的 Interaction RPC 服务
	_, err = l.svcCtx.InteractionRpc.Like(l.ctx, &interaction.LikeRequest{
		UserId: userId,
		PostId: req.PostId,
	})

	if err != nil {
		l.Logger.Errorf("=> [API网关层] RPC返回拦截错误: %v", err)
		// 为了防止前端收到 null，如果出错，这里可以抛出一个直白的错误
		return nil, fmt.Errorf("操作被拦截: %v", err)
	}

	// 3. 返回成功响应
	return &types.LikeResp{}, nil
}
