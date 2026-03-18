package svc

import (
	"miniblog/app/post/model"
	"miniblog/app/post/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 【分布式存储重构】：分离读写模型
	PostMasterModel model.PostModel // 写模型 (连主库)
	PostSlaveModel  model.PostModel // 读模型 (连从库)
	BizRedis        *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 初始化分布式数据库的主从连接池
	masterConn := sqlx.NewMysql(c.MasterDataSource)
	slaveConn := sqlx.NewMysql(c.SlaveDataSource)

	return &ServiceContext{
		Config: c,
		// 2. 将主从连接分别注入到对应的 Model 实例中
		PostMasterModel: model.NewPostModel(masterConn, c.CacheRedis),
		PostSlaveModel:  model.NewPostModel(slaveConn, c.CacheRedis),
		BizRedis:        redis.MustNewRedis(c.BizRedis),
	}
}
