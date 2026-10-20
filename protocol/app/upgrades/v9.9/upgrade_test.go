package v_9_9_test

import (
	"testing"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	v_9_9 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.9"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	canonicalusdc "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestUpgradeAddsCanonicalUsdcStore(t *testing.T) {
	require.Equal(t, []string{canonicalusdctypes.StoreKey}, v_9_9.Upgrade.StoreUpgrades.Added)
}

func TestUpgradeInitializesCanonicalUsdc(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	app := tApp.App

	// Model the state immediately after the store loader adds the empty store:
	// existing modules are initialized, but canonicalusdc has no state or version.
	store := ctx.KVStore(app.GetKVStoreKey()[canonicalusdctypes.StoreKey])
	store.Delete([]byte(canonicalusdctypes.ControlsKey))
	store.Delete([]byte(canonicalusdctypes.LedgerKey))
	iterator := store.Iterator(nil, nil)
	require.False(t, iterator.Valid(), "the newly added store must be empty")
	require.NoError(t, iterator.Close())
	// SetModuleVersionMap only upserts entries; remove the stored version directly.
	upgradeStore := ctx.KVStore(app.GetKVStoreKey()[upgradetypes.StoreKey])
	upgradeStore.Delete(append([]byte{upgradetypes.VersionMapByte}, []byte(canonicalusdctypes.ModuleName)...))
	versions, err := app.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.NotContains(t, versions, canonicalusdctypes.ModuleName)

	plan := upgradetypes.Plan{Name: v_9_9.UpgradeName, Height: ctx.BlockHeight()}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))

	// Getters fall back to defaults on an empty store. Require persisted keys so
	// this test fails if the upgrade handler skips RunMigrations/InitGenesis.
	require.True(t, store.Has([]byte(canonicalusdctypes.ControlsKey)))
	require.True(t, store.Has([]byte(canonicalusdctypes.LedgerKey)))
	genesis := canonicalusdc.ExportGenesis(ctx, *app.CanonicalUsdcKeeper)
	require.Equal(t, canonicalusdctypes.DefaultControls(), genesis.Controls)
	require.Equal(t, canonicalusdctypes.Mode_MODE_DISABLED, genesis.Controls.Mode)
	require.Equal(t, canonicalusdctypes.DefaultLedger(), genesis.Ledger)
	require.Empty(t, genesis.PendingSettlements)
	require.NoError(t, genesis.Validate())
	updatedVersions, err := app.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.Equal(t, app.ModuleManager.GetVersionMap(), updatedVersions)
	require.Equal(t, uint64(1), updatedVersions[canonicalusdctypes.ModuleName])
	require.Equal(t, uint64(1), canonicalusdc.AppModule{}.ConsensusVersion())
}

func TestUpgradeIsNoOpWhenCanonicalUsdcAlreadyInitialized(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	app := tApp.App

	ledger := app.CanonicalUsdcKeeper.GetLedger(ctx)
	ledger.NobleBacking = dtypes.NewInt(123)
	require.NoError(t, app.CanonicalUsdcKeeper.SetLedger(ctx, ledger))
	versions, err := app.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), versions[canonicalusdctypes.ModuleName])

	plan := upgradetypes.Plan{Name: v_9_9.UpgradeName, Height: ctx.BlockHeight()}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))

	require.Equal(t, ledger, app.CanonicalUsdcKeeper.GetLedger(ctx))
	require.Equal(t, canonicalusdctypes.Mode_MODE_DISABLED, app.CanonicalUsdcKeeper.GetControls(ctx).Mode)
	updatedVersions, err := app.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), updatedVersions[canonicalusdctypes.ModuleName])
}
