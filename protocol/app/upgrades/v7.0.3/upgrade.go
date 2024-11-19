package v_7_0_3

import (
	"context"
	"fmt"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

const (
	ID_NUM = 200
)

// Set market, perpetual, and clob ids to a set number
// This is done so that the ids are consistent for convenience
func setMarketListingBaseIds(
	ctx sdk.Context,
	pricesKeeper pricestypes.PricesKeeper,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	clobKeeper clobtypes.ClobKeeper,
) {
	// Set all ids to a set number
	pricesKeeper.SetNextMarketID(ctx, ID_NUM)

	perpetualsKeeper.SetNextPerpetualID(ctx, ID_NUM)

	clobKeeper.SetNextClobPairID(ctx, ID_NUM)
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	pricesKeeper pricestypes.PricesKeeper,
	perpetualsKeeper perptypes.PerpetualsKeeper,
	clobKeeper clobtypes.ClobKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Set market, perpetual, and clob ids to a set number
		setMarketListingBaseIds(sdkCtx, pricesKeeper, perpetualsKeeper, clobKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
