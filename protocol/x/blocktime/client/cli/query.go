package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryDowntimeParams())
	cmd.AddCommand(CmdQueryAllDowntimeInfo())
	cmd.AddCommand(CmdQueryPreviousBlockInfo())
	cmd.AddCommand(CmdQuerySynchronyParams())

	return cmd
}

func CmdQueryDowntimeParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-downtime-params",
		Short: "get the DowntimeParams",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.DowntimeParams(
				context.Background(),
				&types.QueryDowntimeParamsRequest{},
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

func CmdQueryAllDowntimeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all-downtime-info",
		Short: "get all downtime info",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AllDowntimeInfo(
				context.Background(),
				&types.QueryAllDowntimeInfoRequest{},
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

func CmdQueryPreviousBlockInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-previous-block-info",
		Short: "get previous block info",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PreviousBlockInfo(
				context.Background(),
				&types.QueryPreviousBlockInfoRequest{},
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

func CmdQuerySynchronyParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synchrony-params",
		Short: "get synchrony params",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SynchronyParams(
				context.Background(),
				&types.QuerySynchronyParamsRequest{},
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
