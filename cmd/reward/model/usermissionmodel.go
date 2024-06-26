package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserMissionModel = (*customUserMissionModel)(nil)

type (
	// UserMissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserMissionModel.
	UserMissionModel interface {
		userMissionModel
		withSession(session sqlx.Session) UserMissionModel
		GetCreditByUn(ctx context.Context, un string) (int64, error)
	}

	customUserMissionModel struct {
		*defaultUserMissionModel
	}
)

// NewUserMissionModel returns a model for the database table.
func NewUserMissionModel(conn sqlx.SqlConn) UserMissionModel {
	return &customUserMissionModel{
		defaultUserMissionModel: newUserMissionModel(conn),
	}
}

func (m *customUserMissionModel) withSession(session sqlx.Session) UserMissionModel {
	return NewUserMissionModel(sqlx.NewSqlConnFromSession(session))
}

// GetCreditByUn 获取用户社区的任务收益
func (m *customUserMissionModel) GetCreditByUn(ctx context.Context, un string) (int64, error) {
	var credit int64

	if strings.TrimSpace(un) == "" {
		return 0, nil
	}

	query, args, err := squirrel.Select("IFNULL(SUM(credit),0)").From(m.table).Where("username = ?", un).ToSql()
	if err != nil {
		return credit, fmt.Errorf("get sum of credit error:%w", err)
	}

	err = m.conn.QueryRowCtx(ctx, &credit, query, args...)
	switch err {
	case sqlc.ErrNotFound:
		return 0, nil
	case nil:
		return credit, nil
	default:
		return 0, fmt.Errorf("get sum of credit error:%w", err)
	}
}
