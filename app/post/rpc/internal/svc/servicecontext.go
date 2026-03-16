package svc

import (
	"miniblog/app/post/model"
	"miniblog/app/post/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis" // 新增导入
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	PostModel model.PostModel
	BizRedis  *redis.Redis // 新增：业务 Redis 客户端实例
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:    c,
		PostModel: model.NewPostModel(conn, c.CacheRedis),
		// 初始化并注入业务 Redis 客户端
		BizRedis: redis.MustNewRedis(c.BizRedis),
	}
}
