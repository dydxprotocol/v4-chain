package prices

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/shopspring/decimal"
	"math/big"
)

/*
 * Price conversion functions
 *
 * At this time, the protocol represents prices as a tuple of (price, exponent), where the raw price the asset can be
 * inferred as
 *
 * 	rawPrice = price * 10 ^ exponent
 *
 * Price exponents are chosen on a per-market basis, so price conversions that involve multiple market prices require
 * converting to and from raw prices in order to perform the conversion math.
 *
 * For price conversion math, we use the decimal package found here:
 * https://pkg.go.dev/github.com/shopspring/decimal#section-readme
 *
 * The decimal package has some benefits over big.Rat that can make it more suitable for representing money - see
 * the readme. However, it is not as performant as big.Rat, so we only use it for price conversion math.
 */

// Invert inverts a price, returning the inverted price multiplied by 10^-exponent.
// This method is meant to be used for inverting stablecoin ticker prices. For the sake of precision, price inversion
// is only intended to be used for pricing markets that are close to 1:1 price-wise, e.g. USD-USDT. Inverting a price
// that is <<>> 1 could result in a loss of precision.
func Invert(price uint64, exponent types.Exponent) uint64 {
	// Invert 0 to 0. Zero prices are invalid and are prevented by the price fetcher, but we insert this guard just
	// in case to prevent any panics.
	if price == 0 {
		return 0
	}

	decimalPrice := decimal.NewFromBigInt(new(big.Int).SetUint64(price), exponent)
	invertedPrice := decimal.NewFromInt(1).Div(decimalPrice).Mul(
		decimal.NewFromBigInt(new(big.Int).SetUint64(1), -exponent),
	).BigInt().Uint64()
	return invertedPrice
}

// Multiply multiplies two prices, returning the resulting price as a uint64 multiplied by the first exponent.
//
// Formula: rawPrice = price * 10 ^ exponent
//
//	rawAdjustByPrice = adjustByPrice * 10 ^ adjustByExponent
//	rawAdjustedPrice = rawPrice * rawAdjustByPrice
//	adjustedPrice = rawAdjustedPrice * 10 ^ -exponent
//
// The most common use case of multiply will be to convert a market price from one stablecoin to another.
// For example, 1INCH-USD = 1INCH_USDT * USDT-USD.
func Multiply(price uint64, exponent int32, adjustByPrice uint64, adjustByExponent int32) (adjustedPrice uint64) {
	decimalPrice := decimal.NewFromBigInt(new(big.Int).SetUint64(price), exponent)
	decimalAdjustByPrice := decimal.NewFromBigInt(new(big.Int).SetUint64(adjustByPrice), adjustByExponent)
	adjustedPrice = decimalPrice.Mul(decimalAdjustByPrice).Mul(
		decimal.NewFromBigInt(new(big.Int).SetUint64(1), -exponent),
	).BigInt().Uint64()
	return adjustedPrice
}

// Divide divides two prices, returning the resulting price as a uint64 multiplied by the divisor price's exponent.
//
// Formula:   rawPrice = price * 10 ^ exponent
//
//	rawAdjustByPrice = adjustByPrice * 10 ^ adjustByExponent
//	rawAdjustedPrice = rawAdjustByPrice / rawPrice
//	adjustedPrice = rawAdjustedPrice * 10 ^ -exponent
//
// This price conversion method is typically used in practice to derive stablecoin prices by dividing crypto asset
// prices in two different stablecoin quote currencies: for example, USDT-USD = BTC-USD / BTC-USDT.
func Divide(
	adjustByPrice uint64,
	adjustByExponent types.Exponent,
	price uint64,
	exponent types.Exponent,
) (adjustedPrice uint64) {
	// Zero prices are invalid and are prevented by the price fetcher, but we insert this guard just in case to prevent
	// any panics.
	if price == 0 {
		return 0
	}

	decimalPrice := decimal.NewFromBigInt(new(big.Int).SetUint64(price), exponent)
	decimalAdjustByPrice := decimal.NewFromBigInt(new(big.Int).SetUint64(adjustByPrice), adjustByExponent)
	adjustedPrice = decimalAdjustByPrice.Div(decimalPrice).Mul(
		decimal.NewFromBigInt(new(big.Int).SetUint64(1), -exponent),
	).BigInt().Uint64()
	return adjustedPrice
}

// PriceToFloat32ForLogging converts a price, exponent to a float32 for logging purposes. This is not meant to be used
// for price calucations within the protocol, as it could result in an arbitrary loss of precision.
func PriceToFloat32ForLogging(price uint64, exponent types.Exponent) float32 {
	// We're not concerned about truncation here.
	priceFloat32, _ := decimal.NewFromBigInt(new(big.Int).SetUint64(price), exponent).BigFloat().Float32()
	return priceFloat32
}
