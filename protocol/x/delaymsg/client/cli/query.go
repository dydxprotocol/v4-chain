package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// TODO(CORE-437): Implement query commands
	//cmd.AddCommand(CmdQueryNumMessages())
	//cmd.AddCommand(CmdQueryMessage())
	//cmd.AddCommand(CmdQueryBlockMessageIds())

	return cmd
}

// TODO(CORE-437): Implement query commands
//func CmdQueryNumMessages() *cobra.Command {}
//func CmdQueryMessage() *cobra.Command {}
//func CmdQueryBlockMessageIds() *cobra.Command {}
