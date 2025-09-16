package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func CmdUpdateLeverage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-leverage [address] [subaccount-number] [leverage-map]",
		Short: "Update leverage for perpetuals",
		Long: `Update leverage for perpetuals. The leverage-map should be a JSON string with perpetual IDs as keys and leverage amounts as values.
Example: '{"0": 5, "1": 10}' sets leverage of 5 for perpetual 0 and leverage of 10 for perpetual 1.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse address and subaccount number
			address := args[0]
			subaccountNumber, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid subaccount number %s: %w", args[1], err)
			}

			// Parse leverage map
			var leverageMap map[string]uint32
			if err := json.Unmarshal([]byte(args[2]), &leverageMap); err != nil {
				return fmt.Errorf("invalid leverage map JSON: %w", err)
			}

			// Convert string keys to uint32 and create LeverageEntry slice
			var perpetualLeverage []*types.LeverageEntry
			for perpetualIdStr, leverage := range leverageMap {
				perpetualId, err := strconv.ParseUint(perpetualIdStr, 10, 32)
				if err != nil {
					return fmt.Errorf("invalid perpetual ID %s: %w", perpetualIdStr, err)
				}
				perpetualLeverage = append(perpetualLeverage, &types.LeverageEntry{
					PerpetualId: uint32(perpetualId),
					Leverage:    leverage,
				})
			}

			msg := &types.MsgUpdateLeverage{
				SubaccountId: &satypes.SubaccountId{
					Owner:  address,
					Number: uint32(subaccountNumber),
				},
				PerpetualLeverage: perpetualLeverage,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
