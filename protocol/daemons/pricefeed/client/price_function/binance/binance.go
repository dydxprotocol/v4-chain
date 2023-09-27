package binance

import (
	"encoding/json"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// BinanceTicker is our representation of ticker information returned in Binance response.
// It implements interface `Ticker` in util.go.
type BinanceTicker struct {
	Pair      string `json:"symbol" validate:"required"`
	AskPrice  string `json:"askPrice" validate:"required,positive-float-string"`
	BidPrice  string `json:"bidPrice" validate:"required,positive-float-string"`
	LastPrice string `json:"lastPrice" validate:"required,positive-float-string"`
}

// Ensure that BinanceTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*BinanceTicker)(nil)

func (t BinanceTicker) GetPair() string {
	// needs to be wrapped in quotes to be consistent with the API request format.
	return t.Pair
}

func (t BinanceTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t BinanceTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t BinanceTicker) GetLastPrice() string {
	return t.LastPrice
}

// BinancePriceFunction transforms an API response from Binance into a map of tickers to prices that have been
// shifted by a market specific exponent.
func BinancePriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body into a list of tickers.
	var binanceTickers []BinanceTicker
	err = json.NewDecoder(response.Body).Decode(&binanceTickers)
	if err != nil {
		return nil, nil, err
	}

	return price_function.GetMedianPricesFromTickers(
		binanceTickers,
		tickerToExponent,
		resolver,
	)
}
