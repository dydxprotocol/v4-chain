package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/spf13/cobra"
	"strconv"
)

func CmdShowMarketPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-market-price [id]",
		Short: "shows a market price",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argId := args[0]

			id, err := strconv.ParseUint(argId, 10, 32)
			if err != nil {
				return err
			}

			price := &types.QueryMarketPriceRequest{
				Id: uint32(id),
			}

			res, err := queryClient.MarketPrice(context.Background(), price)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListMarketPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-market-price",
		Short: "list all market prices",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllMarketPricesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllMarketPrices(context.Background(), params)
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
