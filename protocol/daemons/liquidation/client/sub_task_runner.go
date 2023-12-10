package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	perpkeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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

	lastCommittedBlockHeight, err := daemonClient.GetPreviousBlockInfo(ctx)
	if err != nil {
		return err
	}

	// 1. Fetch all information needed to calculate total net collateral and margin requirements.
	subaccounts,
		marketPrices,
		perpetuals,
		liquidityTiers,
		err := daemonClient.FetchApplicationStateAtBlockHeight(
		ctx,
		lastCommittedBlockHeight,
		liqFlags,
	)
	if err != nil {
		return err
	}

	// 2. Check collateralization statuses of subaccounts with at least one open position.
	liquidatableSubaccountIds, err := daemonClient.GetLiquidatableSubaccountIds(
		subaccounts,
		marketPrices,
		perpetuals,
		liquidityTiers,
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
	marketPricesMap map[uint32]pricestypes.MarketPrice,
	perpetualsMap map[uint32]perptypes.Perpetual,
	liquidityTiersMap map[uint32]perptypes.LiquidityTier,
	err error,
) {
	// Execute all queries at the given block height.
	queryCtx := newContextWithQueryBlockHeight(ctx, blockHeight)

	// Subaccounts
	subaccounts, err = c.GetAllSubaccounts(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Market prices
	marketPrices, err := c.GetAllMarketPrices(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	marketPricesMap = lib.UniqueSliceToMap(marketPrices, func(m pricestypes.MarketPrice) uint32 {
		return m.Id
	})

	// Perpetuals
	perpetuals, err := c.GetAllPerpetuals(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	perpetualsMap = lib.UniqueSliceToMap(perpetuals, func(p perptypes.Perpetual) uint32 {
		return p.Params.Id
	})

	// Liquidity tiers
	liquidityTiers, err := c.GetAllLiquidityTiers(queryCtx, liqFlags.QueryPageLimit)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	liquidityTiersMap = lib.UniqueSliceToMap(liquidityTiers, func(l perptypes.LiquidityTier) uint32 {
		return l.Id
	})

	return subaccounts, marketPricesMap, perpetualsMap, liquidityTiersMap, nil
}

// GetLiquidatableSubaccountIds verifies collateralization statuses of subaccounts with
// at least one open position and returns a list of unique and potentially liquidatable subaccount ids.
func (c *Client) GetLiquidatableSubaccountIds(
	subaccounts []satypes.Subaccount,
	marketPrices map[uint32]pricestypes.MarketPrice,
	perpetuals map[uint32]perptypes.Perpetual,
	liquidityTiers map[uint32]perptypes.LiquidityTier,
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

	numSubaccountsWithOpenPositions := 0
	liquidatableSubaccountIds = make([]satypes.SubaccountId, 0)
	for _, subaccount := range subaccounts {
		// Skip subaccounts with no open positions.
		if len(subaccount.PerpetualPositions) == 0 {
			numSubaccountsWithOpenPositions++
			continue
		}

		// Check if the subaccount is liquidatable.
		isLiquidatable, err := c.CheckSubaccountCollateralization(
			subaccount,
			marketPrices,
			perpetuals,
			liquidityTiers,
		)
		if err != nil {
			return nil, err
		}

		if isLiquidatable {
			liquidatableSubaccountIds = append(liquidatableSubaccountIds, *subaccount.Id)
		}
	}

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(numSubaccountsWithOpenPositions),
		metrics.SubaccountsWithOpenPositions,
		metrics.Count,
	)

	return liquidatableSubaccountIds, nil
}

// CheckSubaccountCollateralization queries a gRPC server using `AreSubaccountsLiquidatable`
// and returns a list of collateralization statuses for the given list of subaccount ids.
func (c *Client) CheckSubaccountCollateralization(
	subaccount satypes.Subaccount,
	marketPrices map[uint32]pricestypes.MarketPrice,
	perpetuals map[uint32]perptypes.Perpetual,
	liquidityTiers map[uint32]perptypes.LiquidityTier,
) (
	isLiquidatable bool,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.CheckCollateralizationForSubaccounts,
		metrics.Latency,
	)

	bigTotalNetCollateral := big.NewInt(0)
	bigTotalMaintenanceMargin := big.NewInt(0)

	// Calculate the net collateral and maintenance margin for each of the asset positions.
	// Note that we only expect USDC before multi-collateral support is added.
	for _, assetPosition := range subaccount.AssetPositions {
		if assetPosition.AssetId != assetstypes.AssetUsdc.Id {
			panic("liquidation daemon only supports USDC collateral")
		}
		// Net collateral for USDC is the quantums of the position.
		// Margin requirements for USDC are zero.
		bigTotalNetCollateral.Add(bigTotalNetCollateral, assetPosition.GetBigQuantums())
	}

	// Calculate the net collateral and maintenance margin for each of the perpetual positions.
	for _, perpetualPosition := range subaccount.PerpetualPositions {
		perpetual, ok := perpetuals[perpetualPosition.PerpetualId]
		if !ok {
			panic(
				fmt.Sprintf(
					"Perpetual not found for perpetual id %d",
					perpetualPosition.PerpetualId,
				),
			)
		}

		marketPrice, ok := marketPrices[perpetual.Params.MarketId]
		if !ok {
			panic(
				fmt.Sprintf(
					"MarketPrice not found for perpetual %+v",
					perpetual,
				),
			)
		}

		bigQuantums := perpetualPosition.GetBigQuantums()

		// Get the net collateral for the position.
		bigNetCollateral := perpkeeper.GetNetNotional(perpetual, marketPrice, bigQuantums)
		bigTotalNetCollateral.Add(bigTotalNetCollateral, bigNetCollateral)

		liquidityTier, ok := liquidityTiers[perpetual.Params.LiquidityTier]
		if !ok {
			panic(
				fmt.Sprintf(
					"LiquidityTier not found for perpetual %+v",
					perpetual,
				),
			)
		}

		// Get the maintenance margin requirement for the position.
		_, bigMaintenanceMargin := perpkeeper.GetMarginRequirements(
			perpetual,
			marketPrice,
			liquidityTier,
			bigQuantums,
		)
		bigTotalMaintenanceMargin.Add(bigTotalMaintenanceMargin, bigMaintenanceMargin)
	}

	return clobkeeper.IsLiquidatable(bigTotalNetCollateral, bigTotalMaintenanceMargin), nil
}
