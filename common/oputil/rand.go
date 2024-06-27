package oputil

import (
	"crypto/rand"
	"io"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

// GenerateNonce 生成指定长度的随机数
func GenerateNonce(len int) string {
	randBytes := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	nonce := make([]byte, len)
	io.ReadFull(rand.Reader, nonce)

	for i, v := range nonce {
		nonce[i] = randBytes[v%10]
	}

	return string(nonce)
}

// DecRound 向上取小数
func DecRound(dc decimal.Decimal, round int, ups ...bool) float64 {
	var f float64
	up := true
	// 判断当前位数是否已经小于等于需要保留的位数
	_, dec, _ := strings.Cut(dc.String(), ".")
	if utf8.RuneCountInString(dec) <= round {
		f, _ = dc.Float64()
		return f
	}
	fp := decimal.NewFromFloat(5).Div(decimal.NewFromFloat(math.Pow10(round + 1)))
	if dc.Equal(decimal.NewFromInt(0)) {
		return 0
	}

	if len(ups) > 0 {
		up = ups[0]
	}

	if up {
		f, _ = dc.Add(fp).Round(int32(round)).Float64()
	} else {
		f, _ = dc.Sub(fp).Round(int32(round)).Float64()
	}

	if f < 0 {
		return 0
	}

	return f
}
