// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userFieldNames          = builder.RawFieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheUserIdPrefix   = "cache:user:id:"
	cacheUserUuidPrefix = "cache:user:uuid:"
)

type (
	userModel interface {
		Insert(ctx context.Context, session sqlx.Session, datas ...*User) (sql.Result, error)
		FindOne(ctx context.Context, session sqlx.Session, id int64) (*User, error)
		FindOneByUuid(ctx context.Context, session sqlx.Session, uuid string) (*User, error)
		Update(ctx context.Context, session sqlx.Session, data *User) error
		Delete(ctx context.Context, session sqlx.Session, id ...int64) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id         int64  `db:"id"`
		Uuid       string `db:"uuid"`        // 用户uuid
		Email      string `db:"email"`       // 用户邮箱
		WalletAddr string `db:"wallet_addr"` // 用户钱包地址
		Address    string `db:"address"`     // titan的钱包地址
		CreatedAt  int64  `db:"created_at"`  // 创建时间
		DeletedAt  int64  `db:"deleted_at"`  // 删除时间
		Status     int64  `db:"status"`      // 提现状态:0-未提现 1-提现中 2-已提现
	}
)

func newUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultUserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`user`",
	}
}

func (m *defaultUserModel) Delete(ctx context.Context, session sqlx.Session, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}
	var keys []string
	for _, id := range ids {
		data, _ := m.FindOne(ctx, session, id)
		if data == nil {
			continue
		}
		userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
		userUuidKey := fmt.Sprintf("%s%v", cacheUserUuidPrefix, data.Uuid)
		keys = append(keys, userIdKey, userUuidKey)
	}
	vb, _ := json.Marshal(ids)
	values := strings.TrimSuffix(strings.TrimPrefix(string(vb), "["), "]")
	query := fmt.Sprintf("delete from %s where `id` IN (%s)", m.table, values)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query)
		}
		return conn.ExecCtx(ctx, query)
	}, keys...)
	return err
}

func (m *defaultUserModel) FindOne(ctx context.Context, session sqlx.Session, id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, userIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
		if session != nil {
			return session.QueryRowCtx(ctx, v, query, id)
		}
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultUserModel) FindOneByUuid(ctx context.Context, session sqlx.Session, uuid string) (*User, error) {
	userUuidKey := fmt.Sprintf("%s%v", cacheUserUuidPrefix, uuid)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userUuidKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `uuid` = ? limit 1", userRows, m.table)
		if session != nil {
			if err := session.QueryRowCtx(ctx, &resp, query, uuid); err != nil {
				return nil, err
			}
		}
		if err := conn.QueryRowCtx(ctx, &resp, query, uuid); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultUserModel) Insert(ctx context.Context, session sqlx.Session, datas ...*User) (sql.Result, error) {
	sq := squirrel.Insert(m.table).Columns(userRowsExpectAutoSet)

	var keys []string
	for _, data := range datas {
		sq = sq.Values(data.Uuid, data.Email, data.WalletAddr, data.Address, data.DeletedAt, data.Status)
		userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
		userUuidKey := fmt.Sprintf("%s%v", cacheUserUuidPrefix, data.Uuid)
		keys = append(keys, userIdKey, userUuidKey)
	}
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, keys...)

	return ret, err
}
func (m *defaultUserModel) Update(ctx context.Context, session sqlx.Session, newData *User) error {
	data, err := m.FindOne(ctx, session, newData.Id)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userUuidKey := fmt.Sprintf("%s%v", cacheUserUuidPrefix, data.Uuid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, newData.Uuid, newData.Email, newData.WalletAddr, newData.Address, newData.DeletedAt, newData.Status, newData.Id)
		}
		return conn.ExecCtx(ctx, query, newData.Uuid, newData.Email, newData.WalletAddr, newData.Address, newData.DeletedAt, newData.Status, newData.Id)
	}, userIdKey, userUuidKey)
	return err
}

func (m *defaultUserModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserModel) tableName() string {
	return m.table
}
