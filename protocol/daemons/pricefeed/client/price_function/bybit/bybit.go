package bybit

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// BybitResponseBody defines the overall Bybit response.
type BybitResponseBody struct {
	RetCode uint32              `json:"retCode" validate:"required"`
	Result  BybitResponseResult `json:"result" validate:"required"`
}

// BybitResponseResult defines the `result` field of Bybit response.
type BybitResponseResult struct {
	Tickers []BybitTicker `json:"list" validate:"required"`
}

// BybitTicker is our representation of ticker information returned in Bybit response.
// It implements the Ticker interface in util.go.
type BybitTicker struct {
	Pair      string `json:"symbol" validate:"required"`
	AskPrice  string `json:"ask1Price" validate:"required,positive-float-string"`
	BidPrice  string `json:"bid1Price" validate:"required,positive-float-string"`
	LastPrice string `json:"lastPrice" validate:"required,positive-float-string"`
}

// Ensure that BybitTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*BybitTicker)(nil)

func (t BybitTicker) GetPair() string {
	return t.Pair
}

func (t BybitTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t BybitTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t BybitTicker) GetLastPrice() string {
	return t.LastPrice
}

// BybitPriceFunction transforms an API response from Bybit into a map of tickers to prices that have been
// shifted by a market specific exponent.
func BybitPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body into a list of tickers.
	var bybitResponseBody BybitResponseBody
	err = json.NewDecoder(response.Body).Decode(&bybitResponseBody)
	if err != nil {
		return nil, nil, err
	}

	if bybitResponseBody.RetCode != 0 {
		return nil, nil, errors.New("response code is not 0")
	}

	return price_function.GetMedianPricesFromTickers(
		bybitResponseBody.Result.Tickers,
		tickerToExponent,
		resolver,
	)
}
