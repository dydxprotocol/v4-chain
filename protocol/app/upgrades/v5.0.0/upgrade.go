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
	ctx.Logger().Info(
		fmt.Sprintf(
			"Combining the short term order placement and cancellation limits of previous config: %+v\n",
			oldBlockRateLimitConfig,
		),
	)
	numAllowedShortTermOrderPlacementsInOneBlock := 0
	numAllowedShortTermOrderCancellationsInOneBlock := 0
	oldShortTermOrderRateLimits := oldBlockRateLimitConfig.MaxShortTermOrdersPerNBlocks
	for _, limit := range oldShortTermOrderRateLimits {
		if limit.NumBlocks == 1 {
			numAllowedShortTermOrderPlacementsInOneBlock += int(limit.Limit)
			break
		}
	}
	if numAllowedShortTermOrderPlacementsInOneBlock == 0 {
		panic("Failed to find MaxShortTermOrdersPerNBlocks with window 1.")
	}

	oldShortTermOrderCancellationRateLimits := oldBlockRateLimitConfig.MaxShortTermOrderCancellationsPerNBlocks
	for _, limit := range oldShortTermOrderCancellationRateLimits {
		if limit.NumBlocks == 1 {
			numAllowedShortTermOrderCancellationsInOneBlock += int(limit.Limit)
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
	ctx.Logger().Info(
		fmt.Sprintf(
			"Attempting to set rate limiting config to newly combined config: %+v\n",
			blockRateLimitConfig,
		),
	)
	if err := clobKeeper.InitializeBlockRateLimit(ctx, blockRateLimitConfig); err != nil {
		panic(fmt.Sprintf("failed to update the block rate limit configuration: %s", err))
	}
	ctx.Logger().Info(
		"Successfully upgraded block rate limit configuration to: %+v\n",
		clobKeeper.GetBlockRateLimitConfiguration(ctx),
	)
}

// Initialize soft and upper caps for OIMF
func initializeOIMFCaps(
	ctx sdk.Context,
	perpetualsKeeper perptypes.PerpetualsKeeper,
) {
	allLiquidityTiers := perpetualsKeeper.GetAllLiquidityTiers(ctx)
	for _, tier := range allLiquidityTiers {
		if tier.Id == 0 {
			// For large cap, no OIMF caps
			tier.OpenInterestLowerCap = 0
			tier.OpenInterestUpperCap = 0
		} else if tier.Id == 1 {
			// Mid cap
			tier.OpenInterestLowerCap = 25_000_000_000_000 // 25 million USDC
			tier.OpenInterestUpperCap = 50_000_000_000_000 // 50 million USDC
		} else if tier.Id == 2 {
			// Long tail
			tier.OpenInterestLowerCap = 10_000_000_000_000 // 10 million USDC
			tier.OpenInterestUpperCap = 20_000_000_000_000 // 20 million USDC
		} else {
			// Safety
			tier.OpenInterestLowerCap = 500_000_000_000   // 0.5 million USDC
			tier.OpenInterestUpperCap = 1_000_000_000_000 // 1 million USDC
		}

		lt, err := perpetualsKeeper.SetLiquidityTier(
			ctx,
			tier.Id,
			tier.Name,
			tier.InitialMarginPpm,
			tier.MaintenanceFractionPpm,
			tier.ImpactNotional,
			tier.OpenInterestLowerCap,
			tier.OpenInterestUpperCap,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to set liquidity tier: %+v,\n err: %s", tier.Id, err))
		}
		// TODO(OTE-248): Optional - emit indexer events that for updated liquidity tier
		ctx.Logger().Info(
			fmt.Sprintf(
				"Successfully set liqiquidity tier with `OpenInterestLower/UpperCap`: %+v\n",
				lt,
			),
		)
		// TODO(OTE-249): Add upgrade test that checks if the OIMF caps are set correctly
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

		// Migrate pruneable orders to new format
		clobKeeper.MigratePruneableOrders(sdkCtx)
		sdkCtx.Logger().Info("Successfully migrated pruneable orders")

		// Set all perpetuals to cross market type
		perpetualsUpgrade(sdkCtx, perpetualsKeeper)

		// Set block rate limit configuration
		blockRateLimitConfigUpdate(sdkCtx, clobKeeper)

		// Initialize liquidity tier with lower and upper OI caps.
		initializeOIMFCaps(sdkCtx, perpetualsKeeper)

		// TODO(TRA-93): Initialize `x/vault` module.

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
