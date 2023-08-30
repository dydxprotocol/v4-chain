package cmd

import (
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	clobflags "github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	"github.com/spf13/cobra"
)

// GetOptionWithCustomStartCmd returns a root command option with custom start commands.
func GetOptionWithCustomStartCmd() *RootCmdOption {
	option := newRootCmdOption()
	f := func(cmd *cobra.Command) {
		// Add app flags.
		appflags.AddFlagsToCmd(cmd)

		// Add daemon flags.
		daemonflags.AddDaemonFlagsToCmd(cmd)

		// Add indexer flags.
		indexer.AddIndexerFlagsToCmd(cmd)

		// Add clob flags.
		clobflags.AddClobFlagsToCmd(cmd)
	}
	option.setCustomizeStartCmd(f)
	return option
}
