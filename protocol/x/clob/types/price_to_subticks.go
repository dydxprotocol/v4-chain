package types

import (
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

// PriceToSubticks converts price value from Prices module into subticks.
// By construction of the Clob module:
// `sizeQuoteQuantums = subticks * baseQuantums * 10^quantumConversionExponent`
// Substituting `baseQuantums` with a full coin of the base currency:
//
// `tdaiPrice * 10^(-quoteAtomicResolution) = subticks * 10^(-baseAtomicResolution) *
// 10^quantumConversionExponent` (A)
//
// By construction of Prices module:
//
// `tdaiPrice = marketPrice.Price * 10^marketPrice.Exponent` (B)
//
// Combining equations (A) & (B), we get:
//
// `subticks = marketPrice.Price * 10^(marketPrice.Exponent - quantumConversionExponent +
// baseAtomicResolution - quoteAtomicResolution)`
func PnlPriceToSubticks(
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
	return lib.BigMulPow10(
		// TODO(DEC-1256): Use daemon price from the price daemon, instead of oracle price.
		new(big.Int).SetUint64(marketPrice.PnlPrice),
		exponent,
	)
}

func SpotPriceToSubticks(
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
	return lib.BigMulPow10(
		// TODO(DEC-1256): Use daemon price from the price daemon, instead of oracle price.
		new(big.Int).SetUint64(marketPrice.SpotPrice),
		exponent,
	)
}

// SubticksToPrice converts subticks into price value from Prices module.
// By construction of the Clob module:
// `sizeQuoteQuantums = subticks * baseQuantums * 10^quantumConversionExponent`
// Substituting `baseQuantums` with a full coin of the base currency:
//
// `tdaiPrice * 10^(-quoteAtomicResolution) = subticks * 10^(-baseAtomicResolution) *
// 10^quantumConversionExponent` (A)
//
// By construction of Prices module:
//
// `tdaiPrice = marketPrice.Price * 10^marketPrice.Exponent` (B)
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
	return lib.BigRatRound(
		lib.BigMulPow10(
			// TODO(DEC-1256): Use daemon price from the price daemon, instead of oracle price.
			new(big.Int).SetUint64(uint64(subticks)),
			exponent,
		),
		false,
	).Uint64()
}
