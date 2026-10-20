package v_9_9_test

import (
	"testing"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	v_9_9 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.9"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicalusdc "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/module"
)

// The upgrade on a mainnet-like state (logical USDC in supply, an open Noble
// transfer channel) must initialize canonicalusdc as DISABLED and leave every
// other module's state untouched.
func TestUpgradeOnLiveStateOnlyInitializesCanonicalUsdc(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	tApp.InitChain()
	app := tApp.App
	// The InitChain context wraps a store without iteration, which export needs.
	ctx := app.NewUncachedContext(false, cmtproto.Header{Height: app.LastBlockHeight()})

	supplyBefore := app.BankKeeper.GetSupply(ctx, assettypes.UusdcDenom).Amount
	require.True(t, supplyBefore.IsPositive(), "test genesis must carry logical USDC")
	app.IBCKeeper.ConnectionKeeper.SetConnection(ctx, "connection-0", connectiontypes.NewConnectionEnd(
		connectiontypes.OPEN,
		"07-tendermint-0",
		connectiontypes.NewCounterparty("07-tendermint-0", "connection-0", commitmenttypes.NewMerklePrefix([]byte("ibc"))),
		connectiontypes.GetCompatibleVersions(),
		0,
	))
	app.IBCKeeper.ChannelKeeper.SetChannel(ctx, transfertypes.PortID, "channel-0", channeltypes.NewChannel(
		channeltypes.OPEN,
		channeltypes.UNORDERED,
		channeltypes.NewCounterparty(transfertypes.PortID, "channel-0"),
		[]string{"connection-0"},
		transfertypes.Version,
	))

	// The test app never sets ICA host params, so that export panics.
	var exportable []string
	for _, name := range app.ModuleManager.ModuleNames() {
		if name != icatypes.ModuleName {
			exportable = append(exportable, name)
		}
	}
	before, err := app.ModuleManager.ExportGenesisForModules(ctx, app.AppCodec(), exportable)
	require.NoError(t, err)

	fromVM := app.ModuleManager.GetVersionMap()
	delete(fromVM, canonicalusdctypes.ModuleName)
	configurator := module.NewConfigurator(app.AppCodec(), app.MsgServiceRouter(), app.GRPCQueryRouter())
	toVM, err := v_9_9.CreateUpgradeHandler(app.ModuleManager, configurator)(
		ctx, upgradetypes.Plan{Name: v_9_9.UpgradeName}, fromVM,
	)
	require.NoError(t, err)
	require.Equal(t, app.ModuleManager.GetVersionMap(), toVM)

	after, err := app.ModuleManager.ExportGenesisForModules(ctx, app.AppCodec(), exportable)
	require.NoError(t, err)
	require.Equal(t, len(before), len(after))
	for name, state := range before {
		if name == canonicalusdctypes.ModuleName {
			continue
		}
		require.JSONEq(t, string(state), string(after[name]), "module %s changed", name)
	}

	exported := canonicalusdc.ExportGenesis(ctx, *app.CanonicalUsdcKeeper)
	require.Equal(t, canonicalusdctypes.Mode_MODE_DISABLED, exported.Controls.Mode)
	require.Equal(t, canonicalusdctypes.DefaultLedger(), exported.Ledger)
	require.Empty(t, exported.PendingSettlements)
	require.Equal(t, supplyBefore, app.BankKeeper.GetSupply(ctx, assettypes.UusdcDenom).Amount)
}
