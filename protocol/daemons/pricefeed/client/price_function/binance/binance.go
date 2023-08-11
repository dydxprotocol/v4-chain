package binance

import (
	"encoding/json"
	"fmt"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4/lib"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

var (
	validate *validator.Validate
)

// BinanceResponseBody is the response body for the request to the GET
// https://api.binance.us/api/v3/ticker/24hr?symbol=$ Binance API
// Note that not all response fields are included here. See the API response docs for more information.
// https://binance-docs.github.io/apidocs/spot/en/#24hr-ticker-price-change-statistics
type BinanceResponseBody struct {
	// Only relevant fields of the response are included
	AskPrice  string `json:"askPrice" validate:"required,positive-float-string"`
	BidPrice  string `json:"bidPrice" validate:"required,positive-float-string"`
	LastPrice string `json:"lastPrice" validate:"required,positive-float-string"`
}

// unmarshalBinanceResponse converts a raw JSON string representation of the ticker REST API response from
// Binance into a strongly typed struct representation of relevant response fields.
func unmarshalBinanceResponse(body io.ReadCloser) (*BinanceResponseBody, error) {
	var responseBody BinanceResponseBody
	err := json.NewDecoder(body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	if validate == nil {
		validate, err = price_function.GetApiResponseValidator()
		if err != nil {
			return nil, fmt.Errorf("Error creating API response validator (%w)", err)
		}
	}

	err = validate.Struct(responseBody)
	if err != nil {
		return nil, err
	}
	return &responseBody, nil
}

// BinanceFunction transforms an API response from Binance or BinanceUS into a map of market symbols
// to prices that have been shifted by a market specific exponent.
func BinancePriceFunction(
	response *http.Response,
	marketPriceExponent map[string]int32,
	medianizer lib.Medianizer,
) (marketSymbolsToPrice map[string]uint64, unavailableSymbols map[string]error, err error) {
	// Get market symbol and value of exponent. The API response should only contain information
	// for one market.
	marketSymbol, exponent, err := price_function.GetOnlyMarketSymbolAndExponent(
		marketPriceExponent,
		exchange_common.EXCHANGE_NAME_BINANCE,
	)
	if err != nil {
		return nil, nil, err
	}

	responseBody, err := unmarshalBinanceResponse(response.Body)
	// The most likely failure here would be due to a missing ticker, so this will be
	// percolated up as a missing symbol to maintain price function semantics.
	if err != nil {
		return nil, map[string]error{marketSymbol: err}, nil
	}

	// Get big float values from transformed response prices.
	bigFloatSlice, err := lib.ConvertStringSliceToBigFloatSlice(
		[]string{responseBody.AskPrice, responseBody.BidPrice, responseBody.LastPrice})
	if err != nil {
		return nil, nil, err
	}

	// Get the median uint64 value from the slice of big float price values.
	medianPrice, err := price_function.GetUint64MedianFromReverseShiftedBigFloatValues(
		bigFloatSlice,
		exponent,
		medianizer,
	)
	if err != nil {
		return nil, nil, err
	}

	return map[string]uint64{marketSymbol: medianPrice}, nil, nil
}
