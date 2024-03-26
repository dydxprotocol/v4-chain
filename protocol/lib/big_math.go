package lib

import (
	"fmt"
	"math"
	"math/big"
)

// BigMulPow10 returns the result of `val * 10^exponent`, in *big.Rat.
func BigMulPow10(
	val *big.Int,
	exponent int32,
) (
	result *big.Rat,
) {
	ratPow10 := RatPow10(exponent)
	return ratPow10.Mul(
		new(big.Rat).SetInt(val),
		ratPow10,
	)
}

// bigPow10Memo is a cache of the most common exponent value requests. Since bigPow10Memo will be
// accessed from different go-routines, the map should only ever be read from or collision
// could occur.
var bigPow10Memo = warmCache()

// BigPow10 returns the result of `10^exponent`. Caches all calculated values and
// re-uses cached values in any following calls to BigPow10.
func BigPow10(exponent uint64) *big.Int {
	result := bigPow10Helper(exponent)
	// Copy the result, such that no values can be modified by reference in the
	// `bigPow10Memo` cache.
	copy := new(big.Int).Set(result)
	return copy
}

// RatPow10 returns the result of `10^exponent`. Re-uses the cached values by
// calling bigPow10Helper.
func RatPow10(exponent int32) *big.Rat {
	result := new(big.Rat).SetInt(bigPow10Helper(uint64(AbsInt32(exponent))))
	if exponent < 0 {
		result.Inv(result)
	}
	return result
}

// BigIntMulPpm takes a `big.Int` and returns the result of `input * ppm / 1_000_000`. This method rounds towards
// negative infinity.
func BigIntMulPpm(input *big.Int, ppm uint32) *big.Int {
	result := new(big.Int)
	result.Mul(input, big.NewInt(int64(ppm)))
	return result.Div(result, big.NewInt(int64(OneMillion)))
}

// BigIntMulSignedPpm takes a `big.Int` and returns the result of `input * ppm / 1_000_000`.
func BigIntMulSignedPpm(input *big.Int, ppm int32, roundUp bool) *big.Int {
	result := new(big.Rat)
	result.Mul(
		new(big.Rat).SetInt(input),
		new(big.Rat).SetInt64(int64(ppm)),
	)
	result.Quo(result, BigRatOneMillion())
	return BigRatRound(result, roundUp)
}

// BigMin takes two `big.Int` as parameters and returns the smaller one.
func BigMin(a, b *big.Int) *big.Int {
	result := new(big.Int)
	// If `a` is greater than `b`, return `b` since it is smaller.
	// Else, return `a` since it is smaller than or equal to `b`.
	if a.Cmp(b) > 0 {
		result.Set(b)
	} else {
		result.Set(a)
	}
	return result
}

// BigMax takes two `big.Int` as parameters and returns the larger one.
func BigMax(a, b *big.Int) *big.Int {
	result := new(big.Int)
	// If `a` is greater than `b`, return `a` since it is larger.
	// Else, return `b` since it is greater than or equal to `a`.
	if a.Cmp(b) > 0 {
		result.Set(a)
	} else {
		result.Set(b)
	}
	return result
}

// BigRatMulPpm takes a `big.Rat` and returns the result of `input * ppm / 1_000_000`.
func BigRatMulPpm(input *big.Rat, ppm uint32) *big.Rat {
	return new(big.Rat).Mul(
		input,
		new(big.Rat).SetFrac64(
			int64(ppm),
			int64(OneMillion),
		),
	)
}

// bigGenericClamp is a helper function for BigRatClamp and BigIntClamp
// takes an input, upper bound, and lower bound. It returns the result
// bounded within the upper and lower bound, inclusive.
// Note that if there is overlap between the bounds (`lower > upper`), this
// function will do the following:
// - If `n < lower`, the lower bound is returned.
// - Else, the upper bound is returned (since `n >= lower`, then `n > upper` must be true).
func bigGenericClamp[T big.Int | big.Rat, P interface {
	Cmp(P) int
	Set(P) P
	*T
}](n P, lowerBound P, upperBound P) P {
	// If `n` is less than the lower bound, copy and return the lower bound.
	result := P(new(T))
	if n.Cmp(lowerBound) == -1 {
		result.Set(lowerBound)
		return result
	}

	// If `n` is greater than the upper bound, copy and return the upper bound.
	if n.Cmp(upperBound) == 1 {
		result.Set(upperBound)
		return result
	}

	// `n` is between the lower and upper bound, therefore copy and return `n`.
	result.Set(n)
	return result
}

// See `bigGenericClamp` for specification.
func BigRatClamp(n *big.Rat, lowerBound *big.Rat, upperBound *big.Rat) *big.Rat {
	return bigGenericClamp(n, lowerBound, upperBound)
}

// See `bigGenericClamp` for specification.
func BigIntClamp(n *big.Int, lowerBound *big.Int, upperBound *big.Int) *big.Int {
	return bigGenericClamp(n, lowerBound, upperBound)
}

// BigRatRound takes an input and a direction to round (true for up, false for down).
// It returns the result rounded to a `*big.Int` in the specified direction.
func BigRatRound(n *big.Rat, roundUp bool) *big.Int {
	numeratorBig := n.Num()
	denominatorBig := n.Denom()
	resultBig, remainderBig := new(big.Int).DivMod(numeratorBig, denominatorBig, new(big.Int))
	// If the remainder is non-zero, then round up by adding 1.
	// Note this works for negative numbers due to the following reasons:
	// - In euclidean division, the remainder is always positive so the resulting division rounds
	//   down instead of towards zero.
	// - The denominator of `big.Rat` is always positive. Therefore if `n` is negative, that means
	//   the numerator is negative.
	if remainderBig.Sign() > 0 && roundUp {
		resultBig.Add(resultBig, big.NewInt(1))
	}
	return resultBig
}

