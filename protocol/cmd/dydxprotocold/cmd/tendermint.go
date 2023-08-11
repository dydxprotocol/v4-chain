package cmd

import (
	"os"

	"github.com/dydxprotocol/v4/command/genprivkey"
	"github.com/spf13/cobra"
)

// AddTendermintSubcommands adds custom Tendermint subcommands.
func AddTendermintSubcommands(rootCmd *cobra.Command) {
	// Fetch Tendermint subcommand.
	tmCmd, _, err := rootCmd.Find([]string{"tendermint"})
	if err != nil {
		os.Exit(1)
	}

	// Add "gen-priv-key" command to Tendermint subcommand.
	// TODO(DEC-1079): Remove this command after updating to Cosmos `0.46.X` in favor of `init --recover`.
	tmCmd.AddCommand(genprivkey.Command())
}
