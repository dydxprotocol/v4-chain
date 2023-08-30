// Package big provides testing utility methods for the "math/big" library.
package big

import (
	"math/big"
)

// MustFirst is used for returning the first value of the SetString
// method on a `*big.Int` or `*big.Rat`. This will panic if the conversion fails.
func MustFirst[T *big.Int | *big.Rat](n T, success bool) T {
	if !success {
		panic("Conversion failed")
	}
	return n
}
