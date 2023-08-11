package binance_test

import (
	"errors"
	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/binance"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	binanceResponseString = `{"askPrice": "1368.5100", "bidPrice": "1368.0800", "lastPrice": "1368.2100"}`
)

func TestBinancePriceFunction_Mixed(t *testing.T) {
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
		"Failure - Empty market price exponent map": {
			responseJsonString: binanceResponseString,
			exponentMap:        map[string]int32{},
			expectedError: errors.New(
				"Invalid market price exponent map for Binance price function of length: 0, expected length 1",
			),
		},
		"Failure - Two values in the market price exponent map": {
			responseJsonString: binanceResponseString,
			exponentMap: map[string]int32{
				testutil.ETHUSDC: 2,
				testutil.BTCUSDC: 3,
			},
			expectedError: errors.New(
				"Invalid market price exponent map for Binance price function of length: 2, expected length 1",
			),
		},
		"Unavailable - invalid response": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "1368.0800", "lastPrice": "1368.2100",}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New("invalid character '}' looking for beginning of object key string"),
			},
		},
		"Unavailable - invalid type in response: number": {
			// Invalid due to integer bidPrice when string was expected.
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": 1368.0800, "lastPrice": "1368.2100"}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New("json: cannot unmarshal number into Go struct field " +
					"BinanceResponseBody.bidPrice of type string"),
			},
		},
		"Unavailable - invalid type in response: malformed string": {
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "Not a Number", "lastPrice": "1368.2100"}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New("Key: 'BinanceResponseBody.BidPrice' Error:Field validation for " +
					"'BidPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: `{}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New(
					"Key: 'BinanceResponseBody.AskPrice' Error:Field validation for 'AskPrice' failed on the 'required' tag\n" +
						"Key: 'BinanceResponseBody.BidPrice' Error:Field validation for 'BidPrice' failed on the 'required' tag\n" +
						"Key: 'BinanceResponseBody.LastPrice' Error:Field validation for 'LastPrice' failed on the 'required' tag",
				),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "1368.0800"}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New(
					"Key: 'BinanceResponseBody.LastPrice' Error:Field validation for 'LastPrice' failed on the 'required' tag",
				),
			},
		},
		"Failure - overflow due to massively negative exponent": {
			responseJsonString: binanceResponseString,
			exponentMap:        map[string]int32{testutil.ETHUSDC: -3000},
			expectedError:      errors.New("value overflows uint64"),
		},
		"Failure - medianization error": {
			responseJsonString:  binanceResponseString,
			exponentMap:         testutil.ExponentSymbolMap,
			medianFunctionFails: true,
			expectedError:       testutil.MedianizationError,
		},
		"Success - extra fields": {
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "1368.0800", "lastPrice": "1368.2100", "extra": false}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedPriceMap: map[string]uint64{
				testutil.ETHUSDC: uint64(1_368_210_000),
			},
		},
		"Success - integers": {
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "1368", "lastPrice": "1368.2100"}`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedPriceMap: map[string]uint64{
				testutil.ETHUSDC: uint64(1_368_210_000),
			},
		},
		"Success - negative exponent": {
			responseJsonString: binanceResponseString,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedPriceMap: map[string]uint64{
				testutil.ETHUSDC: uint64(1_368_210_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			responseJsonString: `
			{"askPrice": "1368.5100", "bidPrice": "1368.0800", "lastPrice": "1368.211234656788"}
			`,
			exponentMap: testutil.ExponentSymbolMap,
			expectedPriceMap: map[string]uint64{
				testutil.ETHUSDC: uint64(1_368_211_234),
			},
		},
		"Success - positive exponent": {
			responseJsonString: binanceResponseString,
			exponentMap: map[string]int32{
				testutil.ETHUSDC: 2,
			},
			expectedPriceMap: map[string]uint64{
				testutil.ETHUSDC: uint64(13),
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
				prices, unavailable, err = binance.BinancePriceFunction(response, tc.exponentMap, medianizer)
			} else {
				prices, unavailable, err = binance.BinancePriceFunction(response, tc.exponentMap, &lib.MedianizerImpl{})
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
