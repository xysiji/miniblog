// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	// JWT 鉴权配置 (因为我们的 api 文件里写了 jwt: Auth)
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	// 新增：MinIO 配置项
	Minio struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
	}
}
