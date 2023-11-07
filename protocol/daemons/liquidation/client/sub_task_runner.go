package client

import (
	"context"
	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"time"
)

// SubTaskRunner provides an interface that encapsulates the liquidations daemon logic to gather and report
// potentially liquidatable subaccount ids. This interface is used to mock the daemon logic in tests.
type SubTaskRunner interface {
	RunLiquidationDaemonTaskLoop(
		client *Client,
		ctx context.Context,
		liqFlags flags.LiquidationFlags,
		subaccountQueryClient satypes.QueryClient,
		clobQueryClient clobtypes.QueryClient,
		liquidationServiceClient api.LiquidationServiceClient,
	) error
}

type SubTaskRunnerImpl struct{}

// Ensure SubTaskRunnerImpl implements the SubTaskRunner interface.
var _ SubTaskRunner = (*SubTaskRunnerImpl)(nil)

// RunLiquidationDaemonTaskLoop contains the logic to communicate with various gRPC services
// to find the liquidatable subaccount ids.
func (s *SubTaskRunnerImpl) RunLiquidationDaemonTaskLoop(
	client *Client,
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

	// 1. Fetch all subaccounts from query service.
	subaccounts, err := GetAllSubaccounts(
		client,
		ctx,
		subaccountQueryClient,
		liqFlags.SubaccountPageLimit,
	)
	if err != nil {
		return err
	}

	// 2. Check collateralization statuses of subaccounts with at least one open position.
	liquidatableSubaccountIds, err := GetLiquidatableSubaccountIds(
		client,
		ctx,
		clobQueryClient,
		liqFlags,
		subaccounts,
	)
	if err != nil {
		return err
	}

	// 3. Send the list of liquidatable subaccount ids to the daemon server.
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

// CheckCollateralizationForSubaccounts queries a gRPC server using `AreSubaccountsLiquidatable`
// and returns a list of collateralization statuses for the given list of subaccount ids.
func CheckCollateralizationForSubaccounts(
	daemon *Client,
	ctx context.Context,
	client clobtypes.QueryClient,
	subaccountIds []satypes.SubaccountId,
) (
	results []clobtypes.AreSubaccountsLiquidatableResponse_Result,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.CheckCollateralizationForSubaccounts,
		metrics.Latency,
	)

	query := &clobtypes.AreSubaccountsLiquidatableRequest{
		SubaccountIds: subaccountIds,
	}
	response, err := client.AreSubaccountsLiquidatable(ctx, query)
	if err != nil {
		return nil, err
	}

	// For the purposes of the health check, log the successful request as an indicator of daemon health.
	daemon.ReportSuccess()

	return response.Results, nil
}

// SendLiquidatableSubaccountIds sends a list of unique and potentially liquidatable
// subaccount ids to a gRPC server via `LiquidateSubaccounts`.
func SendLiquidatableSubaccountIds(
	ctx context.Context,
	client api.LiquidationServiceClient,
	subaccountIds []satypes.SubaccountId,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.SendLiquidatableSubaccountIds,
		metrics.Latency,
	)

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(subaccountIds)),
		metrics.LiquidatableSubaccountIds,
		metrics.Count,
	)

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
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.LiquidationDaemon,
		[]string{metrics.GetSubaccountsFromKey, metrics.Latency},
		time.Now(),
		[]gometrics.Label{
			metrics.GetLabelForIntValue(metrics.PageLimit, int(limit)),
		},
	)

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

// GetAllSubaccounts queries a gRPC server and returns a list of subaccounts and
// their balances and open positions.
func GetAllSubaccounts(
	daemon *Client,
	ctx context.Context,
	client satypes.QueryClient,
	limit uint64,
) (
	subaccounts []satypes.Subaccount,
	err error,
) {
	defer telemetry.ModuleMeasureSince(metrics.LiquidationDaemon, time.Now(), metrics.GetAllSubaccounts, metrics.Latency)
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

		// For the purposes of the health check, log the successful request as an indicator of daemon health.
		daemon.ReportSuccess()

		subaccounts = append(subaccounts, subaccountsFromKey...)
		nextKey = next

		if len(nextKey) == 0 {
			break
		}
	}

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(subaccounts)),
		metrics.GetAllSubaccounts,
		metrics.Count,
	)

	return subaccounts, nil
}

// GetLiquidatableSubaccountIds verifies collateralization statuses of subaccounts with
// at least one open position and returns a list of unique and potentially liquidatable subaccount ids.
func GetLiquidatableSubaccountIds(
	daemon *Client,
	ctx context.Context,
	client clobtypes.QueryClient,
	liqFlags flags.LiquidationFlags,
	subaccounts []satypes.Subaccount,
) (
	liquidatableSubaccountIds []satypes.SubaccountId,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.GetLiquidatableSubaccountIds,
		metrics.Latency,
	)

	// Filter out subaccounts with no open positions.
	subaccountsToCheck := make([]satypes.SubaccountId, 0)
	for _, subaccount := range subaccounts {
		if len(subaccount.PerpetualPositions) > 0 {
			subaccountsToCheck = append(subaccountsToCheck, *subaccount.Id)
		}
	}

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(subaccountsToCheck)),
		metrics.SubaccountsWithOpenPositions,
		metrics.Count,
	)

	// Query the gRPC server in chunks of size `liqFlags.RequestChunkSize`.
	liquidatableSubaccountIds = make([]satypes.SubaccountId, 0)
	for start := 0; start < len(subaccountsToCheck); start += int(liqFlags.RequestChunkSize) {
		end := lib.Min(start+int(liqFlags.RequestChunkSize), len(subaccountsToCheck))

		results, err := CheckCollateralizationForSubaccounts(
			daemon,
			ctx,
			client,
			subaccountsToCheck[start:end],
		)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			if result.IsLiquidatable {
				liquidatableSubaccountIds = append(liquidatableSubaccountIds, result.SubaccountId)
			}
		}
	}
	return liquidatableSubaccountIds, nil
}
