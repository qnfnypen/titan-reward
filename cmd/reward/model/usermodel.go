package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		QueryBuilder(fields ...string) squirrel.SelectBuilder
		FindOneByEmail(ctx context.Context, email string) (*User, error)
		FindOneByWalletAddr(ctx context.Context, walletAddr string) (*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *customUserModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customUserModel) QueryBuilder(fields ...string) squirrel.SelectBuilder {
	var queryFields string

	if len(fields) <= 0 {
		queryFields = "*"
	} else {
		queryFields = strings.Join(fields, ",")
	}

	return squirrel.Select(queryFields).From(m.table)
}

// FindOneByEmail 通过email获取用户信息
func (m *customUserModel) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	var resp User

	query, args, err := m.QueryBuilder().Where("email = ? limit 1", email).ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate sql of get info by email error:%w", err)
	}

	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	case nil:
		return &resp, nil
	default:
		return nil, fmt.Errorf("get info by email error:%w", err)
	}
}

// FindOneByWalletAddr 通过小狐狸钱包地址获取用户信息
func (m *customUserModel) FindOneByWalletAddr(ctx context.Context, walletAddr string) (*User, error) {
	var resp User

	query, args, err := m.QueryBuilder().Where("wallet_addr = ? limit 1", walletAddr).ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate sql of get info by wallet_addr error:%w", err)
	}

	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	case nil:
		return &resp, nil
	default:
		return nil, fmt.Errorf("get info by wallet_addr error:%w", err)
	}
}
