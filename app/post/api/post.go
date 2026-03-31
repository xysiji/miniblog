// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"miniblog/app/post/api/internal/config"
	"miniblog/app/post/api/internal/handler"
	"miniblog/app/post/api/internal/middleware" // 新增引入
	"miniblog/app/post/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/post-api.yaml", "the config file")

// AccessLog 定义我们用于数据分析的结构化日志格式
type AccessLog struct {
	Timestamp string `json:"timestamp"` // 访问时间
	Method    string `json:"method"`    // 请求方法 (GET/POST)
	Path      string `json:"path"`      // 请求路径
	Duration  int64  `json:"duration"`  // 接口耗时 (微秒)，用于分析多级缓存的性能
}

// DataAnalysisMiddleware 全局数据分析拦截器 (埋点)
func DataAnalysisMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		duration := time.Since(startTime).Microseconds()

		logData := AccessLog{
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Method:    r.Method,
			Path:      r.URL.Path,
			Duration:  duration,
		}

		go func(data AccessLog) {
			logBytes, _ := json.Marshal(data)
			file, err := os.OpenFile("data_analysis.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				defer file.Close()
				file.WriteString(string(logBytes) + "\n")
			}
		}(logData)
	}
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	// ===============================================
	// 1. 挂载数据分析埋点中间件
	server.Use(DataAnalysisMiddleware)

	// 2. 挂载我们新写的全局弱认证中间件 (解析 Token 但不阻拦)
	softAuth := middleware.NewSoftAuthMiddleware(c.Auth.AccessSecret)
	server.Use(softAuth.Handle)
	// ===============================================

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
