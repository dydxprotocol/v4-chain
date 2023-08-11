package kraken_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/kraken"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"
	"github.com/dydxprotocol/v4/testutil/stringutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	ETHUSDC_SYMBOL = "XETHZUSD"
	BTCUSDC_SYMBOL = "XXBTZUSD"
)

var (
	ValidSymbolMap = map[string]int32{
		ETHUSDC_SYMBOL: constants.StaticMarketPriceExponent[exchange_common.MARKET_ETH_USD],
		BTCUSDC_SYMBOL: constants.StaticMarketPriceExponent[exchange_common.MARKET_BTC_USD],
	}
	EthSymbolMap = map[string]int32{
		ETHUSDC_SYMBOL: constants.StaticMarketPriceExponent[exchange_common.MARKET_ETH_USD],
	}
)

// Take a test file with human-readable JSON, load it, strip all whitespace / newlines, and return a string
func readJsonTestFile(t *testing.T, fileName string) string {
	fileBytes, err := os.ReadFile(fmt.Sprintf("testdata/%v", fileName))
	require.NoError(t, err)
	return stringutils.StripSpaces(fileBytes)
}

func TestKrakenPriceFunction_Mixed(t *testing.T) {
	krakenValidResponseString := readJsonTestFile(t, "kraken_2_ticker_response.json")
	tests := map[string]struct {
		// parameters
		responseJsonString  string
		exponentMap         map[string]int32
		medianFunctionFails bool

		// expectations
		expectedPriceMap       map[string]uint64
		expectedUnavailableMap map[string]error
		expectedError          error
	}{
		"Failure - invalid response, not JSON": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{,}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedError:      errors.New("invalid character ',' looking for beginning of object key string"),
		},
		"Unavailable - float instead of string data type, missing": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{"a":[2105.8]}}}`,
			exponentMap:        ValidSymbolMap,
			expectedPriceMap:   map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: errors.New("expected nonempty string value for field a[0], but found 2105.8"),
				BTCUSDC_SYMBOL: errors.New("no ticker found for market symbol XXBTZUSD"),
			},
		},
		"Unavailable - invalid response: returned values aren't parsable as floats": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{"a":["cat"],"b":["1"],"c":["2"]}}}`,
			exponentMap:        EthSymbolMap,
			expectedPriceMap:   map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: errors.New("invalid, value is not a number: cat"),
			},
		},
		"Unavailable - invalid response: underflow due to invalid negative": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{"a":["-1234.56"],"b":["1"],"c":["2"]}}}`,
			exponentMap:        EthSymbolMap,
			expectedPriceMap:   map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: errors.New("value underflows uint64"),
			},
		},
		"Unavailable - invalid response: overflow due to negative exponent": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{"a":["1"],"b":["1"],"c":["2"]}}}`,
			exponentMap:        map[string]int32{ETHUSDC_SYMBOL: -3000},
			expectedPriceMap:   map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: errors.New("value overflows uint64"),
			},
		},
		"Unavailable - invalid response: missing expected response field": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{}}}`,
			exponentMap:        EthSymbolMap,
			expectedPriceMap:   map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: errors.New("expected non-empty list for fieldname 'a'"),
			},
		},
		"Mixed success, unavailable - one ticker invalid": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{"a":["1"],"b":"abc","c":["2"]},` +
				`"XXBTZUSD":{"a":["1"],"b":["1"],"c":["2"]}}}`,
			exponentMap: ValidSymbolMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_SYMBOL: uint64(100000),
			},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: errors.New("expected non-empty list for fieldname 'b'"),
			},
		},
		"Unavailable - fails on medianization error": {
			responseJsonString:  krakenValidResponseString,
			exponentMap:         ValidSymbolMap,
			medianFunctionFails: true,
			expectedPriceMap:    map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: testutil.MedianizationError,
				BTCUSDC_SYMBOL: testutil.MedianizationError,
			},
		},
		"Failure - Kraken API Error response": {
			responseJsonString: `{"error":["EQuery:Unknown asset pair"]}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedError:      errors.New("kraken API call error: [EQuery:Unknown asset pair]"),
		},
		"Failure - Kraken API Empty response": {
			responseJsonString: `{}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedError:      errors.New("kraken API call error: map[]"),
		},
		"Success: one market response": {
			responseJsonString: krakenValidResponseString,
			exponentMap:        EthSymbolMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_SYMBOL: uint64(1_888_000_000),
			},
		},
		"Success: two market response": {
			responseJsonString: krakenValidResponseString,
			exponentMap:        ValidSymbolMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_SYMBOL: uint64(1_888_000_000),
				BTCUSDC_SYMBOL: uint64(2_920_750_000),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			response := testutil.CreateResponseFromJson(tc.responseJsonString)

			var prices map[string]uint64
			var unavailable map[string]error
			var err error
			if tc.medianFunctionFails {
				medianizer := &mocks.Medianizer{}
				medianizer.On("MedianUint64", mock.Anything).Return(uint64(0), testutil.MedianizationError)
				prices, unavailable, err = kraken.KrakenPriceFunction(response, tc.exponentMap, medianizer)
			} else {
				prices, unavailable, err = kraken.KrakenPriceFunction(response, tc.exponentMap, &lib.MedianizerImpl{})
			}

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Nil(t, prices)
				require.Nil(t, unavailable)
			} else {
				require.Equal(t, tc.expectedPriceMap, prices)
				pricefeed.ErrorMapsEqual(t, tc.expectedUnavailableMap, unavailable)
				require.NoError(t, err)
			}
		})
	}
}
