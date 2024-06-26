package user

import (
	"context"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
)

var defaultRate = 14.52

// Sctx 类型再定义
type sctx svc.ServiceContext

func getTTNT(num *big.Int) float64 {
	numFloat := new(big.Float)

	numFloat = numFloat.SetInt(num)
	nf, _ := numFloat.Quo(numFloat, big.NewFloat(math.Pow10(6))).Float64()

	return nf
}

// getRate 获取验证节点的预期年收益率
func (s *sctx) getRate(ctx context.Context) float64 {
	rate, err := s.RedisCli.GetCtx(ctx, types.RateKey)
	if err != nil || rate == "" {
		return defaultRate
	}

	rf, err := strconv.ParseFloat(strings.TrimSpace(rate), 10)
	if rf != 0 {
		defaultRate = rf
	}

	return defaultRate
}
