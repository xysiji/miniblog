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

type CommentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentDeleteLogic {
	return &CommentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentDeleteLogic) CommentDelete(req *types.CommentDeleteReq) (resp *types.CommentDeleteResp, err error) {
	// 1. 【终极兼容版】解析 JWT 获取 userId (用于底层鉴权，确保只能删自己的)
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

	l.Logger.Infof("=> [API网关层] 接收到删除评论请求: 操作人 UserId=%d, 目标 CommentId=%d", userId, req.CommentId)

	// 2. 调用底层的 Interaction RPC 服务
	_, err = l.svcCtx.InteractionRpc.CommentDelete(l.ctx, &interaction.CommentDeleteRequest{
		UserId:    userId,
		CommentId: req.CommentId,
	})

	if err != nil {
		l.Logger.Errorf("=> [API网关层] 删除评论RPC返回拦截错误: %v", err)
		return nil, fmt.Errorf("操作被拦截: %v", err)
	}

	return &types.CommentDeleteResp{}, nil
}
