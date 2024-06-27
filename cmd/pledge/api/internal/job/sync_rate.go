package job

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ctx = context.Background()
)

// syncRateByInf 通过通货膨胀率计算验证者预期年收益
func syncRateByInf(sctx *svc.ServiceContext) func() {
	return func() {
		var raddr = sctx.Config.TitanClientConf.RateAddr
		// 获取当前时间验证者节点的总余额
		balance, err := sctx.TitanCli.GetBalance(ctx, raddr)
		if err != nil {
			logx.Error(fmt.Errorf("get balance error:%w", err))
			return
		}
		// 获取通货膨胀率
		inf, err := sctx.TitanCli.GetMintInflation(ctx)
		if err != nil {
			logx.Error(err)
			return
		}
		// 获取质押总金额
		delegate, err := sctx.TitanCli.GetDelegations(ctx, raddr)
		if err != nil {
			logx.Error(err)
			return
		}
		bf := new(big.Float).SetInt(balance.Amount.BigInt())
		df := new(big.Float).SetInt(delegate.Amount.BigInt())
		bf = bf.Mul(bf, inf)
		bff, _ := bf.Quo(bf, df).Float64()
		bff, _ = decimal.NewFromFloat(bff).Round(2).Float64()
		sctx.RedisCli.Set(types.RateKey, fmt.Sprintf("%v", bff))
	}
}

// syncRate 同步验证者预期年收益
func syncRate(sctx *svc.ServiceContext) func() {
	return func() {
		// 获取上一次记录的质押详情
		flag, dc, ldcr, err := getLastDelegateInfo(sctx)
		if err != nil {
			logx.Error(err)
			return
		}
		// 如果为获取上一次质押详情则不进行计算，直接退出
		if !flag {
			return
		}
		// 获取当前质押收益
		ndcr, err := getCurrentRewards(sctx)
		if err != nil {
			logx.Error(err)
			return
		}
		// 换算为年收益
		dcr := ndcr.Sub(ndcr, ldcr)
		dcr = dcr.Mul(dcr, big.NewFloat(36.5))
		dcrf, _ := dcr.Quo(dcr, dc).Float64()
		dcrf, _ = decimal.NewFromFloat(dcrf).Round(2).Float64()
		sctx.RedisCli.Set(types.RateKey, fmt.Sprintf("%v", dcrf))
	}
}

// getLastDelegateInfo 获取上一次的质押详情
func getLastDelegateInfo(sctx *svc.ServiceContext) (bool, *big.Float, *big.Float, error) {
	var (
		flag      = true
		dcf, dcrf *big.Float
		raddr     = sctx.Config.TitanClientConf.RateAddr
	)

	// 获取质押的金额
	dc, err := sctx.RedisCli.Get(types.DelegateCoinKey)
	if err != nil || dc == "" {
		flag = false
		// 没有获取到或者为空，则获取当前的质押金额进行存储
		dcs, err := sctx.TitanCli.GetDelegations(ctx, raddr)
		if err != nil {
			return false, nil, nil, merror.NewError(fmt.Errorf("get delegation's coins error:%w", err))
		}
		// 转换为ttnt进行计算
		dcsf := new(big.Float).SetInt(dcs.Amount.BigInt())
		dcsf = dcsf.Quo(dcsf, big.NewFloat(math.Pow10(6)))
		dc = dcsf.String()
		// 存储到redis中
		sctx.RedisCli.Setex(types.DelegateCoinKey, dc, 146*60)
		dcf = dcsf
	} else {
		dcf, _ = new(big.Float).SetString(dc)
	}
	// 获取上一次记录的质押收益
	dcr, err := sctx.RedisCli.Get(types.LastDelegateReward)
	if err != nil || dcr == "" {
		// 没有获取到或者为空，则获取当前的质押收益进行存储
		dcrsf, err := getCurrentRewards(sctx)
		if err != nil {
			return false, nil, nil, err
		}
		dcr = dcrsf.String()
		// 存储到redis中
		sctx.RedisCli.Setex(types.LastDelegateReward, dc, 146*60)
		dcrf = dcrsf
	} else {
		dcrf, _ = new(big.Float).SetString(dcr)
	}

	return flag, dcf, dcrf, nil
}

func getCurrentRewards(sctx *svc.ServiceContext) (*big.Float, error) {
	dcrf, err := sctx.TitanCli.GetRewards(ctx, sctx.Config.TitanClientConf.RateAddr)
	if err != nil {
		return nil, merror.NewError(fmt.Errorf("get delegation's rewards error:%w", err))
	}
	// 转换为ttnt进行计算
	dcrf = dcrf.Quo(dcrf, big.NewFloat(math.Pow10(6)))

	return dcrf, nil
}
