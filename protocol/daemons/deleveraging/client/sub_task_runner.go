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

// SubTaskRunner provides an interface that encapsulates the deleveraging daemon logic to gather subaccounts with open positions for each perp.
// This interface is used to mock the daemon logic in tests.
type SubTaskRunner interface {
	RunDeleveragingDaemonTaskLoop(
		ctx context.Context,
		client *Client,
		deleveragingFlags flags.DeleveragingFlags,
	) error
}

type SubTaskRunnerImpl struct{}

// Ensure SubTaskRunnerImpl implements the SubTaskRunner interface.
var _ SubTaskRunner = (*SubTaskRunnerImpl)(nil)

// RunDeleveragingDaemonTaskLoop contains the logic to communicate with various gRPC services
// to generate the list of subaccounts with open positions for each perpetual.
func (s *SubTaskRunnerImpl) RunDeleveragingDaemonTaskLoop(
	ctx context.Context,
	daemonClient *Client,
	deleveragingFlags flags.DeleveragingFlags,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.DeleveragingDaemon,
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
		deleveragingFlags,
	)
	if err != nil {
		return err
	}

	// Build a map of perpetual id to subaccounts with open positions in that perpetual.
	subaccountOpenPositionInfo := daemonClient.GetSubaccountOpenPositionInfo(subaccounts)

	// 3. Send the list of deleveraging subaccount ids to the daemon server.
	err = daemonClient.SendDeleveragingSubaccountIds(
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
	deleveragingFlags flags.DeleveragingFlags,
) (
	subaccounts []satypes.Subaccount,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.DeleveragingDaemon,
		time.Now(),
		metrics.FetchSubaccountsAtBlockHeight,
		metrics.Latency,
	)

	// Execute all queries at the given block height.
	queryCtx := newContextWithQueryBlockHeight(ctx, blockHeight)

	// Subaccounts
	subaccounts, err = c.GetAllSubaccounts(queryCtx, deleveragingFlags.QueryPageLimit)
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
		metrics.DeleveragingDaemon,
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
		metrics.DeleveragingDaemon,
		float32(numSubaccountsWithOpenPositions),
		metrics.SubaccountsWithOpenPositions,
		metrics.Count,
	)

	return subaccountOpenPositionInfo
}
