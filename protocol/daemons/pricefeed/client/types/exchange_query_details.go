package types

import (
	"net/http"

	"github.com/dydxprotocol/v4/lib"
)

// ExchangeQueryDetails represents the information needed to query a specific exchange.
type ExchangeQueryDetails struct {
	Exchange ExchangeId
	Url      string // url to query exchange
	// function to get a map of tickers to prices from an exchange's response
	PriceFunction func(
		response *http.Response,
		tickerToPriceExponent map[string]int32,
		medianizer lib.Medianizer,
	) (
		tickerToPrice map[string]uint64,
		unavailableTickers map[string]error,
		err error,
	)
	IsMultiMarket bool
}
