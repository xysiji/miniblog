package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	CommentModel interface {
		commentModel
		// 覆盖基础方法，支持分库分表与读写分离
		InsertShard(ctx context.Context, data *Comment) (sql.Result, error)
		FindPageListByPostIdShard(ctx context.Context, postId int64, page, pageSize int) ([]*Comment, error)
		CountByPostIdShard(ctx context.Context, postId int64) (int64, error)
	}

	customCommentModel struct {
		*defaultCommentModel
		readConn sqlx.SqlConn // 读库连接（Slave）
	}
)

// 注入 readConn，实现读写分离
func NewCommentModel(conn sqlx.SqlConn, readConn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
		readConn:            readConn,
	}
}

// ==========================================
// 核心算法：根据 PostId 进行 Hash 路由分表
// ==========================================
func (m *customCommentModel) getShardTable(postId int64) string {
	shardIndex := postId % 4 // 分为4张表: comment_0, comment_1, comment_2, comment_3
	return fmt.Sprintf("`comment_%d`", shardIndex)
}

// 1. 分表写入 (走主库)
func (m *customCommentModel) InsertShard(ctx context.Context, data *Comment) (sql.Result, error) {
	tableName := m.getShardTable(data.PostId)
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", tableName, commentRowsExpectAutoSet)

	// 【核心修复】：go-zero 的缓存模型没有暴露 conn，必须使用框架提供的 ExecNoCacheCtx 来操作主库
	return m.ExecNoCacheCtx(ctx, query, data.Id, data.PostId, data.UserId, data.Content)
}

// 2. 分表读取列表 (走从库 m.readConn)
func (m *customCommentModel) FindPageListByPostIdShard(ctx context.Context, postId int64, page, pageSize int) ([]*Comment, error) {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("select %s from %s where `post_id` = ? order by create_time desc limit ?, ?", commentRows, tableName)

	var resp []*Comment
	offset := (page - 1) * pageSize

	// 读操作使用读库 readConn
	err := m.readConn.QueryRowsCtx(ctx, &resp, query, postId, offset, pageSize)
	return resp, err
}

// 3. 分表读取总数 (走从库 m.readConn)
func (m *customCommentModel) CountByPostIdShard(ctx context.Context, postId int64) (int64, error) {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("select count(*) from %s where `post_id` = ?", tableName)

	var count int64

	// 读操作使用读库 readConn
	err := m.readConn.QueryRowCtx(ctx, &count, query, postId)
	return count, err
}
