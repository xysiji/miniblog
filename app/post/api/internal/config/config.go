// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	PostRpc        zrpc.RpcClientConf
	UserRpc        zrpc.RpcClientConf // 【新增】映射 YAML 里的 UserRpc
	InteractionRpc zrpc.RpcClientConf // 【新增】映射 YAML 里的 InteractionRpc
	Auth           struct {           // 映射 JWT 配置
		AccessSecret string
		AccessExpire int64
	}
}
