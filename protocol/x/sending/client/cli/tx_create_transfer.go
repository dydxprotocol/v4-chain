package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	customflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	aptypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-transfer sender_owner sender_number recipient_owner recipient_number quantums",
		Short: "Broadcast message CreateTransfer",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSenderOwner := args[0]
			argSenderNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			argRecipientOwner := args[2]
			argRecipientNumber, err := cast.ToUint32E(args[3])
			if err != nil {
				return err
			}

			argAmount, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateTransfer(
				&types.Transfer{
					Sender: satypes.SubaccountId{
						Owner:  argSenderOwner,
						Number: argSenderNumber,
					},
					Recipient: satypes.SubaccountId{
						Owner:  argRecipientOwner,
						Number: argRecipientNumber,
					},
					AssetId: assettypes.AssetUsdc.Id,
					Amount:  argAmount,
				},
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			authenticatorIds, err := customflags.GetPermisionedKeyAuthenticatorsForExtOptions(cmd)
			if err == nil && len(authenticatorIds) > 0 {
				value, err := codectypes.NewAnyWithValue(
					&aptypes.TxExtension{
						SelectedAuthenticators: authenticatorIds,
					},
				)
				if err != nil {
					return err
				}
				txf = txf.WithNonCriticalExtensionOptions(value)
			}
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	customflags.AddTxPermissionedKeyFlagsToCmd(cmd)

	return cmd
}
