package svc

import (
	"log"
	"miniblog/app/oss/api/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ServiceContext struct {
	Config      config.Config
	MinioClient *minio.Client // 新增：全局 MinIO 客户端实例
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MinIO 客户端
	minioClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyID, c.Minio.SecretAccessKey, ""),
		Secure: c.Minio.UseSSL,
	})
	if err != nil {
		log.Fatalf("初始化 MinIO 客户端失败: %v", err)
	}

	return &ServiceContext{
		Config:      c,
		MinioClient: minioClient,
	}
}
