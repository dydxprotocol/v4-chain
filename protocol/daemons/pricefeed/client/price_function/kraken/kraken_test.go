package kraken_test

import (
	"errors"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/kraken"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestKrakenPriceFunction_Mixed(t *testing.T) {
	krakenValidResponseString := pricefeed.ReadJsonTestFile(t, "kraken_2_ticker_response.json")
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
			expectedError: errors.New(
				"kraken API response JSON parse error (invalid character ',' looking for beginning of object " +
					"key string)",
			),
		},
		"Failure - invalid response, float instead of string data type, missing": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{"a":[2105.8]}}}`,
			exponentMap:        ValidSymbolMap,
			expectedError: errors.New(
				"kraken API response JSON parse error (json: cannot unmarshal number into Go struct field " +
					"KrakenTickerResult.result.a of type string)",
			),
		},
		"Unavailable - overflow due to negative exponent": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: krakenValidResponseString,
			exponentMap:        map[string]int32{ETHUSDC_SYMBOL: -3000},
			expectedPriceMap:   map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_SYMBOL: errors.New("value overflows uint64"),
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
			expectedError:      errors.New("kraken API call error: EQuery:Unknown asset pair"),
		},
		"Success - one market response": {
			responseJsonString: krakenValidResponseString,
			exponentMap:        EthSymbolMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_SYMBOL: uint64(1_888_000_000),
			},
		},
		"Success - two market response": {
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
