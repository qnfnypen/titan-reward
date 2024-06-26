package user

import (
	"context"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/shopspring/decimal"
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

// convertTimestamp 将时间戳转换为剩余天数
func (s *sctx) convertTimestamp(ts int64) float64 {
	tnum := s.Config.TitanClientConf.UnbindTime

	if ts == 0 {
		return tnum
	}

	tn := time.Now().Unix()
	if ts > tn {
		tnum, _ = decimal.NewFromInt(ts - tn).Div(decimal.NewFromInt(86400)).Round(1).Float64()
	}

	return tnum
}
