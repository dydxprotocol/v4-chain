package marketmap

import (
	dydx "github.com/dydxprotocol/slinky/providers/apis/dydx"
	dydxtypes "github.com/dydxprotocol/slinky/providers/apis/dydx/types"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// Construct a MarketMap struct from a slice of MarketParams
func ConstructMarketMapFromParams(
	allMarketParams []pricestypes.MarketParam,
) (marketmaptypes.MarketMap, error) {
	var mpr dydxtypes.QueryAllMarketParamsResponse
	for _, mp := range allMarketParams {
		mpr.MarketParams = append(mpr.MarketParams, dydxtypes.MarketParam{
			Id:                 mp.Id,
			Pair:               mp.Pair,
			Exponent:           mp.Exponent,
			MinExchanges:       mp.MinExchanges,
			MinPriceChangePpm:  mp.MinPriceChangePpm,
			ExchangeConfigJson: mp.ExchangeConfigJson,
		})
	}
	mm, err := dydx.ConvertMarketParamsToMarketMap(mpr)
	if err != nil {
		return marketmaptypes.MarketMap{}, err
	}

	return mm.MarketMap, nil
}
