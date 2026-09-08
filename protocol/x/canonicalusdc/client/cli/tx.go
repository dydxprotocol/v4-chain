package cli

import (
	"fmt"
	"os"
	"strconv"

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
	cmd.AddCommand(executeBackingSwapCmd(), cancelBackingSwapCmd(), updateControlsCmd(), updateParticipantsCmd())
	return cmd
}

func executeBackingSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-backing-swap [amount] [source-port] [timeout-timestamp]",
		Short: "release Noble USDC for funded Injective backing",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			timeout, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			msg := &types.MsgExecuteBackingSwap{
				Controller:       clientCtx.FromAddress.String(),
				Amount:           args[0],
				SourcePort:       args[1],
				TimeoutTimestamp: timeout,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cancelBackingSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-backing-swap [amount]",
		Short: "return uncommitted Injective backing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgCancelBackingSwap{Controller: clientCtx.FromAddress.String(), Amount: args[0]}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
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

func updateParticipantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-participants [participants-json-file]",
		Short: "replace canonical USDC participants",
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
			var request types.MsgUpdateParticipants
			if err := clientCtx.Codec.UnmarshalJSON(contents, &request); err != nil {
				return err
			}
			request.Authority = clientCtx.FromAddress.String()
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &request)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
