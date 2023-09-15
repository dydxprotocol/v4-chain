package kraken_test

import (
	"errors"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/kraken"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/require"
)

const (
	ETHUSDC_TICKER = "XETHZUSD"
	BTCUSDC_TICKER = "XXBTZUSD"
)

var (
	EthExponentMap = map[string]int32{
		ETHUSDC_TICKER: constants.EthUsdExponent,
	}
	BtcAndEthExponentMap = map[string]int32{
		BTCUSDC_TICKER: constants.BtcUsdExponent,
		ETHUSDC_TICKER: constants.EthUsdExponent,
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
			exponentMap:        EthExponentMap,
			expectedError: errors.New(
				"invalid character ',' looking for beginning of object key string",
			),
		},
		"Failure - invalid response, float instead of string data type, missing": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"result":{"XETHZUSD":{"a":[2105.8]}}}`,
			exponentMap:        BtcAndEthExponentMap,
			expectedError: errors.New(
				"json: cannot unmarshal number into Go struct field KrakenTickerResult.result.a of type string",
			),
		},
		"Unavailable - overflow due to negative exponent": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: krakenValidResponseString,
			exponentMap:        map[string]int32{ETHUSDC_TICKER: -3000},
			expectedPriceMap:   map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New("value overflows uint64"),
			},
		},
		"Unavailable - fails on medianization error": {
			responseJsonString:  krakenValidResponseString,
			exponentMap:         BtcAndEthExponentMap,
			medianFunctionFails: true,
			expectedPriceMap:    map[string]uint64{},
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: testutil.MedianizationError,
				BTCUSDC_TICKER: testutil.MedianizationError,
			},
		},
		"Failure - Kraken API Error response": {
			responseJsonString: `{"error":["EQuery:Unknown asset pair"]}`,
			exponentMap:        EthExponentMap,
			expectedError:      errors.New("kraken API call error: EQuery:Unknown asset pair"),
		},
		"Success - one market response": {
			responseJsonString: krakenValidResponseString,
			exponentMap:        EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_888_000_000),
			},
		},
		"Success - two market response": {
			responseJsonString: krakenValidResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_888_000_000),
				BTCUSDC_TICKER: uint64(2_920_750_000),
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
				prices, unavailable, err = kraken.KrakenPriceFunction(response, tc.exponentMap, testutil.MedianErr)
			} else {
				prices, unavailable, err = kraken.KrakenPriceFunction(response, tc.exponentMap, lib.Median[uint64])
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
