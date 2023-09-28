package types

import (
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// ExchangeQueryDetails represents the information needed to query a specific exchange.
type ExchangeQueryDetails struct {
	Exchange ExchangeId
	// Url is the url to query the exchange.
	Url string
	// PriceFunction computes a map of tickers to prices from an exchange's response
	PriceFunction func(
		response *http.Response,
		tickerToPriceExponent map[string]int32,
		resolver types.Resolver,
	) (
		tickerToPrice map[string]uint64,
		unavailableTickers map[string]error,
		err error,
	)
	// IsMultiMarket indicates whether the url query response contains multiple tickers.
	IsMultiMarket bool
}
