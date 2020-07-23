package util

import (
	"crypto/rand"
	"math/big"
)

// RandomInt 获取min ~ max之间的随机值，包括min和max
func RandomInt(min, max int) int {
	maxInt := big.NewInt(int64(max) - int64(min) + 1)
	i, _ := rand.Int(rand.Reader, maxInt)
	return min + int(i.Int64())
}
