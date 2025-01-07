package cli

import (
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

// CmdWithdrawFromSubaccount initiates a transfer from sender (an `x/subaccounts` subaccount)
// to a recipient (an `x/banks` account).
func CmdWithdrawFromSubaccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-from-subaccount [sender_key_or_address] [sender_subaccount_number] [recipient_address] [quantums]",
		Short: "Withdraw funds from a subaccount to an account.",
		Long: `Withdraw funds from a subaccount to an account.
Note, the '--from' flag is ignored as it is implied from [sender_key_or_address].
[sender_key_or_address] and [sender_subaccount_number] together specify the sender subaccount.
[quantums] specifies the amount to withdraw.
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSenderOwner := args[0]
			err = cmd.Flags().Set(flags.FlagFrom, argSenderOwner)
			if err != nil {
				return err
			}
			argSenderNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			// Recipient address validation done in `ValidateBasic()` below.
			argRecipient := args[2]

			argAmount, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawFromSubaccount(
				satypes.SubaccountId{
					Owner:  clientCtx.GetFromAddress().String(),
					Number: argSenderNumber,
				},
				argRecipient,
				assettypes.AssetUsdc.Id,
				argAmount,
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
