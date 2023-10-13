package kucoin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// KucoinResponseBody defines the overall Kucoin response.
type KucoinResponseBody struct {
	Code string             `json:"code" validate:"required"`
	Data KucoinResponseData `json:"data" validate:"required"`
}

// KucoinResponseData defines the `data` field of Kucoin response.
type KucoinResponseData struct {
	Tickers []KucoinTicker `json:"ticker" validate:"required"`
}

// KucoinTicker is our representation of ticker information returned in Kucoin response.
// KucoinTicker implements interface `Ticker` in util.go.
type KucoinTicker struct {
	Pair      string `json:"symbol" validate:"required"`
	AskPrice  string `json:"sell" validate:"required,positive-float-string"`
	BidPrice  string `json:"buy" validate:"required,positive-float-string"`
	LastPrice string `json:"last" validate:"required,positive-float-string"`
}

// Ensure that KucoinTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*KucoinTicker)(nil)

func (t KucoinTicker) GetPair() string {
	return t.Pair
}

func (t KucoinTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t KucoinTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t KucoinTicker) GetLastPrice() string {
	return t.LastPrice
}

// KucoinPriceFunction transforms an API response from Kucoin into a map of tickers to prices that have been
// shifted by a market specific exponent.
func KucoinPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body.
	var kucoinResponseBody KucoinResponseBody
	err = json.NewDecoder(response.Body).Decode(&kucoinResponseBody)
	if err != nil {
		return nil, nil, err
	}

	if kucoinResponseBody.Code != "200000" {
		return nil, nil, errors.New(`kucoin response code is not "200000"`)
	}

	return price_function.GetMedianPricesFromTickers(
		kucoinResponseBody.Data.Tickers,
		tickerToExponent,
		resolver,
	)
}
