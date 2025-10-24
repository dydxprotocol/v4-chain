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
	cmd.AddCommand(CmdQueryMarketFeeDiscountParams())
	cmd.AddCommand(CmdQueryStakingTiers())
	cmd.AddCommand(CmdQueryUserStakingTier())

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

func CmdQueryMarketFeeDiscountParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-market-fee-discount-params [clob_pair_id]",
		Short: "get the fee discount parameters for all markets or a specific CLOB pair",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				// Query all market fee discount params
				res, err := queryClient.AllMarketFeeDiscountParams(
					context.Background(),
					&types.QueryAllMarketFeeDiscountParamsRequest{},
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

				// Query specific market fee discount params
				res, err := queryClient.PerMarketFeeDiscountParams(
					context.Background(),
					&types.QueryPerMarketFeeDiscountParamsRequest{
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

func CmdQueryStakingTiers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking-tiers",
		Short: "get all staking tiers",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.StakingTiers(
				context.Background(),
				&types.QueryStakingTiersRequest{},
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

func CmdQueryUserStakingTier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-staking-tier [address]",
		Short: "get the staking tier and discount of a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.UserStakingTier(
				context.Background(),
				&types.QueryUserStakingTierRequest{
					Address: args[0],
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
