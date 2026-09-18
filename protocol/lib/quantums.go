package lib

import (
	"math/big"
)

// BaseToQuoteQuantums converts an amount denoted in base quantums, to an equivalent amount denoted in quote
// quantums. To determine the equivalent amount, an oracle price is used.
//
//   - `priceValue * 10^priceExponent` represents the conversion rate from one full coin of base currency
//     to one full coin of quote currency.
//   - `10^baseCurrencyAtomicResolution` represents the amount of one full coin that a base quantum is equal to.
//   - `10^quoteCurrencyAtomicResolution` represents the amount of one full coin that a quote quantum is equal to.
//
// To convert from base to quote quantums, we use the following equation:
//
//	quoteQuantums =
//	  (baseQuantums * 10^baseCurrencyAtomicResolution) *
//	  (priceValue * 10^priceExponent) /
//	  (10^quoteCurrencyAtomicResolution)
//	=
//	  baseQuantums * priceValue *
//	  10^(priceExponent + baseCurrencyAtomicResolution - quoteCurrencyAtomicResolution) [expression 1]
//
// The result is rounded down.
func BaseToQuoteQuantums(
	bigBaseQuantums *big.Int,
	baseCurrencyAtomicResolution int32,
	priceValue uint64,
	priceExponent int32,
) (bigNotional *big.Int) {
	// Multiply all numerators.
	numResult := new(big.Int).SetUint64(priceValue)
	numResult.Mul(numResult, bigBaseQuantums)
	exponent := priceExponent + baseCurrencyAtomicResolution - QuoteCurrencyAtomicResolution

	// Special case: if the exponent is zero, we can return early.
	if exponent == 0 {
		return numResult
	}

	// Otherwise multiply or divide by the 1e^exponent.
	pow10, inverse := BigPow10(exponent)
	if inverse {
		// Trucated division (towards zero) instead of Euclidean division.
		return numResult.Quo(numResult, pow10)
	} else {
		return numResult.Mul(numResult, pow10)
	}
}

// QuoteToBaseQuantums converts an amount denoted in quote quantums, to an equivalent amount denoted in base
// quantums. To determine the equivalent amount, an oracle price is used.
//
//   - `priceValue * 10^priceExponent` represents the conversion rate from one full coin of base currency
//     to one full coin of quote currency.
//   - `10^baseCurrencyAtomicResolution` represents the amount of one full coin that a base quantum is equal to.
//   - `10^quoteCurrencyAtomicResolution` represents the amount of one full coin that a quote quantum is equal to.
//
// To convert from quote to base quantums, we use the following equation:
//
//	baseQuantums =
//	  quoteQuantums / priceValue /
//	  10^(priceExponent + baseCurrencyAtomicResolution - quoteCurrencyAtomicResolution)
//
// The result is rounded towards zero.
func QuoteToBaseQuantums(
	bigQuoteQuantums *big.Int,
	baseCurrencyAtomicResolution int32,
	priceValue uint64,
	priceExponent int32,
) (bigNotional *big.Int) {
	// Initialize result to quoteQuantums.
	result := new(big.Int).Set(bigQuoteQuantums)

	// Divide result (towards zero) by 10^(exponent).
	exponent := priceExponent + baseCurrencyAtomicResolution - QuoteCurrencyAtomicResolution
	p10, inverse := BigPow10(exponent)
	if inverse {
		result.Mul(result, p10)
	} else {
		result.Quo(result, p10)
	}

	// Divide result (towards zero) by priceValue.
	// If there are two divisions, it is okay to do them separately as the result is the same.
	result.Quo(result, new(big.Int).SetUint64(priceValue))

	return result
}

// BigRatRoundToMultiple rounds a big.Rat to the nearest multiple of the given number.
// If roundUp is true, it rounds up to the nearest multiple, otherwise it rounds down.
// This function always expects the `multiple` parameter to be positive, otherwise it will panic.
func BigRatRoundToMultiple(
	value *big.Rat,
	multiple *big.Int,
	roundUp bool,
) *big.Int {
	// A zero multiple divides by zero and a negative one silently reverses the
	// rounding direction, so reject both up front like BigIntRoundToMultiple does.
	if multiple.Sign() <= 0 {
		panic("BigRatRoundToMultiple: multiple must be positive")
	}

	// Round `value / multiple` to an integer and scale it back up, so that the
	// fractional part of `value` takes part in the rounding decision. Rounding
	// `value` to an integer first would discard it, and a value such as `10.5`
	// would then never round up past the multiple its floor already sits on.
	quotient := new(big.Rat).Quo(value, new(big.Rat).SetInt(multiple))
	rounded := BigRatRound(quotient, roundUp)
	return rounded.Mul(rounded, multiple)
}
