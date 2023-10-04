package coinbase_pro

import (
	"encoding/json"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// CoinbaseProTicker is our representation of ticker information returned in CoinbasePro response.
// CoinbaseProTicker implements interface `Ticker` in util.go.
type CoinbaseProTicker struct {
	// `Pair` is not part of API response but can be set manually to reuse existing helper functions.
	Pair      string `validate:"required"`
	AskPrice  string `json:"ask" validate:"required,positive-float-string"`
	BidPrice  string `json:"bid" validate:"required,positive-float-string"`
	LastPrice string `json:"price" validate:"required,positive-float-string"`
}

// Ensure that CoinbaseProTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*CoinbaseProTicker)(nil)

func (t CoinbaseProTicker) GetPair() string {
	return t.Pair
}

func (t CoinbaseProTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t CoinbaseProTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t CoinbaseProTicker) GetLastPrice() string {
	return t.LastPrice
}

// CoinbaseProPriceFunction transforms an API response from CoinbasePro into a map of tickers
// to prices that have been shifted by a market specific exponent.
func CoinbaseProPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Get ticker. The API response should only contain information for one market.
	ticker, _, err := price_function.GetOnlyTickerAndExponent(
		tickerToExponent,
		exchange_common.EXCHANGE_ID_COINBASE_PRO,
	)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal response body.
	var coinbaseProTicker CoinbaseProTicker
	err = json.NewDecoder(response.Body).Decode(&coinbaseProTicker)
	if err != nil {
		return nil, nil, err
	}

	// Invoke `GetMedianPricesFromTickers` on a list of one ticker whose `Pair`
	// matches the only ticker in `tickerToExponent`.
	coinbaseProTicker.Pair = ticker
	return price_function.GetMedianPricesFromTickers(
		[]CoinbaseProTicker{coinbaseProTicker},
		tickerToExponent,
		resolver,
	)
}
