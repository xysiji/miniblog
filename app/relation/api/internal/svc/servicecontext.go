package svc

import (
	"miniblog/app/relation/api/internal/config"
	"miniblog/app/relation/rpc/relationclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	// 挂载 RPC 客户端实例
	RelationRpc relationclient.Relation
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		RelationRpc: relationclient.NewRelation(zrpc.MustNewClient(c.RelationRpc)),
	}
}
