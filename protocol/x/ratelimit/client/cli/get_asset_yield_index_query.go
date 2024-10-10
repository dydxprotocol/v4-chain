package cli

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetAssetYieldIndexQuery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-asset-yield-index",
		Short: "Get the asset yield index for tDAI",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.GetAssetYieldIndexQuery(cmd.Context(), &types.GetAssetYieldIndexQueryRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
