// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"miniblog/app/notice/api/internal/config"
	"miniblog/app/notice/model"
	"miniblog/app/notice/rpc/noticerpc" // 引入 RPC 客户端代码

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	NoticeModel model.NoticeModel   // 【新增】：挂载数据库模型（读数据用）
	NoticeRpc   noticerpc.NoticeRpc // 【保留】：挂载 RPC 客户端（未来写数据用）
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:      c,
		NoticeModel: model.NewNoticeModel(conn, c.CacheRedis),
		NoticeRpc:   noticerpc.NewNoticeRpc(zrpc.MustNewClient(c.NoticeRpc)),
	}
}
