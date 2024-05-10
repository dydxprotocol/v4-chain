package app

import (
	"fmt"

	v5_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v5.0.0"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
)

var (
	// `Upgrades` defines the upgrade handlers and store loaders for the application.
	// New upgrades should be added to this slice after they are implemented.
	Upgrades = []upgrades.Upgrade{
		v5_0_0.Upgrade,
	}
	Forks = []upgrades.Fork{}
)

// setupUpgradeHandlers registers the upgrade handlers to perform custom upgrade
// logic and state migrations for software upgrades.
func (app *App) setupUpgradeHandlers() {
	if app.UpgradeKeeper.HasHandler(v5_0_0.UpgradeName) {
		panic(fmt.Sprintf("Cannot register duplicate upgrade handler '%s'", v5_0_0.UpgradeName))
	}
	app.UpgradeKeeper.SetUpgradeHandler(
		v5_0_0.UpgradeName,
		v5_0_0.CreateUpgradeHandler(
			app.ModuleManager,
			app.configurator,
			app.PerpetualsKeeper,
			app.ClobKeeper,
			app.SubaccountsKeeper,
			app.ConsensusParamsKeeper,
			app.VaultKeeper,
		),
	)
}

// setupUpgradeStoreLoaders sets custom store loaders to customize the rootMultiStore
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
