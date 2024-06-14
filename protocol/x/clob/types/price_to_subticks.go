package types

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PriceToSubticks converts price value from Prices module into subticks.
// By construction of the Clob module:
// `sizeQuoteQuantums = subticks * baseQuantums * 10^quantumConversionExponent`
// Substituting `baseQuantums` with a full coin of the base currency:
//
// `usdcPrice * 10^(-quoteAtomicResolution) = subticks * 10^(-baseAtomicResolution) *
// 10^quantumConversionExponent` (A)
//
// By construction of Prices module:
//
// `usdcPrice = marketPrice.Price * 10^marketPrice.Exponent` (B)
//
// Combining equations (A) & (B), we get:
//
// `subticks = marketPrice.Price * 10^(marketPrice.Exponent - quantumConversionExponent +
// baseAtomicResolution - quoteAtomicResolution)`
func PriceToSubticks(
	marketPrice pricestypes.MarketPrice,
	clobPair ClobPair,
	baseAtomicResolution int32,
	quoteAtomicResolution int32,
) (
	ratSubticks *big.Rat,
) {
	exponent := int32(
		marketPrice.Exponent - clobPair.QuantumConversionExponent + baseAtomicResolution - quoteAtomicResolution,
	)
	// TODO(DEC-1256): Use index price from the price daemon, instead of oracle price.
	bigPrice := new(big.Int).SetUint64(marketPrice.Price)
	p10, inverse := lib.BigPow10(exponent)
	if inverse {
		return new(big.Rat).SetFrac(bigPrice, p10)
	} else {
		return new(big.Rat).SetInt(bigPrice.Mul(bigPrice, p10))
	}
}

// SubticksToPrice converts subticks into price value from Prices module.
// By construction of the Clob module:
// `sizeQuoteQuantums = subticks * baseQuantums * 10^quantumConversionExponent`
// Substituting `baseQuantums` with a full coin of the base currency:
//
// `usdcPrice * 10^(-quoteAtomicResolution) = subticks * 10^(-baseAtomicResolution) *
// 10^quantumConversionExponent` (A)
//
// By construction of Prices module:
//
// `usdcPrice = marketPrice.Price * 10^marketPrice.Exponent` (B)
//
// Combining equations (A) & (B), we get:
//
// `marketPrice.Price = subticks * 10^(-marketPrice.Exponent + quantumConversionExponent -
// baseAtomicResolution + quoteAtomicResolution)`
// Note this function rounds down in order to typecast into an int. It should really only be used
// in testing with well-defined integer subticks.
func SubticksToPrice(
	subticks Subticks,
	marketPriceExponent int32,
	clobPair ClobPair,
	baseAtomicResolution int32,
	quoteAtomicResolution int32,
) (
	price uint64,
) {
	exponent := int32(
		-marketPriceExponent + clobPair.QuantumConversionExponent - baseAtomicResolution + quoteAtomicResolution,
	)
	result := lib.BigIntMulPow10(new(big.Int).SetUint64(uint64(subticks)), exponent, false)
	if !result.IsUint64() {
		panic("SubticksToPrice: result is not a uint64")
	}
	return result.Uint64()
}
