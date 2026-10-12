package v_9_8

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper *clobkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// One-shot: build the TWAP trigger index and drop the orphan and duplicate trigger entries
		// left by pre-upgrade bugs. Reads state directly; this runs before the clob module's
		// PreBlocker hydrates the memclob, so no in-memory TWAP state exists yet.
		indexed, deleted := clobKeeper.IndexTwapTriggers(sdkCtx)
		sdkCtx.Logger().Info(fmt.Sprintf(
			"IndexTwapTriggers: indexed %d trigger(s), deleted %d orphan/duplicate trigger(s)",
			indexed,
			deleted,
		))

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
