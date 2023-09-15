package crypto_com_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/crypto_com"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/require"
)

// Test tickers for CryptoCom.
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
	btcTicker = `{"i":"BTC_USD","h":"26437.90","l":"25722.02","a":"25988.88","v":"617.7597",
		"vv":"16062822.01","c":"-0.0011","b":"25989.77","k":"25990.15","t":1686716261608}`
	ethTicker = `{"i":"ETH_USD","h":"1767.34","l":"1724.35","a":"1745.33","v":"4854.0452",
		"vv":"8458590.29","c":"-0.0005","b":"1745.21","k":"1745.33","t":1686716257576}`
	ResponseStringTemplate  = `{"code":0,"result":{"data":[%s]}}`
	BtcResponseString       = fmt.Sprintf(ResponseStringTemplate, btcTicker)
	EthResponseString       = fmt.Sprintf(ResponseStringTemplate, ethTicker)
	BtcAndEthResponseString = fmt.Sprintf(`{"code":0,"result":{"data":[%s,%s]}}`, btcTicker, ethTicker)
)

func TestCryptoComPriceFunction_Mixed(t *testing.T) {
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
				`{"i":"BTC_USD","a":"25988.88","b":"25989.77","k":"25990.15",}`),
			exponentMap:   BtcExponentMap,
			expectedError: errors.New("invalid character '}' looking for beginning of object key string"),
		},
		"Unavailable - invalid type in response: number": {
			// Invalid due to integer askPrice when string was expected.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"i":"BTC_USD","a":"25988.88","b":"25989.77","k":25990.15}`),
			exponentMap: BtcExponentMap,
			expectedError: errors.New("json: cannot unmarshal number into Go struct field " +
				"CryptoComTicker.result.data.k of type string"),
		},
		"Unavailable - invalid type in response: malformed string": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"i":"BTC_USD","a":"25988.88","b":"malformed number","k":"25990.15"},
				{"i":"ETH_USD","a":"1745.33","b":"1745.21","k":"1745.33"}`),
			exponentMap: BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_745_330_000), // ETH_USD should still be available despite BTC_USD error.
			},
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("Key: 'CryptoComTicker.BidPrice' Error:Field validation for 'BidPrice' " +
					"failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, `{}`),
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC_USD"),
			},
		},
		"Unavailable - empty list response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, ``),
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
				ETHUSDC_TICKER: uint64(1_745_330_000),
			},
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker BTC_USD"),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"i":"BTC_USD","a":"25988.88","k":"25990.15"}`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(
					"Key: 'CryptoComTicker.BidPrice' Error:Field validation for 'BidPrice' failed on the 'required' tag",
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
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"i":"BTC_USD","a":"25988","b":"25989","k":"25990.15"}`),
			exponentMap: BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_598_900_000),
			},
		},
		"Success - negative exponent": {
			responseJsonString: BtcResponseString,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_598_977_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"i":"BTC_USD","a":"25988.88","b":"25989.775294213","k":"25990.15"}`),
			exponentMap: BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_598_977_529),
			},
		},
		"Success - positive exponent": {
			responseJsonString: BtcResponseString,
			exponentMap: map[string]int32{
				BTCUSDC_TICKER: 1,
			},
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_598),
			},
		},
		"Success - null bid price": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"i":"BTC_USD","a":"25990.88","b":null,"k":"25988.15"}`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(
					"Key: 'CryptoComTicker.BidPrice' Error:Field validation for 'BidPrice' failed on the 'required' tag",
				),
			},
		},
		"Success - null ask prices": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"i":"BTC_USD","a":"25988.88","b":"25984.23","k":null}`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(
					"Key: 'CryptoComTicker.AskPrice' Error:Field validation for 'AskPrice' failed on the 'required' tag",
				),
			},
		},
		"Success - two tickers in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_598_977_000),
				ETHUSDC_TICKER: uint64(1_745_330_000),
			},
		},
		"Success - one ticker in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_745_330_000),
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
				prices, unavailable, err = crypto_com.CryptoComPriceFunction(response, tc.exponentMap, testutil.MedianErr)
			} else {
				prices, unavailable, err = crypto_com.CryptoComPriceFunction(response, tc.exponentMap, lib.Median[uint64])
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
