package logic

import (
	"context"
	"fmt"

	"miniblog/app/user/rpc/internal/svc"
	"miniblog/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	// 查询数据库获取用户信息
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Errorf("查询用户信息失败: %v", err)
		return nil, fmt.Errorf("用户不存在")
	}

	return &user.UserInfoResponse{
		UserId:   userInfo.Id,
		Username: userInfo.Username,
		Avatar:   userInfo.Avatar, // 数据库里读取头像
		Bio:      userInfo.Bio,    // 数据库里读取简介
	}, nil
}
