package bitfinex

import (
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func TestUnmarshalBitfinexResponse_Mixed(t *testing.T) {
	tests := map[string]struct {
		responseJsonString   string
		expectedResponseBody BitfinexResponseBody
		expectedError        error
	}{
		"Failure - Response too short": {
			responseJsonString: `[1234.4]`,
			expectedError:      errors.New("Invalid response body length for Bitfinex with length of: 1, expected length 10"),
		},
		"Failure - Response not a slice": {
			responseJsonString: `1234.5`,
			expectedError:      errors.New("json: cannot unmarshal number into Go value of type bitfinex.bitfinexRawResponse"),
		},
		"Failure - Response contains incorrect types": {
			responseJsonString: `["abcdef"]`,
			expectedError:      errors.New("json: cannot unmarshal string into Go value of type float64"),
		},
		"Failure - Bid < 0": {
			responseJsonString: `[-1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000]`,
			expectedError: errors.New("Key: 'BitfinexResponseBody.BidPrice' Error:Field validation for " +
				"'BidPrice' failed on the 'gt' tag"),
		},
		"Failure - Bid == 0": {
			responseJsonString: `[0, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000]`,
			expectedError: errors.New("Key: 'BitfinexResponseBody.BidPrice' Error:Field validation for " +
				"'BidPrice' failed on the 'gt' tag"),
		},
		"Failure - Ask < 0": {
			responseJsonString: `[1.000, 1.000, -1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000]`,
			expectedError: errors.New("Key: 'BitfinexResponseBody.AskPrice' Error:Field validation for 'AskPrice' " +
				"failed on the 'gt' tag"),
		},
		"Failure - Ask == 0": {
			responseJsonString: `[1.000, 1.000, 0, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 1.000]`,
			expectedError: errors.New("Key: 'BitfinexResponseBody.AskPrice' Error:Field validation for " +
				"'AskPrice' failed on the 'gt' tag"),
		},
		"Failure - Close < 0": {
			responseJsonString: `[1.000, 1.000, 1.000, 1.000, 1.000, 1.000, -1.000, 1.000, 1.000, 1.000]`,
			expectedError: errors.New("Key: 'BitfinexResponseBody.LastPrice' Error:Field validation for " +
				"'LastPrice' failed on the 'gt' tag"),
		},
		"Failure - Close == 0": {
			responseJsonString: `[1.000, 1.000, 1.000, 1.000, 1.000, 1.000, 0.000, 1.000, 1.000, 1.000]`,
			expectedError: errors.New("Key: 'BitfinexResponseBody.LastPrice' Error:Field validation for " +
				"'LastPrice' failed on the 'gt' tag"),
		},
		"Success": {
			responseJsonString: `[1.234, 1.000, 5.678, 1.000, 1.000, 1.000, 3.904, 1.000, 1.000, 1.000]`,
			expectedResponseBody: BitfinexResponseBody{
				BidPrice:  1.234,
				AskPrice:  5.678,
				LastPrice: 3.904,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody, err := unmarshalBitfinexResponse(io.NopCloser(strings.NewReader(tc.responseJsonString)))
			if tc.expectedError == nil {
				require.Nil(t, err)
				require.Equal(t, tc.expectedResponseBody, *responseBody)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}
}
