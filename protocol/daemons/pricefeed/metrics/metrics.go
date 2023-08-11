package metrics

import (
	gometrics "github.com/armon/go-metrics"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/lib/metrics"
)

const (
	INVALID = "INVALID"
)

// GetLabelForMarketId converts a marketId uint32 into a human-readable symbol and then returns the
// label with the market symbol.
func GetLabelForMarketId(marketId types.MarketId) gometrics.Label {
	marketSymbol, exists := StaticMarketSymbols[marketId]
	if !exists {
		return metrics.GetLabelForStringValue(metrics.MarketId, INVALID)
	}

	return metrics.GetLabelForStringValue(metrics.MarketId, marketSymbol)
}

// GetLabelForExchangeFeedId converts an exchangeFeedId uint32 into a name and then
// returns the label with the name.
func GetLabelForExchangeFeedId(exchangeFeedId types.ExchangeFeedId) gometrics.Label {
	exchangeFeedName, exists := StaticExchangeNames[exchangeFeedId]
	if !exists {
		return metrics.GetLabelForStringValue(metrics.ExchangeFeedId, INVALID)
	}

	return metrics.GetLabelForStringValue(metrics.ExchangeFeedId, exchangeFeedName)
}
