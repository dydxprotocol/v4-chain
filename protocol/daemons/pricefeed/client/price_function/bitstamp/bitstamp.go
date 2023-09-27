package bitstamp

import (
	"encoding/json"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// BitstampTicker is our representation of ticker information returned in Bitstamp response.
// Need to implement interface `Ticker` in util.go.
type BitstampTicker struct {
	Pair      string `json:"pair" validate:"required"`
	AskPrice  string `json:"ask" validate:"required,positive-float-string"`
	BidPrice  string `json:"bid" validate:"required,positive-float-string"`
	LastPrice string `json:"last" validate:"required,positive-float-string"`
}

// Ensure that BitstampTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*BitstampTicker)(nil)

func (t BitstampTicker) GetPair() string {
	return t.Pair
}

func (t BitstampTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t BitstampTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t BitstampTicker) GetLastPrice() string {
	return t.LastPrice
}

// BitstampPriceFunction transforms an API response from Bitstamp into a map of tickers to prices that have been
// shifted by a market specific exponent.
func BitstampPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body into a list of tickers.
	var bitstampTickers []BitstampTicker
	err = json.NewDecoder(response.Body).Decode(&bitstampTickers)
	if err != nil {
		return nil, nil, err
	}

	return price_function.GetMedianPricesFromTickers(
		bitstampTickers,
		tickerToExponent,
		resolver,
	)
}
