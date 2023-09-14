package bitstamp_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitstamp"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/require"
)

// Test tickers for Bitstamp.
const (
	BTCUSDC_TICKER = "BTC/USD"
	ETHUSDC_TICKER = "ETH/USD"
)

// Test exponent maps.
var (
	// Test exponent maps.
	BtcExponentMap = map[string]int32{
		BTCUSDC_TICKER: constants.BtcUsdExponent,
	}
	EthExponentMap = map[string]int32{
		ETHUSDC_TICKER: constants.EthUsdExponent,
	}
	BtcAndEthExponentMap = map[string]int32{
		BTCUSDC_TICKER: constants.BtcUsdExponent,
		ETHUSDC_TICKER: constants.EthUsdExponent,
	}
)

// Test response strings.
var (
	btcTicker = `{"timestamp": "1686600672", "open": "25940", "high": "26209", "low": "25634", "last": "25846",
		"volume": "1903.95560640", "vwap": "25868", "bid": "25841", "ask": "25842", "open_24": "26183",
		"percent_change_24": "-1.29", "pair": "BTC/USD"}`
	ethTicker = `{"timestamp": "1686600672", "open": "1753.4", "high": "1777.7", "low": "1720.1", "last": "1734.9",
		"volume": "6462.32622552", "vwap": "1738.6", "bid": "1734.3", "ask": "1734.9", "open_24": "1777.0",
		"percent_change_24": "-2.37", "pair": "ETH/USD"}`

	BtcResponseString       = fmt.Sprintf("[%s]", btcTicker)
	EthResponseString       = fmt.Sprintf("[%s]", ethTicker)
	BtcAndEthResponseString = fmt.Sprintf("[%s,%s]", ethTicker, btcTicker)
)

func TestBitstampPriceFunction_Mixed(t *testing.T) {
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
		"Unavailable - invalid response": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `[{"pair":"BTC/USD", "last":"26452", "ask":26453.23, "bid":"26449.38",}]`,
			exponentMap:        BtcExponentMap,
			expectedError:      errors.New("invalid character '}' looking for beginning of object key string"),
		},
		"Unavailable - invalid type in response: number": {
			// Invalid due to integer bidPrice when string was expected.
			responseJsonString: `[{"pair":"BTC/USD", "last":"26452", "ask":26453.23, "bid":"26449.38"}]`,
			exponentMap:        BtcExponentMap,
			expectedError: errors.New("json: cannot unmarshal number into Go struct field " +
				"BitstampTicker.ask of type string"),
		},
		"Unavailable - invalid type in response: malformed string": {
			responseJsonString: `[{"pair":"BTC/USD", "last":"26452", "ask":"not a number", "bid":"26449.38"}]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("Key: 'BitstampTicker.AskPrice' Error:Field validation for " +
					"'AskPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: `[{}]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC/USD"),
			},
		},
		"Unavailable - empty list response": {
			responseJsonString: `[]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC/USD"),
			},
		},
		"Unavailable - missing btc response": {
			responseJsonString: EthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_734_900_000),
			},
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC/USD"),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: `[{"pair":"BTC/USD", "last":"26452", "bid":"26449.38"}]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(
					"Key: 'BitstampTicker.AskPrice' Error:Field validation for 'AskPrice' failed on the 'required' tag",
				),
			},
		},
		"Failure - overflow due to massively negative exponent": {
			responseJsonString: BtcResponseString,
			exponentMap:        map[string]int32{BTCUSDC_TICKER: -3000},
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("value overflows uint64"),
			},
		},
		"Failure - medianization error": {
			responseJsonString:  BtcResponseString,
			exponentMap:         BtcExponentMap,
			medianFunctionFails: true,
			expectedPriceMap:    make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: testutil.MedianizationError,
			},
		},
		"Success - integers": {
			responseJsonString: `[{"pair":"BTC/USD", "last":"26452", "ask":"26455", "bid":"26449.38"}]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_645_200_000),
			},
		},
		"Success - negative exponent": {
			responseJsonString: BtcResponseString,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_584_200_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			responseJsonString: `[{"pair":"BTC/USD", "last":"26452.4415621293", "ask":"26455.51", "bid":"26449.38"}]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_645_244_156),
			},
		},
		"Success - positive exponent": {
			responseJsonString: BtcResponseString,
			exponentMap: map[string]int32{
				BTCUSDC_TICKER: 1,
			},
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2584),
			},
		},
		"Success - two tickers in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_584_200_000),
				ETHUSDC_TICKER: uint64(1_734_900_000),
			},
		},
		"Success - one ticker in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_734_900_000),
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
				prices, unavailable, err = bitstamp.BitstampPriceFunction(response, tc.exponentMap, testutil.MedianErr)
			} else {
				prices, unavailable, err = bitstamp.BitstampPriceFunction(response, tc.exponentMap, lib.Median[uint64])
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
