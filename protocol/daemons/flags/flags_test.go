package flags_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddDaemonFlagsToCmd(t *testing.T) {
	cmd := cobra.Command{}

	flags.AddDaemonFlagsToCmd(&cmd)
	tests := []string{
		flags.FlagUnixSocketAddress,

		flags.FlagBridgeDaemonEnabled,
		flags.FlagBridgeDaemonLoopDelayMs,

		flags.FlagLiquidationDaemonEnabled,
		flags.FlagLiquidationDaemonLoopDelayMs,
		flags.FlagLiquidationDaemonSubaccountPageLimit,

		flags.FlagPriceDaemonEnabled,
		flags.FlagPriceDaemonLoopDelayMs,
		flags.FlagPriceDaemonExchangeConfigOverride,
	}

	for _, v := range tests {
		testName := fmt.Sprintf("Has %s flag", v)
		t.Run(testName, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), v)
		})
	}
}

func TestGetDaemonFlagValuesFromOptions_Custom(t *testing.T) {
	optsMap := make(map[string]interface{})

	optsMap[flags.FlagUnixSocketAddress] = "test-socket-address"

	optsMap[flags.FlagBridgeDaemonEnabled] = true
	optsMap[flags.FlagBridgeDaemonLoopDelayMs] = uint32(1111)
	optsMap[flags.FlagBridgeDaemonEthRpcEndpoint] = "test-eth-rpc-endpoint"

	optsMap[flags.FlagLiquidationDaemonEnabled] = true
	optsMap[flags.FlagLiquidationDaemonLoopDelayMs] = uint32(2222)
	optsMap[flags.FlagLiquidationDaemonSubaccountPageLimit] = uint64(3333)
	optsMap[flags.FlagLiquidationDaemonRequestChunkSize] = uint64(4444)

	optsMap[flags.FlagPriceDaemonEnabled] = true
	optsMap[flags.FlagPriceDaemonLoopDelayMs] = uint32(4444)
	optsMap[flags.FlagPriceDaemonExchangeConfigOverride] = `{"exchange_query_configs":[]}`

	mockOpts := mocks.AppOptions{}
	mockOpts.On("Get", mock.Anything).
		Return(func(key string) interface{} {
			return optsMap[key]
		})

	r := flags.GetDaemonFlagValuesFromOptions(&mockOpts)

	// Shared.
	require.Equal(t, optsMap[flags.FlagUnixSocketAddress], r.Shared.SocketAddress)

	// Bridge Daemon.
	require.Equal(t, optsMap[flags.FlagBridgeDaemonEnabled], r.Bridge.Enabled)
	require.Equal(t, optsMap[flags.FlagBridgeDaemonLoopDelayMs], r.Bridge.LoopDelayMs)
	require.Equal(t, optsMap[flags.FlagBridgeDaemonEthRpcEndpoint], r.Bridge.EthRpcEndpoint)

	// Liquidation Daemon.
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonEnabled], r.Liquidation.Enabled)
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonLoopDelayMs], r.Liquidation.LoopDelayMs)
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonSubaccountPageLimit], r.Liquidation.SubaccountPageLimit)
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonRequestChunkSize], r.Liquidation.RequestChunkSize)

	// Price Daemon.
	require.Equal(t, optsMap[flags.FlagPriceDaemonEnabled], r.Price.Enabled)
	require.Equal(t, optsMap[flags.FlagPriceDaemonLoopDelayMs], r.Price.LoopDelayMs)
	require.Equal(t, optsMap[flags.FlagPriceDaemonExchangeConfigOverride], r.Price.ExchangeConfigOverride)
}

func TestGetDaemonFlagValuesFromOptions_Default(t *testing.T) {
	mockOpts := mocks.AppOptions{}
	mockOpts.On("Get", mock.Anything).
		Return(func(key string) interface{} {
			return nil
		})

	r := flags.GetDaemonFlagValuesFromOptions(&mockOpts)
	d := flags.GetDefaultDaemonFlags()
	require.Equal(t, d, r)
}

func TestParseExchangeConfigOverride(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    types.ClientExchangeQueryConfigs
		expectedErr error
	}{
		"invalid: invalid json": {
			input:       "",
			expected:    types.ClientExchangeQueryConfigs{},
			expectedErr: fmt.Errorf("Error unmarshalling exchange config override: unexpected end of JSON input"),
		},
		"valid: empty json object": {
			input:       `{}`,
			expected:    types.ClientExchangeQueryConfigs{},
			expectedErr: nil,
		},
		"valid: empty exchange query configs": {
			input: `{"exchange_query_configs":[]}`,
			expected: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{},
			},
			expectedErr: nil,
		},
		"invalid: invalid exchange id": {
			input:       `{"exchange_query_configs":[{"exchange_id":"invalid"}]}`,
			expected:    types.ClientExchangeQueryConfigs{},
			expectedErr: fmt.Errorf("Error validating exchange config override: invalid exchange id invalid"),
		},
		"valid client exchange query configs - disable some exchanges": {
			input: `{"exchange_query_configs":[{"exchange_id":"Binance","disabled":true},{"exchange_id":"CoinbasePro","disabled":true}]}`,
			expected: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						Disabled:   true,
					},
					{
						ExchangeId: "CoinbasePro",
						Disabled:   true,
					},
				},
			},
			expectedErr: nil,
		},
		"valid client exchange query configs - multiple updates": {
			input: `{"exchange_query_configs":[` +
				`{"exchange_id":"Binance","interval_ms":1000,"timeout_ms":2000,"max_queries":4},` +
				`{"exchange_id":"CoinbasePro","disabled":true},` +
				`{"exchange_id":"Huobi","interval_ms":6000},` +
				`{"exchange_id":"Bybit","disabled":true}]}`,
			expected: types.ClientExchangeQueryConfigs{
				ExchangeQueryConfigs: []*types.ExchangeQueryConfig{
					{
						ExchangeId: "Binance",
						IntervalMs: 1000,
						TimeoutMs:  2000,
						MaxQueries: 4,
					},
					{
						ExchangeId: "CoinbasePro",
						Disabled:   true,
					},
					{
						ExchangeId: "Huobi",
						IntervalMs: 6000,
					},
					{
						ExchangeId: "Bybit",
						Disabled:   true,
					},
				},
			},
			expectedErr: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := flags.ParseExchangeConfigOverride(tc.input)
			// require.Equal(t, tc.expected, actual)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}
