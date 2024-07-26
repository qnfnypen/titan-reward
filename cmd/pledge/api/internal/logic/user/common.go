package user

import (
	"context"
	"fmt"
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
func (s *sctx) convertTimestamp(ts int64, lan string) (float64, string) {
	tnum := s.Config.TitanClientConf.UnbindTime

	if ts != 0 {
		tn := time.Now().Unix()
		if ts > tn {
			tnum, _ = decimal.NewFromInt(ts - tn).Div(decimal.NewFromInt(3600)).Round(1).Float64()
		}
	}

	switch lan {
	case "en":
		if tnum >= 24 {
			tnum, _ = decimal.NewFromFloat(tnum).Div(decimal.NewFromInt(24)).Round(1).Float64()
			return tnum, "DAY"
		}
		return tnum, "HOUR"
	default:
		if tnum >= 24 {
			tnum, _ = decimal.NewFromFloat(tnum).Div(decimal.NewFromInt(24)).Round(1).Float64()
			return tnum, "天"
		}
		return tnum, "小时"
	}
}

func converTimeDur(ut time.Duration, lan string) string {
	uth := ut.Hours()

	d := int(uth / 24)
	h := math.Mod(uth, 24)

	switch lan {
	case "en":
		if d > 0 {
			return fmt.Sprintf("%dd%gh", d, h)
		}
		return fmt.Sprintf("%gh", h)
	case "cn":
		if d > 0 {
			return fmt.Sprintf("%d天%g小时", d, h)
		}
		return fmt.Sprintf("%g小时", h)
	}

	return ut.String()
}
