package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/spf13/cobra"
)

func CmdQueryRewardShare() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-share [address]",
		Short: "shows the reward share for the specified address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			argAddress := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.RewardShare(cmd.Context(), &types.QueryRewardShareRequest{
				Address: argAddress,
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
