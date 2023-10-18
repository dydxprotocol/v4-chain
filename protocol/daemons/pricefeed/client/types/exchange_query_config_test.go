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

func TestValidateDelta(t *testing.T) {
	tests := map[string]struct {
		exchangeQueryConfig *types.ExchangeQueryConfig
		expectedError       error
	}{
		"valid": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 0, // valid for a delta
				TimeoutMs:  0, // valid for a delta
				MaxQueries: 0, // valid for a delta
			},
		},
		"invalid - invalid exchange id": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "abc", // invalid
			},
			expectedError: fmt.Errorf("invalid exchange id abc"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.exchangeQueryConfig.ValidateDelta(validExchanges)
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
			}
		})
	}
}

func TestApplyDeltaAndValidate(t *testing.T) {
	tests := map[string]struct {
		exchangeQueryConfig *types.ExchangeQueryConfig
		delta               *types.ExchangeQueryConfig
		expected            *types.ExchangeQueryConfig
		expectedError       error
	}{
		"valid, applies all fields": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 1,
				TimeoutMs:  1,
				MaxQueries: 1,
			},
			delta: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				Disabled:   true,
				IntervalMs: 2,
				TimeoutMs:  2,
				MaxQueries: 2,
			},
			expected: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				Disabled:   true,
				IntervalMs: 2,
				TimeoutMs:  2,
				MaxQueries: 2,
			},
		},
		"invalid - mismatched exchange id": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 1,
				TimeoutMs:  1,
				MaxQueries: 1,
			},
			delta: &types.ExchangeQueryConfig{
				ExchangeId: "CoinbasePro", // invalid - does not match above
			},
			expectedError: fmt.Errorf("exchange id mismatch: CoinbasePro, Binance"),
		},
		"valid, enables disabled exchange": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				Disabled:   true,
				IntervalMs: 1,
				TimeoutMs:  1,
				MaxQueries: 1,
			},
			delta: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				Disabled:   false, // even though this is zero, expect it to always be applied.
			},
			expected: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 1,
				TimeoutMs:  1,
				MaxQueries: 1,
			},
		},
		"invalid - invalid exchange config": {
			exchangeQueryConfig: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
				IntervalMs: 0, // invalid
				TimeoutMs:  1,
				MaxQueries: 1,
			},
			delta: &types.ExchangeQueryConfig{
				ExchangeId: "Binance",
			},
			expectedError: fmt.Errorf("Error:Field validation for 'IntervalMs' failed on the 'gt' tag"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			updatedConfig, err := tc.exchangeQueryConfig.ApplyDeltaAndValidate(tc.delta, validExchanges)
			if tc.expectedError == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expected, updatedConfig)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
				require.Zero(t, updatedConfig)
			}
		})
	}
}

func TestExchangeQueryConfig_Copy(t *testing.T) {
	// Make a struct with all non-zero values to validate that all values are propagated to the copy.
	ecq := &types.ExchangeQueryConfig{
		ExchangeId: "Binance",
		Disabled:   true,
		IntervalMs: 1,
		TimeoutMs:  2,
		MaxQueries: 3,
	}
	require.Equal(t, ecq, ecq.Copy())
}
