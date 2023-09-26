package huobi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// HuobiResponseBody defines the overall Huobi response.
type HuobiResponseBody struct {
	Status  string        `json:"status" validate:"required"`
	Tickers []HuobiTicker `json:"data" validate:"required"`
}

// HuobiTicker is our representation of ticker information returned in Huobi response.
// HuobiTicker implements interface `Ticker` in util.go.
type HuobiTicker struct {
	Pair      string  `json:"symbol" validate:"required"`
	AskPrice  float64 `json:"ask" validate:"required,gt=0"`
	BidPrice  float64 `json:"bid" validate:"required,gt=0"`
	LastPrice float64 `json:"close" validate:"required,gt=0"`
}

// Ensure that HuobiTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*HuobiTicker)(nil)

func (t HuobiTicker) GetPair() string {
	return t.Pair
}

func (t HuobiTicker) GetAskPrice() string {
	return price_function.ConvertFloat64ToString(t.AskPrice)
}

func (t HuobiTicker) GetBidPrice() string {
	return price_function.ConvertFloat64ToString(t.BidPrice)
}

func (t HuobiTicker) GetLastPrice() string {
	return price_function.ConvertFloat64ToString(t.LastPrice)
}

// HuobiPriceFunction transforms an API response from Huobi into a map of tickers to prices that have been
// shifted by a market specific exponent.
func HuobiPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body.
	var huobiResponseBody HuobiResponseBody
	err = json.NewDecoder(response.Body).Decode(&huobiResponseBody)
	if err != nil {
		return nil, nil, err
	}

	if huobiResponseBody.Status != "ok" {
		return nil, nil, errors.New(`huobi response status is not "ok"`)
	}

	return price_function.GetMedianPricesFromTickers(
		huobiResponseBody.Tickers,
		tickerToExponent,
		resolver,
	)
}
