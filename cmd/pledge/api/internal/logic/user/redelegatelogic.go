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

// ReDelegateLogic 质押转移
type ReDelegateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewReDelegateLogic 新建 质押转移
func NewReDelegateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReDelegateLogic {
	return &ReDelegateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ReDelegate 实现 质押转移
func (l *ReDelegateLogic) ReDelegate(req *types.ReDelegateReq) error {
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

	err := l.svcCtx.TitanCli.ReDelegate(l.ctx, wallet, req.SrcValidator, req.DstValidator, amount)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return gzErr
	}

	return nil
}
