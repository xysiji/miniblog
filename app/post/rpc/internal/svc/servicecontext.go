package svc

import (
	"miniblog/app/post/model" // 引入之前生成的 post model
	"miniblog/app/post/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	PostModel model.PostModel // 新增模型实例
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 连接
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config: c,
		// 实例化带有 Redis 多级缓存的 PostModel
		PostModel: model.NewPostModel(conn, c.CacheRedis),
	}
}
