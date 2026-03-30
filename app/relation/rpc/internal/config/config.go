package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// 新增：数据源和缓存配置
	DataSource string
	CacheRedis cache.CacheConf
}
