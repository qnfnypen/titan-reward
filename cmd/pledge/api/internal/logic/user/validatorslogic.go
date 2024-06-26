package user

import (
	"context"
	"fmt"
	"math/big"
	"unsafe"

	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

// ValidatorsLogic 获取验证者信息
type ValidatorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewValidatorsLogic 新建 获取验证者信息
func NewValidatorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidatorsLogic {
	return &ValidatorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Validators 实现 获取验证者信息
func (l *ValidatorsLogic) Validators(req *types.GetValidatorReq) (resp *types.Validators, err error) {
	var (
		gzErr  merror.GzErr
		comctx = (*sctx)(unsafe.Pointer(l.svcCtx))
	)
	resp = new(types.Validators)
	resp.List = make([]types.ValidatorInfo, 0)

	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.GetValitorErrCode, lan)

	// 获取所有的验证者的token
	tokens, total, err := l.getAllTokens(req.Key)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return nil, gzErr
	}
	var validators []staking.Validator
	switch req.Kind {
	case 0:
		resp.Total = total
		validators, err = l.svcCtx.TitanCli.QueryValidators(l.ctx, req.Page, req.Size, req.Key)
		if err != nil {
			gzErr.LogErr = merror.NewError(fmt.Errorf("get all validators error:%w", err)).Error()
			return nil, gzErr
		}
	case 1:
		resp.Total, err = l.getDelgatorVidatorNums(wallet)
		if err != nil {
			gzErr.LogErr = merror.NewError(err).Error()
			return nil, gzErr
		}
		validators, err = l.svcCtx.TitanCli.QueryDelgatorVlidators(l.ctx, wallet, req.Page, req.Size)
		if err != nil {
			gzErr.LogErr = merror.NewError(fmt.Errorf("get delgator validators error:%w", err)).Error()
			return nil, gzErr
		}
	}
	for i, v := range validators {
		token := v.Tokens.BigInt()
		info := types.ValidatorInfo{}
		info.Name = v.OperatorAddress
		info.Validator = v.OperatorAddress
		if req.Kind == 1 {
			del, err := l.svcCtx.TitanCli.QueryDelegation(l.ctx, wallet, v.OperatorAddress)
			if err != nil {
				gzErr.LogErr = merror.NewError(fmt.Errorf("get delgator's delegation error:%w", err)).Error()
				return nil, gzErr
			}
			token = del.DelegationResponse.Balance.Amount.BigInt()
		}
		info.ID = int64(int(req.Page*req.Size) + i + 1)
		info.StakedTokens = getTTNT(token)
		// rf, _ := new(big.Float).Quo(new(big.Float).SetInt(v.DelegatorShares.BigInt()), new(big.Float).SetInt(v.Tokens.BigInt())).Float64()
		// info.Rate = rf
		info.Rate = comctx.getRate(l.ctx)
		vpf, _ := new(big.Float).Quo(new(big.Float).SetInt(v.Tokens.BigInt()), new(big.Float).SetInt(tokens)).Float64()
		info.VotingPower, _ = decimal.NewFromFloat(vpf).Round(4).Mul(decimal.NewFromInt(100)).Float64()
		info.UnbindingPeriod = comctx.convertTimestamp(v.UnbondingTime.Unix())
		dc, _ := decimal.NewFromString(v.Commission.Rate.String())
		info.HandlingFees, _ = dc.Round(4).Mul(decimal.NewFromInt(100)).Float64()
		resp.List = append(resp.List, info)
	}

	return resp, nil
}

func (l *ValidatorsLogic) getAllTokens(key string) (*big.Int, int64, error) {
	var (
		totalTokens = new(big.Int)
		count       int64
	)

	validators, err := l.svcCtx.TitanCli.QueryValidators(l.ctx, 0, 0, key)
	if err != nil {
		return nil, 0, fmt.Errorf("get all tokens of validators error:%w", err)
	}
	for _, v := range validators {
		totalTokens = totalTokens.Add(totalTokens, v.Tokens.BigInt())
		count++
	}

	return totalTokens, count, nil
}

func (l *ValidatorsLogic) getDelgatorVidatorNums(addr string) (int64, error) {
	vs, err := l.svcCtx.TitanCli.QueryDelgatorVlidators(l.ctx, addr, 0, 0)
	if err != nil {
		return 0, fmt.Errorf("get total of delgator vidators error:%w", err)
	}

	return int64(len(vs)), nil
}
