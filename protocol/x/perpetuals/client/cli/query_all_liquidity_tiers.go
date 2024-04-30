package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/spf13/cobra"
)

func CmdQueryAllLiquidityTiers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all-liquidity-tiers",
		Short: "get all liquidity tiers",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AllLiquidityTiers(
				context.Background(),
				&types.QueryAllLiquidityTiersRequest{},
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
