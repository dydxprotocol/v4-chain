package bitfinex_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

// Test tickers for Bitfinex.
const (
	BTCUSDC_TICKER = "tBTCUSD"
	ETHUSDC_TICKER = "tETHUSD"
)

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

	// Test response strings.
	btcTicker = `["tBTCUSD",30169,24.43585233,30170,19.64058018,1812,0.06388605,30175,3180.89036382,30801,28271]`
	ethTicker = `["tETHUSD",1898.4,567.73204365,1898.5,479.12571969,104.6,0.05831196,1898.4,12095.61872826,1903.4,1787.7]`

	ResponseStringTemplate  = `[%s]`
	BtcResponseString       = fmt.Sprintf(ResponseStringTemplate, btcTicker)
	EthResponseString       = fmt.Sprintf(ResponseStringTemplate, ethTicker)
	BtcAndEthResponseString = fmt.Sprintf(`[%s,%s]`, btcTicker, ethTicker)
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
		"Unavailable - invalid response": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169,24.43585233,30170,19.64058018,1812,0.06388605,30175,3180.89036382,30801,28271,]`),
			exponentMap:   BtcExponentMap,
			expectedError: errors.New("invalid character ']' looking for beginning of value"),
		},
		"Unavailable - invalid pair type in response: number": {
			// Invalid due to number pair when string was expected.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`[12,30169,24.43585233,30170,19.64058018,1812,0.06388605,30175,3180.89036382,30801,28271]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker tBTCUSD"),
			},
		},
		"Unavailable - invalid bid price type in response: string": {
			// Invalid due to string bid price when float was expected.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD","30169",24.43585233,30170,19.64058018,1812,0.06388605,30175,3180.89036382,30801,28271]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("invalid bid price in response - not a float64"),
			},
		},
		"Unavailable - invalid ask price type in response: string": {
			// Invalid due to string ask price when float was expected.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169.2,24.43585233,"30170",19.64058018,1812,0.06388605,30175,3180.89036382,30801,28271]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("invalid ask price in response - not a float64"),
			},
		},
		"Unavailable - invalid last price type in response: string": {
			// Invalid due to string last price when float was expected.
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169.2,24.43585233,30170,19.64058018,1812,0.06388605,"30175",3180.89036382,30801,28271]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("invalid last price in response - not a float64"),
			},
		},
		"Unavailable - bid price is 0": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",0,24.43585233,30170,19.64058018,1812,0.06388605,30175,3180.89036382,30801,28271]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("Key: 'BitfinexTicker.BidPrice' Error:Field validation for " +
					"'BidPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - ask price is negative": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169,24.43585233,-1,19.64058018,1812,0.06388605,30175,3180.89036382,30801,28271]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("Key: 'BitfinexTicker.AskPrice' Error:Field validation for " +
					"'AskPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - last price is negative": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169,24.43585233,30170,19.64058018,1812,0.06388605,-2,3180.89036382,30801,28271]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("Key: 'BitfinexTicker.LastPrice' Error:Field validation for " +
					"'LastPrice' failed on the 'positive-float-string' tag"),
			},
		},
		"Unavailable - empty response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, `[]`),
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker tBTCUSD"),
			},
		},
		"Unavailable - empty list response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate, ``),
			exponentMap:        BtcExponentMap,
			expectedPriceMap:   make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker tBTCUSD"),
			},
		},
		"Unavailable - incomplete response": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169,24.43585233,30170,19.64058018,1812,0.06388605,30175,3180.89036382,30801]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker tBTCUSD"),
			},
		},
		"Unavailable - non-requested ticker doesn't get marked as unavailable": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBBBUSD",30169,24.43585233,"30170",19.64058018,1812,0.06388605,30175,3180.89036382,30801]`),
			exponentMap:      BtcExponentMap,
			expectedPriceMap: make(map[string]uint64),
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker tBTCUSD"),
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
				ETHUSDC_TICKER: uint64(1_898_400_000),
			},
			expectedUnavailableMap: map[string]error{
				BTCUSDC_TICKER: errors.New("no listing found for ticker tBTCUSD"),
			},
		},
		"Success - integers": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169,24.43585233,30170,19.64058018,1812,0.06388605,30175.43,3180.89036382,30801,28271]`),
			exponentMap: BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(3_017_000_000),
			},
		},
		"Success - negative exponent": {
			responseJsonString: BtcResponseString,
			exponentMap:        BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(3_017_000_000),
			},
		},
		"Success - decimals beyond supported precision ignored": {
			responseJsonString: fmt.Sprintf(ResponseStringTemplate,
				`["tBTCUSD",30169.23994231423,24.43,30169,19.64058018,1812,0.06388605,30175,3180.890,30801,28271]`),
			exponentMap: BtcExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(3_016_923_994),
			},
		},
		"Success - positive exponent": {
			responseJsonString: BtcResponseString,
			exponentMap: map[string]int32{
				BTCUSDC_TICKER: 1,
			},
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(3_017),
			},
		},
		"Success - two tickers in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        BtcAndEthExponentMap,
			expectedPriceMap: map[string]uint64{
				BTCUSDC_TICKER: uint64(3_017_000_000),
				ETHUSDC_TICKER: uint64(1_898_400_000),
			},
		},
		"Success - one ticker in request, two tickers in response": {
			responseJsonString: BtcAndEthResponseString,
			exponentMap:        EthExponentMap,
			expectedPriceMap: map[string]uint64{
				ETHUSDC_TICKER: uint64(1_898_400_000),
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
				prices, unavailable, err = bitfinex.BitfinexPriceFunction(response, tc.exponentMap, testutil.MedianErr)
			} else {
				prices, unavailable, err = bitfinex.BitfinexPriceFunction(response, tc.exponentMap, lib.Median[uint64])
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
