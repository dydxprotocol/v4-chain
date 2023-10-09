package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdListClobPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-clob-pair",
		Short: "list all clob_pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllClobPairRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ClobPairAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowClobPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-clob-pair [index]",
		Short: "shows a clob_pair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argId, err := cast.ToUint32E(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryGetClobPairRequest{
				Id: argId,
			}

			res, err := queryClient.ClobPair(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
