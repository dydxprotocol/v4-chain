package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/spf13/cobra"
)

func CmdGetBlockRateLimitConfiguration() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block-rate-limit-config",
		Short: "get the block rate limit configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBlockRateLimitConfigurationRequest{}

			res, err := queryClient.BlockRateLimitConfiguration(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
