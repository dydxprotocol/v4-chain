package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group clob queries under a subcommand.
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListClobPair())
	cmd.AddCommand(CmdShowClobPair())
	cmd.AddCommand(CmdGetBlockRateLimitConfiguration())
	cmd.AddCommand(CmdGetEquityTierLimitConfig())
	cmd.AddCommand(CmdGetLiquidationsConfiguration())
	cmd.AddCommand(CmdQueryStatefulOrder())
	cmd.AddCommand(CmdQueryLeverage())

	return cmd
}
