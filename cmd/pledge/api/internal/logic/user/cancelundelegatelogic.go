package user

import (
	"context"
	"math"
	"math/big"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"

	"github.com/zeromicro/go-zero/core/logx"
)

// CancelUnDelegateLogic 取消解除质押
type CancelUnDelegateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCancelUnDelegateLogic 新建 取消解除质押
func NewCancelUnDelegateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelUnDelegateLogic {
	return &CancelUnDelegateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CancelUnDelegate 实现 取消解除质押
func (l *CancelUnDelegateLogic) CancelUnDelegate(req *types.CancelUnDelegateReq) error {
	var (
		gzErr  merror.GzErr
		amount = new(big.Int)
	)

	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.CancelUnDelegateTokenErrCode, lan)

	// 处理ttnt
	amountFloat := new(big.Float).SetFloat64(req.Amount)
	amount, _ = amountFloat.Mul(amountFloat, big.NewFloat(math.Pow10(6))).Int(amount)

	err := l.svcCtx.TitanCli.CancelUnbondingDelegation(l.ctx, wallet, req.Validator, req.Height, amount)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return gzErr
	}

	return nil
}
