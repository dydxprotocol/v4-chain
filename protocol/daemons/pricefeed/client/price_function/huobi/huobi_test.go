package huobi_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/huobi"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/require"
)

// Test tickers for Huobi.
const (
	BTCUSDC_TICKER = "btcusdt"
	ETHUSDC_TICKER = "ethusdt"
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

func TestHuobiPriceFunction_Mixed(t *testing.T) {
	// Test response strings.
	var (
		btcTicker = pricefeed.ReadJsonTestFile(t, "btc_ticker.json")
		ethTicker = pricefeed.ReadJsonTestFile(t, "eth_ticker.json")

		ResponseStringTemplate  = `{"status":"ok","data":[%s]}`
		BtcResponseString       = fmt.Sprintf(ResponseStringTemplate, btcTicker)
		EthResponseString       = fmt.Sprintf(ResponseStringTemplate, ethTicker)
		BtcAndEthResponseString = fmt.Sprintf(`{"status":"ok","data":[%s,%s]}`, btcTicker, ethTicker)
	)

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
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ethusdt","close":1648.74,"bid":1648.95,"ask":1648.96,}`),
			exponentMap:   EthExponentMap,
			expectedError: errors.New("invalid character '}' looking for beginning of object key string"),
		},
		"Unavailable - invalid type in response: string": {
			// Invalid due to integer askPrice when string was expected.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ethusdt","close":1648.74,"bid":"1648.95","ask":1648.96}`),
			exponentMap: EthExponentMap,
			expectedError: errors.New("json: cannot unmarshal string into Go struct field " +
				"HuobiTicker.data.bid of type float64"),
		},
		"Unavailable - bid price is 0": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ethusdt","close":1648.74,"bid":0,"ask":1648.96}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New("Key: 'HuobiTicker.BidPrice' Error:Field validation for " +
					"'BidPrice' failed on the 'required' tag"),
			},
		},
		"Unavailable - ask price is negative": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ethusdt","close":1648.74,"bid":1648.96,"ask":-0.2}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New("Key: 'HuobiTicker.AskPrice' Error:Field validation for " +
					"'AskPrice' failed on the 'gt' tag"),
			},
		},
		"Unavailable - last price is negative": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ethusdt","close":-0.5,"bid":1648.96,"ask":1648.24}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New("Key: 'HuobiTicker.LastPrice' Error:Field validation for " +
					"'LastPrice' failed on the 'gt' tag"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, `{}`),
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker btcusdt"),
			},
		},
		"Unavailable - empty list response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, ``),
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker btcusdt"),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ethusdt","bid":1648.95,"ask":1648.96}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New(
					"Key: 'HuobiTicker.LastPrice' Error:Field validation for 'LastPrice' failed on the 'required' tag",
				),
			},
		},
		"Failure - response status is not ok": {
			responseJsonString: `{"status":"err","data":[]}`,
			exponentMap:        EthExponentMap,
			expectedError:      errors.New(`huobi response status is not "ok"`),
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
		"Mixed - missing btc response and has eth response": {
			responseJsonString: EthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_648_950_000),
			},
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker btcusdt"),
			},
		},
		"Success - integers": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ethusdt","close":1648,"bid":1649,"ask":1648}`),
			exponentMap: EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_648_000_000),
			},
		},
		"Success - negative exponent": {
			responseJsonString: BtcResponseString,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_507_701_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"btcusdt","close":25077.99,"bid":25077.0,"ask":25077.0175230123572}`),
			exponentMap: BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_507_701_752),
			},
		},
		"Success - positive exponent": {
			responseJsonString: BtcResponseString,
			exponentMap: map[string]int32{
				BTCUSDC_TICKER: 1,
			},
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_507),
			},
		},
		"Success - two tickers in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_507_701_000),
				ETHUSDC_TICKER: uint64(1_648_950_000),
			},
		},
		"Success - one ticker in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_648_950_000),
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
				prices, unavailable, err = huobi.HuobiPriceFunction(response, tc.exponentMap, testutil.MedianErr)
			} else {
				prices, unavailable, err = huobi.HuobiPriceFunction(response, tc.exponentMap, lib.Median[uint64])
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
