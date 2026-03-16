package svc

import (
	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	LikeRecordModel model.LikeRecordModel
	CommentModel    model.CommentModel
	BizRedis        *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 连接
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:          c,
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
		CommentModel:    model.NewCommentModel(conn, c.CacheRedis),
		BizRedis:        redis.MustNewRedis(c.BizRedis),
	}
}
