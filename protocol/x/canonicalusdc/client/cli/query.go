package cli

import (
	"fmt"
	"strconv"

	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(queryStateCmd(), queryPendingSettlementsCmd(), queryPendingSettlementCmd())
	return cmd
}

func queryStateCmd() *cobra.Command {
	return queryCmd("state", "query controls and backing ledger", cobra.NoArgs, func(
		clientCtx client.Context,
		queryClient types.QueryClient,
		cmd *cobra.Command,
		_ []string,
	) error {
		response, err := queryClient.State(cmd.Context(), &types.QueryStateRequest{})
		if err != nil {
			return err
		}
		return clientCtx.PrintProto(response)
	})
}

func queryPendingSettlementsCmd() *cobra.Command {
	cmd := queryCmd("pending-settlements", "query pending settlements", cobra.NoArgs, func(
		clientCtx client.Context,
		queryClient types.QueryClient,
		cmd *cobra.Command,
		_ []string,
	) error {
		pageRequest, err := client.ReadPageRequest(cmd.Flags())
		if err != nil {
			return err
		}
		response, err := queryClient.PendingSettlements(
			cmd.Context(),
			&types.QueryPendingSettlementsRequest{Pagination: pageRequest},
		)
		if err != nil {
			return err
		}
		return clientCtx.PrintProto(response)
	})
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	return cmd
}

func queryPendingSettlementCmd() *cobra.Command {
	return queryCmd(
		"pending-settlement [source-channel] [sequence]",
		"query one pending settlement by packet source channel and sequence",
		cobra.ExactArgs(2),
		func(clientCtx client.Context, queryClient types.QueryClient, cmd *cobra.Command, args []string) error {
			sequence, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			response, err := queryClient.PendingSettlement(
				cmd.Context(),
				&types.QueryPendingSettlementRequest{SourceChannel: args[0], Sequence: sequence},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(response)
		},
	)
}

func queryCmd(
	use,
	short string,
	args cobra.PositionalArgs,
	run func(client.Context, types.QueryClient, *cobra.Command, []string) error,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  args,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			return run(clientCtx, types.NewQueryClient(clientCtx), cmd, args)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
