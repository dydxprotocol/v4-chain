package client

import (
	"context"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
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

	lastCommittedBlockHeight, err := daemonClient.GetPreviousBlockInfo(ctx)
	if err != nil {
		return err
	}

	// 1. Fetch all information needed to calculate total net collateral and margin requirements.
	subaccounts, err := daemonClient.FetchSubaccountsAtBlockHeight(
		ctx,
		lastCommittedBlockHeight,
		liqFlags,
	)
	if err != nil {
		return err
	}

	// Build a map of perpetual id to subaccounts with open positions in that perpetual.
	subaccountOpenPositionInfo := daemonClient.GetSubaccountOpenPositionInfo(subaccounts)

	// 3. Send the list of liquidatable subaccount ids to the daemon server.
	err = daemonClient.SendLiquidatableSubaccountIds(
		ctx,
		subaccountOpenPositionInfo,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) FetchSubaccountsAtBlockHeight(
	ctx context.Context,
	blockHeight uint32,
	liqFlags flags.LiquidationFlags,
) (
	subaccounts []satypes.Subaccount,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.FetchSubaccountsAtBlockHeight,
		metrics.Latency,
	)

	// Execute all queries at the given block height.
	queryCtx := newContextWithQueryBlockHeight(ctx, blockHeight)

	// Subaccounts
	subaccounts, err = c.GetAllSubaccounts(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, err
	}

	return subaccounts, nil
}

// GetSubaccountOpenPositionInfo iterates over the given subaccounts and returns a map of
// perpetual id to open position info.
func (c *Client) GetSubaccountOpenPositionInfo(
	subaccounts []satypes.Subaccount,
) (
	subaccountOpenPositionInfo map[uint32]*clobtypes.SubaccountOpenPositionInfo,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.GetSubaccountOpenPositionInfo,
		metrics.Latency,
	)

	numSubaccountsWithOpenPositions := 0
	subaccountOpenPositionInfo = make(map[uint32]*clobtypes.SubaccountOpenPositionInfo)
	for _, subaccount := range subaccounts {
		// Skip subaccounts with no open positions.
		if len(subaccount.PerpetualPositions) == 0 {
			continue
		}

		for _, perpetualPosition := range subaccount.PerpetualPositions {
			openPositionInfo, ok := subaccountOpenPositionInfo[perpetualPosition.PerpetualId]
			if !ok {
				openPositionInfo = &clobtypes.SubaccountOpenPositionInfo{
					PerpetualId:                  perpetualPosition.PerpetualId,
					SubaccountsWithLongPosition:  make([]satypes.SubaccountId, 0),
					SubaccountsWithShortPosition: make([]satypes.SubaccountId, 0),
				}
				subaccountOpenPositionInfo[perpetualPosition.PerpetualId] = openPositionInfo
			}

			if perpetualPosition.GetIsLong() {
				openPositionInfo.SubaccountsWithLongPosition = append(
					openPositionInfo.SubaccountsWithLongPosition,
					*subaccount.Id,
				)
			} else {
				openPositionInfo.SubaccountsWithShortPosition = append(
					openPositionInfo.SubaccountsWithShortPosition,
					*subaccount.Id,
				)
			}
		}

		numSubaccountsWithOpenPositions++
	}

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(numSubaccountsWithOpenPositions),
		metrics.SubaccountsWithOpenPositions,
		metrics.Count,
	)

	return subaccountOpenPositionInfo
}
