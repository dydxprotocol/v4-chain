package kraken

import (
	"errors"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func TestUnmarshalKrakenResponse_Mixed(t *testing.T) {
	krakenValidResponseString := pricefeed.ReadJsonTestFile(t, "kraken_2_ticker_response.json")
	tests := map[string]struct {
		responseJsonString   string
		expectedResponseBody KrakenResponseBody
		expectedError        error
	}{
		"invalid response - float instead of string": {
			responseJsonString: `{"result":{"XETHZUSD":{"a":[2105.8]}}}`,
			expectedError: errors.New(
				"kraken API response JSON parse error (json: cannot unmarshal number into Go struct field " +
					"KrakenTickerResult.result.a of type string)"),
		},
		"invalid response - non-numeric string": {
			responseJsonString: `{"result":{"XETHZUSD":{"a":["cat","2","3"],"b":["1","2","3"],"c":["1","2"]}}}`,
			expectedError: errors.New(
				"kraken API response validation error (Key: " +
					"'KrakenResponseBody.Tickers[XETHZUSD].AskPriceStats[0]' Error:Field validation for " +
					"'AskPriceStats[0]' failed on the 'positive-float-string' tag)",
			),
		},
		"invalid response - response field too short": {
			responseJsonString: `{"result":{"XETHZUSD":{"a":["2","3"],"b":["1","2","3"],"c":["1","2"]}}}`,
			expectedError: errors.New(
				"kraken API response validation error (Key: 'KrakenResponseBody.Tickers[XETHZUSD].AskPriceStats' " +
					"Error:Field validation for 'AskPriceStats' failed on the 'len' tag)",
			),
		},
		"invalid response - response field too long": {
			responseJsonString: `{"result":{"XETHZUSD":{"a":["1","2","3","4"],"b":["1","2","3"],"c":["1","2"]}}}`,
			expectedError: errors.New(
				"kraken API response validation error (Key: 'KrakenResponseBody.Tickers[XETHZUSD].AskPriceStats' " +
					"Error:Field validation for 'AskPriceStats' failed on the 'len' tag)",
			),
		},
		"invalid response - negative price": {
			responseJsonString: `{"result":{"XETHZUSD":{"a":["-1234.56","2","3"],"b":["1","2","3"],"c":["1","2"]}}}`,
			expectedError: errors.New(
				"kraken API response validation error (Key: " +
					"'KrakenResponseBody.Tickers[XETHZUSD].AskPriceStats[0]' Error:Field validation for " +
					"'AskPriceStats[0]' failed on the 'positive-float-string' tag)",
			),
		},
		"invalid response - missing expected ticker content": {
			responseJsonString: `{"result":{"XETHZUSD":{}}}`,
			expectedError: errors.New(
				"kraken API response validation error (Key: 'KrakenResponseBody.Tickers[XETHZUSD].AskPriceStats' " +
					"Error:Field validation for 'AskPriceStats' failed on the 'len' tag\nKey: " +
					"'KrakenResponseBody.Tickers[XETHZUSD].BidPriceStats' Error:Field validation for 'BidPriceStats' " +
					"failed on the 'len' tag\nKey: 'KrakenResponseBody.Tickers[XETHZUSD].ClosePriceStats' " +
					"Error:Field validation for 'ClosePriceStats' failed on the 'len' tag)",
			),
		},
		"invalid response - empty JSON object": {
			responseJsonString: `{}`,
			expectedError: errors.New(
				"kraken API response validation error (Key: 'KrakenResponseBody.Tickers' Error:Field validation " +
					"for 'Tickers' failed on the 'required_without' tag)",
			),
		},
		"invalid response - empty string": {
			responseJsonString: "",
			expectedError:      errors.New("kraken API response JSON parse error (EOF)"),
		},
		"invalid response - both tickers and errors defined": {
			responseJsonString: `{"error":["abcdefg"],"result":{"XETHZUSD":{"a":["1","2","3"],"b":["1","2","3"],` +
				`"c":["1","2","3"]}}}`,
			expectedError: errors.New(
				"kraken API response validation error (Key: 'KrakenResponseBody.Tickers' Error:Field validation " +
					"for 'Tickers' failed on the 'excluded_with' tag)",
			),
		},
		"success": {
			responseJsonString: krakenValidResponseString,
			expectedResponseBody: KrakenResponseBody{
				Tickers: map[string]KrakenTickerResult{
					"XETHZUSD": {
						AskPriceStats:   []string{"1888.00000", "59", "59.000"},
						BidPriceStats:   []string{"1887.99000", "62", "62.000"},
						ClosePriceStats: []string{"1888.00000", "0.65587578"},
					},
					"XXBTZUSD": {
						AskPriceStats:   []string{"29207.50000", "1", "1.000"},
						BidPriceStats:   []string{"29204.50000", "2", "2.000"},
						ClosePriceStats: []string{"29207.50000", "0.01327170"},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody, err := unmarshalKrakenResponse(io.NopCloser(strings.NewReader(tc.responseJsonString)))
			if tc.expectedError == nil {
				require.Nil(t, err)
				require.Equal(t, tc.expectedResponseBody, *responseBody)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}
}
