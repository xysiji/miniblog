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
		// 【新增】：一键将用户的所有未读通知置为已读
		ReadAllByUserId(ctx context.Context, userId int64) error
		// 【新增】：单条已读
		ReadById(ctx context.Context, id int64) error
	}

	customNoticeModel struct {
		*defaultNoticeModel
	}
)

// NewNoticeModel returns a model for the database table.
func NewNoticeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NoticeModel {
	return &customNoticeModel{
		defaultNoticeModel: newNoticeModel(conn, c, opts...),
	}
}

// 实现查询列表的具体 SQL
func (m *customNoticeModel) FindListByUserId(ctx context.Context, userId int64, page, pageSize int) ([]*Notice, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by create_time desc limit ?, ?", noticeRows, m.table)
	var resp []*Notice
	offset := (page - 1) * pageSize
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, offset, pageSize)
	return resp, err
}

// 实现查询总数的具体 SQL
func (m *customNoticeModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, userId)
	return count, err
}

// 实现一键已读的具体 SQL
func (m *customNoticeModel) ReadAllByUserId(ctx context.Context, userId int64) error {
	// is_read = 0 表示未读，将其更新为 1 已读
	query := fmt.Sprintf("update %s set `is_read` = 1 where `user_id` = ? and `is_read` = 0", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, userId)
	return err
}

// 实现单条已读的具体 SQL
func (m *customNoticeModel) ReadById(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set `is_read` = 1 where `id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, id)
	return err
}
