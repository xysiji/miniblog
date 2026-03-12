package logic

import (
	"context"
	"fmt"

	"miniblog/app/user/model"
	"miniblog/app/user/rpc/internal/svc"
	"miniblog/app/user/rpc/user"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 1. 检查用户名是否已存在 (这里展示了如何调用 model 层)
	_, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 2. 密码加密 (Bcrypt 算法)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 3. 生成雪花算法全局唯一 ID (分布式架构标配)
	node, err := snowflake.NewNode(1) // 节点ID为1，生产环境应从配置读取
	if err != nil {
		return nil, fmt.Errorf("生成用户ID失败: %v", err)
	}
	userId := node.Generate().Int64()

	// 4. 组装数据并写入 MySQL
	newUser := &model.User{
		Id:       userId,
		Username: in.Username,
		Password: string(hashPassword),
	}
	// 之前是：_, err = l.svcCtx.UserModel.Insert(l.ctx, newUser)
	// 修改为：
	_, err = l.svcCtx.UserModel.InsertWithId(l.ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("插入数据库失败: %v", err)
	}

	// 5. 返回生成的 UserID
	return &user.RegisterResponse{
		UserId: userId,
	}, nil
}
