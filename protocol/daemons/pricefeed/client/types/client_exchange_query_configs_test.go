package types_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientExchangeQueryConfigs_Validate(t *testing.T) {
	tests := map[string]struct {
		configs     types.ClientExchangeQueryConfigs
		expectedErr error
	}{
		"valid: empty configs": {
			configs:     types.ClientExchangeQueryConfigs{},
			expectedErr: nil,
		},
		"valid: populated_configs": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
					{
						ExchangeId: "CoinbasePro",
						IntervalMs: 2,
						TimeoutMs:  2,
						MaxQueries: 2,
					},
					{
						ExchangeId: "Bybit",
						IntervalMs: 3,
						TimeoutMs:  3,
						MaxQueries: 3,
					},
				},
			},
			expectedErr: nil,
		},
		"invalid: duplicate exchange id": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
				},
			},
			expectedErr: fmt.Errorf("duplicate exchange id Binance"),
		},
		"invalid: invalid config (invalid exchange id)": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "InvalidExchange",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
				},
			},
			expectedErr: fmt.Errorf("invalid exchange id InvalidExchange"),
		},
		"invalid: interval_ms = 0": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 0, // invalid
						TimeoutMs:  1,
						MaxQueries: 1,
					},
				},
			},
			expectedErr: fmt.Errorf("validation for 'IntervalMs' failed on the 'gt' tag"),
		},
		"invalid: timeout_ms = 0": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  0, // invalid
						MaxQueries: 1,
					},
				},
			},
			expectedErr: fmt.Errorf("validation for 'TimeoutMs' failed on the 'gt' tag"),
		},
		"invalid: max_queries = 0": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 0, // invalid
					},
				},
			},
			expectedErr: fmt.Errorf("validation for 'MaxQueries' failed on the 'gt' tag"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.configs.Validate(constants.GetValidExchanges())
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func TestClientExchangeQueryConfigs_ValidateDelta(t *testing.T) {
	tests := map[string]struct {
		configs     types.ClientExchangeQueryConfigs
		expectedErr error
	}{
		"valid: empty configs": {
			configs:     types.ClientExchangeQueryConfigs{},
			expectedErr: nil,
		},
		"valid: populated_configs": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
					{
						ExchangeId: "CoinbasePro",
						IntervalMs: 2,
						TimeoutMs:  2,
						MaxQueries: 2,
					},
					{
						ExchangeId: "Bybit",
						IntervalMs: 3,
						TimeoutMs:  3,
						MaxQueries: 3,
					},
				},
			},
			expectedErr: nil,
		},
		"invalid: duplicate exchange id": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
				},
			},
			expectedErr: fmt.Errorf("duplicate exchange id Binance"),
		},
		"invalid: invalid config (invalid exchange id)": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "InvalidExchange",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 1,
					},
				},
			},
			expectedErr: fmt.Errorf("invalid exchange id InvalidExchange"),
		},
		"valid: interval_ms = 0": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 0, // valid for a delta
						TimeoutMs:  1,
						MaxQueries: 1,
					},
				},
			},
			expectedErr: nil,
		},
		"invalid: timeout_ms = 0": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  0, // valid for a delta
						MaxQueries: 1,
					},
				},
			},
			expectedErr: nil,
		},
		"invalid: max_queries = 0": {
			configs: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1,
						TimeoutMs:  1,
						MaxQueries: 0, // valid for a delta
					},
				},
			},
			expectedErr: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.configs.ValidateDelta(constants.GetValidExchanges())
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func TestApplyClientExchangeQueryConfigOverride(t *testing.T) {
	tests := map[string]struct {
		exchangeQueryConfigs map[types.ExchangeId]*types.ExchangeQueryConfig
		overrideConfigs      *types.ClientExchangeQueryConfigs
		expected             map[types.ExchangeId]*types.ExchangeQueryConfig
		expectedErr          error
	}{
		"valid: no overrides": {
			exchangeQueryConfigs: map[types.ExchangeId]*types.ExchangeQueryConfig{},
			overrideConfigs:      &types.ClientExchangeQueryConfigs{},
			expected:             map[types.ExchangeId]*types.ExchangeQueryConfig{},
		},
		"invalid: invalid override exchange id": {
			exchangeQueryConfigs: map[types.ExchangeId]*types.ExchangeQueryConfig{},
			overrideConfigs: &types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "InvalidExchange", // invalid
					},
				},
			},
			expectedErr: fmt.Errorf("invalid exchange id InvalidExchange"),
		},
		"valid: disable some exchanges": {
			exchangeQueryConfigs: map[types.ExchangeId]*types.ExchangeQueryConfig{
				"Binance": {
					ExchangeId: "Binance",
					IntervalMs: 1,
					TimeoutMs:  1,
					MaxQueries: 1,
				},
				"CoinbasePro": {
					ExchangeId: "CoinbasePro",
					IntervalMs: 2,
					TimeoutMs:  2,
					MaxQueries: 2,
				},
				"Huobi": {
					ExchangeId: "Huobi",
					IntervalMs: 3,
					TimeoutMs:  3,
					MaxQueries: 3,
				},
			},
			overrideConfigs: &types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						Disabled:   true,
					},
					{
						ExchangeId: "Huobi",
						Disabled:   true,
					},
				},
			},
			expected: map[types.ExchangeId]*types.ExchangeQueryConfig{
				"CoinbasePro": {
					ExchangeId: "CoinbasePro",
					IntervalMs: 2,
					TimeoutMs:  2,
					MaxQueries: 2,
				},
			},
		},
		"valid: multiple updates": {
			exchangeQueryConfigs: map[types.ExchangeId]*types.ExchangeQueryConfig{
				"Binance": {
					ExchangeId: "Binance",
					IntervalMs: 1,
					TimeoutMs:  1,
					MaxQueries: 1,
				},
				"CoinbasePro": {
					ExchangeId: "CoinbasePro",
					IntervalMs: 2,
					TimeoutMs:  2,
					MaxQueries: 2,
				},
				"Huobi": {
					ExchangeId: "Huobi",
					IntervalMs: 3,
					TimeoutMs:  3,
					MaxQueries: 3,
				},
			},
			overrideConfigs: &types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 111,
						TimeoutMs:  222,
						MaxQueries: 99,
					},
					{
						ExchangeId: "Huobi",
						IntervalMs: 333,
					},
				},
			},
			expected: map[types.ExchangeId]*types.ExchangeQueryConfig{
				"Binance": {
					ExchangeId: "Binance",
					IntervalMs: 111,
					TimeoutMs:  222,
					MaxQueries: 99,
				},
				"CoinbasePro": {
					ExchangeId: "CoinbasePro",
					IntervalMs: 2,
					TimeoutMs:  2,
					MaxQueries: 2,
				},
				"Huobi": {
					ExchangeId: "Huobi",
					IntervalMs: 333,
					TimeoutMs:  3,
					MaxQueries: 3,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := types.ApplyClientExchangeQueryConfigOverride(
				tc.exchangeQueryConfigs,
				tc.overrideConfigs,
			)
			// require.Equal(t, tc.expected, actual)
			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Zero(t, actual)
			}
		})
	}
}
