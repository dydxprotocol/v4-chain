package bitfinex

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"math/big"
	"net/http"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4/lib"
)

// These indices into the REST API response are defined in https://docs.bitfinex.com/reference/rest-public-ticker
const (
	BidPriceIndex  = 0
	AskPriceIndex  = 2
	LastPriceIndex = 6
	// We don't need all 10 fields, but 10 fields indicates this is a valid API response. See above link
	// for documentation on the response format.
	BitfinexResponseLength = 10
)

// BitfinexResponseBody is our representation of the response body for the request to the GET
// The raw response is a slice of floats. We use this constructed response to enable stricter
// validation.
// https://api.bitfinex.com/v2/ticker/$ Bitfinex API
// Note that not all response values are included here. See the API response docs for more information.
// Documentation: https://docs.bitfinex.com/reference/rest-public-ticker
type BitfinexResponseBody struct {
	BidPrice  float64 `validate:"gt=0"`
	AskPrice  float64 `validate:"gt=0"`
	LastPrice float64 `validate:"gt=0"`
}

// bitfinexRawResponse is the go representation of the raw response format
type bitfinexRawResponse []float64

// unmarshalBinfinexResponse converts a raw JSON string representation of the ticker REST API response from
// Bitfinex into a strongly typed struct representation of relevant response fields.
func unmarshalBitfinexResponse(body io.ReadCloser) (*BitfinexResponseBody, error) {
	var responseBody BitfinexResponseBody
	var rawResponse bitfinexRawResponse
	err := json.NewDecoder(body).Decode(&rawResponse)
	if err != nil {
		return nil, err
	}

	// Verify the API response is the expected length.
	if len(rawResponse) != BitfinexResponseLength {
		return nil, fmt.Errorf(
			"Invalid response body length for %s with length of: %v, expected length %v",
			exchange_common.EXCHANGE_NAME_BITFINEX,
			len(rawResponse),
			BitfinexResponseLength,
		)
	}
	// Manually assign relevant slice fields to an annotated struct and validate field values for stricter validation.
	responseBody.BidPrice = rawResponse[BidPriceIndex]
	responseBody.AskPrice = rawResponse[AskPriceIndex]
	responseBody.LastPrice = rawResponse[LastPriceIndex]

	validate := validator.New()
	err = validate.Struct(responseBody)
	if err != nil {
		return nil, err
	}
	return &responseBody, nil
}

// BitfinexPriceFunction transforms an API response from Bitfinex into a map of market symbols
// to prices that have been shifted by a market specific exponent.
func BitfinexPriceFunction(
	response *http.Response,
	marketPriceExponent map[string]int32,
	medianizer lib.Medianizer,
) (marketSymbolsToPrice map[string]uint64, unavailableSymbols map[string]error, err error) {
	// Get market symbol and value of exponent. The API response should only contain information
	// for one market.
	marketSymbol, exponent, err := price_function.GetOnlyMarketSymbolAndExponent(
		marketPriceExponent,
		exchange_common.EXCHANGE_NAME_BITFINEX,
	)
	if err != nil {
		return nil, nil, err
	}

	// Parse response body.
	responseBody, err := unmarshalBitfinexResponse(response.Body)
	// The most likely failure here would be due to a missing ticker, so this will be
	// percolated up as a missing symbol to maintain price function semantics.
	if err != nil {
		return nil, map[string]error{marketSymbol: err}, nil
	}

	// Get big float values from transformed response prices.
	bigFloatSlice := []*big.Float{
		new(big.Float).SetFloat64(responseBody.BidPrice),
		new(big.Float).SetFloat64(responseBody.AskPrice),
		new(big.Float).SetFloat64(responseBody.LastPrice),
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
