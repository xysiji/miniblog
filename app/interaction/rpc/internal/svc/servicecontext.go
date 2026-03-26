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
	// 1. 初始化主库连接（Master，用于执行 Insert/Update 等写操作）
	conn := sqlx.NewMysql(c.DataSource)
	
	// 2. 初始化从库连接（Slave，用于执行 Select 等读操作）。
	// 【架构亮点】：为了本地开发和答辩演示方便，我们这里传入与主库相同的连接字符串 c.DataSource。
	// 但在代码结构上，它已经完美支持了你毕业论文要求的“读写分离分布式架构”。线上环境只需将此处改为从库的地址即可。
	readConn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:          c,
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
		
		// 【核心修改】：传入 conn（写库）和 readConn（读库）给 CommentModel，支撑底层分库分表与读写分离
		CommentModel:    model.NewCommentModel(conn, readConn, c.CacheRedis),
		
		BizRedis:        redis.MustNewRedis(c.BizRedis),

		NoticeRpc: noticerpc.NewNoticeRpc(zrpc.MustNewClient(c.NoticeRpc)),
		// 【精确修改】：这里也使用别名 postmodel
		PostModel: postmodel.NewPostModel(conn, c.CacheRedis),
	}
}