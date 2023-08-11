package cmd

import (
	"github.com/dydxprotocol/v4/app/flags"
	"github.com/dydxprotocol/v4/daemons/pricefeed"
	"github.com/dydxprotocol/v4/indexer"
	"github.com/spf13/cobra"
)

// GetOptionWithCustomStartCmd returns a root command option with custom start commands.
func GetOptionWithCustomStartCmd() *RootCmdOption {
	option := newRootCmdOption()
	f := func(cmd *cobra.Command) {
		// Add app flags.
		flags.AddFlagsToCmd(cmd)

		// Add pricefeed flags.
		pricefeed.AddSharedPriceFeedFlagsToCmd(cmd)
		pricefeed.AddClientPriceFeedFlagsToCmd(cmd)

		// Add indexer flags.
		indexer.AddIndexerFlagsToCmd(cmd)
	}
	option.setCustomizeStartCmd(f)
	return option
}
