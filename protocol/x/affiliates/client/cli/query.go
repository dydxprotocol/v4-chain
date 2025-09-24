package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group x/affiliates queries under a subcommand.
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryAffiliateTiers(),
		GetCmdQueryAffiliateInfo(),
		GetCmdQueryReferredBy(),
		GetCmdQueryAffiliateWhitelist(),
		GetCmdQueryAffiliateOverrides(),
		GetCmdQueryAffiliateParameters(),
	)
	return cmd
}

func GetCmdQueryAffiliateTiers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affiliate-tiers",
		Short: "Query affiliate tiers",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AllAffiliateTiers(context.Background(), &types.AllAffiliateTiersRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func GetCmdQueryAffiliateInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affiliate-info [affiliate-address]",
		Short: "Query affiliate info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AffiliateInfo(context.Background(), &types.AffiliateInfoRequest{
				Address: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func GetCmdQueryReferredBy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "referred-by [address]",
		Short: "Query the referee that referred the given addresss",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ReferredBy(context.Background(), &types.ReferredByRequest{
				Address: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func GetCmdQueryAffiliateWhitelist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affiliate-whitelist",
		Short: "Query affiliate whitelist",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AffiliateWhitelist(context.Background(), &types.AffiliateWhitelistRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func GetCmdQueryAffiliateOverrides() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affiliate-overrides",
		Short: "Query affiliate overrides",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AffiliateOverrides(context.Background(), &types.AffiliateOverridesRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func GetCmdQueryAffiliateParameters() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affiliate-parameters",
		Short: "Query affiliate parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AffiliateParameters(context.Background(), &types.AffiliateParametersRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}
