package job

import (
	"context"
	"fmt"
	"log"
	"math/big"

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
		var (
			di = new(big.Int)
		)
		// 获取当前时间验证者节点的总余额
		balance, err := sctx.TitanCli.GetTotalBalance(ctx)
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
		validators, err := sctx.TitanCli.QueryValidators(ctx, 0, 0, "")
		if err != nil {
			logx.Error(err)
			return
		}
		for _, v := range validators {
			di.Add(di, v.Tokens.BigInt())
		}
		// 总金额为所有可用余额+所有质押金额
		bf := new(big.Float).SetInt(balance.Amount.BigInt())
		df := new(big.Float).SetInt(di)
		bf.Mul(bf, inf)
		bff, _ := bf.Quo(bf, df).Float64()
		bff, _ = decimal.NewFromFloat(bff * 100).Round(2).Float64()
		log.Println(bff)
		sctx.RedisCli.Set(types.RateKey, fmt.Sprintf("%v", bff))
	}
}
