package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
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

	cmd.AddCommand(CmdCreateMarketPermissionless())

	return cmd
}

// CmdCreateMarketPermissionless is the CLI command for creating a permissionless market.
func CmdCreateMarketPermissionless() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-market [ticker] [address]",
		Short: "Create new market with permissionless access",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			ticker, address := args[0], args[1]

			// Create MsgCreateMarketPermissionless.
			msg := &types.MsgCreateMarketPermissionless{
				Ticker: ticker,
				SubaccountId: &satypes.SubaccountId{
					Owner:  address,
					Number: 0,
				},
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
