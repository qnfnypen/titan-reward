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

// DelegateLogic 质押token
type DelegateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDelegateLogic 新建 质押token
func NewDelegateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelegateLogic {
	return &DelegateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Delegate 实现 质押token
func (l *DelegateLogic) Delegate(req *types.DelegateReq) error {
	var (
		gzErr  merror.GzErr
		amount = new(big.Int)
	)

	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.DelegateTokenErrCode, lan)

	// 处理ttnt
	amountFloat := new(big.Float).SetFloat64(req.Amount)
	amount, _ = amountFloat.Mul(amountFloat, big.NewFloat(math.Pow10(6))).Int(amount)

	err := l.svcCtx.TitanCli.Delegate(l.ctx, wallet, req.Validator, amount)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return gzErr
	}

	return nil
}
