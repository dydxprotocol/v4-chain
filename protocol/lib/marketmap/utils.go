package marketmap

import (
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

// Construct a MarketMap struct from a slice of MarketParams
func ConstructMarketMapFromParams(
	allMarketParams []pricestypes.MarketParam,
) (marketmaptypes.MarketMap, error) {
	// fill out config with dummy variables to pass validation.  This handler is only used to run the
	// ConvertMarketParamsToMarketMap member function.
	h, err := dydx.NewAPIHandler(zap.NewNop(), config.APIConfig{
		Enabled:          true,
		Timeout:          1,
		Interval:         1,
		ReconnectTimeout: 1,
		MaxQueries:       1,
		Atomic:           false,
		Endpoints:        []config.Endpoint{{URL: "upgrade"}},
		BatchSize:        0,
		Name:             dydx.Name,
	})
	if err != nil {
		return marketmaptypes.MarketMap{}, err
	}

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
	mm, err := h.ConvertMarketParamsToMarketMap(mpr)
	if err != nil {
		return marketmaptypes.MarketMap{}, err
	}

	return mm.MarketMap, nil
}
