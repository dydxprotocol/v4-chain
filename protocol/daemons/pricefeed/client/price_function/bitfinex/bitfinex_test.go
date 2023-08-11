package bitfinex_test

import (
	"errors"
	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	bitfinexResponseString = "[" +
		"1247.8," + // 0) Bid - used to get median
		"701.99719833," + // 1) Bid Period
		"1248," + // 2) Ask - used to get median
		"636.3844446100001," + // 3) Ask Period
		"-84," + // 4) Ask Size
		"-0.0631," + // 5) Daily Change
		"1248," + // 6) Last Price - used to get median
		"84500.66677816," + // 7) Volume
		"1412.76385631," + // 8) High
		"1220.11135781" + // 9) Low
		"]"
)

func TestBitfinexPriceFunction_Mixed(t *testing.T) {
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
			responseJsonString: bitfinexResponseString,
			exponentMap:        map[string]int32{},
			expectedError: errors.New(
				"Invalid market price exponent map for Bitfinex price function of length: 0, expected length 1",
			),
		},
		"Failure - Two values in the market price exponent map": {
			responseJsonString: bitfinexResponseString,
			exponentMap: map[string]int32{
				testutil.ETHUSDC: 2,
				testutil.BTCUSDC: 3,
			},
			expectedError: errors.New(
				"Invalid market price exponent map for Bitfinex price function of length: 2, expected length 1",
			),
		},
		"Unavailable - response is invalid JSON": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `[
				1247.8,
				701.99719833,
				1248,
				636.3844446100001,
				-84,
				-0.0631,
				1248,
				84500.66677816,
				1412.76385631,
				1220.11135781,
				]`,
			exponentMap: testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New("invalid character ']' looking for beginning of value"),
			},
		},
		"Unavailable - response contains invalid type": {
			// Invalid due to first value being a string when all values should be floating point
			// numbers.
			responseJsonString: `[
				"1247.8",
				701.99719833,
				1248,
				636.3844446100001,
				-84,
				-0.0631,
				1248,
				84500.66677816,
				1412.76385631,
				1220.11135781
				]`,
			exponentMap: testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New("json: cannot unmarshal string into Go value of type float64"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: `[]`,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New(
					"Invalid response body length for Bitfinex with length of: 0, expected length 10",
				),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: `[
				1247.8,
				701.99719833
				]`,
			exponentMap: testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New(
					"Invalid response body length for Bitfinex with length of: 2, expected length 10",
				),
			},
		},
		"Unavailable - response contains invalid negative": {
			// Bid price (index 0) is negative which causes underflow.
			// Other negative values are fine as they are not used to derive the median price.
			responseJsonString: `[
				-1247.8,
				701.99719833,
				1248,
				636.3844446100001,
				-84,
				-0.0631,
				1248,
				84500.66677816,
				1412.76385631,
				1220.11135781
				]`,
			exponentMap: testutil.ExponentSymbolMap,
			expectedUnavailableMap: map[string]error{
				testutil.ETHUSDC: errors.New("Key: 'BitfinexResponseBody.BidPrice' Error:Field validation for " +
					"'BidPrice' failed on the 'gt' tag"),
			},
		},
		"Failure - overflow due to massively negative exponent": {
			responseJsonString: bitfinexResponseString,
			exponentMap:        map[string]int32{testutil.ETHUSDC: -3000},
			expectedError:      errors.New("value overflows uint64"),
		},
		"Success - negative exponent": {
			responseJsonString: bitfinexResponseString,
			exponentMap:        testutil.ExponentSymbolMap,
			expectedPriceMap: map[string]uint64{
				"ETHUSDC": uint64(1_248_000_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			// Last price (index 6) has a very small decimal value which is ignored in the final
			// medianized price.
			responseJsonString: `[
				1247.8,
				701.99719833,
				1248,
				636.3844446100001,
				-84,
				-0.0631,
				1248.00000000000000000123,
				84500.66677816,
				1412.76385631,
				1220.11135781
				]`,
			exponentMap: testutil.ExponentSymbolMap,
			expectedPriceMap: map[string]uint64{
				testutil.ETHUSDC: uint64(1_248_000_000),
			},
		},
		"Success - positive exponent": {
			responseJsonString: bitfinexResponseString,
			exponentMap: map[string]int32{
				testutil.ETHUSDC: 2,
			},
			expectedPriceMap: map[string]uint64{
				"ETHUSDC": uint64(12),
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
				prices, unavailable, err = bitfinex.BitfinexPriceFunction(response, tc.exponentMap, medianizer)
			} else {
				prices, unavailable, err = bitfinex.BitfinexPriceFunction(response, tc.exponentMap, &lib.MedianizerImpl{})
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
