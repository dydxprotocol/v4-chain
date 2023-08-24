package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strconv"

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

	cmd.AddCommand(CmdQueryNumMessages())
	cmd.AddCommand(CmdQueryMessage())
	cmd.AddCommand(CmdQueryBlockMessageIds())

	return cmd
}

func CmdQueryNumMessages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-num-messages",
		Short: "get the number of delayed messages",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.NumMessages(
				context.Background(),
				&types.QueryNumMessagesRequest{},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-message",
		Short: "get the delayed message with the given id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			argId := args[0]

			id, err := strconv.ParseUint(argId, 10, 32)
			if err != nil {
				return err
			}

			res, err := queryClient.Message(
				context.Background(),
				&types.QueryMessageRequest{
					Id: uint32(id),
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryBlockMessageIds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block-message-ids",
		Short: "get the ids of the message to be executed at a given block",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			argId := args[0]

			id, err := strconv.ParseInt(argId, 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.BlockMessageIds(
				context.Background(),
				&types.QueryBlockMessageIdsRequest{
					BlockHeight: id,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
