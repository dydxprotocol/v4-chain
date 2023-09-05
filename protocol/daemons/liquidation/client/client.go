package client

import (
	"context"
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Start begins a job that periodically:
// 1) Queries a gRPC server for all subaccounts including their open positions.
// 2) Checks collateralization statuses of subaccounts with at least one open position.
// 3) Sends a list of subaccount ids that potentially need to be liquidated to the application.
func Start(
	ctx context.Context,
	flags flags.DaemonFlags,
	logger log.Logger,
	grpcClient lib.GrpcClient,
) error {
	// Make a connection to the Cosmos gRPC query services.
	queryConn, err := grpcClient.NewTcpConnection(ctx, flags.Shared.GrpcServerAddress)
	if err != nil {
		logger.Error("Failed to establish gRPC connection to Cosmos gRPC query services", "error", err)
		return err
	}
	defer func() {
		if connErr := grpcClient.CloseConnection(queryConn); connErr != nil {
			err = connErr
		}
	}()

	// Make a connection to the private daemon gRPC server.
	daemonConn, err := grpcClient.NewGrpcConnection(ctx, flags.Shared.SocketAddress)
	if err != nil {
		logger.Error("Failed to establish gRPC connection to socket address", "error", err)
		return err
	}
	defer func() {
		if connErr := grpcClient.CloseConnection(daemonConn); connErr != nil {
			err = connErr
		}
	}()

	subaccountQueryClient := satypes.NewQueryClient(queryConn)
	clobQueryClient := clobtypes.NewQueryClient(queryConn)
	liquidationServiceClient := api.NewLiquidationServiceClient(daemonConn)

	ticker := time.NewTicker(time.Duration(flags.Liquidation.LoopDelayMs) * time.Millisecond)
	for ; true; <-ticker.C {
		if err := RunLiquidationDaemonTaskLoop(
			ctx,
			flags.Liquidation,
			subaccountQueryClient,
			clobQueryClient,
			liquidationServiceClient,
		); err != nil {
			// TODO(DEC-947): Move daemon shutdown to application.
			logger.Error("Liquidations daemon returned error", "error", err)
		}
	}

	return nil
}

// RunLiquidationDaemonTaskLoop contains the logic to communicate with various gRPC services
// to find the liquidatable subaccount ids.
func RunLiquidationDaemonTaskLoop(
	ctx context.Context,
	liqFlags flags.LiquidationFlags,
	subaccountQueryClient satypes.QueryClient,
	clobQueryClient clobtypes.QueryClient,
	liquidationServiceClient api.LiquidationServiceClient,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.MainTaskLoop,
		metrics.Latency,
	)

	// Fetch all subaccounts from query service.
	subaccounts, err := GetAllSubaccounts(
		ctx,
		subaccountQueryClient,
		liqFlags.SubaccountPageLimit,
	)
	if err != nil {
		return err
	}
	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(subaccounts)),
		metrics.AllSubaccounts,
		metrics.Count,
	)

	// Filter out subaccounts with no open positions.
	subaccountsWithOpenPositions := make([]satypes.SubaccountId, 0)
	for _, subaccount := range subaccounts {
		if len(subaccount.PerpetualPositions) > 0 {
			subaccountsWithOpenPositions = append(subaccountsWithOpenPositions, *subaccount.Id)
		}
	}
	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(subaccountsWithOpenPositions)),
		metrics.SubaccountsWithOpenPositions,
		metrics.Count,
	)
	liquidatableSubaccountIds := make([]satypes.SubaccountId, 0)
	if len(subaccountsWithOpenPositions) > 0 {
		// Check collateralization statuses of subaccounts with at least one open position.
		collateralizationCheckResults, err :=
			CheckCollateralizationForSubaccounts(
				ctx,
				clobQueryClient,
				subaccountsWithOpenPositions,
			)
		if err != nil {
			return err
		}

		// Append all liquidatable subaccount ids to a new slice.
		for _, result := range collateralizationCheckResults {
			if result.IsLiquidatable {
				liquidatableSubaccountIds = append(liquidatableSubaccountIds, result.SubaccountId)
			}
		}
	}

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(liquidatableSubaccountIds)),
		metrics.LiquidatableSubaccountIds,
		metrics.Count,
	)

	// Send the list of liquidatable subaccount ids to the daemon server.
	err = SendLiquidatableSubaccountIds(
		ctx,
		liquidationServiceClient,
		liquidatableSubaccountIds,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAllSubaccounts queries a gRPC server and returns a list of subaccounts and
// their balances and open positions.
func GetAllSubaccounts(
	ctx context.Context,
	client satypes.QueryClient,
	limit uint64,
) (
	subaccounts []satypes.Subaccount,
	err error,
) {
	subaccounts = make([]satypes.Subaccount, 0)

	var nextKey []byte
	for {
		subaccountsFromKey, next, err := getSubaccountsFromKey(
			ctx,
			client,
			limit,
			nextKey,
		)
		if err != nil {
			return nil, err
		}

		subaccounts = append(subaccounts, subaccountsFromKey...)
		nextKey = next

		if len(nextKey) == 0 {
			break
		}
	}
	return subaccounts, nil
}

// CheckCollateralizationForSubaccounts queries a gRPC server using `AreSubaccountsLiquidatable`
// and returns a list of collateralization statuses for the given list of subaccount ids.
func CheckCollateralizationForSubaccounts(
	ctx context.Context,
	client clobtypes.QueryClient,
	subaccountIds []satypes.SubaccountId,
) (
	results []clobtypes.AreSubaccountsLiquidatableResponse_Result,
	err error,
) {
	query := &clobtypes.AreSubaccountsLiquidatableRequest{
		SubaccountIds: subaccountIds,
	}
	response, err := client.AreSubaccountsLiquidatable(ctx, query)
	if err != nil {
		return nil, err
	}
	return response.Results, nil
}

// SendLiquidatableSubaccountIds sends a list of unique and potentially liquidatable
// subaccount ids to a gRPC server via `LiquidateSubaccounts`.
func SendLiquidatableSubaccountIds(
	ctx context.Context,
	client api.LiquidationServiceClient,
	subaccountIds []satypes.SubaccountId,
) error {
	request := &api.LiquidateSubaccountsRequest{
		SubaccountIds: subaccountIds,
	}

	if _, err := client.LiquidateSubaccounts(ctx, request); err != nil {
		return err
	}
	return nil
}

func getSubaccountsFromKey(
	ctx context.Context,
	client satypes.QueryClient,
	limit uint64,
	pageRequestKey []byte,
) (
	subaccounts []satypes.Subaccount,
	nextKey []byte,
	err error,
) {
	query := &satypes.QueryAllSubaccountRequest{
		Pagination: &query.PageRequest{
			Key:   pageRequestKey,
			Limit: limit,
		},
	}

	response, err := client.SubaccountAll(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	if response.Pagination != nil {
		nextKey = response.Pagination.NextKey
	}
	return response.Subaccount, nextKey, nil
}
