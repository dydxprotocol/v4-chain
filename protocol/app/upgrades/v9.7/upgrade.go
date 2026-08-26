package v_9_7

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

		newVm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVm, err
		}

		// Activate the bounded conditional-order trigger path. This is a state-breaking change:
		// build the trigger-price secondary index from the full resting untriggered conditional
		// set in one pass (the O(N) cost is paid once, during this upgrade block, never per block),
		// and persist the default governance-tunable per-block budgets and caps. After this height
		// the EndBlocker uses the bounded index path for every resting conditional order — including
		// those placed before this upgrade. The per-subaccount / global memstore counters are
		// rehydrated by InitMemStore on restart.
		clobKeeper.BuildConditionalOrderTriggerPriceIndex(sdkCtx)
		// Persist the default budgets/caps explicitly so they are queryable and governance has a
		// concrete starting point. Zero values normalize to the package defaults.
		clobKeeper.SetConditionalOrderTriggerConfigParams(sdkCtx, 0, 0, 0, 0)

		return newVm, nil
	}
}
