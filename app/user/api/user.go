// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"miniblog/app/user/api/internal/config"
	"miniblog/app/user/api/internal/handler"
	"miniblog/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	//	server := rest.MustNewServer(c.RestConf)
	// 加上 rest.WithCors() 允许所有跨域请求
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
