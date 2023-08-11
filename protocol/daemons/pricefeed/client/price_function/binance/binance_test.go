package binance_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4/testutil/constants"
	"testing"

	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/binance"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test tickers for Binance.
const (
	BTCUSDC_TICKER = `"BTCUSDT"`
	ETHUSDC_TICKER = `"ETHUSDT"`
)

// Test exponent maps.
var (
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

func TestBinancePriceFunction_Mixed(t *testing.T) {
	// Test response strings.
	var (
		btcTicker = pricefeed.ReadJsonTestFile(t, "btc_ticker_binance.json")
		ethTicker = pricefeed.ReadJsonTestFile(t, "eth_ticker_binance.json")

		ResponseStringTemplate  = `[%s]`
		BtcResponseString       = fmt.Sprintf(ResponseStringTemplate, btcTicker)
		EthResponseString       = fmt.Sprintf(ResponseStringTemplate, ethTicker)
		BtcAndEthResponseString = fmt.Sprintf(`[%s,%s]`, btcTicker, ethTicker)
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
				`{"symbol":"ETHUSDT","lastPrice":"1780.29000000","bidPrice":"1780.24000000","askPrice":"1780.25000000",}`),
			exponentMap:   EthExponentMap,
			expectedError: errors.New("invalid character '}' looking for beginning of object key string"),
		},
		"Unavailable - invalid type in response: number": {
			// Invalid due to number askPrice when string was expected.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ETHUSDT","lastPrice":"1780.29000000","bidPrice":"1780.24000000","askPrice":1780.25000000}`),
			exponentMap: EthExponentMap,
			expectedError: errors.New("json: cannot unmarshal number into Go struct field " +
				"BinanceTicker.askPrice of type string"),
		},
		"Unavailable - bid price is 0": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ETHUSDT","lastPrice":"1780.29000000","bidPrice":"0","askPrice":"1780.25000000"}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New("Key: 'BinanceTicker.BidPrice' Error:Field validation for " +
					"'BidPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - ask price is negative": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ETHUSDT","lastPrice":"1780.29000000","bidPrice":"1780.24000000","askPrice":"-1780.25000000"}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New("Key: 'BinanceTicker.AskPrice' Error:Field validation for " +
					"'AskPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - last price is negative": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ETHUSDT","lastPrice":"-1780.29000000","bidPrice":"1780.24000000","askPrice":"1780.25000000"}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New("Key: 'BinanceTicker.LastPrice' Error:Field validation for " +
					"'LastPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, `{}`),
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(`no listing found for ticker "BTCUSDT"`),
			},
		},
		"Unavailable - empty list response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, ``),
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(`no listing found for ticker "BTCUSDT"`),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ETHUSDT","lastPrice":"1780.29000000","askPrice":"1780.25000000"}`),
			exponentMap:      EthExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				ETHUSDC_TICKER: errors.New(
					"Key: 'BinanceTicker.BidPrice' Error:Field validation for 'BidPrice' failed on the 'required' tag",
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
		"Mixed - missing btc response and has eth response": {
			responseJsonString: EthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_780_250_000),
			},
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New(`no listing found for ticker "BTCUSDT"`),
			},
		},
		"Success - integers": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ETHUSDT","lastPrice":"1780","bidPrice":"1780","askPrice":"1780.25000000"}`),
			exponentMap: EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_780_000_000),
			},
		},
		"Success - negative exponent": {
			responseJsonString: BtcResponseString,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_794_470_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`{"symbol":"ETHUSDT","lastPrice":"1780.2900","bidPrice":"1780.240","askPrice":"1780.25123752942"}`),
			exponentMap: EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_780_251_237),
			},
		},
		"Success - positive exponent": {
			responseJsonString: BtcResponseString,
			exponentMap: map[string]int32{
				BTCUSDC_TICKER: 1,
			},
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_794),
			},
		},
		"Success - two tickers in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(2_794_470_000),
				ETHUSDC_TICKER: uint64(1_780_250_000),
			},
		},
		"Success - one ticker in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_780_250_000),
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
