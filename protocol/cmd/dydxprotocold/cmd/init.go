package cmd

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/configs"
	"github.com/spf13/cobra"
)

// AddInitCmdPostRunE adds a PostRunE to the `init` subcommand.
func AddInitCmdPostRunE(rootCmd *cobra.Command) {
	// Fetch init subcommand.
	initCmd, _, err := rootCmd.Find([]string{"init"})
	if err != nil {
		os.Exit(1)
	}

	// Add PostRun to configure required setups after `init`.
	initCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		// Get home directory.
		clientCtx := client.GetClientContextFromCmd(cmd)

		// Add default pricefeed exchange config toml file if it does not exist.
		configs.WriteDefaultPricefeedExchangeToml(clientCtx.HomeDir)
		return nil
	}
}
