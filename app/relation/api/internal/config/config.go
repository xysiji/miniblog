package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// JWT 配置
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	// RPC 客户端配置
	RelationRpc zrpc.RpcClientConf
}
