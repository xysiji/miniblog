package svc

import (
	"miniblog/app/interaction/model"
	"miniblog/app/interaction/rpc/internal/config"
	"miniblog/app/notice/rpc/noticerpc"

	// 【精确修改】：明确引入 post 服务的 model，并强制命名别名为 postmodel
	postmodel "miniblog/app/post/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	LikeRecordModel model.LikeRecordModel
	CommentModel    model.CommentModel
	BizRedis        *redis.Redis

	NoticeRpc noticerpc.NoticeRpc
	// 【精确修改】：这里使用别名 postmodel
	PostModel postmodel.PostModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:          c,
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
		CommentModel:    model.NewCommentModel(conn, c.CacheRedis),
		BizRedis:        redis.MustNewRedis(c.BizRedis),

		NoticeRpc: noticerpc.NewNoticeRpc(zrpc.MustNewClient(c.NoticeRpc)),
		// 【精确修改】：这里也使用别名 postmodel
		PostModel: postmodel.NewPostModel(conn, c.CacheRedis),
	}
}
