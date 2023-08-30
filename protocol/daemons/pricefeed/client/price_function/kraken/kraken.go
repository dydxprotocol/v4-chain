package kraken

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"strings"
)

var (
	validate *validator.Validate
)

// https://api.kraken.com/0/public/Ticker?pair=$
// https://docs.kraken.com/rest/#tag/Market-Data/operation/getTickerInformation
type KrakenTickerResult struct {
	pair            string
	AskPriceStats   []string `json:"a" validate:"len=3,dive,positive-float-string"`
	BidPriceStats   []string `json:"b" validate:"len=3,dive,positive-float-string"`
	ClosePriceStats []string `json:"c" validate:"len=2,dive,positive-float-string"`
}

func (ktr KrakenTickerResult) WithPair(pair string) KrakenTickerResult {
	ktr.pair = pair
	return ktr
}

func (ktr KrakenTickerResult) GetPair() string {
	return ktr.pair
}

func (ktr KrakenTickerResult) GetAskPrice() string {
	return ktr.AskPriceStats[0]
}

func (ktr KrakenTickerResult) GetBidPrice() string {
	return ktr.BidPriceStats[0]
}

func (ktr KrakenTickerResult) GetLastPrice() string {
	return ktr.ClosePriceStats[0]
}

type KrakenResponseBody struct {
	// As of this time, the Kraken API response is all-or-nothing - either valid ticker data, or one or more errors,
	// but not both. We enforce this expectation by defining mutual exclusivity in the validation tags of the Errors
	// field so that any validated API result always meets our expectation in the response parsing logic.
	Errors  []string                      `json:"error" validate:"omitempty"`
	Tickers map[string]KrakenTickerResult `validate:"required_without=Errors,excluded_with=Errors,dive" json:"result"`
}

// unmarshalKrakenResponse converts a raw JSON string representation of the ticker REST API response from
// Kraken into a strongly typed struct representation of relevant response fields.
func unmarshalKrakenResponse(body io.ReadCloser) (*KrakenResponseBody, error) {
	var responseBody KrakenResponseBody
	err := json.NewDecoder(body).Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("kraken API response JSON parse error (%w)", err)
	}

	// The Kraken API will return an empty list of errors with an API result containing valid tickers. However, it's
	// easier for us to validate that there were no errors if this field is set to nil whenever it's empty.
	if len(responseBody.Errors) == 0 {
		responseBody.Errors = nil
	}

	if validate == nil {
		validate, err = price_function.GetApiResponseValidator()
		if err != nil {
			return nil, fmt.Errorf("Error creating API response validator (%w)", err)
		}
	}

	err = validate.Struct(responseBody)
	if err != nil {
		return nil, fmt.Errorf("kraken API response validation error (%w)", err)
	}
	return &responseBody, nil
}

// KrakenPriceFunction transforms an API response from Kraken into a map of tickers to prices that have been
// shifted by a market-specific exponent.
func KrakenPriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	medianizer lib.Medianizer,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	responseBody, err := unmarshalKrakenResponse(response.Body)
	if err != nil {
		return nil, nil, err
	}

	if len(responseBody.Errors) > 0 {
		// TODO(CORE-185): Update to Go 1.20 and replace strings.Join with errors.Join.
		apiCallError := fmt.Errorf(
			"kraken API call error: %w", errors.New(strings.Join(responseBody.Errors, ", ")),
		)
		return nil, nil, apiCallError
	}

	tickers := make([]KrakenTickerResult, 0, len(responseBody.Tickers))
	for pair, ticker := range responseBody.Tickers {
		tickers = append(tickers, ticker.WithPair(pair))
	}

	return price_function.GetMedianPricesFromTickers(
		tickers,
		tickerToExponent,
		medianizer,
	)
}
