package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/spf13/cobra"
)

func CmdQueryPremiumSamples() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-premium-samples",
		Short: "Get PremiumSamples from the current funding-tick epoch",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PremiumSamples(
				context.Background(),
				&types.QueryPremiumSamplesRequest{},
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

func CmdQueryPremiumVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-premium-votes",
		Short: "Get PremiumVotes from the current funding-sample epoch",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PremiumVotes(
				context.Background(),
				&types.QueryPremiumVotesRequest{},
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
