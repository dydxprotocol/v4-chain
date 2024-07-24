package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group ratelimit queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListLimitParams())
	cmd.AddCommand(CmdQueryCapacityByDenom())
	cmd.AddCommand(CmdPendingSendPackets())
	cmd.AddCommand(CmdGetSDAIPriceQuery())
	cmd.AddCommand(CmdGetAssetYieldIndexQuery())

	return cmd
}
