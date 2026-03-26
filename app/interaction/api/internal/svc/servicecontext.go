package svc

import (
	"miniblog/app/interaction/api/internal/config"
	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/interactionclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	InteractionRpc  interactionclient.Interaction
	LikeRecordModel model.LikeRecordModel // 新增
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource) // 新增建立连接
	return &ServiceContext{
		Config:          c,
		InteractionRpc:  interactionclient.NewInteraction(zrpc.MustNewClient(c.InteractionRpc)),
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis), // 新增挂载模型
	}
}
