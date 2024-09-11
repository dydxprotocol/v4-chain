package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/spf13/cobra"
)

func CmdQueryListingVaultDepositParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listing-vault-deposit-params",
		Short: "listing vault deposit params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ListingVaultDepositParams(
				context.Background(),
				&types.QueryListingVaultDepositParams{},
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
