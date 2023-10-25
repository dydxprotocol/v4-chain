package prices

import (
	"math/big"

	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type MarketParamPriceModifierOption func(cp *pricestypes.MarketParamPrice)

func WithId(id uint32) MarketParamPriceModifierOption {
	return func(cp *pricestypes.MarketParamPrice) {
		cp.Param.Id = id
		cp.Price.Id = id
	}
}

func WithPair(pair string) MarketParamPriceModifierOption {
	return func(cp *pricestypes.MarketParamPrice) {
		cp.Param.Pair = pair
	}
}

func WithExponent(exp int32) MarketParamPriceModifierOption {
	return func(cp *pricestypes.MarketParamPrice) {
		cp.Param.Exponent = exp
		cp.Price.Exponent = exp
	}
}

func WithPriceValue(price uint64) MarketParamPriceModifierOption {
	return func(cp *pricestypes.MarketParamPrice) {
		cp.Price.Price = price
	}
}

func WithMinExchanges(minExchanges uint32) MarketParamPriceModifierOption {
	return func(cp *pricestypes.MarketParamPrice) {
		cp.Param.MinExchanges = minExchanges
	}
}

func WithExchangeConfigJson(configJson string) MarketParamPriceModifierOption {
	return func(cp *pricestypes.MarketParamPrice) {
		cp.Param.ExchangeConfigJson = configJson
	}
}

// GenerateMarketParamPrice returns a `MarketParamPrice` object set to default values.
// Passing in `MarketParamPriceModifierOption` methods alters the value of the `MarketParamPrice` returned.
// It will start with the default, valid `MarketParamPrice` value defined within the method
// and make the requested modifications before returning the object.
//
// Example usage:
// `GenerateMarketParamPrice(WithId(10))`
// This will start with the default `MarketParamPrice` object defined within the method and
// return the newly-created object after overriding the values of
// `Id` to 10.
func GenerateMarketParamPrice(optionalModifications ...MarketParamPriceModifierOption) *pricestypes.MarketParamPrice {
	marketParamPrice := &pricestypes.MarketParamPrice{
		Param: pricestypes.MarketParam{
			Id:                 0,
			Pair:               "BTC-USDC",
			MinExchanges:       3,
			MinPriceChangePpm:  100,
			Exponent:           -8,
			ExchangeConfigJson: "{}",
		},
		Price: pricestypes.MarketPrice{
			Id:       0,
			Exponent: -8,
			Price:    100000,
		},
	}

	for _, opt := range optionalModifications {
		opt(marketParamPrice)
	}

	return marketParamPrice
}

func MustHumanPriceToMarketPrice(
	humanPrice string,
	exponent int32,
) (marketPrice uint64) {
	// Ensure the exponent is negative
	if exponent >= 0 {
		panic("Only negative exponents are supported")
	}

	// Parse the humanPrice string to a big rational
	ratValue, ok := new(big.Rat).SetString(humanPrice)
	if !ok {
		panic("Failed to parse humanPrice to big.Rat")
	}

	// Convert exponent to its absolute value for calculations
	absResolution := int64(-exponent)

	// Create a multiplier which is 10 raised to the power of the absolute exponent
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(absResolution), nil)

	// Multiply the parsed humanPrice with the multiplier
	ratValue.Mul(ratValue, new(big.Rat).SetInt(multiplier))

	// Convert the result to an unsigned 64-bit integer
	resultInt := ratValue.Num() // Get the numerator which now represents the whole value

	return resultInt.Uint64()
}
