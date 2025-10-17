package client

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	salib "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/lib"
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

type SubTaskRunnerImpl struct {
	lastLoopBlockHeight uint32
}

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

	// Skip the loop if no new block has been committed.
	// Note that lastLoopBlockHeight is initialized to 0, so the first loop will always run.
	if lastCommittedBlockHeight == s.lastLoopBlockHeight {
		daemonClient.logger.Info(
			"Skipping liquidation daemon task loop as no new block has been committed",
			"blockHeight", lastCommittedBlockHeight,
		)
		return nil
	}

	// Update the last loop block height.
	s.lastLoopBlockHeight = lastCommittedBlockHeight

	// 1. Fetch all information needed to calculate total net collateral and margin requirements.
	subaccounts,
		perpInfos,
		err := daemonClient.FetchApplicationStateAtBlockHeight(
		ctx,
		lastCommittedBlockHeight,
		liqFlags,
	)
	if err != nil {
		return err
	}

	// 2. Check collateralization statuses of subaccounts with at least one open position.
	liquidatableSubaccountIds,
		negativeTncSubaccountIds,
		err := daemonClient.GetLiquidatableSubaccountIds(
		subaccounts,
		perpInfos,
	)
	if err != nil {
		return err
	}

	// Build a map of perpetual id to subaccounts with open positions in that perpetual.
	subaccountOpenPositionInfo := daemonClient.GetSubaccountOpenPositionInfo(subaccounts)

	// 3. Send the list of liquidatable subaccount ids to the daemon server.
	err = daemonClient.SendLiquidatableSubaccountIds(
		ctx,
		lastCommittedBlockHeight,
		liquidatableSubaccountIds,
		negativeTncSubaccountIds,
		subaccountOpenPositionInfo,
		liqFlags.ResponsePageLimit,
	)
	if err != nil {
		return err
	}

	return nil
}

// FetchApplicationStateAtBlockHeight queries a gRPC server and fetches the following information given a block height:
// - Last committed block height.
// - Subaccounts including their open positions.
// - Market prices.
// - Perpetuals.
// - Liquidity tiers.
func (c *Client) FetchApplicationStateAtBlockHeight(
	ctx context.Context,
	blockHeight uint32,
	liqFlags flags.LiquidationFlags,
) (
	subaccounts []satypes.Subaccount,
	perpInfos perptypes.PerpInfos,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.FetchApplicationStateAtBlockHeight,
		metrics.Latency,
	)

	// Execute all queries at the given block height.
	queryCtx := newContextWithQueryBlockHeight(ctx, blockHeight)

	// Subaccounts
	subaccounts, err = c.GetAllSubaccounts(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, errorsmod.Wrapf(
			err,
			"failed to fetch subaccounts at block height %d",
			blockHeight,
		)
	}

	// Market prices
	marketPrices, err := c.GetAllMarketPrices(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, errorsmod.Wrapf(
			err,
			"failed to fetch market prices at block height %d",
			blockHeight,
		)
	}
	marketPricesMap := lib.UniqueSliceToMap(marketPrices, func(m pricestypes.MarketPrice) uint32 {
		return m.Id
	})

	// Perpetuals
	perpetuals, err := c.GetAllPerpetuals(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, errorsmod.Wrapf(
			err,
			"failed to fetch perpetuals at block height %d",
			blockHeight,
		)
	}

	// Liquidity tiers
	liquidityTiers, err := c.GetAllLiquidityTiers(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, errorsmod.Wrapf(
			err,
			"failed to fetch liquidity tiers at block height %d",
			blockHeight,
		)
	}
	liquidityTiersMap := lib.UniqueSliceToMap(liquidityTiers, func(l perptypes.LiquidityTier) uint32 {
		return l.Id
	})

	perpInfos = make(perptypes.PerpInfos, len(perpetuals))
	for _, perp := range perpetuals {
		price, ok := marketPricesMap[perp.Params.MarketId]
		if !ok {
			return nil, nil, errorsmod.Wrapf(
				pricestypes.ErrMarketPriceDoesNotExist,
				"%d",
				perp.Params.MarketId,
			)
		}
		liquidityTier, ok := liquidityTiersMap[perp.Params.LiquidityTier]
		if !ok {
			return nil, nil, errorsmod.Wrapf(
				perptypes.ErrLiquidityTierDoesNotExist,
				"%d",
				perp.Params.LiquidityTier,
			)
		}
		perpInfos[perp.Params.Id] = perptypes.PerpInfo{
			Perpetual:     perp,
			Price:         price,
			LiquidityTier: liquidityTier,
		}
	}

	return subaccounts, perpInfos, nil
}

// GetLiquidatableSubaccountIds verifies collateralization statuses of subaccounts with
// at least one open position and returns a list of unique and potentially liquidatable subaccount ids.
func (c *Client) GetLiquidatableSubaccountIds(
	subaccounts []satypes.Subaccount,
	perpInfos perptypes.PerpInfos,
) (
	liquidatableSubaccountIds []satypes.SubaccountId,
	negativeTncSubaccountIds []satypes.SubaccountId,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.GetLiquidatableSubaccountIds,
		metrics.Latency,
	)

	liquidatableSubaccountIds = make([]satypes.SubaccountId, 0)
	negativeTncSubaccountIds = make([]satypes.SubaccountId, 0)
	for _, subaccount := range subaccounts {
		// Skip subaccounts with no open positions.
		if len(subaccount.PerpetualPositions) == 0 {
			continue
		}

		// Check if the subaccount is liquidatable.
		isLiquidatable, hasNegativeTnc, err := c.CheckSubaccountCollateralization(
			subaccount,
			perpInfos,
		)
		if err != nil {
			c.logger.Error("Error checking collateralization status", "error", err)
			return nil, nil, err
		}

		if isLiquidatable {
			liquidatableSubaccountIds = append(liquidatableSubaccountIds, *subaccount.Id)
		}
		if hasNegativeTnc {
			negativeTncSubaccountIds = append(negativeTncSubaccountIds, *subaccount.Id)
		}
	}

	return liquidatableSubaccountIds, negativeTncSubaccountIds, nil
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

// CheckSubaccountCollateralization performs the same collateralization check as the application
// using the provided market prices, perpetuals, and liquidity tiers.
//
// Note that current implementation assumes that the only asset is USDC and multi-collateral support
// is not yet implemented.
func (c *Client) CheckSubaccountCollateralization(
	unsettledSubaccount satypes.Subaccount,
	perpInfos perptypes.PerpInfos,
) (
	isLiquidatable bool,
	hasNegativeTnc bool,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.CheckCollateralizationForSubaccounts,
		metrics.Latency,
	)

	// Funding payments are lazily settled, so get the settled subaccount
	// to ensure that the funding payments are included in the net collateral calculation.
	settledSubaccount, _ := salib.GetSettledSubaccountWithPerpetuals(
		unsettledSubaccount,
		perpInfos,
	)

	risk, err := salib.GetRiskForSubaccount(
		settledSubaccount,
		perpInfos,
		nil, // No leverage needed for liquidation calculations
	)

	return risk.IsLiquidatable(), risk.NC.Sign() < 0, nil
}
