package bitfinex

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// These indices into the REST API response are defined in https://docs.bitfinex.com/reference/rest-public-tickers
const (
	PairIndex      = 0
	BidPriceIndex  = 1
	AskPriceIndex  = 3
	LastPriceIndex = 7
	// We don't need all 11 fields, but 11 fields indicates this is a valid API response. See above link
	// for documentation on the response format.
	BitfinexResponseLength = 11
)

// BitfinexTicker is our representation of the ticker information in Bitfinex API response.
// The raw response is a slice of floats. We use this constructed response to enable stricter
// validation.
// BitfinexTicker implements interface `Ticker` in util.go.
type BitfinexTicker struct {
	Pair      string `validate:"required"`
	BidPrice  string `validate:"required,positive-float-string"`
	AskPrice  string `validate:"required,positive-float-string"`
	LastPrice string `validate:"required,positive-float-string"`
}

// Ensure that BitfinexTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*BitfinexTicker)(nil)

func (t BitfinexTicker) GetPair() string {
	return t.Pair
}

func (t BitfinexTicker) GetAskPrice() string {
	return t.AskPrice
}

func (t BitfinexTicker) GetBidPrice() string {
	return t.BidPrice
}

func (t BitfinexTicker) GetLastPrice() string {
	return t.LastPrice
}

// BitfinexPriceFunction transforms an API response from Bitfinex into a map of tickers
// to prices that have been shifted by a market specific exponent.
func BitfinexPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Unmarshal response body into raw format first.
	var rawResponse [][]interface{}
	err = json.NewDecoder(response.Body).Decode(&rawResponse)
	if err != nil {
		return nil, nil, err
	}

	// Convert raw tickers in response into a list of `BitfinexTicker`.
	bitfinexTickers := []BitfinexTicker{}
	invalidRawTickers := map[string]error{}
	for _, rawTicker := range rawResponse {
		// Verify raw ticker is the expected length. If not, continue to next.
		if len(rawTicker) != BitfinexResponseLength {
			continue
		}

		// Get `pair` as `string`. If invalid, continue on to next raw ticker.
		pair, ok := rawTicker[PairIndex].(string)
		if !ok {
			continue
		}
		// Get `bidPrice` as `float64`. If invalid, mark pair as invalid with bid price error.
		bidPrice, ok := rawTicker[BidPriceIndex].(float64)
		if !ok {
			invalidRawTickers[pair] = errors.New("invalid bid price in response - not a float64")
			continue
		}
		// Get `askPrice` as `float64`. If invalid, mark pair as invalid with ask price error.
		askPrice, ok := rawTicker[AskPriceIndex].(float64)
		if !ok {
			invalidRawTickers[pair] = errors.New("invalid ask price in response - not a float64")
			continue
		}
		// Get `lastPrice`. If invalid, mark pair as invalid with last price error.
		lastPrice, ok := rawTicker[LastPriceIndex].(float64)
		if !ok {
			invalidRawTickers[pair] = errors.New("invalid last price in response - not a float64")
			continue
		}
		bitfinexTickers = append(bitfinexTickers, BitfinexTicker{
			Pair:      pair,
			BidPrice:  price_function.ConvertFloat64ToString(bidPrice),
			AskPrice:  price_function.ConvertFloat64ToString(askPrice),
			LastPrice: price_function.ConvertFloat64ToString(lastPrice),
		})
	}

	// Calculate median price of each ticker in `tickerToExponent`.
	tickerToPrice, unavailableTickers, err = price_function.GetMedianPricesFromTickers(
		bitfinexTickers,
		tickerToExponent,
		resolver,
	)

	// Mark as unavailable requested tickers whose raw ticker response was invalid.
	for ticker, err := range invalidRawTickers {
		if _, exists := tickerToExponent[ticker]; exists {
			unavailableTickers[ticker] = err
		}
	}

	return tickerToPrice, unavailableTickers, err
}
