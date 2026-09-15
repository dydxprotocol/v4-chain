package cli

import (
	"fmt"
	"os"

	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transaction subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(updateControlsCmd())
	return cmd
}

func updateControlsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-controls [controls-json-file]",
		Short: "replace canonical USDC controls",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			contents, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			var controls types.Controls
			if err := clientCtx.Codec.UnmarshalJSON(contents, &controls); err != nil {
				return err
			}
			msg := &types.MsgUpdateControls{Authority: clientCtx.FromAddress.String(), Controls: controls}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
