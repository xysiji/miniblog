package svc

import (
	"miniblog/app/relation/rpc/relationclient" // 引入 relation 的 client
	"miniblog/app/user/api/internal/config"
	"miniblog/app/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User
	// 新增：挂载实例
	RelationRpc relationclient.Relation
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		UserRpc:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		RelationRpc: relationclient.NewRelation(zrpc.MustNewClient(c.RelationRpc)),
	}
}
