// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	NoticeRpc zrpc.RpcClientConf // 【保留】：RPC 配置

	DataSource string          // 【新增】：直连数据库
	CacheRedis cache.CacheConf // 【新增】：Redis 缓存
}
