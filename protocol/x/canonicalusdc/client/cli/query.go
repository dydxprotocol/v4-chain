package cli

import (
	"fmt"

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
	cmd.AddCommand(queryStateCmd(), queryPendingCmd())
	return cmd
}

func queryStateCmd() *cobra.Command {
	return queryCmd("state", "query controls and backing ledger", func(
		clientCtx client.Context,
		queryClient types.QueryClient,
		cmd *cobra.Command,
	) error {
		response, err := queryClient.State(cmd.Context(), &types.QueryStateRequest{})
		if err != nil {
			return err
		}
		return clientCtx.PrintProto(response)
	})
}

func queryPendingCmd() *cobra.Command {
	return queryCmd("pending-settlements", "query pending settlements", func(
		clientCtx client.Context,
		queryClient types.QueryClient,
		cmd *cobra.Command,
	) error {
		response, err := queryClient.PendingSettlements(cmd.Context(), &types.QueryPendingSettlementsRequest{})
		if err != nil {
			return err
		}
		return clientCtx.PrintProto(response)
	})
}

func queryCmd(
	use,
	short string,
	run func(client.Context, types.QueryClient, *cobra.Command) error,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			return run(clientCtx, types.NewQueryClient(clientCtx), cmd)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
