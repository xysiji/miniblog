// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"miniblog/app/interaction/rpc/interactionclient"
	"miniblog/app/post/api/internal/config"
	"miniblog/app/post/rpc/postclient"
	"miniblog/app/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	PostRpc        postclient.Post
	UserRpc        userclient.User               // 【新增】用户服务客户端实例
	InteractionRpc interactionclient.Interaction // 【新增】互动服务客户端实例
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 初始化三个底层的 RPC 客户端，供 logic 层调用
		PostRpc:        postclient.NewPost(zrpc.MustNewClient(c.PostRpc)),
		UserRpc:        userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),                      // 【新增】
		InteractionRpc: interactionclient.NewInteraction(zrpc.MustNewClient(c.InteractionRpc)), // 【新增】
	}
}
