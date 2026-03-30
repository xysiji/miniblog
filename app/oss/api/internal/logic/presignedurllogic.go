package logic

import (
	"context"
	"fmt"
	"time"

	"miniblog/app/oss/api/internal/svc"
	"miniblog/app/oss/api/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type PresignedUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPresignedUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PresignedUrlLogic {
	return &PresignedUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PresignedUrlLogic) PresignedUrl(req *types.PresignedUrlReq) (resp *types.PresignedUrlResp, err error) {
	// 1. 规范化桶名 (Bucket)
	bucketName := "miniblog" // 我们将所有博客文件统一放在 miniblog 这个桶里

	// 2. 根据前端传入的 Type 划分子目录
	var dir string
	switch req.Type {
	case 1:
		dir = "avatar" // 头像目录
	case 2:
		dir = "post" // 帖子配图目录
	default:
		dir = "other"
	}

	// 3. 生成全局唯一的随机文件名，防止文件被覆盖
	// 例: avatar/f47ac10b-58cc-0372-8567-0e02b2c3d479.png
	objectName := fmt.Sprintf("%s/%s%s", dir, uuid.New().String(), req.FileExt)

	// 4. 设置预签名链接的过期时间 (例如 10 分钟内上传有效，过期作废)
	expires := time.Minute * 10

	// 5. 向 MinIO 请求生成直传 PUT 链接
	presignedURL, err := l.svcCtx.MinioClient.PresignedPutObject(
		l.ctx,
		bucketName,
		objectName,
		expires,
	)
	if err != nil {
		l.Logger.Errorf("生成 MinIO 预签名链接失败: %v", err)
		return nil, err
	}

	// 6. 拼接文件上传成功后的永久下载/访问地址
	// 格式: http://{Endpoint}/{BucketName}/{ObjectName}
	downloadUrl := fmt.Sprintf("http://%s/%s/%s", l.svcCtx.Config.Minio.Endpoint, bucketName, objectName)

	return &types.PresignedUrlResp{
		UploadUrl:   presignedURL.String(),
		DownloadUrl: downloadUrl,
	}, nil
}
