package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/spf13/cobra"
)

func CmdShowMarketParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-market-param [id]",
		Short: "shows a market param",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argId := args[0]

			id, err := strconv.ParseUint(argId, 10, 32)
			if err != nil {
				return err
			}

			param := &types.QueryMarketParamRequest{
				Id: uint32(id),
			}

			res, err := queryClient.MarketParam(context.Background(), param)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListMarketParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-market-param",
		Short: "list all market params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllMarketParamsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllMarketParams(context.Background(), params)
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
