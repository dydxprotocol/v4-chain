package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdQueryStatefulOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stateful-order subaccount_owner subaccount_number client_id clob_pair_id order_flags",
		Short: "queries a stateful order by id",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			owner := args[0]

			number, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			clientId, err := cast.ToUint32E(args[2])
			if err != nil {
				return err
			}

			clobPairId, err := cast.ToUint32E(args[3])
			if err != nil {
				return err
			}

			orderFlag, err := cast.ToUint32E(args[4])
			if err != nil {
				return err
			}

			req := &types.QueryStatefulOrderRequest{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{
						Owner:  owner,
						Number: number,
					},
					ClientId:   clientId,
					ClobPairId: clobPairId,
					OrderFlags: orderFlag,
				},
			}

			res, err := queryClient.StatefulOrder(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
