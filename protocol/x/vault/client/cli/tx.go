package cli

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdDepositToVault())

	return cmd
}

func CmdDepositToVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-to-vault [vault_type] [vault_number] [depositor_owner] [depositor_number] [quantums]",
		Short: "Broadcast message DepositToVault",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Parse vault type.
			vaultType, err := GetVaultTypeFromString(args[0])
			if err != nil {
				return err
			}

			// Parse vault number.
			vaultNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			// Parse depositor number.
			depositorNumber, err := cast.ToUint32E(args[3])
			if err != nil {
				return err
			}

			// Parse quantums.
			quantums, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Create MsgDepositToVault.
			msg := &types.MsgDepositToVault{
				VaultId: &types.VaultId{
					Type:   vaultType,
					Number: vaultNumber,
				},
				SubaccountId: &satypes.SubaccountId{
					Owner:  args[2],
					Number: depositorNumber,
				},
				QuoteQuantums: dtypes.NewIntFromUint64(quantums),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
