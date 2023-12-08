package client

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// SubTaskRunner provides an interface that encapsulates the liquidations daemon logic to gather and report
// potentially liquidatable subaccount ids. This interface is used to mock the daemon logic in tests.
type SubTaskRunner interface {
	RunLiquidationDaemonTaskLoop(
		ctx context.Context,
		client *Client,
		liqFlags flags.LiquidationFlags,
	) error
}

type SubTaskRunnerImpl struct{}

// Ensure SubTaskRunnerImpl implements the SubTaskRunner interface.
var _ SubTaskRunner = (*SubTaskRunnerImpl)(nil)

// RunLiquidationDaemonTaskLoop contains the logic to communicate with various gRPC services
// to find the liquidatable subaccount ids.
func (s *SubTaskRunnerImpl) RunLiquidationDaemonTaskLoop(
	ctx context.Context,
	daemonClient *Client,
	liqFlags flags.LiquidationFlags,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.MainTaskLoop,
		metrics.Latency,
	)

	// 1. Fetch all subaccounts from query service.
	subaccounts, err := daemonClient.GetAllSubaccounts(ctx, liqFlags.SubaccountPageLimit)
	if err != nil {
		return err
	}

	// 2. Check collateralization statuses of subaccounts with at least one open position.
	liquidatableSubaccountIds, err := daemonClient.GetLiquidatableSubaccountIds(
		ctx,
		liqFlags,
		subaccounts,
	)
	if err != nil {
		return err
	}

	// 3. Send the list of liquidatable subaccount ids to the daemon server.
	err = daemonClient.SendLiquidatableSubaccountIds(ctx, liquidatableSubaccountIds)
	if err != nil {
		return err
	}

	return nil
}

// GetLiquidatableSubaccountIds verifies collateralization statuses of subaccounts with
// at least one open position and returns a list of unique and potentially liquidatable subaccount ids.
func (c *Client) GetLiquidatableSubaccountIds(
	ctx context.Context,
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

		results, err := c.CheckCollateralizationForSubaccounts(
			ctx,
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
