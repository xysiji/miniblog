package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string          // MySQL 链接
	CacheRedis cache.CacheConf // go-zero 内置行缓存
	BizRedis   redis.RedisConf // 业务 Redis
}