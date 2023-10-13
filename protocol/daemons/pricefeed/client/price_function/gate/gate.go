package gate

import (
	"encoding/json"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// GateTicker is our representation of ticker information returned in Gate response.
// Need to implement interface `Ticker` in util.go.
type GateTicker struct {
	Pair      string `json:"currency_pair" validate:"required"`
	AskPrice  string `json:"lowest_ask" validate:"required,positive-float-string"`
	BidPrice  string `json:"highest_bid" validate:"required,positive-float-string"`
	LastPrice string `json:"last" validate:"required,positive-float-string"`
}

// Ensure that GateTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*GateTicker)(nil)

func (t GateTicker) GetPair() string {
	return t.Pair
}

func (t GateTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t GateTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t GateTicker) GetLastPrice() string {
	return t.LastPrice
}

// GatePriceFunction transforms an API response from Gate into a map of tickers to prices that have been
// shifted by a market specific exponent.
func GatePriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body into a list of tickers.
	var gateTickers []GateTicker
	err = json.NewDecoder(response.Body).Decode(&gateTickers)
	if err != nil {
		return nil, nil, err
	}

	return price_function.GetMedianPricesFromTickers(
		gateTickers,
		tickerToExponent,
		resolver,
	)
}
