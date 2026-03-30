package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PostModel = (*customPostModel)(nil)

type (
	PostModel interface {
		postModel
		InsertWithId(ctx context.Context, data *Post) (sql.Result, error)

		// --- 分页查询相关底层方法 ---
		FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Post, error)
		Count(ctx context.Context) (int64, error)

		// --- 终极闭环：统计数原子增减（自带缓存清理） ---
		IncrLikeCount(ctx context.Context, id int64) error
		DecrLikeCount(ctx context.Context, id int64) error
		IncrCommentCount(ctx context.Context, id int64) error
		DecrCommentCount(ctx context.Context, id int64) error
	}

	customPostModel struct {
		*defaultPostModel
	}
)

func NewPostModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PostModel {
	return &customPostModel{
		defaultPostModel: newPostModel(conn, c, opts...),
	}
}

func (m *customPostModel) InsertWithId(ctx context.Context, data *Post) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`id`, `user_id`, `content`, `images`) values (?, ?, ?, ?)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.Id, data.UserId, data.Content, data.Images)
}

func (m *customPostModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Post, error) {
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select %s from %s order by id desc limit ? offset ?", postRows, m.table)
	var resp []*Post
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, pageSize, offset)
	return resp, err
}

func (m *customPostModel) Count(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query)
	return count, err
}

// ===============================================
// 以下为点赞数、评论数原子操作，ExecCtx 会自动清理缓存
// ===============================================

func (m *customPostModel) IncrLikeCount(ctx context.Context, id int64) error {
	postIdKey := fmt.Sprintf("%s%v", cachePostIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set like_count = like_count + 1 where id = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, postIdKey)
	return err
}

func (m *customPostModel) DecrLikeCount(ctx context.Context, id int64) error {
	postIdKey := fmt.Sprintf("%s%v", cachePostIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		// 保证点赞数不出现负数
		query := fmt.Sprintf("update %s set like_count = like_count - 1 where id = ? and like_count > 0", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, postIdKey)
	return err
}

func (m *customPostModel) IncrCommentCount(ctx context.Context, id int64) error {
	postIdKey := fmt.Sprintf("%s%v", cachePostIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set comment_count = comment_count + 1 where id = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, postIdKey)
	return err
}

func (m *customPostModel) DecrCommentCount(ctx context.Context, id int64) error {
	postIdKey := fmt.Sprintf("%s%v", cachePostIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set comment_count = comment_count - 1 where id = ? and comment_count > 0", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, postIdKey)
	return err
}
