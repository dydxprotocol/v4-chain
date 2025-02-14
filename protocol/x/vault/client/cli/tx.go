package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

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

	cmd.AddCommand(CmdDepositToMegavault())
	cmd.AddCommand(CmdSetVaultParams())
	cmd.AddCommand(CmdAllocateToVault())
	cmd.AddCommand(CmdRetrieveFromVault())
	cmd.AddCommand(CmdWithdrawFromMegavault())
	cmd.AddCommand(CmdUpdateDefaultQuotingParams())

	return cmd
}

func CmdDepositToMegavault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-to-megavault [depositor_owner] [depositor_number] [quantums]",
		Short: "Broadcast message DepositToMegavault",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Parse depositor number.
			depositorNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			// Parse quantums.
			quantums, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Create MsgDepositToMegavault.
			msg := &types.MsgDepositToMegavault{
				SubaccountId: &satypes.SubaccountId{
					Owner:  args[0],
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

func CmdSetVaultParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-vault-params [authority] [vault_type] [vault_number] [status] [quoting_params_json]",
		Short: "Broadcast message SetVaultParams",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Parse vault type.
			vaultType, err := GetVaultTypeFromString(args[1])
			if err != nil {
				return err
			}

			// Parse vault number.
			vaultNumber, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			// Parse status.
			status, err := GetVaultStatusFromString(args[3])
			if err != nil {
				return err
			}

			// Parse quoting_params (optional).
			var quotingParams *types.QuotingParams
			if args[4] != "" {
				if err := json.Unmarshal([]byte(args[4]), &quotingParams); err != nil {
					return fmt.Errorf("invalid quoting params JSON: %w", err)
				}
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Create MsgSetVaultParams.
			msg := &types.MsgSetVaultParams{
				Authority: args[0],
				VaultId: types.VaultId{
					Type:   vaultType,
					Number: uint32(vaultNumber),
				},
				VaultParams: types.VaultParams{
					Status:        status,
					QuotingParams: quotingParams, // nil if not provided.
				},
			}

			// Validate vault params.
			if err := msg.VaultParams.Validate(); err != nil {
				return err
			}

			// Broadcast or generate the transaction.
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add the necessary flags.
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdAllocateToVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocate-to-vault [authority] [vault_type] [vault_number] [quote_quantums]",
		Short: "Broadcast message AllocateToVault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Parse vault type.
			vaultType, err := GetVaultTypeFromString(args[1])
			if err != nil {
				return err
			}

			// Parse vault number.
			vaultNumber, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			// Parse quantums.
			quantums, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Create MsgAllocateToVault.
			msg := &types.MsgAllocateToVault{
				Authority: args[0],
				VaultId: types.VaultId{
					Type:   vaultType,
					Number: uint32(vaultNumber),
				},
				QuoteQuantums: dtypes.NewIntFromUint64(quantums),
			}

			// Broadcast or generate the transaction.
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add the necessary flags.
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRetrieveFromVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retrieve-from-vault [authority] [vault_type] [vault_number] [quote_quantums]",
		Short: "Broadcast message RetrieveFromVault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Parse vault type.
			vaultType, err := GetVaultTypeFromString(args[1])
			if err != nil {
				return err
			}

			// Parse vault number.
			vaultNumber, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			// Parse quantums.
			quantums, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Create MsgRetrieveFromVault.
			msg := &types.MsgRetrieveFromVault{
				Authority: args[0],
				VaultId: types.VaultId{
					Type:   vaultType,
					Number: uint32(vaultNumber),
				},
				QuoteQuantums: dtypes.NewIntFromUint64(quantums),
			}

			// Broadcast or generate the transaction.
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add the necessary flags.
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdWithdrawFromMegavault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-from-megavault [withdrawer_owner] [withdrawer_number] [shares] [min_quote_quantums]",
		Short: "Broadcast message WithdrawFromMegavault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse withdrawer owner and number
			withdrawerOwner := args[0]
			withdrawerNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			// Parse shares.
			shares, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			// Parse min quote quantums.
			minQuoteQuantums, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			// Create MsgWithdrawFromMegavault.
			msg := &types.MsgWithdrawFromMegavault{
				SubaccountId: satypes.SubaccountId{
					Owner:  withdrawerOwner,
					Number: withdrawerNumber,
				},
				Shares:           types.NumShares{NumShares: dtypes.NewIntFromUint64(shares)},
				MinQuoteQuantums: dtypes.NewIntFromUint64(minQuoteQuantums),
			}

			// Broadcast or generate the transaction.
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add the necessary flags.
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateDefaultQuotingParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-default-quoting-params [authority] [quoting_params_json]",
		Short: "Broadcast message UpdateDefaultQuotingParams",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse quoting params.
			var quotingParams types.QuotingParams
			if err := json.Unmarshal([]byte(args[1]), &quotingParams); err != nil {
				return fmt.Errorf("invalid quoting params JSON: %w", err)
			}
			if err := quotingParams.Validate(); err != nil {
				return err
			}

			// Create MsgUpdateDefaultQuotingParams.
			msg := &types.MsgUpdateDefaultQuotingParams{
				Authority:            args[0],
				DefaultQuotingParams: quotingParams,
			}

			// Broadcast or generate the transaction.
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add the necessary flags.
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
