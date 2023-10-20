package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CmdListStatefulOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-stateful-orders",
		Short: "list all stateful orders",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAllStatefulOrdersRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllStatefulOrders(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetStatefulOrderCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-stateful-order-count [owner] [account-number]",
		Short: "shows stateful order count for a subaccount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argOwner := args[0]
			var argNumber uint32
			if argNumber, err = cast.ToUint32E(args[1]); err != nil {
				return err
			}

			if _, err := sdk.AccAddressFromBech32(argOwner); err != nil {
				return status.Error(
					codes.InvalidArgument,
					fmt.Sprintf("Invalid owner address: %v", err),
				)
			}

			params := &types.QueryStatefulOrderCountRequest{
				SubaccountId: &satypes.SubaccountId{
					Owner:  argOwner,
					Number: argNumber,
				},
			}

			var res *types.QueryStatefulOrderCountResponse
			if res, err = queryClient.StatefulOrderCount(context.Background(), params); err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
