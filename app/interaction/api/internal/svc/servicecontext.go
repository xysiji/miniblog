// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"miniblog/app/interaction/api/internal/config"
	"miniblog/app/interaction/rpc/interactionclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	InteractionRpc interactionclient.Interaction // 新增：RPC 客户端实例
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 初始化 RPC 客户端
		InteractionRpc: interactionclient.NewInteraction(zrpc.MustNewClient(c.InteractionRpc)),
	}
}
