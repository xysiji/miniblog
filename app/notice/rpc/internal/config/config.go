package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// 【新增】：声明数据库和缓存配置，与 yaml 文件对应
	DataSource string
	CacheRedis cache.CacheConf
}
