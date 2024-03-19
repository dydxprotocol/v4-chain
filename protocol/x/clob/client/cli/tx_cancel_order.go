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

func CmdCancelOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-order owner subaccount_number clientId clobPairId goodTilBlock",
		Short: "Broadcasts message cancel_order. Assumes short term order cancellation.",
		Args:  cobra.ExactArgs(5),
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

			argGoodTilBlock, err := cast.ToUint32E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelOrderShortTerm(
				types.OrderId{
					ClobPairId: argClobPairId,
					ClientId:   argClientId,
					SubaccountId: satypes.SubaccountId{
						Owner:  argOwner,
						Number: argSubaccountNumber,
					},
				},
				argGoodTilBlock,
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
