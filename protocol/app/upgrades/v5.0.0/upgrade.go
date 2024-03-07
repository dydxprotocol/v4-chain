package v_5_0_0

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// Set all existing perpetuals to cross market type
func perpetualsUpgrade(
	ctx sdk.Context,
	perpetualsKeeper perptypes.PerpetualsKeeper,
) error {

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

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	perpetualsKeeper perptypes.PerpetualsKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Set all perpetuals to cross market type
		perpetualsUpgrade(sdkCtx, perpetualsKeeper)

		// TODO(TRA-93): Initialize `x/vault` module.

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
