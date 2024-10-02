package marketmap

import (
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	dydx "github.com/skip-mev/connect/v2/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/connect/v2/providers/apis/dydx/types"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
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
