package v_9_4

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	affiliatekeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

func setDefaultAffiliateTiersForSlidingAffiliates(ctx sdk.Context, affiliateKeeper affiliatekeeper.Keeper) {
	err := affiliateKeeper.UpdateAffiliateTiers(ctx, UpdatedAffiliateTiers)

	if err != nil {
		panic(fmt.Sprintf("failed to set default affiliate tiers: %s", err))
	}
}

func setDefaultAffiliateParameters(ctx sdk.Context, affiliateKeeper affiliatekeeper.Keeper) {
	err := affiliateKeeper.UpdateAffiliateParameters(ctx, &affiliatetypes.MsgUpdateAffiliateParameters{
		AffiliateParameters: UpdatedAffiliateParameters,
	})

	if err != nil {
		panic(fmt.Sprintf("failed to set default affiliate parameters: %s", err))
	}
}

func migrateAffiliateOverrides(ctx sdk.Context, affiliateKeeper affiliatekeeper.Keeper) {
	// Get all whitelist
	whitelist, err := affiliateKeeper.GetAffiliateWhitelist(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to get affiliate whitelist: %s", err))
	}
	// Create overrides for all whitelist addresses
	overrides := affiliatetypes.AffiliateOverrides{}
	var overridesList []string
	for _, addr := range whitelist.Tiers {
		overridesList = append(overridesList, addr.Addresses...)
	}
	overrides.Addresses = overridesList
	// Update affiliate overrides
	err = affiliateKeeper.SetAffiliateOverrides(ctx, overrides)
	if err != nil {
		panic(fmt.Sprintf("failed to set affiliate overrides: %s", err))
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	affiliateKeeper affiliatekeeper.Keeper,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Set default affiliate tiers and parameters.
		setDefaultAffiliateTiersForSlidingAffiliates(sdkCtx, affiliateKeeper)

		// Set default affiliate parameters.
		setDefaultAffiliateParameters(sdkCtx, affiliateKeeper)

		// Migrate affiliate overrides.
		migrateAffiliateOverrides(sdkCtx, affiliateKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
