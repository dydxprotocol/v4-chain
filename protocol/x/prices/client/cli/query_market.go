package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/spf13/cobra"
)

func CmdListMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-market",
		Short: "list all market",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllMarketsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllMarkets(context.Background(), params)
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

func CmdShowMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-market [id]",
		Short: "shows a market",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argId := args[0]

			id, err := strconv.ParseUint(argId, 10, 32)
			if err != nil {
				return err
			}

			params := &types.QueryMarketRequest{
				Id: uint32(id),
			}

			res, err := queryClient.Market(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
