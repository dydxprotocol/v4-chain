package v_9_7

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	subaccountskeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
)

// SweepIsolatedInsuranceFunds transfers the isolated insurance fund balance of every isolated
// perpetual whose market is in final settlement to the cross insurance fund. Markets that entered
// final settlement before this upgrade never had their insurance funds swept; markets closing
// after this upgrade are swept automatically on transition. An error aborts the upgrade so the
// one-shot migration can be fixed and retried rather than silently stranding funds.
//
// Clob pairs are read from state rather than via GetClobPairIdForPerpetual: the upgrade module's
// PreBlocker runs before the clob module's PreBlocker hydrates the in-memory
// perpetual-to-clob-pair mapping.
func SweepIsolatedInsuranceFunds(
	ctx sdk.Context,
	clobKeeper *clobkeeper.Keeper,
	subaccountsKeeper *subaccountskeeper.Keeper,
) error {
	for _, clobPair := range clobKeeper.GetAllClobPairs(ctx) {
		if clobPair.Status != clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT {
			continue
		}

		perpetualId, err := clobPair.GetPerpetualId()
		if err != nil {
			return fmt.Errorf("failed to get perpetual id for clob pair %d: %w", clobPair.Id, err)
		}
		if err := subaccountsKeeper.TransferIsolatedInsuranceFundToCross(ctx, perpetualId); err != nil {
			return fmt.Errorf("failed to sweep insurance fund for perpetual %d: %w", perpetualId, err)
		}

		ctx.Logger().Info(fmt.Sprintf(
			"SweepIsolatedInsuranceFunds: processed final settlement perpetual %d",
			perpetualId,
		))
	}
	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper *clobkeeper.Keeper,
	subaccountsKeeper *subaccountskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		if err := SweepIsolatedInsuranceFunds(sdkCtx, clobKeeper, subaccountsKeeper); err != nil {
			return vm, err
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
