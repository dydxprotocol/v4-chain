package v_0_3_0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
)

var (
	Upgrade = upgrades.Upgrade{
		UpgradeName: UpgradeName,
		StoreUpgrades: store.StoreUpgrades{
			Added: []string{
				evidencetypes.ModuleName,
			},
		},
	}
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running v0.3.0 Upgrade...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
