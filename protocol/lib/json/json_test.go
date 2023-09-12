package json_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib/json"
)

func TestIsValidJSON(t *testing.T) {
	type testCase struct {
		name    string
		input   string
		wantErr bool
	}

	testCases := []testCase{
		{
			name:    "Basic JSON string",
			input:   `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "Invalid string",
			input:   `{"key": "value`,
			wantErr: true,
		},
		{
			name:    "Nested well formed JSON object",
			input:   `{"key": {"nestedKey": "nestedValue"}}`,
			wantErr: false,
		},
		{
			name:    "Typo in nested JSON object",
			input:   `{"key": {"nestedKey": "nestedValue}`,
			wantErr: true,
		},
		{
			name: "Nested well formed JSON object with array",
			input: `{
				"exchanges": [
				  {
					"exchangeName": "Binance",
					"ticker": "LINKUSDT",
					"adjustByMarket": "USDT-USD"
				  },
				  {
					"exchangeName": "Bitstamp",
					"ticker": "LINK/USD"
				  },
				  {
					"exchangeName": "Bybit",
					"ticker": "LINKUSDT",
					"adjustByMarket": "USDT-USD"
				  }
				]
			}`,
			wantErr: false,
		},
		{
			name: "Typo in JSON object with array",
			input: `{
				"exchanges": [
				  {
					"exchangeName": "Binance",
					"ticker": "LINKUSDT",
					"adjustByMarket": "USDT-USD"
				  },
				  {
					"exchangeName": "Bitstamp",
					"ticker": "LINK/USD"
				  },
				  {
					"exchangeName": "Bybit",
					"ticker": "LINKUSDT",
					"adjustByMarket": "USDT-USD"
				  },
				]
			}`,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := json.IsValidJSON(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("IsValidJSON() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
