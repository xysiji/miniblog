package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoticeModel = (*customNoticeModel)(nil)

type (
	// NoticeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNoticeModel.
	NoticeModel interface {
		noticeModel

		// 【新增】：根据用户ID查询通知列表 (支持分页)
		FindListByUserId(ctx context.Context, userId int64, page, pageSize int) ([]*Notice, error)
		// 【新增】：查询某个用户的通知总数 (用于前端分页组件)
		CountByUserId(ctx context.Context, userId int64) (int64, error)
	}

	customNoticeModel struct {
		*defaultNoticeModel
	}
)

// NewNoticeModel returns a model for the database table.
// 【⚠️精确修复】：保留你原版代码中的 opts ...cache.Option，防止编译报错
func NewNoticeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NoticeModel {
	return &customNoticeModel{
		defaultNoticeModel: newNoticeModel(conn, c, opts...),
	}
}

// 【新增】：实现查询列表的具体 SQL
func (m *customNoticeModel) FindListByUserId(ctx context.Context, userId int64, page, pageSize int) ([]*Notice, error) {
	// noticeRows 变量在 _gen.go 中已自动生成，包含了所有字段名
	// 按 create_time 降序排列，最新的通知在最前面
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by create_time desc limit ?, ?", noticeRows, m.table)

	var resp []*Notice
	// 计算分页 offset
	offset := (page - 1) * pageSize

	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, offset, pageSize)
	return resp, err
}

// 【新增】：实现查询总数的具体 SQL
func (m *customNoticeModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, userId)
	return count, err
}
