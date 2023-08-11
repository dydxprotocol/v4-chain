package kraken

import (
	"encoding/json"
	"fmt"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4/lib"
	"net/http"
)

const (
	askPriceField   = "a"
	bidPriceField   = "b"
	closePriceField = "c"

	resultField = "result"
	errorField  = "error"
)

// TODO(DEC-753): add strict type
type KrakenResponseBody map[string]interface{}

// extractPriceFromTicker takes a map representation of the Kraken GetTicker response for a single
// ticker and computes the market price based on the median of ask price, bid price and last trade
// close price, shifted by the market-specific exponent.
func extractPriceFromTicker(
	ticker map[string]interface{},
	exponent int32,
	medianizer lib.Medianizer,
) (uint64, error) {
	askPriceStr, err := price_function.ExtractFirstStringFromSliceField(ticker, askPriceField)
	if err != nil {
		return 0, err
	}
	bidPriceStr, err := price_function.ExtractFirstStringFromSliceField(ticker, bidPriceField)
	if err != nil {
		return 0, err
	}
	lastPriceStr, err := price_function.ExtractFirstStringFromSliceField(ticker, closePriceField)
	if err != nil {
		return 0, err
	}

	bigFloatSlice, err := lib.ConvertStringSliceToBigFloatSlice([]string{askPriceStr, bidPriceStr, lastPriceStr})
	if err != nil {
		return 0, err
	}

	// Get the median uint64 value from the slice of big float price values.
	medianPrice, err := price_function.GetUint64MedianFromReverseShiftedBigFloatValues(
		bigFloatSlice,
		exponent,
		medianizer,
	)
	if err != nil {
		return 0, err
	}

	return medianPrice, nil
}

// KrakenPriceFunction transforms an API response from Kraken into a map of market symbols
// to prices that have been shifted by a market-specific exponent.
func KrakenPriceFunction(
	response *http.Response,
	marketPriceExponent map[string]int32,
	medianizer lib.Medianizer,
) (marketSymbolsToPrice map[string]uint64, unavailableSymbols map[string]error, err error) {
	var responseBody KrakenResponseBody
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, nil, err
	}

	results, ok := responseBody[resultField].(map[string]interface{})
	if !ok {
		errors, ok := responseBody[errorField]
		if ok {
			return nil, nil, fmt.Errorf("kraken API call error: %v", errors)
		} else {
			return nil, nil, fmt.Errorf("kraken API call error: %v", responseBody)
		}
	}

	marketSymbolsToPrice = make(map[string]uint64, len(marketPriceExponent))
	unavailableSymbols = make(map[string]error)
	for symbol, exponent := range marketPriceExponent {
		result, ok := results[symbol]
		if !ok {
			unavailableSymbols[symbol] = fmt.Errorf("no ticker found for market symbol %v", symbol)
			continue
		}
		resultMap := result.(map[string]interface{})
		medianPrice, err := extractPriceFromTicker(resultMap, exponent, medianizer)
		if err != nil {
			unavailableSymbols[symbol] = err
			continue
		}
		marketSymbolsToPrice[symbol] = medianPrice
	}
	return marketSymbolsToPrice, unavailableSymbols, nil
}
