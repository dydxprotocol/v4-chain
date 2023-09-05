package app

import (
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// setupUpgradeHandlers registers the upgrade handlers to perform custom upgrade
// logic and state migrations for software upgrades.
func (app *App) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		if app.UpgradeKeeper.HasHandler(upgrade.UpgradeName) {
			panic(fmt.Sprintf("Cannot register duplicate upgrade handler '%s'", upgrade.UpgradeName))
		}
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(app.ModuleManager, app.configurator),
		)
	}
}

// setUpgradeStoreLoaders sets custom store loaders to customize the rootMultiStore
// initialization for software upgrades.
func (app *App) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
		}
	}
}

// ScheduleForkUpgrade executes any necessary fork logic for based upon the current
// block height. It sets an upgrade plan once the chain reaches the pre-defined upgrade height.
//
// CONTRACT: for this logic to work properly it is required to:
//
//  1. Release a non-breaking patch version so that the chain can set the scheduled upgrade plan at upgrade-height.
//  2. Release the software defined in the upgrade-info.
func (app *App) scheduleForkUpgrade(ctx sdk.Context) {
	currentBlockHeight := ctx.BlockHeight()
	for _, fork := range Forks {
		if currentBlockHeight == fork.UpgradeHeight {
			upgradePlan := upgradetypes.Plan{
				Height: currentBlockHeight,
				Name:   fork.UpgradeName,
				Info:   fork.UpgradeInfo,
			}

			// schedule the upgrade plan to the current block height, effectively performing
			// a hard fork that uses the upgrade handler to manage the migration.
			if err := app.UpgradeKeeper.ScheduleUpgrade(ctx, upgradePlan); err != nil {
				panic(
					fmt.Errorf(
						"Hard Fork: failed to schedule upgrade %s during BeginBlock at height %d: %w",
						upgradePlan.Name,
						ctx.BlockHeight(),
						err,
					),
				)
			}
		}
	}
}
