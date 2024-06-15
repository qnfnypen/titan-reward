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

// UnDelegateLogic 解除质押
type UnDelegateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUnDelegateLogic 新建 解除质押
func NewUnDelegateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnDelegateLogic {
	return &UnDelegateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UnDelegate 实现 解除质押
func (l *UnDelegateLogic) UnDelegate(req *types.DelegateReq) error {
	var (
		gzErr  merror.GzErr
		amount = new(big.Int)
	)

	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.UnDelegateTokenErrCode, lan)

	// 处理ttnt
	amountFloat := new(big.Float).SetFloat64(req.Amount)
	amount, _ = amountFloat.Mul(amountFloat, big.NewFloat(math.Pow10(6))).Int(amount)

	err := l.svcCtx.TitanCli.UnDelegate(l.ctx, wallet, req.Validator, amount)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return gzErr
	}

	return nil
}
