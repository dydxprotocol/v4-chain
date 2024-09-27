package v_8_0_0

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	listingtypes "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	listingkeeper "github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
)

func initListingModuleState(ctx sdk.Context, listingKeeper listingkeeper.Keeper) {
	// Set hard cap on listed markets
	err := listingKeeper.SetMarketsHardCap(ctx, 500)
	if err != nil {
		panic(fmt.Sprintf("failed to set markets hard cap: %s", err))
	}

	// Set listing vault deposit params
	err = listingKeeper.SetListingVaultDepositParams(
		ctx,
		listingtypes.DefaultParams(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to set listing vault deposit params: %s", err))
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	listingKeeper listingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		initListingModuleState(sdkCtx, listingKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
