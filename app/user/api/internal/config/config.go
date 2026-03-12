// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpc zrpc.RpcClientConf // 新增：RPC 客户端配置
	// 新增 JWT 配置结构体映射
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
