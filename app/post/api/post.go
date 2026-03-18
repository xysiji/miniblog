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
		// 1. 记录请求开始时间
		startTime := time.Now()

		// 2. 放行请求，执行实际的业务逻辑 (比如拉取列表、发布博文)
		next(w, r)

		// 3. 计算接口耗时 (微秒)
		duration := time.Since(startTime).Microseconds()

		// 4. 组装日志数据
		logData := AccessLog{
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Method:    r.Method,
			Path:      r.URL.Path,
			Duration:  duration,
		}

		// 5. 将日志转为 JSON 并异步写入本地文件 (不阻塞正常请求)
		go func(data AccessLog) {
			logBytes, _ := json.Marshal(data)
			// 打开或创建埋点日志文件 (追加模式)
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

	//server := rest.MustNewServer(c.RestConf)
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	// ===============================================
	// 新增：挂载全局数据分析中间件 (所有经过 post-api 的请求都会被埋点采集)
	// ===============================================
	server.Use(DataAnalysisMiddleware)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
