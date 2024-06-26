package user

import (
	"context"
	"unsafe"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetUnbindingDelegateLogic 获取进行中的解除质押
type GetUnbindingDelegateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUnbindingDelegateLogic 新建 获取进行中的解除质押
func NewGetUnbindingDelegateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnbindingDelegateLogic {
	return &GetUnbindingDelegateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUnbindingDelegate 实现 获取进行中的解除质押
func (l *GetUnbindingDelegateLogic) GetUnbindingDelegate() (resp []types.UnbindingDelegateInfo, err error) {
	var (
		gzErr  merror.GzErr
		comctx = (*sctx)(unsafe.Pointer(l.svcCtx))
	)
	resp = make([]types.UnbindingDelegateInfo, 0)

	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.GetValitorErrCode, lan)

	des, err := l.svcCtx.TitanCli.UnbondingDelegations(l.ctx, wallet, 0, 0)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return nil, gzErr
	}

	for i, v := range des {
		for ii, vv := range v.Entries {
			info := types.UnbindingDelegateInfo{}
			info.ID = int64((i+1)*ii + 1)
			info.Name = v.ValidatorAddress
			info.Validator = v.ValidatorAddress
			info.Height = vv.CreationHeight
			info.Tokens = getTTNT(vv.Balance.BigInt())
			info.UnbindingPeriod = comctx.convertTimestamp(vv.CompletionTime.Unix())

			resp = append(resp, info)
		}
	}

	return resp, nil
}
