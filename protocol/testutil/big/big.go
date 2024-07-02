// Package big provides testing utility methods for the "math/big" library.
package big

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// MustFirst is used for returning the first value of the SetString
// method on a `*big.Int` or `*big.Rat`. This will panic if the conversion fails.
func MustFirst[T *big.Int | *big.Rat](n T, success bool) T {
	if !success {
		panic("Conversion failed")
	}
	return n
}

// Int64MulPow10 returns the result of `val * 10^exponent`, in *big.Int.
func Int64MulPow10(
	val int64,
	exponent uint64,
) (
	result *big.Int,
) {
	return lib.BigIntMulPow10(big.NewInt(val), exponent, false)
}
