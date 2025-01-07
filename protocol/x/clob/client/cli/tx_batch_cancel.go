package cli

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	customflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	aptypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdBatchCancel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-cancel owner subaccount_number clobPairId goodTilBlock --clientIds=\"<list of ids>\"",
		Short: "Broadcast message batch cancel for a specific clobPairId",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argOwner := args[0]
			clientIds, err := cmd.Flags().GetString("clientIds")
			if err != nil {
				return err
			}
			argClientIds := []uint32{}
			for _, idString := range strings.Fields(clientIds) {
				idUint64, err := strconv.ParseUint(idString, 10, 32)
				if err != nil {
					return err
				}
				argClientIds = append(argClientIds, uint32(idUint64))
			}

			argNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			argClobPairId, err := cast.ToUint32E(args[2])
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

			msg := types.NewMsgBatchCancel(
				satypes.SubaccountId{
					Owner:  argOwner,
					Number: argNumber,
				},
				[]types.OrderBatch{
					{
						ClobPairId: argClobPairId,
						ClientIds:  argClientIds,
					},
				},
				argGoodTilBlock,
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
