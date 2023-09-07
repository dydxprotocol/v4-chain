package prices

import (
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
			Id:                0,
			Pair:              "BTC-USDC",
			MinExchanges:      3,
			MinPriceChangePpm: 100,
			Exponent:          -8,
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
