package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	UserModel interface {
		userModel
		InsertWithId(ctx context.Context, data *User) (sql.Result, error)
		// 终极杀招：绕过被污染的缓存和框架 SQL，直接查全表
		FindUserByUsername(ctx context.Context, username string) (*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *customUserModel) InsertWithId(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`id`, `username`, `password`) values (?, ?, ?)", m.table)
	return m.ExecNoCacheCtx(ctx, query, data.Id, data.Username, data.Password)
}

// 终极杀招实现：使用 select * 确保拿出所有字段，并使用 QueryRowNoCacheCtx 彻底无视 Redis 缓存
func (m *customUserModel) FindUserByUsername(ctx context.Context, username string) (*User, error) {
	var resp User
	query := fmt.Sprintf("select * from %s where `username` = ? limit 1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, username)
	return &resp, err
}
