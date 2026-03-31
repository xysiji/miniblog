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
		FindPageListByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64, offset, pageSize int) ([]*Comment, error)
		CountByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64) (int64, error)
		FindAllByRootIdsShard(ctx context.Context, postId int64, rootIds []int64) ([]*Comment, error)
		SoftDeleteShard(ctx context.Context, postId int64, commentId int64) error
		// 新增：分表查询单条评论 (用于删除前的越权和防重校验)
		FindOneShard(ctx context.Context, postId int64, commentId int64) (*Comment, error)
	}

	customCommentModel struct {
		*defaultCommentModel
		readConn sqlx.SqlConn
	}
)

func NewCommentModel(conn sqlx.SqlConn, readConn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
		readConn:            readConn,
	}
}

func (m *customCommentModel) getShardTable(postId int64) string {
	shardIndex := postId % 4
	return fmt.Sprintf("`comment_%d`", shardIndex)
}

func (m *customCommentModel) InsertShard(ctx context.Context, data *Comment) (sql.Result, error) {
	tableName := m.getShardTable(data.PostId)
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", tableName, commentRowsExpectAutoSet)
	return m.ExecNoCacheCtx(ctx, query, data.Id, data.PostId, data.RootId, data.ParentId, data.UserId, data.ReplyToUserId, data.Content, data.Status)
}

func (m *customCommentModel) FindPageListByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64, offset, pageSize int) ([]*Comment, error) {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("select %s from %s where `post_id` = ? and `root_id` = ? and `status` = 1 order by create_time desc limit ?, ?", commentRows, tableName)

	var resp []*Comment
	err := m.readConn.QueryRowsCtx(ctx, &resp, query, postId, rootId, offset, pageSize)
	return resp, err
}

func (m *customCommentModel) CountByPostIdRootIdShard(ctx context.Context, postId int64, rootId int64) (int64, error) {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("select count(*) from %s where `post_id` = ? and `root_id` = ? and `status` = 1", tableName)

	var count int64
	err := m.readConn.QueryRowCtx(ctx, &count, query, postId, rootId)
	return count, err
}

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

func (m *customCommentModel) SoftDeleteShard(ctx context.Context, postId int64, commentId int64) error {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("update %s set `status` = 0 where `id` = ?", tableName)
	_, err := m.ExecNoCacheCtx(ctx, query, commentId)
	return err
}

// 核心实现：为了删除前校验，去正确的物理表里查数据
func (m *customCommentModel) FindOneShard(ctx context.Context, postId int64, commentId int64) (*Comment, error) {
	tableName := m.getShardTable(postId)
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", commentRows, tableName)
	var resp Comment
	// 查主库，防止从库延迟导致校验通过
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, commentId)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
