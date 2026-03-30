package logic

import (
	"context"
	"fmt"

	"miniblog/app/user/rpc/internal/svc"
	"miniblog/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserUpdateLogic) UserUpdate(in *user.UserUpdateRequest) (*user.UserUpdateResponse, error) {
	// 1. 先查出原有的用户信息
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 2. 只有当前端传了对应的值才更新（按需更新）
	if in.Avatar != "" {
		userInfo.Avatar = in.Avatar
	}
	if in.Bio != "" {
		userInfo.Bio = in.Bio
	}

	// 3. 执行数据库更新操作
	err = l.svcCtx.UserModel.Update(l.ctx, userInfo)
	if err != nil {
		l.Logger.Errorf("更新用户资料失败: %v", err)
		return nil, fmt.Errorf("资料更新失败，请稍后重试")
	}

	return &user.UserUpdateResponse{}, nil
}
