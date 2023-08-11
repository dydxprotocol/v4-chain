package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/dydxprotocol/v4/x/bridge/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group bridge queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryEventParams())
	cmd.AddCommand(CmdQueryProposeParams())
	cmd.AddCommand(CmdQuerySafetyParams())
	cmd.AddCommand(CmdQueryNextAcknowledgedEventId())

	return cmd
}

func CmdQueryEventParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-event-params",
		Short: "get the EventParams",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.EventParams(
				context.Background(),
				&types.QueryEventParamsRequest{},
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

func CmdQueryProposeParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-propose-params",
		Short: "get the ProposeParams",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ProposeParams(
				context.Background(),
				&types.QueryProposeParamsRequest{},
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

func CmdQuerySafetyParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-safety-params",
		Short: "get the SafetyParams",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SafetyParams(
				context.Background(),
				&types.QuerySafetyParamsRequest{},
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

func CmdQueryNextAcknowledgedEventId() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-next-acknowledged-event-id",
		Short: "get the NextAcknowledgedEventId",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.NextAcknowledgedEventId(
				context.Background(),
				&types.QueryNextAcknowledgedEventIdRequest{},
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
