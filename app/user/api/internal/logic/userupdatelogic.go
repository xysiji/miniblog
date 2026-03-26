// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"

	"miniblog/app/user/api/internal/svc"
	"miniblog/app/user/api/internal/types"

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
	// 假设从 JWT ctx 中获取了当前用户 ID
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	// 查询原用户信息
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, err
	}

	// 更新有变动的字段
	if req.Avatar != "" {
		userInfo.Avatar = req.Avatar
	}
	if req.Bio != "" {
		userInfo.Bio = req.Bio
	}

	err = l.svcCtx.UserModel.Update(l.ctx, userInfo)
	return &types.UserUpdateResp{}, err
}
