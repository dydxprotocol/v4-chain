package client

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

// ExchangeFeedIdMarketPriceTimestamp contains an `ExchangeFeedId` and an associated
// `types.MarketPriceTimestamp`. This type exists for convenience and clarity in testing the
// pricefeed client.
type ExchangeFeedIdMarketPriceTimestamp struct {
	ExchangeFeedId       types.ExchangeFeedId
	MarketPriceTimestamp *types.MarketPriceTimestamp
}
