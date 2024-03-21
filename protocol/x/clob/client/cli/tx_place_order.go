package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdPlaceOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-order owner subaccount_number clientId clobPairId side quantums subticks goodTilBlock",
		Short: "Broadcast message place_order. Assumes short term order placement.",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argOwner := args[0]

			argSubaccountNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			argClientId, err := cast.ToUint32E(args[2])
			if err != nil {
				return err
			}

			argClobPairId, err := cast.ToUint32E(args[3])
			if err != nil {
				return err
			}

			argSide, err := cast.ToInt32E(args[4])
			if err != nil {
				return err
			}

			argQuantums, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			argSubticks, err := cast.ToUint64E(args[6])
			if err != nil {
				return err
			}

			argGoodTilBlock, err := cast.ToUint32E(args[7])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceOrder(
				types.Order{
					OrderId: types.OrderId{
						ClientId: argClientId,
						SubaccountId: satypes.SubaccountId{
							Owner:  argOwner,
							Number: argSubaccountNumber,
						},
						ClobPairId: argClobPairId,
					},
					Side:         types.Order_Side(argSide),
					Quantums:     argQuantums,
					Subticks:     argSubticks,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: argGoodTilBlock},
				},
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
