package user

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"

	"github.com/zeromicro/go-zero/core/logx"
)

// InfoLogic 获取用户信息详情
type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewInfoLogic 新建 获取用户信息详情
func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Info 实现 获取用户信息详情
func (l *InfoLogic) Info() (resp *types.UserInfo, err error) {
	var (
		gzErr merror.GzErr
	)
	resp = new(types.UserInfo)
	
	lan := l.ctx.Value(types.LangKey).(string)
	wallet := l.ctx.Value("wallet").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.GetUserInfoErrCode, lan)

	// 获取用户可用余额
	balance, err := l.svcCtx.TitanCli.GetBalance(l.ctx, wallet)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get balance of user error:%w", err)).Error()
		return nil, gzErr
	}
	balanceFloat := new(big.Float).SetInt(balance.Amount.BigInt())
	resp.AvailableToken, _ = balanceFloat.Quo(balanceFloat, big.NewFloat(math.Pow10(6))).Float64()
	// 获取质押token的数量
	stakedTokens, err := l.svcCtx.TitanCli.GetDelegations(l.ctx, wallet)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get staked tokens of user error:%w", err)).Error()
		return nil, gzErr
	}
	stakedTokensFloat := new(big.Float).SetInt(stakedTokens.Amount.BigInt())
	resp.StakedToken, _ = stakedTokensFloat.Quo(stakedTokensFloat, big.NewFloat(math.Pow10(6))).Float64()
	// 获取质押的收益
	rewards, err := l.svcCtx.TitanCli.GetRewards(l.ctx, wallet)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get rewards of user error:%w", err)).Error()
		return nil, gzErr
	}
	rewardsFloat := new(big.Float).SetInt(rewards)
	resp.Reward, _ = rewardsFloat.Quo(rewardsFloat, big.NewFloat(math.Pow10(6))).Float64()
	// 获取质押锁仓的token数量
	unStakedTokens, err := l.svcCtx.TitanCli.GetUnBondingDelegations(l.ctx, wallet)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get unstaked tokens of user error:%w", err)).Error()
		return nil, gzErr
	}
	unStakedTokensFloat := new(big.Float).SetInt(unStakedTokens.Amount.BigInt())
	resp.UnstakedToken, _ = unStakedTokensFloat.Quo(unStakedTokensFloat, big.NewFloat(math.Pow10(6))).Float64()

	return resp, nil
}
