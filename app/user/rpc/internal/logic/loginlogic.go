package logic

import (
	"context"
	"fmt"

	"miniblog/app/user/rpc/internal/svc"
	"miniblog/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增：处理登录的 RPC 方法.
func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	// 1. 根据用户名从 MySQL 查询用户记录
	// 之前是：userInfo, err := l.svcCtx.UserModel.FindOneByUsername(...)
	// 修改为：
	userInfo, err := l.svcCtx.UserModel.FindUserByUsername(l.ctx, in.Username)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 2. 比对哈希密码 (微服务安全标准操作)
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(in.Password))
	if err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	// 3. 密码正确，返回 UserId 给 API 网关
	return &user.LoginResponse{
		UserId: userInfo.Id,
	}, nil
}
