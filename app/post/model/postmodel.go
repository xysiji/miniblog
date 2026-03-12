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

		// --- B阶段新增：分页查询相关底层方法 ---
		// 1. 查询某一页的博文列表 (按时间倒序)
		FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Post, error)
		// 2. 查询博文总条数 (用于分页计算)
		Count(ctx context.Context) (int64, error)
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
	query := fmt.Sprintf("insert into %s (`id`, `user_id`, `content`) values (?, ?, ?)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.Id, data.UserId, data.Content)
}

// B阶段新增实现：利用 QueryRowsNoCacheCtx 绕过行级缓存，直接查 MySQL 列表
func (m *customPostModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Post, error) {
	// 核心逻辑：计算分页偏移量 (Offset)
	// 如果是第 1 页，每页 10 条，跳过 (1-1)*10 = 0 条
	// 如果是第 2 页，每页 10 条，跳过 (2-1)*10 = 10 条
	offset := (page - 1) * pageSize

	// 编写原生 SQL：按 ID 倒序（最新的在最前），限制条数和偏移量
	query := fmt.Sprintf("select %s from %s order by id desc limit ? offset ?", postRows, m.table)

	var resp []*Post
	// 注意：查询多行数据必须用 QueryRowsNoCacheCtx
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, pageSize, offset)
	return resp, err
}

// B阶段新增实现：查询数据总数
func (m *customPostModel) Count(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s", m.table)
	var count int64
	// 注意：查询单行单列数据用 QueryRowNoCacheCtx
	err := m.QueryRowNoCacheCtx(ctx, &count, query)
	return count, err
}
