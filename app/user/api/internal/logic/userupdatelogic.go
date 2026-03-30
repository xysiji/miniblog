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

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户信息
func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) (resp *types.UserUpdateResp, err error) {
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

	// 2. 呼叫底层 RPC 执行更新
	_, err = l.svcCtx.UserRpc.UserUpdate(l.ctx, &userclient.UserUpdateRequest{
		UserId: userId,
		Avatar: req.Avatar,
		Bio:    req.Bio,
	})

	if err != nil {
		return nil, err
	}

	return &types.UserUpdateResp{}, nil
}
