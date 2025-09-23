package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/spf13/cobra"
)

func CmdQueryLeverage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leverage [address] [subaccount-number]",
		Short: "Query leverage for a subaccount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			// Parse address and subaccount number
			address := args[0]
			subaccountNumber, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid subaccount number %s: %w", args[1], err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLeverageRequest{
				Owner:  address,
				Number: uint32(subaccountNumber),
			}

			res, err := queryClient.Leverage(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
