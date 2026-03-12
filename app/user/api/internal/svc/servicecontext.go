// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"miniblog/app/user/api/internal/config"
	"miniblog/app/user/rpc/userclient" // 引入生成的 rpc 客户端包

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User // 新增：RPC 客户端接口
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 初始化 RPC 客户端并注入
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
