package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
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
	cmd.AddCommand(CmdRegisterAffiliate())
	cmd.AddCommand(CmdRegisterBrokerAffiliate())
	return cmd
}

func CmdRegisterAffiliate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-affiliate [affiliate] [referee]",
		Short: "Register an affiliate",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.MsgRegisterAffiliate{
				Affiliate: args[0],
				Referee:   args[1],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRegisterBrokerAffiliate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-broker-affiliate [broker-id] [broker-address] [broker-fee-share-ppm]",
		Short: "Register a broker affiliate",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			brokerId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			brokerAddress := args[1]
			brokerFeeSharePpm, err := cast.ToUint32E(args[2])
			if err != nil {
				return err
			}
			msg := types.MsgRegisterBrokerAffiliate{
				BrokerAffiliate: types.BrokerAffiliate{
					BrokerId:            brokerId,
					BrokerAddress:       brokerAddress,
					BrokerFeeSharePpm:   brokerFeeSharePpm,
				},
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

