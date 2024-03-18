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

func blockRateLimitConfigUpdate(
	ctx sdk.Context,
	clobKeeper clobtypes.ClobKeeper,
) {
	// Based off of https://docs.dydx.exchange/trading/rate_limits
	blockRateLimitConfig := clobtypes.BlockRateLimitConfiguration{
		// Kept the same
		MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 2,
				Limit:     1,
			},
			{
				NumBlocks: 20,
				Limit:     100,
			},
		},
		MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 5,
				Limit:     2000,
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
