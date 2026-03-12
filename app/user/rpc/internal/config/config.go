package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string          // 对应 yaml 的 DataSource
	CacheRedis cache.CacheConf // 对应 yaml 的 CacheRedis
}
