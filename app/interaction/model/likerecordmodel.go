package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeRecordModel = (*customLikeRecordModel)(nil)

type (
	// LikeRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeRecordModel.
	LikeRecordModel interface {
		likeRecordModel
		// 新增：获取用户所有处于点赞状态的博文ID
		FindLikedPostIdsByUserId(ctx context.Context, userId int64) ([]int64, error)
	}

	customLikeRecordModel struct {
		*defaultLikeRecordModel
	}
)

// NewLikeRecordModel returns a model for the database table.
func NewLikeRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeRecordModel {
	return &customLikeRecordModel{
		defaultLikeRecordModel: newLikeRecordModel(conn, c, opts...),
	}
}

// 实现获取用户点赞列表的具体 SQL
func (m *customLikeRecordModel) FindLikedPostIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	// 只需要查出 post_id，且状态 status = 1 表示已点赞
	query := fmt.Sprintf("select post_id from %s where `user_id` = ? and `status` = 1", m.table)

	var resp []int64
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
