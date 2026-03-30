package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	CommentModel interface {
		commentModel
		InsertShard(ctx context.Context, data *Comment) (sql.Result, error)
		// 改造：原有的查询改为只查一级主评论 (root_id = 0)
		FindPageListByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64, offset, pageSize int) ([]*Comment, error)
		// 改造：原有的统计改为只统计一级主评论 (root_id = 0)
		CountByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64) (int64, error)
		// 新增：防 N+1 核心，批量查询子评论（必须传入 postId 以便找到正确的分表）
		FindAllByRootIdsShard(ctx context.Context, postId int64, rootIds []int64) ([]*Comment, error)
	}

	customCommentModel struct {
		*defaultCommentModel
		readConn sqlx.SqlConn // 读库连接（Slave）
	}
)

func NewCommentModel(conn sqlx.SqlConn, readConn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
		readConn:            readConn,
	}
}

// 核心算法：根据 PostId 进行 Hash 路由分表
func (m *customCommentModel) getShardTable(postId int64) string {
	shardIndex := postId % 4 // 分为4张表: comment_0, comment_1, comment_2, comment_3
	return fmt.Sprintf("`comment_%d`", shardIndex)
}

// 1. 分表写入 (走主库)
func (m *customCommentModel) InsertShard(ctx context.Context, data *Comment) (sql.Result, error) {
	tableName := m.getShardTable(data.PostId)
	// 补充了嵌套所需的 root_id, parent_id, reply_to_user_id, status 字段
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", tableName, commentRowsExpectAutoSet)
	return m.ExecNoCacheCtx(ctx, query, data.Id, data.PostId, data.RootId, data.ParentId, data.UserId, data.ReplyToUserId, data.Content, data.Status)
}

// 2. 分表读取主评论列表 (走从库)
func (m *customCommentModel) FindPageListByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64, offset, pageSize int) ([]*Comment, error) {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("select %s from %s where `post_id` = ? and `root_id` = ? and `status` = 1 order by create_time desc limit ?, ?", commentRows, tableName)

	var resp []*Comment
	err := m.readConn.QueryRowsCtx(ctx, &resp, query, postId, rootId, offset, pageSize)
	return resp, err
}

// 3. 分表读取主评论总数 (走从库)
func (m *customCommentModel) CountByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64) (int64, error) {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("select count(*) from %s where `post_id` = ? and `root_id` = ? and `status` = 1", tableName)

	var count int64
	err := m.readConn.QueryRowCtx(ctx, &count, query, postId, rootId)
	return count, err
}

// 4. 新增：分表批量读取子评论 (走从库)
func (m *customCommentModel) FindAllByRootIdsShard(ctx context.Context, postId int64, rootIds []int64) ([]*Comment, error) {
	if len(rootIds) == 0 {
		return nil, nil
	}
	tableName := m.getShardTable(postId)

	var placeholders []string
	var args []interface{}
	for _, id := range rootIds {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf("select %s from %s where `root_id` in (%s) and `status` = 1 order by id asc", commentRows, tableName, strings.Join(placeholders, ","))
	var resp []*Comment
	err := m.readConn.QueryRowsCtx(ctx, &resp, query, args...)
	return resp, err
}
