package binance

import (
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func TestUnmarshalBinanceResponse_Mixed(t *testing.T) {
	tests := map[string]struct {
		// parameters
		responseJsonString string
		// expectations
		expectedResponseBody BinanceResponseBody
		expectedError        error
	}{
		"Invalid response - invalid JSON": {
			// Invalid due to trailing comma in JSON.
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "1368.0800", "lastPrice": "1368.2100",}`,
			expectedError:      errors.New("invalid character '}' looking for beginning of object key string"),
		},
		"Invalid response - expected string field type": {
			// Invalid due to integer bidPrice when string was expected.
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": 1368.0800, "lastPrice": "1368.2100"}`,
			expectedError: errors.New("json: cannot unmarshal number into Go struct field " +
				"BinanceResponseBody.bidPrice of type string"),
		},
		"Invalid response - expected numeric string field value": {
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "Not a Number", "lastPrice": "1368.2100"}`,
			expectedError: errors.New("Key: 'BinanceResponseBody.BidPrice' Error:Field validation for " +
				"'BidPrice' failed on the 'numeric' tag"),
		},
		"Invalid response - empty": {
			responseJsonString: `{}`,
			expectedError: errors.New(
				"Key: 'BinanceResponseBody.AskPrice' Error:Field validation for 'AskPrice' failed on the 'required' tag\n" +
					"Key: 'BinanceResponseBody.BidPrice' Error:Field validation for 'BidPrice' failed on the 'required' tag\n" +
					"Key: 'BinanceResponseBody.LastPrice' Error:Field validation for 'LastPrice' failed on the 'required' tag",
			),
		},
		"Invalid response - missing field": {
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "1368.0800"}`,
			expectedError: errors.New("Key: 'BinanceResponseBody.LastPrice' Error:Field validation for " +
				"'LastPrice' failed on the 'required' tag"),
		},
		"Success": {
			responseJsonString: `{"askPrice": "1368.5100", "bidPrice": "1368.0800", "lastPrice": "1368.2100"}`,
			expectedResponseBody: BinanceResponseBody{
				AskPrice:  "1368.5100",
				BidPrice:  "1368.0800",
				LastPrice: "1368.2100",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody, err := unmarshalBinanceResponse(io.NopCloser(strings.NewReader(tc.responseJsonString)))
			if tc.expectedError == nil {
				require.Nil(t, err)
				require.Equal(t, tc.expectedResponseBody, *responseBody)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}
}
