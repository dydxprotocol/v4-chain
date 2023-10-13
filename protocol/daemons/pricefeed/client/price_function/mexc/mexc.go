package mexc

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// MexcResponseBody defines the overall Mexc response.
type MexcResponseBody struct {
	Code    uint32       `json:"code" validate:"required"`
	Tickers []MexcTicker `json:"data" validate:"required"`
}

// MexcTicker is our representation of ticker information returned in Mexc response.
// MexcTicker implements interface `Ticker` in util.go.
type MexcTicker struct {
	Pair      string `json:"symbol" validate:"required"`
	AskPrice  string `json:"ask" validate:"required,positive-float-string"`
	BidPrice  string `json:"bid" validate:"required,positive-float-string"`
	LastPrice string `json:"last" validate:"required,positive-float-string"`
}

// Ensure that MexcTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*MexcTicker)(nil)

func (t MexcTicker) GetPair() string {
	return t.Pair
}

func (t MexcTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t MexcTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t MexcTicker) GetLastPrice() string {
	return t.LastPrice
}

// MexcPriceFunction transforms an API response from Mexc into a map of tickers to prices that have been
// shifted by a market specific exponent.
func MexcPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body.
	var mexcResponseBody MexcResponseBody
	err = json.NewDecoder(response.Body).Decode(&mexcResponseBody)
	if err != nil {
		return nil, nil, err
	}

	if mexcResponseBody.Code != 200 {
		return nil, nil, errors.New(`mexc response code is not 200`)
	}

	return price_function.GetMedianPricesFromTickers(
		mexcResponseBody.Tickers,
		tickerToExponent,
		resolver,
	)
}
