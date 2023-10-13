package crypto_com

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// CryptoComResponseBody defines the overall CryptoCom response.
type CryptoComResponseBody struct {
	Code   uint32                  `json:"code" validate:"required"`
	Result CryptoComResponseResult `json:"result" validate:"required"`
}

// CryptoComResponseResult defines the `result` field of CryptoCom response.
type CryptoComResponseResult struct {
	Tickers []CryptoComTicker `json:"data" validate:"required"`
}

// CryptoComTicker is our representation of ticker information returned in CryptoCom response.
// Need to implement interface `Ticker` in util.go.
// Note: CryptoCom returns `null` for bids and asks if there are none on the orderbook, in which
// case we mark the ticker as unavailable.
type CryptoComTicker struct {
	Pair      string `json:"i" validate:"required"`
	AskPrice  string `json:"k" validate:"required,positive-float-string"`
	BidPrice  string `json:"b" validate:"required,positive-float-string"`
	LastPrice string `json:"a" validate:"required,positive-float-string"`
}

// Ensure that CryptoComTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*CryptoComTicker)(nil)

func (t CryptoComTicker) GetPair() string {
	return t.Pair
}

func (t CryptoComTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t CryptoComTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t CryptoComTicker) GetLastPrice() string {
	return t.LastPrice
}

// CryptoComPriceFunction transforms an API response from CryptoCom into a map of tickers to prices that have been
// shifted by a market specific exponent.
func CryptoComPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body into a list of tickers.
	var cryptoComResponseBody CryptoComResponseBody
	err = json.NewDecoder(response.Body).Decode(&cryptoComResponseBody)
	if err != nil {
		return nil, nil, err
	}

	if cryptoComResponseBody.Code != 0 {
		return nil, nil, errors.New("response code is not 0")
	}

	return price_function.GetMedianPricesFromTickers(
		cryptoComResponseBody.Result.Tickers,
		tickerToExponent,
		resolver,
	)
}
