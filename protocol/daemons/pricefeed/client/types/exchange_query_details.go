package types

import (
	"net/http"

	"github.com/dydxprotocol/v4/lib"
)

// ExchangeQueryDetails represents the information needed to query a specific exchange.
type ExchangeQueryDetails struct {
	Exchange      ExchangeFeedId
	Url           string              // url to query exchange
	MarketSymbols map[MarketId]string // map of market Id to exchange-specific symbol
	// function to get a map of market symbols to prices from an exchange's response
	PriceFunction func(
		response *http.Response,
		marketSymbolPriceExponentMap map[string]int32,
		medianizer lib.Medianizer,
	) (
		marketSymbolsToPrice map[string]uint64,
		unavailableSymbols map[string]error,
		err error,
	)
	IsMultiMarket bool
}
