// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	PostRpc zrpc.RpcClientConf
	Auth    struct { // 映射 JWT 配置
		AccessSecret string
		AccessExpire int64
	}
}
