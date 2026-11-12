package v_9_9_test

import (
	"testing"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	dydxapp "github.com/dydxprotocol/v4-chain/protocol/app"
	v_9_9 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.9"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	canonicalusdc "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

// openRoute writes an open transfer channel, connection and Tendermint client
// for chainID into the test app's IBC state.
func openRoute(ctx sdk.Context, app *dydxapp.App, chainID, channel, connection, client string) {
	app.IBCKeeper.ClientKeeper.SetClientState(ctx, client, &ibctm.ClientState{ChainId: chainID})
	app.IBCKeeper.ConnectionKeeper.SetConnection(ctx, connection, connectiontypes.NewConnectionEnd(
		connectiontypes.OPEN, client, connectiontypes.Counterparty{}, nil, 0,
	))
	app.IBCKeeper.ChannelKeeper.SetChannel(ctx, transfertypes.PortID, channel, channeltypes.NewChannel(
		channeltypes.OPEN, channeltypes.UNORDERED,
		channeltypes.NewCounterparty(transfertypes.PortID, "channel-remote"),
		[]string{connection}, transfertypes.Version,
	))
}

func TestUpgradeEntersGradualWhenBothRoutesResolve(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	app := tApp.App
	openRoute(ctx, app, v_9_9.NobleChainID, "channel-0", "connection-0", "07-tendermint-0")
	openRoute(ctx, app, v_9_9.InjectiveChainID, "channel-62", "connection-30", "07-tendermint-26")
	supply := app.BankKeeper.GetSupply(ctx, canonicalusdctypes.DefaultControls().LogicalDenom).Amount

	plan := upgradetypes.Plan{Name: v_9_9.UpgradeName, Height: ctx.BlockHeight()}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))

	controls := app.CanonicalUsdcKeeper.GetControls(ctx)
	require.Equal(t, canonicalusdctypes.Mode_MODE_GRADUAL, controls.Mode)
	require.Equal(t, "channel-0", controls.NobleChannel)
	require.Equal(t, "connection-0", controls.NobleConnection)
	require.Equal(t, "07-tendermint-0", controls.NobleClient)
	require.Equal(t, v_9_9.NoblePacketDenom, controls.NoblePacketDenom)
	require.True(t, controls.NobleWithdrawalsEnabled)
	require.False(t, controls.NobleDepositsDisabled)
	require.Equal(t, "channel-62", controls.InjectiveChannel)
	require.Equal(t, "connection-30", controls.InjectiveConnection)
	require.Equal(t, "07-tendermint-26", controls.InjectiveClient)
	require.Equal(t, v_9_9.InjectivePacketDenom, controls.InjectivePacketDenom)
	require.Equal(t,
		transfertypes.ParseDenomTrace("transfer/channel-62/"+v_9_9.InjectivePacketDenom).IBCDenom(),
		controls.InjectiveDenom,
	)
	require.False(t, controls.InjectiveDepositsDisabled)
	require.Equal(t, uint32(v_9_9.MaxPendingSettlements), controls.MaxPendingSettlements)
	require.NoError(t, controls.Validate())
	require.NoError(t, app.CanonicalUsdcKeeper.ValidateRouteBindings(ctx, controls))

	ledger := app.CanonicalUsdcKeeper.GetLedger(ctx)
	require.Equal(t, supply.String(), ledger.NobleBacking.String(), "Noble backing snapshots the logical supply")
	require.Equal(t, "0", ledger.InjectiveBacking.String())
}

func TestUpgradeStaysDisabledWhenARouteIsAmbiguous(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	app := tApp.App
	openRoute(ctx, app, v_9_9.NobleChainID, "channel-0", "connection-0", "07-tendermint-0")
	openRoute(ctx, app, v_9_9.InjectiveChainID, "channel-62", "connection-30", "07-tendermint-26")
	openRoute(ctx, app, v_9_9.InjectiveChainID, "channel-63", "connection-31", "07-tendermint-27")

	plan := upgradetypes.Plan{Name: v_9_9.UpgradeName, Height: ctx.BlockHeight()}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))

	require.Equal(t, canonicalusdctypes.DefaultControls(), app.CanonicalUsdcKeeper.GetControls(ctx))
	require.Equal(t, canonicalusdctypes.DefaultLedger(), app.CanonicalUsdcKeeper.GetLedger(ctx))
}

func TestUpgradeUsesPinnedRoutesOnMainnet(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain().WithChainID("dydx-mainnet-1")
	app := tApp.App
	pinned := v_9_9.PinnedRoutes["dydx-mainnet-1"]
	// A second Injective channel would make discovery ambiguous; pinning ignores it.
	openRoute(ctx, app, v_9_9.NobleChainID, pinned.Noble.Channel, pinned.Noble.Connection, pinned.Noble.Client)
	openRoute(ctx, app, v_9_9.InjectiveChainID, "channel-7", "connection-7", "07-tendermint-7")
	openRoute(ctx, app, v_9_9.InjectiveChainID,
		pinned.Injective.Channel, pinned.Injective.Connection, pinned.Injective.Client)

	plan := upgradetypes.Plan{Name: v_9_9.UpgradeName, Height: ctx.BlockHeight()}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))

	controls := app.CanonicalUsdcKeeper.GetControls(ctx)
	require.Equal(t, canonicalusdctypes.Mode_MODE_GRADUAL, controls.Mode)
	require.Equal(t, pinned.Injective.Channel, controls.InjectiveChannel)
	require.Equal(t, pinned.Injective.Client, controls.InjectiveClient)
}

func TestUpgradeStaysDisabledWhenPinnedRoutesAreMissing(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain().WithChainID("dydx-mainnet-1")
	app := tApp.App
	// Discovery would find these, but mainnet is pinned to other ids.
	openRoute(ctx, app, v_9_9.NobleChainID, "channel-0", "connection-0", "07-tendermint-0")
	openRoute(ctx, app, v_9_9.InjectiveChainID, "channel-1", "connection-1", "07-tendermint-1")

	plan := upgradetypes.Plan{Name: v_9_9.UpgradeName, Height: ctx.BlockHeight()}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))

	require.Equal(t, canonicalusdctypes.Mode_MODE_DISABLED, app.CanonicalUsdcKeeper.GetControls(ctx).Mode)
}
