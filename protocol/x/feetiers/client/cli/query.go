package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group feetiers queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryPerpetualFeeParams())
	cmd.AddCommand(CmdQueryUserFeeTier())
	cmd.AddCommand(CmdQueryFeeDiscountCampaignParams())

	return cmd
}

func CmdQueryPerpetualFeeParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-perpetual-fee-params",
		Short: "get the PerpetualFeeParams",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PerpetualFeeParams(
				context.Background(),
				&types.QueryPerpetualFeeParamsRequest{},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryUserFeeTier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-user-fee-tier",
		Short: "get the fee tier of a User",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.UserFeeTier(
				context.Background(),
				&types.QueryUserFeeTierRequest{
					User: args[0],
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryFeeDiscountCampaignParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-fee-discount-compaign-params [clob_pair_id]",
		Short: "get the FeeDiscountCampaignParams for all or a specific CLOB pair",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				// Query all fee discount params
				res, err := queryClient.AllFeeDiscountCampaignParams(
					context.Background(),
					&types.QueryAllFeeDiscountCampaignParamsRequest{},
				)
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			} else {
				// Parse CLOB pair ID
				var clobPairID uint32
				if _, err := fmt.Sscanf(args[0], "%d", &clobPairID); err != nil {
					return fmt.Errorf("clob_pair_id %s not a valid uint32", args[0])
				}

				// Query specific fee discount params
				res, err := queryClient.FeeDiscountCampaignParams(
					context.Background(),
					&types.QueryFeeDiscountCampaignParamsRequest{
						ClobPairId: clobPairID,
					},
				)
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
