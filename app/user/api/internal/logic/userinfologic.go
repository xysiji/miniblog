// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/user/api/internal/svc"
	"miniblog/app/user/api/internal/types"
	"miniblog/app/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户信息
func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// 1. 从 JWT 中安全提取 UserId
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

	// 2. 呼叫底层 RPC
	rpcResp, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userclient.UserInfoRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	// 3. 组装给前端的 HTTP JSON
	return &types.UserInfoResp{
		UserId:   rpcResp.UserId,
		Username: rpcResp.Username,
		Avatar:   rpcResp.Avatar,
		Bio:      rpcResp.Bio,
	}, nil
}
