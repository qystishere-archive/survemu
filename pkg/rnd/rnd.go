package rnd

import (
	"crypto/rand"
	"math/big"
)

func Next(limit int) int {
	if limit == 0 {
		return 0
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(limit)))
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64())
}

func NextLimit(min int, limit int) int {
	if limit == 0 {
		return 0
	}
	return min + Next(limit-min)
}
