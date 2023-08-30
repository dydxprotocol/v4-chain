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
		Use:   "cancel-order owner number clientId goodTilBlock",
		Short: "Broadcast message cancel_order",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argOwner := args[0]

			argNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			argClientId, err := cast.ToUint32E(args[2])
			if err != nil {
				return err
			}

			argGoodTilBlock, err := cast.ToUint32E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelOrderShortTerm(
				types.OrderId{
					ClientId: argClientId,
					SubaccountId: satypes.SubaccountId{
						Owner:  argOwner,
						Number: argNumber,
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
