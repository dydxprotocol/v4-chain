package lib

import (
	"fmt"
	"math"
	"math/big"
)

// BigU returns a new big.Int from the input unsigned integer.
func BigU[T uint | uint32 | uint64](u T) *big.Int {
	return new(big.Int).SetUint64(uint64(u))
}

// BigI returns a new big.Int from the input signed integer.
func BigI[T int | int32 | int64](i T) *big.Int {
	return big.NewInt(int64(i))
}

// BigMulPpm returns the result of `val * ppm / 1_000_000`, rounding in the direction indicated.
func BigMulPpm(val *big.Int, ppm *big.Int, roundUp bool) *big.Int {
	result := new(big.Int).Mul(val, ppm)
	oneMillion := BigIntOneMillion()
	if roundUp {
		return BigDivCeil(result, oneMillion)
	} else {
		return result.Div(result, oneMillion)
	}
}

// bigPow10Memo is a cache of the most common exponent value requests. Since bigPow10Memo will be
// accessed from different go-routines, the map should only ever be read from or collision
// could occur.
var bigPow10Memo = warmCache()

// BigPow10 returns the result of `10^abs(exponent)` and whether the exponent is non-negative.
func BigPow10[T int | int32 | int64 | uint | uint32 | uint64](
	exponent T,
) (
	result *big.Int,
	inverse bool,
) {
	inverse = exponent < 0
	var absExponent uint64
	if inverse {
		absExponent = uint64(-exponent)
	} else {
		absExponent = uint64(exponent)
	}

	return new(big.Int).Set(bigPow10Helper(absExponent)), inverse
}

// BigIntMulPow10 returns the result of `input * 10^exponent`, rounding in the direction indicated.
// There is no rounding if `exponent` is non-negative.
func BigIntMulPow10[T int | int32 | int64 | uint | uint32 | uint64](
	input *big.Int,
	exponent T,
	roundUp bool,
) *big.Int {
	p10, inverse := BigPow10(exponent)
	if inverse {
		if roundUp {
			return BigDivCeil(input, p10)
		}
		return new(big.Int).Div(input, p10)
	}
	return new(big.Int).Mul(p10, input)
}

// BigIntMulPpm takes a `big.Int` and returns the result of `input * ppm / 1_000_000`. This method rounds towards
// negative infinity.
func BigIntMulPpm(input *big.Int, ppm uint32) *big.Int {
	result := new(big.Int)
	result.Mul(input, big.NewInt(int64(ppm)))
	return result.Div(result, big.NewInt(int64(OneMillion)))
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

// BigRatMin takes two `big.Rat` as parameters and returns the smaller one.
func BigRatMin(a, b *big.Rat) *big.Rat {
	result := new(big.Rat)
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
	num := new(big.Int).Mul(input.Num(), big.NewInt(int64(ppm)))
	den := new(big.Int).Mul(input.Denom(), big.NewInt(int64(OneMillion)))
	return new(big.Rat).SetFrac(num, den)
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

// BigDivCeil returns the ceiling of `a / b`.
func BigDivCeil(a *big.Int, b *big.Int) *big.Int {
	result, remainder := new(big.Int).QuoRem(a, b, new(big.Int))

	// If the value was rounded (i.e. there is a remainder), and the exact result would be positive,
	// then add 1 to the result.
	if remainder.Sign() != 0 && (a.Sign() == b.Sign()) {
		result.Add(result, big.NewInt(1))
	}

	return result
}

// BigDivFloor returns the floor of `a / b`.
func BigDivFloor(a *big.Int, b *big.Int) *big.Int {
	result, remainder := new(big.Int).QuoRem(a, b, new(big.Int))

	// If the value was rounded (i.e. there is a remainder), and the exact result would be negative,
	// then subtract 1 from the result.
	if remainder.Sign() != 0 && (a.Sign() != b.Sign()) {
		result.Sub(result, big.NewInt(1))
	}

	return result
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
