package cli

import (
	"encoding/json"
	"fmt"
	"sort"
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
		Args:  cobra.ExactArgs(3),
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

			// Sort the keys to ensure deterministic ordering
			keys := make([]string, 0, len(leverageMap))
			for clobPairIdStr := range leverageMap {
				keys = append(keys, clobPairIdStr)
			}
			sort.Strings(keys)

			var clobPairLeverage []*types.LeverageEntry
			for _, clobPairIdStr := range keys {
				clobPairId, err := strconv.ParseUint(clobPairIdStr, 10, 32)
				if err != nil {
					return fmt.Errorf("invalid clob pair ID %s: %w", clobPairIdStr, err)
				}

				clobPairLeverage = append(clobPairLeverage, &types.LeverageEntry{
					ClobPairId:   uint32(clobPairId),
					CustomImfPpm: leverageMap[clobPairIdStr],
				})
			}

			msg := &types.MsgUpdateLeverage{
				SubaccountId: &satypes.SubaccountId{
					Owner:  address,
					Number: uint32(subaccountNumber),
				},
				ClobPairLeverage: clobPairLeverage,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
