package types

import (
	"errors"
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

// FillAmountToQuoteQuantums converts a fill amount denoted in base quantums, to a price denoted in quote
// quantums given an order's subticks.
//
// `subticks * 10^quantumConversionExponent` represents the price-per-base quantum of an order.
//
// To convert from base to quote quantums, we use the following equation:
// `sizeQuoteQuantums = subticks * baseQuantums * 10^quantumConversionExponent`.
// Note that the result is rounded down.
//
// Note: If `subticks`, `baseQuantums`, are small enough, `quantumConversionExponent` is negative, it's possible that
// this method returns `0` `quoteQuantums` for a non-zero amount of `baseQuantums`. This could mean that it's possible
// that a maker sell order on the book at a very unfavorable price (subticks) could receive `0` `quoteQuantums` amount.
// This case is unlikely to happen in a production environment as:
//   - The maker orders would need to be placed at a very unfavorable price for the maker. There is no incentive for
//     makers to place such orders.
//   - The taker would need to place hundreds of thousands of transactions, filling the maker orders, in order to
//     profit a single dollar of `quoteQuantums`.
func FillAmountToQuoteQuantums(
	subticks Subticks,
	baseQuantums satypes.BaseQuantums,
	quantumConversionExponent int32,
) (bigNotional *big.Int) {
	bigSubticks := subticks.ToBigInt()
	bigBaseQuantums := baseQuantums.ToBigInt()

	bigSubticksMulBaseQuantums := new(big.Int).Mul(bigSubticks, bigBaseQuantums)

	exponent := int64(quantumConversionExponent)

	// To ensure we are always doing math with integers, we take the absolute
	// value of the exponent. If `exponent` is non-negative, then `10^exponent` is an
	// integer and we can multiply by it. Else, `10^exponent` is less than 1 and we should
	// multiply by `1 / 10^exponent` (which must be an integer if `exponent < 0`).
	bigExponentValue := lib.BigPow10(lib.AbsInt64(exponent))

	bigQuoteQuantums := new(big.Int)
	if exponent < 0 {
		// `1 / 10^exponent` is an integer.
		bigQuoteQuantums.Div(bigSubticksMulBaseQuantums, bigExponentValue)
	} else {
		// `10^exponent` is an integer.
		bigQuoteQuantums.Mul(bigSubticksMulBaseQuantums, bigExponentValue)
	}

	return bigQuoteQuantums
}

// GetAveragePriceSubticks computes the average price (in subticks) of filled
// amount in `quoteQuantums` and `baseQuantums`.
// To calculate quote quantums from base quantums and subticks, we use the
// following equation:
// `sizeQuoteQuantums = subticks * baseQuantums * 10^quantumConversionExponent`.
//
// Thus, to get `subticks`:
// `subticks = sizeQuoteQuantums * 10^(-quantumConversionExponent) / baseQuantums`
//
// This function panics if `bigBaseQuantums == 0`. The result of division is rounded down.
func GetAveragePriceSubticks(
	bigQuoteQuantums *big.Int,
	bigBaseQuantums *big.Int,
	quantumConversionExponent int32,
) (bigSubticks *big.Rat) {
	if bigBaseQuantums.Sign() == 0 {
		panic(errors.New("GetAveragePriceSubticks: bigBaseQuantums = 0"))
	}

	result := lib.BigMulPow10(bigQuoteQuantums, -quantumConversionExponent)
	return result.Quo(
		result,
		new(big.Rat).SetInt(bigBaseQuantums),
	)
}

// NotionalToCoinAmount returns the coin amount (e.g. `uatom`) that has equal worth to the notional (in quote quantums).
// For example, given price of 9.5 TDAI/ATOM, notional of 9_500_000 quote quantums, return 1_000_000 `uatom` (since
// `tokenDenomExp“=-6).
// Note the return value is in coin amount, which is different from base quantums.
//
// Given the below by definitions:
//
//	quote_quantums * 10^quote_atomic_resolution = full_quote_coin_amount (e.g. 2_000_000 quote quantums * 10^-6 = 2 TDAI)
//	coin_amount * 10^denom_exponent = full_coin_amount (e.g. 1_000_000 uatom * 10^-6 = 1 ATOM)
//	full_coin_amount * coin_price = full_quote_coin_amount (e.g. 1 ATOM * 9.5 TDAI/ATOM = 9.5 TDAI)
//
// Therefore:
//
//	coin_amount * 10^denom_exponent * coin_price = quote_quantums * 10^quote_atomic_resolution
//	coin_amount = quote_quantums * 10^(quote_atomic_resolution - denom_exponent) / coin_price
func NotionalToCoinAmount(
	notionalQuoteQuantums *big.Int,
	quoteAtomicResolution int32,
	denomExp int32,
	marketPrice pricestypes.MarketPrice,
) *big.Rat {
	fullCoinPrice := lib.BigMulPow10(
		new(big.Int).SetUint64(marketPrice.PnlPrice),
		marketPrice.Exponent,
	)
	ret := lib.BigMulPow10(notionalQuoteQuantums, quoteAtomicResolution-denomExp)
	return ret.Quo(
		ret,
		fullCoinPrice,
	)
}
