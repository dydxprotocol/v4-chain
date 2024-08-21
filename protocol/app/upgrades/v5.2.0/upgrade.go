package v_5_2_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	vaultkeeper "github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper clobtypes.ClobKeeper,
	vaultKeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Migrate x/clob.
		clobKeeper.UnsafeMigrateOrderExpirationState(sdkCtx)

		// Migrate x/vault.
		setVaultOrderExpirationSecondsToSixty(sdkCtx, vaultKeeper)
		addAllVaultsToVaultAddressStore(sdkCtx, vaultKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// addAllVaultsToVaultAddressStore adds all existing vaults to Vault Address store.
func addAllVaultsToVaultAddressStore(
	ctx sdk.Context,
	vaultKeeper vaultkeeper.Keeper,
) {
	allVaults := vaultKeeper.GetAllVaults(ctx)
	for _, vault := range allVaults {
		vaultKeeper.AddVaultToAddressStore(ctx, vault.VaultId)
	}
}

// Deprecated: setVaultOrderExpirationSecondsToSixty sets vault module param `OrderExpirationSeconds` to 60.
// This function is deprecated because `Params` in `x/vault` is replaced with `DefaultQuotingParams` in v6.x.
func setVaultOrderExpirationSecondsToSixty(
	ctx sdk.Context,
	vaultKeeper vaultkeeper.Keeper,
) {
	// params := vaultKeeper.GetParams(ctx)
	// params.OrderExpirationSeconds = 60
	// err := vaultKeeper.SetParams(
	// 	ctx,
	// 	params,
	// )
	// if err != nil {
	// 	panic(err)
	// }
}
