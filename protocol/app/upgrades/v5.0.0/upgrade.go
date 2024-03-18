package v_5_0_0

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// Set all existing perpetuals to cross market type
func perpetualsUpgrade(
	ctx sdk.Context,
	perpetualsKeeper perptypes.PerpetualsKeeper,
) {
	// Set all perpetuals to cross market type
	perpetuals := perpetualsKeeper.GetAllPerpetuals(ctx)
	for _, p := range perpetuals {
		_, err := perpetualsKeeper.SetPerpetualMarketType(
			ctx, p.GetId(),
			perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS)
		if err != nil {
			panic(fmt.Sprintf("failed to set perpetual market type for perpetual %d: %s", p.GetId(), err))
		}
	}
}

// blockRateLimitConfigUpdate upgrades the block rate limit. It searches for the
// 1-block window limit for short term and cancellations, sums them, and creates a new
// combined rate limit.
func blockRateLimitConfigUpdate(
	ctx sdk.Context,
	clobKeeper clobtypes.ClobKeeper,
) {
	oldBlockRateLimitConfig := clobKeeper.GetBlockRateLimitConfiguration(ctx)
	numAllowedShortTermOrderPlacementsInOneBlock := 0
	numAllowedShortTermOrderCancellationsInOneBlock := 0
	oldShortTermOrderRateLimits := oldBlockRateLimitConfig.MaxShortTermOrdersPerNBlocks
	for _, limit := range oldShortTermOrderRateLimits {
		if limit.NumBlocks == 1 {
			numAllowedShortTermOrderPlacementsInOneBlock += int(limit.NumBlocks)
			break
		}
	}
	if numAllowedShortTermOrderPlacementsInOneBlock == 0 {
		panic("Failed to find MaxShortTermOrdersPerNBlocks with window 1.")
	}

	oldShortTermOrderCancellationRateLimits := oldBlockRateLimitConfig.MaxShortTermOrderCancellationsPerNBlocks
	for _, limit := range oldShortTermOrderCancellationRateLimits {
		if limit.NumBlocks == 1 {
			numAllowedShortTermOrderCancellationsInOneBlock += int(limit.NumBlocks)
			break
		}
	}
	if numAllowedShortTermOrderCancellationsInOneBlock == 0 {
		panic("Failed to find MaxShortTermOrdersPerNBlocks with window 1.")
	}

	allowedNumShortTermPlaceAndCancelInFiveBlocks :=
		(numAllowedShortTermOrderPlacementsInOneBlock + numAllowedShortTermOrderCancellationsInOneBlock) * 5

	blockRateLimitConfig := clobtypes.BlockRateLimitConfiguration{
		// Kept the same
		MaxStatefulOrdersPerNBlocks: oldBlockRateLimitConfig.MaxStatefulOrdersPerNBlocks,
		// Combine place and cancel, gate over 5 blocks to allow burst
		MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 5,
				Limit:     uint32(allowedNumShortTermPlaceAndCancelInFiveBlocks),
			},
		},
	}

	if err := clobKeeper.InitializeBlockRateLimit(ctx, blockRateLimitConfig); err != nil {
		panic(fmt.Sprintf("failed to update the block rate limit configuration: %s", err))
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	clobKeeper clobtypes.ClobKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Set all perpetuals to cross market type
		perpetualsUpgrade(sdkCtx, perpetualsKeeper)

		// Set block rate limit configuration
		blockRateLimitConfigUpdate(sdkCtx, clobKeeper)

		// TODO(TRA-93): Initialize `x/vault` module.

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
