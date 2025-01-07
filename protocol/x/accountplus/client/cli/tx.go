package cli

import (
	"encoding/base64"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
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
	cmd.AddCommand(CmdAddAuthenticator())
	cmd.AddCommand(CmdRemoveAuthenticator())
	return cmd
}

func CmdAddAuthenticator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-authenticator [account] [authenticator_type] [data]",
		Short: "Registers an authenticator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			config, err := base64.StdEncoding.DecodeString(args[2])
			if err != nil {
				return err
			}
			msg := types.MsgAddAuthenticator{
				Sender:            args[0],
				AuthenticatorType: args[1],
				Data:              config,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRemoveAuthenticator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-authenticator [sender] [authenticator_id]",
		Short: "Removes an authenticator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			id, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			msg := types.MsgRemoveAuthenticator{
				Sender: args[0],
				Id:     id,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
