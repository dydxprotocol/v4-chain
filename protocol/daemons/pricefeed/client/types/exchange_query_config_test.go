package types_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	validExchanges = map[string]struct{}{
		"Binance": {},
	}
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		exchangeQueryConfig *types.ExchangeQueryConfig
		expectedError       error
	}{
		"valid": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 1,
				TimeoutMs:  1,
				MaxQueries: 1,
			},
		},
		"invalid - invalid exchange id": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "abc", // invalid
			},
			expectedError: fmt.Errorf("invalid exchange id abc"),
		},
		"invalid: interval ms 0": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 0, // invalid
			},
			expectedError: fmt.Errorf("Error:Field validation for 'IntervalMs' failed on the 'gt' tag"),
		},
		"invalid: timeout ms 0": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 1,
				TimeoutMs:  0, // invalid
			},
			expectedError: fmt.Errorf("Error:Field validation for 'TimeoutMs' failed on the 'gt' tag"),
		},
		"invalid: max queries 0": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 1,
				TimeoutMs:  1,
				MaxQueries: 0, // invalid
			},
			expectedError: fmt.Errorf("Error:Field validation for 'MaxQueries' failed on the 'gt' tag"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.exchangeQueryConfig.Validate(validExchanges)
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
			}
		})
	}
}
