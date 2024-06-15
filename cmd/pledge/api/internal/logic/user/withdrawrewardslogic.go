package user

import (
	"context"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/zeromicro/go-zero/core/logx"
)

// WithdrawRewardsLogic 提取收益
type WithdrawRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewWithdrawRewardsLogic 新建 提取收益
func NewWithdrawRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawRewardsLogic {
	return &WithdrawRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// WithdrawRewards 实现 提取收益
func (l *WithdrawRewardsLogic) WithdrawRewards() error {
	var (
		gzErr merror.GzErr
	)

	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.WithdrawRewardsErrCode, lan)

	err := l.svcCtx.TitanCli.WithdrawRewards(l.ctx, wallet)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return gzErr
	}

	return nil
}
