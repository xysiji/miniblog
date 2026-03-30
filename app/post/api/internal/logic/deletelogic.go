// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/post/api/internal/svc"
	"miniblog/app/post/api/internal/types"
	"miniblog/app/post/rpc/postclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除微型博客
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteReq) (resp *types.DeleteResp, err error) {
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

	// 2. 调用 RPC 删除
	_, err = l.svcCtx.PostRpc.Delete(l.ctx, &postclient.DeleteRequest{
		PostId: req.PostId,
		UserId: userId,
	})

	if err != nil {
		l.Logger.Errorf("网关层删除博文失败: %v", err)
		return nil, err
	}

	return &types.DeleteResp{}, nil
}
