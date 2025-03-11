package v8_1

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	perpetualskeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
)

func createNewLiquidityTier(
	ctx sdk.Context,
	perpetualsKeeper perpetualskeeper.Keeper,
) {
	_, err := perpetualsKeeper.SetLiquidityTier(
		ctx,
		NewTierId,
		NewTierName,
		InitialMarginPpm,
		MaintenanceFractionPpm,
		ImpactNotional,
		OpenInterestLowerCap,
		OpenInterestUpperCap,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create new liquidity tier: %s", err))
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	perpetualsKeeper perpetualskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		createNewLiquidityTier(sdkCtx, perpetualsKeeper)

		sdkCtx.Logger().Info("Successfully created new liquidity tier for instant market listing")

		return mm.RunMigrations(ctx, configurator, vm)
	}
}