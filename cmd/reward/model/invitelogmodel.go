package model

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InviteLogModel = (*customInviteLogModel)(nil)

type (
	// InviteLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInviteLogModel.
	InviteLogModel interface {
		inviteLogModel
		withSession(session sqlx.Session) InviteLogModel
		GetInviteCreditByUn(ctx context.Context, un string) (int64, error)
	}

	customInviteLogModel struct {
		*defaultInviteLogModel
	}
)

// NewInviteLogModel returns a model for the database table.
func NewInviteLogModel(conn sqlx.SqlConn) InviteLogModel {
	return &customInviteLogModel{
		defaultInviteLogModel: newInviteLogModel(conn),
	}
}

func (m *customInviteLogModel) withSession(session sqlx.Session) InviteLogModel {
	return NewInviteLogModel(sqlx.NewSqlConnFromSession(session))
}

// GetInviteCreditByUn 获取用户社区邀请好友的奖励
func (m *customInviteLogModel) GetInviteCreditByUn(ctx context.Context, un string) (int64, error) {
	var credit int64

	query, args, err := squirrel.Select("SUM(credit)").From(m.table).Where("username = ?", un).ToSql()
	if err != nil {
		return credit, fmt.Errorf("get sum of credit error:%w", err)
	}

	err = m.conn.QueryRowCtx(ctx, &credit, query, args...)
	if err != nil {
		return 0, fmt.Errorf("get sum of credit error:%w", err)
	}

	return credit, nil
}
