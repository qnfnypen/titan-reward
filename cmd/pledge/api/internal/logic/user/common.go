package user

import (
	"math"
	"math/big"
)

func getTTNT(num *big.Int) float64 {
	numFloat := new(big.Float)

	numFloat = numFloat.SetInt(num)
	nf, _ := numFloat.Quo(numFloat, big.NewFloat(math.Pow10(6))).Float64()

	return nf
}
