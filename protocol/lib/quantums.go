package lib

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
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
	baseQuantums *int256.Int,
	baseCurrencyAtomicResolution int32,
	priceValue uint64,
	priceExponent int32,
) (notional *int256.Int) {
	return multiplyByPrice(
		baseQuantums,
		baseCurrencyAtomicResolution,
		priceValue,
		priceExponent,
	)
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
// The result is rounded down.
func QuoteToBaseQuantums(
	quoteQuantums *int256.Int,
	baseCurrencyAtomicResolution int32,
	priceValue uint64,
	priceExponent int32,
) (notional *int256.Int) {
	notional = new(int256.Int).Set(quoteQuantums)
	exponent := int64(-(priceExponent + baseCurrencyAtomicResolution - QuoteCurrencyAtomicResolution))
	notional.MulExp10(notional, exponent)
	notional.Div(notional, int256.NewUnsignedInt(priceValue))
	return notional
}

// multiplyByPrice multiples a value by price, factoring in exponents of base
// and quote currencies.
// Given `value`, returns result of the following:
//
// `value * priceValue * 10^(priceExponent + baseAtomicResolution - quoteAtomicResolution)` [expression 2]
//
// Note that both `BaseToQuoteQuantums` and `FundingRateToIndex` directly wrap around this function.
// - For `BaseToQuoteQuantums`, substituing `value` with `baseQuantums` in expression 2 yields expression 1.
// - For `FundingRateToIndex`, substituing `value` with `fundingRatePpm * time` in expression 2 yields expression 3.
func multiplyByPrice(
	value *int256.Int,
	baseCurrencyAtomicResolution int32,
	priceValue uint64,
	priceExponent int32,
) (result *int256.Int) {
	result = int256.NewUnsignedInt(priceValue)
	result.Mul(result, value)
	return result.MulExp10(result, int64(priceExponent+baseCurrencyAtomicResolution-QuoteCurrencyAtomicResolution))
}

// FundingRateToIndex converts funding rate (in ppm) to FundingIndex given the oracle price.
//
// To get funding index from funding rate, we know that:
//
//   - `fundingPaymentQuoteQuantum = fundingRatePpm / 1_000_000 * time * quoteQuantums`
//   - Divide both sides by `baseQuantums`:
//   - Left side: `fundingPaymentQuoteQuantums / baseQuantums = fundingIndexDelta / 1_000_000`
//   - right side:
//     ```
//     fundingRate * time * quoteQuantums / baseQuantums = fundingRatePpm / 1_000_000 *
//     priceValue * 10^(priceExponent + baseCurrencyAtomicResolution - quoteCurrencyAtomicResolution) [expression 3]
//     ```
//
// Hence, further multiplying both sides by 1_000_000, we have:
//
//	fundingIndexDelta =
//	  (fundingRatePpm * time) * priceValue *
//	  10^(priceExponent + baseCurrencyAtomicResolution - quoteCurrencyAtomicResolution)
//
// Arguments:
//
//	proratedFundingRate: prorated funding rate adjusted by time delta, in parts-per-million
//	baseCurrencyAtomicResolution: atomic resolution of the base currency
//	priceValue: index price of the perpetual market according to the pricesKeeper
//	priceExponent: priceExponent of the market according to the pricesKeeper
func FundingRateToIndex(
	proratedFundingRate *int256.Int,
	baseCurrencyAtomicResolution int32,
	priceValue uint64,
	priceExponent int32,
) (fundingIndex *int256.Int) {
	return multiplyByPrice(
		proratedFundingRate,
		baseCurrencyAtomicResolution,
		priceValue,
		priceExponent,
	)
}
