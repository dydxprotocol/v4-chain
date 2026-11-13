package v_9_9_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store/rootmulti"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	v_9_9 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.9"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	canonicalusdckeeper "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
)

// Runs the v9.9 upgrade handler against a real node home (MAINNET_HOME, holding
// data/application.db) without committing, and reports what GRADUAL would look like.
func TestMainnetStateDryRun(t *testing.T) {
	home := os.Getenv("MAINNET_HOME")
	if home == "" {
		t.Skip("MAINNET_HOME not set")
	}
	chainID := os.Getenv("DRYRUN_CHAIN_ID")
	if chainID == "" {
		chainID = "dydx-mainnet-1"
	}
	dataDir := filepath.Join(home, "data")
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, dataDir)
	require.NoError(t, err)
	defer db.Close()

	last := rootmulti.GetLatestVersion(db)
	require.Greater(t, last, int64(0))
	upgradeHeight := last + 1
	info := fmt.Sprintf(`{"name":%q,"height":%d}`, v_9_9.UpgradeName, upgradeHeight)
	require.NoError(t, os.WriteFile(filepath.Join(dataDir, "upgrade-info.json"), []byte(info), 0o644))

	logger := log.NewLogger(os.Stdout, log.LevelOption(zerolog.WarnLevel))
	opened := time.Now()
	dydxApp := app.New(logger, db, nil, true, appoptions.GetDefaultTestAppOptions(home, map[string]interface{}{
		flags.FlagChainID: chainID,
	}))
	loadTime := time.Since(opened)
	require.Equal(t, last, dydxApp.LastBlockHeight())

	ctx := dydxApp.NewUncachedContext(false, cmtproto.Header{ChainID: chainID, Height: upgradeHeight, Time: time.Now().UTC()})
	before, err := dydxApp.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	require.NotContains(t, before, canonicalusdctypes.ModuleName)
	k := dydxApp.CanonicalUsdcKeeper
	supplyBefore := dydxApp.BankKeeper.GetSupply(ctx, canonicalusdctypes.DefaultControls().LogicalDenom).Amount

	started := time.Now()
	require.NoError(t, dydxApp.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: v_9_9.UpgradeName, Height: upgradeHeight}))
	handlerTime := time.Since(started)

	after, err := dydxApp.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.NoError(t, err)
	changed := map[string][2]uint64{}
	for name, v := range after {
		if before[name] != v {
			changed[name] = [2]uint64{before[name], v}
		}
	}
	controls := k.GetControls(ctx)
	ledger := k.GetLedger(ctx)
	supply := dydxApp.BankKeeper.GetSupply(ctx, controls.LogicalDenom).Amount
	invariantMsg, broken := canonicalusdckeeper.BackingInvariant(*k)(ctx)
	clients := map[string]string{}
	for _, c := range []string{controls.NobleClient, controls.InjectiveClient} {
		cs, found := dydxApp.IBCKeeper.ClientKeeper.GetClientState(ctx, c)
		require.True(t, found, c)
		clients[c] = dydxApp.IBCKeeper.ClientKeeper.GetClientStatus(ctx, cs, c).String()
	}
	controlsJSON, err := dydxApp.AppCodec().MarshalJSON(&controls)
	require.NoError(t, err)
	var controlsMap map[string]any
	require.NoError(t, json.Unmarshal(controlsJSON, &controlsMap))
	report := map[string]any{
		"chain_id": chainID, "last_height": last, "upgrade_height": upgradeHeight,
		"load_seconds": loadTime.Seconds(), "handler_seconds": handlerTime.Seconds(),
		"module_versions_changed": changed,
		"controls":                controlsMap,
		"ledger": map[string]string{
			"noble_backing": ledger.NobleBacking.BigInt().String(), "injective_backing": ledger.InjectiveBacking.BigInt().String(),
			"pending_injective": ledger.PendingInjective.BigInt().String(), "pending_noble": ledger.PendingNoble.BigInt().String(),
			"refillable_capacity": ledger.RefillableCapacity.BigInt().String(),
		},
		"supply_before": supplyBefore.String(), "supply_after": supply.String(),
		"pending_count": k.PendingCount(ctx), "invariant_broken": broken, "invariant_msg": invariantMsg,
		"ibc_client_status": clients,
	}
	out, _ := json.MarshalIndent(report, "", "  ")
	fmt.Printf("DRYRUN_REPORT %s\n", out)

	require.Equal(t, canonicalusdctypes.Mode_MODE_GRADUAL, controls.Mode)
	require.Equal(t, "channel-0", controls.NobleChannel)
	require.Equal(t, "channel-62", controls.InjectiveChannel)
	require.True(t, controls.NobleDepositsDisabled)
	require.True(t, controls.NobleWithdrawalsEnabled)
	require.False(t, controls.InjectiveDepositsDisabled)
	require.Equal(t, v_9_9.InjectivePacketDenom, controls.InjectivePacketDenom)
	require.Equal(t, supply.String(), supplyBefore.String())
	require.Equal(t, 0, ledger.NobleBacking.BigInt().Cmp(supply.BigInt()))
	require.Equal(t, int64(0), ledger.InjectiveBacking.BigInt().Int64())
	require.Equal(t, int64(0), ledger.PendingInjective.BigInt().Int64())
	require.Equal(t, int64(0), ledger.PendingNoble.BigInt().Int64())
	require.False(t, broken, invariantMsg)
	require.Equal(t, map[string][2]uint64{canonicalusdctypes.ModuleName: {0, 1}}, changed)
}
