package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/spf13/cobra"
)

func CmdQueryRevShareParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revshare-params",
		Short: "revshare params",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.MarketMapperRevenueShareParams(
				context.Background(),
				&types.QueryMarketMapperRevenueShareParams{},
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

func CmdQueryRevShareDetailsForMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revshare-details-for-market [market-id]",
		Short: "revshare details for a market",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			argId := args[0]

			marketId, err := strconv.ParseUint(argId, 10, 32)
			if err != nil {
				return err
			}

			res, err := queryClient.MarketMapperRevShareDetails(
				context.Background(),
				&types.QueryMarketMapperRevShareDetails{
					MarketId: uint32(marketId),
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

func CmdQueryUnconditionalRevShareConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unconditional-revshare-config",
		Short: "unconditional revshare config",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.UnconditionalRevShareConfig(
				context.Background(),
				&types.QueryUnconditionalRevShareConfig{},
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

func CmdQueryOrderRouterRevShare() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order-router-rev-share [address]",
		Short: "order router rev share",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.OrderRouterRevShare(
				context.Background(),
				&types.QueryOrderRouterRevShare{
					Address: args[0],
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
