package cmd_test

import (
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/cmd/dydxprotocold/cmd"
	"github.com/stretchr/testify/require"
	"testing"
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
		require.NoError(t, svrcmd.Execute(rootCmd, app.AppDaemonName, tempDir))
	}

	// Set the client config to point to a fake address
	{
		option := cmd.GetOptionWithCustomStartCmd()
		rootCmd := cmd.NewRootCmd(option, tempDir)

		cmd.AddTendermintSubcommands(rootCmd)
		cmd.AddInitCmdPostRunE(rootCmd)
		rootCmd.SetArgs([]string{"config", "set", "client", "node", "fakeTestAddress"})
		require.NoError(t, svrcmd.Execute(rootCmd, app.AppDaemonName, tempDir))
	}

	// Run a query command (that will fail) to ensure that we are reading the client config
	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmd(option, tempDir)

	cmd.AddTendermintSubcommands(rootCmd)
	cmd.AddInitCmdPostRunE(rootCmd)
	rootCmd.SetArgs([]string{"query", "auth", "params"})
	require.ErrorContains(t, svrcmd.Execute(rootCmd, app.AppDaemonName, tempDir), "fakeTestAddress")
}
