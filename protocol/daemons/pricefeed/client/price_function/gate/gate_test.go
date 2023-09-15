package gate_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/gate"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/require"
)

// Test tickers for Gate.
const (
	BTCUSDC_TICKER = "BTC_USD"
	ETHUSDC_TICKER = "ETH_USD"
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
	btcTicker = `{"currency_pair":"BTC_USD","last":"26452.44","lowest_ask":"26455.51","highest_bid":"26449.38",
		"change_percentage":"0.25","base_volume":"162.03142","quote_volume":"4294834.9535027","high_24h":"26825.72",
		"low_24h":"26308.05"}`
	ethTicker = `{"currency_pair":"ETH_USD","last":"1757.36","lowest_ask":"1757.96","highest_bid":"1757.71",
		"change_percentage":"-3.94","change_utc0":"0","change_utc8":"0.88","base_volume":"1249.407483",
		"quote_volume":"2197610.5129584","high_24h":"1829.58","low_24h":"1722.33"}`

	BtcResponseString       = fmt.Sprintf("[%s]", btcTicker)
	EthResponseString       = fmt.Sprintf("[%s]", ethTicker)
	BtcAndEthResponseString = fmt.Sprintf("[%s,%s]", ethTicker, btcTicker)
)

func TestGatePriceFunction_Mixed(t *testing.T) {
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
			responseJsonString: `[{"currency_pair":"BTC_USD","last":"26452",
				"lowest_ask":26453.23,"highest_bid":"26449.38",}]`,
			exponentMap:   BtcExponentMap,
			expectedError: errors.New("invalid character '}' looking for beginning of object key string"),
		},
		"Unavailable - invalid type in response: number": {
			// Invalid due to integer bidPrice when string was expected.
			responseJsonString: `[{"currency_pair":"BTC_USD","last":"26452",
				"lowest_ask":26453.23,"highest_bid":"26449.38"}]`,
			exponentMap: BtcExponentMap,
			expectedError: errors.New("json: cannot unmarshal number into Go struct field " +
				"GateTicker.lowest_ask of type string"),
		},
		"Unavailable - invalid type in response: malformed string": {
			responseJsonString: `[{"currency_pair":"BTC_USD","last":"26452",
				"lowest_ask":"not a number","highest_bid":"26449.38"}]`,
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("Key: 'GateTicker.AskPrice' Error:Field validation for " +
					"'AskPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: `[{}]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC_USD"),
			},
		},
		"Unavailable - empty list response": {
			responseJsonString: `[]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC_USD"),
			},
		},
		"Unavailable - missing btc response": {
			responseJsonString: EthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_757_710_000),
			},
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC_USD"),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: `[{"currency_pair":"BTC_USD","last":"26452","highest_bid":"26449.38"}]`,
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(
					"Key: 'GateTicker.AskPrice' Error:Field validation for 'AskPrice' failed on the 'required' tag",
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
			responseJsonString: `[{"currency_pair":"BTC_USD","last":"26452",
				"lowest_ask":"26455","highest_bid":"26449.38"}]`,
			exponentMap: BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_645_200_000),
			},
		},
		"Success - negative exponent": {
			responseJsonString: BtcResponseString,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_645_244_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			responseJsonString: `[{"currency_pair":"BTC_USD","last":"26452.4415621293",
				"lowest_ask":"26455.51","highest_bid":"26449.38"}]`,
			exponentMap: BtcExponentMap,
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
				BTCUSDC_TICKER: uint64(2645),
			},
		},
		"Success - two tickers in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_645_244_000),
				ETHUSDC_TICKER: uint64(1_757_710_000),
			},
		},
		"Success - one ticker in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_757_710_000),
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
				prices, unavailable, err = gate.GatePriceFunction(response, tc.exponentMap, testutil.MedianErr)
			} else {
				prices, unavailable, err = gate.GatePriceFunction(response, tc.exponentMap, lib.Median[uint64])
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
