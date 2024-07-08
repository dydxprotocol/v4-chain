package cli

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// CmdUpdateMarketPrices updates the conversion rate and ethereum block number for sDAI.
func CmdUpdateMarketPrices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-market-prices [sender_key_or_address] [conversion_rate] [ethereum_block_number]",
		Short: "Update the conversion rate and ethereum block number for sDAI.",
		Long: `Update the conversion rate and ethereum block number for sDAI.
Note, the '--from' flag is ignored as it is implied from [sender_key_or_address].
[conversion_rate] is the conversion rate of sDAI to USD.
[ethereum_block_number] is the block number on Ethereum.
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSenderOwner := args[0]
			err = cmd.Flags().Set(flags.FlagFrom, argSenderOwner)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			conversionRate := args[1]
			ethereumBlockNumber := args[2]

			msg := types.NewMsgUpdateSDAIConversionRate(
				clientCtx.GetFromAddress(),
				conversionRate,
				ethereumBlockNumber,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
