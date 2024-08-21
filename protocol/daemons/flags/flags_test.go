package flags_test

import (
	"fmt"
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
		flags.FlagPanicOnDaemonFailureEnabled,
		flags.FlagMaxDaemonUnhealthySeconds,

		flags.FlagBridgeDaemonEnabled,
		flags.FlagBridgeDaemonLoopDelayMs,

		flags.FlagLiquidationDaemonEnabled,
		flags.FlagLiquidationDaemonLoopDelayMs,
		flags.FlagLiquidationDaemonQueryPageLimit,
		flags.FlagLiquidationDaemonResponsePageLimit,

		flags.FlagPriceDaemonEnabled,
		flags.FlagPriceDaemonLoopDelayMs,
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
	optsMap[flags.FlagPanicOnDaemonFailureEnabled] = false
	optsMap[flags.FlagMaxDaemonUnhealthySeconds] = uint32(1234)

	optsMap[flags.FlagBridgeDaemonEnabled] = true
	optsMap[flags.FlagBridgeDaemonLoopDelayMs] = uint32(1111)
	optsMap[flags.FlagBridgeDaemonEthRpcEndpoint] = "test-eth-rpc-endpoint"

	optsMap[flags.FlagLiquidationDaemonEnabled] = true
	optsMap[flags.FlagLiquidationDaemonLoopDelayMs] = uint32(2222)
	optsMap[flags.FlagLiquidationDaemonQueryPageLimit] = uint64(3333)
	optsMap[flags.FlagLiquidationDaemonResponsePageLimit] = uint64(4444)

	optsMap[flags.FlagPriceDaemonEnabled] = true
	optsMap[flags.FlagPriceDaemonLoopDelayMs] = uint32(4444)

	mockOpts := mocks.AppOptions{}
	mockOpts.On("Get", mock.Anything).
		Return(func(key string) interface{} {
			return optsMap[key]
		})

	r := flags.GetDaemonFlagValuesFromOptions(&mockOpts)

	// Shared.
	require.Equal(t, optsMap[flags.FlagUnixSocketAddress], r.Shared.SocketAddress)
	require.Equal(t, optsMap[flags.FlagPanicOnDaemonFailureEnabled], r.Shared.PanicOnDaemonFailureEnabled)
	require.Equal(
		t,
		optsMap[flags.FlagMaxDaemonUnhealthySeconds],
		r.Shared.MaxDaemonUnhealthySeconds,
	)

	// Bridge Daemon.
	require.Equal(t, optsMap[flags.FlagBridgeDaemonEnabled], r.Bridge.Enabled)
	require.Equal(t, optsMap[flags.FlagBridgeDaemonLoopDelayMs], r.Bridge.LoopDelayMs)
	require.Equal(t, optsMap[flags.FlagBridgeDaemonEthRpcEndpoint], r.Bridge.EthRpcEndpoint)

	// Liquidation Daemon.
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonEnabled], r.Liquidation.Enabled)
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonLoopDelayMs], r.Liquidation.LoopDelayMs)
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonQueryPageLimit], r.Liquidation.QueryPageLimit)
	require.Equal(t, optsMap[flags.FlagLiquidationDaemonResponsePageLimit], r.Liquidation.ResponsePageLimit)

	// Price Daemon.
	require.Equal(t, optsMap[flags.FlagPriceDaemonEnabled], r.Price.Enabled)
	require.Equal(t, optsMap[flags.FlagPriceDaemonLoopDelayMs], r.Price.LoopDelayMs)
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
