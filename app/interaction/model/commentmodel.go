package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	CommentModel interface {
		commentModel
		// 新增：根据 PostId 分页查询评论列表 (按时间倒序)
		FindPageListByPostId(ctx context.Context, postId int64, page, pageSize int64) ([]*Comment, error)
		// 新增：查询某篇博文的评论总数
		CountByPostId(ctx context.Context, postId int64) (int64, error)
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

func NewCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
	}
}

// ==========================================
// 自定义方法实现区
// ==========================================

// FindPageListByPostId 根据博文ID拉取评论列表
func (m *customCommentModel) FindPageListByPostId(ctx context.Context, postId int64, page, pageSize int64) ([]*Comment, error) {
	offset := (page - 1) * pageSize
	// 原生 SQL：按 id 倒序（最新评论在最上面），限制查某篇具体的博文
	query := fmt.Sprintf("select %s from %s where post_id = ? order by id desc limit ? offset ?", commentRows, m.table)

	var resp []*Comment
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, postId, pageSize, offset)
	return resp, err
}

// CountByPostId 统计评论总数
func (m *customCommentModel) CountByPostId(ctx context.Context, postId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where post_id = ?", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, postId)
	return count, err
}
