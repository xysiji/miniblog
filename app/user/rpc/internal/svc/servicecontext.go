package svc

import (
	"miniblog/app/user/model" // 引入刚刚生成的 model 包
	"miniblog/app/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel // 注入 UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 连接
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config: c,
		// 实例化带有 Redis 缓存的 UserModel
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
