package okx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// OkxResponseBody defines the overall Okx response.
type OkxResponseBody struct {
	Code    string      `json:"code" validate:"required"`
	Tickers []OkxTicker `json:"data" validate:"required"`
}

// OkxTicker is our representation of ticker information returned in Okx response.
// OkxTicker implements interface `Ticker` in util.go.
type OkxTicker struct {
	Pair      string `json:"instId" validate:"required"`
	AskPrice  string `json:"askPx" validate:"required,positive-float-string"`
	BidPrice  string `json:"bidPx" validate:"required,positive-float-string"`
	LastPrice string `json:"last" validate:"required,positive-float-string"`
}

// Ensure that OkxTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*OkxTicker)(nil)

func (t OkxTicker) GetPair() string {
	return t.Pair
}

func (t OkxTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t OkxTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t OkxTicker) GetLastPrice() string {
	return t.LastPrice
}

// OkxPriceFunction transforms an API response from Okx into a map of tickers
// to prices that have been shifted by a market specific exponent.
func OkxPriceFunction(
	response *http.Response,
	marketPriceExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body.
	var okxResponseBody OkxResponseBody
	err = json.NewDecoder(response.Body).Decode(&okxResponseBody)
	if err != nil {
		return nil, nil, err
	}

	if okxResponseBody.Code != "0" {
		return nil, nil, errors.New(`okx response code is not "0"`)
	}

	return price_function.GetMedianPricesFromTickers(
		okxResponseBody.Tickers,
		marketPriceExponent,
		resolver,
	)
}
