package svc

import (
	"miniblog/app/notice/model" // 引入我们刚刚生成的通知数据库模型
	"miniblog/app/notice/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	NoticeModel model.NoticeModel // 【新增】：挂载 Notice 的数据库模型
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 连接
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:      c,
		NoticeModel: model.NewNoticeModel(conn, c.CacheRedis), // 【新增】：实例化模型
	}
}
