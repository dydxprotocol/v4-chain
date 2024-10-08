package flags_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
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

		flags.FlagDeleveragingDaemonEnabled,
		flags.FlagDeleveragingDaemonLoopDelayMs,
		flags.FlagDeleveragingDaemonQueryPageLimit,

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

	optsMap[flags.FlagDeleveragingDaemonEnabled] = true
	optsMap[flags.FlagDeleveragingDaemonLoopDelayMs] = uint32(2222)
	optsMap[flags.FlagDeleveragingDaemonQueryPageLimit] = uint64(3333)

	optsMap[flags.FlagSDAIDaemonEnabled] = true
	optsMap[flags.FlagSDAIDaemonMockEnabled] = true
	optsMap[flags.FlagSDAIDaemonMockNoYield] = true
	optsMap[flags.FlagSDAIDaemonEthRpcEndpoint] = "test-eth-rpc-endpoint"
	optsMap[flags.FlagSDAIDaemonLoopDelayMs] = uint32(5555)

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

	// Deleveraging Daemon.
	require.Equal(t, optsMap[flags.FlagDeleveragingDaemonEnabled], r.Deleveraging.Enabled)
	require.Equal(t, optsMap[flags.FlagDeleveragingDaemonLoopDelayMs], r.Deleveraging.LoopDelayMs)
	require.Equal(t, optsMap[flags.FlagDeleveragingDaemonQueryPageLimit], r.Deleveraging.QueryPageLimit)

	// SDAI Daemon.
	require.Equal(t, optsMap[flags.FlagSDAIDaemonEnabled], r.SDAI.Enabled)
	require.Equal(t, optsMap[flags.FlagSDAIDaemonMockEnabled], r.SDAI.MockEnabled)
	require.Equal(t, optsMap[flags.FlagSDAIDaemonMockNoYield], r.SDAI.MockNoYield)
	require.Equal(t, optsMap[flags.FlagSDAIDaemonEthRpcEndpoint], r.SDAI.EthRpcEndpoint)
	require.Equal(t, optsMap[flags.FlagSDAIDaemonLoopDelayMs], r.SDAI.LoopDelayMs)

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
