package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// 【分布式存储重构】：主从读写分离架构
	MasterDataSource string          // 主库：专门负责写入 (Publish)
	SlaveDataSource  string          // 从库：专门负责读取 (List)
	CacheRedis       cache.CacheConf // go-zero 内置行缓存 (支持 Redis 集群)
	BizRedis         redis.RedisConf // 业务 Redis (Timeline 缓存)
}
