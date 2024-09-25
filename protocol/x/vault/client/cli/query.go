package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group x/vault queries under a subcommand.
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryVault())
	cmd.AddCommand(CmdQueryListVault())
	cmd.AddCommand(CmdQueryTotalShares())
	cmd.AddCommand(CmdQueryListOwnerShares())
	cmd.AddCommand(CmdQueryMegavaultWithdrawalInfo())
	cmd.AddCommand(CmdQueryOwnerShares())

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-params",
		Short: "get x/vault parameters",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
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

func CmdQueryVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-vault [type] [number]",
		Short: "get a vault by its type and number",
		Long:  "get a vault by its type and number. Current support types are: clob.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// Parse vault type.
			vaultType, err := GetVaultTypeFromString(args[0])
			if err != nil {
				return err
			}

			// Parse vault number.
			vaultNumber, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			res, err := queryClient.Vault(
				context.Background(),
				&types.QueryVaultRequest{
					Type:   vaultType,
					Number: uint32(vaultNumber),
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

func CmdQueryListVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-vault",
		Short: "list all vaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryAllVaultsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllVaults(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryTotalShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-shares",
		Short: "get total shares of megavault",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.MegavaultTotalShares(
				context.Background(),
				&types.QueryMegavaultTotalSharesRequest{},
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

func CmdQueryListOwnerShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-owner-shares",
		Short: "list owner shares",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryMegavaultAllOwnerSharesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.MegavaultAllOwnerShares(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryMegavaultWithdrawalInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "megavault-withdrawal-info [shares_to_withdraw]",
		Short: "get megavault withdrawal info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			shares, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.MegavaultWithdrawalInfo(
				context.Background(),
				&types.QueryMegavaultWithdrawalInfoRequest{
					SharesToWithdraw: types.NumShares{
						NumShares: dtypes.NewIntFromUint64(shares),
					},
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

func CmdQueryOwnerShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-shares [address]",
		Short: "get owner shares by their address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.MegavaultOwnerShares(
				context.Background(),
				&types.QueryMegavaultOwnerSharesRequest{
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
