package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/spf13/cobra"
)

func CmdPendingSendPackets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-send-packets",
		Short: "gets all pending send packets",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AllPendingSendPackets(cmd.Context(), &types.QueryAllPendingSendPacketsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
