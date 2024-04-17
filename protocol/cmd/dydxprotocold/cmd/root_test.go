package cmd_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/config"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/cmd/dydxprotocold/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/stretchr/testify/require"
)

func TestNewRootCmd_UsesClientConfig(t *testing.T) {
	tempDir := t.TempDir()

	config.SetupConfig()

	// Set the client config to point to a fake chain id since this is a required option
	{
		option := cmd.GetOptionWithCustomStartCmd()
		rootCmd := cmd.NewRootCmd(option, tempDir)

		cmd.AddTendermintSubcommands(rootCmd)
		cmd.AddInitCmdPostRunE(rootCmd)
		rootCmd.SetArgs([]string{"config", "set", "client", "chain-id", "fakeChainId"})
		require.NoError(t, svrcmd.Execute(rootCmd, constants.AppDaemonName, tempDir))
	}

	// Set the client config to point to a fake address
	{
		option := cmd.GetOptionWithCustomStartCmd()
		rootCmd := cmd.NewRootCmd(option, tempDir)

		cmd.AddTendermintSubcommands(rootCmd)
		cmd.AddInitCmdPostRunE(rootCmd)
		rootCmd.SetArgs([]string{"config", "set", "client", "node", "fakeTestAddress"})
		require.NoError(t, svrcmd.Execute(rootCmd, constants.AppDaemonName, tempDir))
	}

	// Run a query command (that will fail) to ensure that we are reading the client config
	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmd(option, tempDir)

	cmd.AddTendermintSubcommands(rootCmd)
	cmd.AddInitCmdPostRunE(rootCmd)
	rootCmd.SetArgs([]string{"query", "auth", "params"})
	require.ErrorContains(t, svrcmd.Execute(rootCmd, constants.AppDaemonName, tempDir), "fakeTestAddress")
}
