package svc

import (
	"miniblog/app/relation/model"
	"miniblog/app/relation/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 新增：关系表的 Model 实例
	RelationModel model.RelationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 连接
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config: c,
		// 注入数据库连接和 Redis 缓存配置
		RelationModel: model.NewRelationModel(conn, c.CacheRedis),
	}
}