// BigIntRoundToMultiple takes an input, a multiple, and a direction to round (true for up,
// false for down). It returns a rounded result such that it is evenly divided by `multiple`.
// This function always expects the `multiple` parameter to be positive, otherwise it will panic.
func BigIntRoundToMultiple(
	n *big.Int,
	multiple *big.Int,
	roundUp bool,
) *big.Int {
	if multiple.Sign() <= 0 {
		panic("BigIntRoundToMultiple: multiple must be positive")
	}

	result, remainder := new(big.Int).DivMod(n, multiple, new(big.Int))
	if roundUp && remainder.Sign() > 0 {
		result = result.Add(result, big.NewInt(1))
	}
	return result.Mul(result, multiple)
}

// BigInt32Clamp takes a `big.Int` as input, and `int32` upper and lower bounds. It returns
// `int32` bounded within the upper and lower bound, inclusive.
// Note that if there is overlap between the bounds (`lower > upper`), this
// function will do the following:
// - If `n < lower`, the lower bound is returned.
// - Else, the upper bound is returned (since `n >= lower`, then `n > upper` must be true).
func BigInt32Clamp(n *big.Int, lowerBound, upperBound int32) int32 {
	// If `n` is less than the lower bound, return the lower bound.
	if n.Cmp(new(big.Int).SetInt64(int64(lowerBound))) == -1 {
		return lowerBound
	}

	// If `n` is greater than the upper bound, return the upper bound.
	if n.Cmp(new(big.Int).SetInt64(int64(upperBound))) == 1 {
		return upperBound
	}

	// `n` is between the lower and upper bound, which also means it must fit in a `int32`.
	// Therefore return `n`.
	return int32(n.Int64())
}

// BigUint64Clamp takes a `big.Int` as input, and `uint64` upper and lower bounds. It returns
// `uint64` bounded within the upper and lower bound, inclusive.
// Note that if there is overlap between the bounds (`lower > upper`), this
// function will do the following:
// - If `n < lower`, the lower bound is returned.
// - Else, the upper bound is returned (since `n >= lower`, then `n > upper` must be true).
func BigUint64Clamp(n *big.Int, lowerBound, upperBound uint64) uint64 {
	// If `n` is less than the lower bound, return the lower bound.
	if n.Cmp(new(big.Int).SetUint64(lowerBound)) == -1 {
		return lowerBound
	}

	// If `n` is greater than the upper bound, return the upper bound.
	if n.Cmp(new(big.Int).SetUint64(upperBound)) == 1 {
		return upperBound
	}

	// `n` is between the lower and upper bound, which also means it must fit in a `uint64`.
	// Therefore return `n`.
	return n.Uint64()
}

// `MustConvertBigIntToInt32` converts a `big.Int` to an `int32` and panics if the input value overflows
// or underflows `int32`.
func MustConvertBigIntToInt32(n *big.Int) int32 {
	// If `n` is greater than maxInt32 or less than minInt32, panic.
	if n.Cmp(new(big.Int).SetInt64(math.MaxInt32)) > 0 || n.Cmp(new(big.Int).SetInt64(math.MinInt32)) < 0 {
		panic("MustConvertBigIntToInt32: input value overflows or underflows int32")
	}
	return int32(n.Int64())
}

func bigPow10Helper(exponent uint64) *big.Int {
	m, ok := bigPow10Memo[exponent]
	if ok {
		return m
	}

	// Subdivide the exponent and recursively calculate each result, then multiply
	// both results together (given that `10^exponent = 10^(exponent / 2) *
	// 10^(exponent - (exponent / 2))`.
	e1 := exponent / 2
	e2 := exponent - e1
	return new(big.Int).Mul(bigPow10Helper(e1), bigPow10Helper(e2))
}

// warmCache is used to populate `bigPow10Memo` with the most common exponent requests. Since,
// none of the exponents should ever be invalid - panic immediately if an exponent is cannot be
// parsed.
func warmCache() map[uint64]*big.Int {
	exponentString := "1"
	bigExponentValues := make(map[uint64]*big.Int, 100)
	for i := 0; i < 100; i++ {
		bigValue, ok := new(big.Int).SetString(exponentString, 0)

		if !ok {
			panic(fmt.Sprintf("Failed to get big from string for exponent memo: %v", exponentString))
		}

		bigExponentValues[uint64(i)] = bigValue
		exponentString = exponentString + "0"
	}

	return bigExponentValues
}

// BigRatRoundToNearestMultiple rounds `value` up/down to the nearest multiple of `base`.
// Returns 0 if `base` is 0.
func BigRatRoundToNearestMultiple(
	value *big.Rat,
	base uint32,
	up bool,
) uint64 {
	if base == 0 {
		return 0
	}

	quotient := new(big.Rat).Quo(
		value,
		new(big.Rat).SetUint64(uint64(base)),
	)
	quotientFloored := new(big.Int).Div(quotient.Num(), quotient.Denom())

	if up && quotientFloored.Cmp(quotient.Num()) != 0 {
		return (quotientFloored.Uint64() + 1) * uint64(base)
	}

	return quotientFloored.Uint64() * uint64(base)
}
