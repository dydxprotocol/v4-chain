package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdPlaceOrder())
	cmd.AddCommand(CmdCancelOrder())
	batchCancelCmd := CmdBatchCancel()
	batchCancelCmd.PersistentFlags().String("clientIds", "", "A list of client ids to to batch cancel")
	cmd.AddCommand(batchCancelCmd)
	cmd.AddCommand(CmdUpdateLeverage())
	// this line is used by starport scaffolding # 1

	return cmd
}
