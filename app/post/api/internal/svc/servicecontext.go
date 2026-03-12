// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"miniblog/app/post/api/internal/config"
	"miniblog/app/post/rpc/postclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	PostRpc postclient.Post // 新增 RPC 客户端
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 初始化 RPC 客户端
		PostRpc: postclient.NewPost(zrpc.MustNewClient(c.PostRpc)),
	}
}
