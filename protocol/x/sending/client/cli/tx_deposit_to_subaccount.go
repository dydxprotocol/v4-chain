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

// CmdDepositToSubaccount initiates a transfer from sender (an `x/banks` account)
// to a recipient (an `x/subaccounts` subaccount).
func CmdDepositToSubaccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-to-subaccount [sender_key_or_address] [recipient_address] [recipient_subaccount_number] [quantums]",
		Short: "Deposit funds from an account to a subaccount.",
		Long: `Deposit funds from an account to a subaccount.
Note, the '--from' flag is ignored as it is implied from [sender_key_or_address].
[recipient_address] and [recipient_subaccount_number] together specify the recipient subaccount.
[quantums] specifies the amount to deposit.
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Sender address validation done in `ValidateBasic()` below.
			argSender := args[0]
			err = cmd.Flags().Set(flags.FlagFrom, argSender)
			if err != nil {
				return err
			}

			argRecipientOwner := args[1]
			argRecipientNumber, err := cast.ToUint32E(args[2])
			if err != nil {
				return err
			}

			argAmount, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositToSubaccount(
				clientCtx.GetFromAddress().String(),
				satypes.SubaccountId{
					Owner:  argRecipientOwner,
					Number: argRecipientNumber,
				},
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
