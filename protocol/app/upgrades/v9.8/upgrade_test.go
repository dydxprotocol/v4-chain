package v_9_8_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// TestUpgradeMigration_IndexesTwapTriggers exercises the v9.8 migration on the real application
// keeper: an orphan trigger-store entry and a duplicate entry for a live suborder are removed, while
// the live parent's first-due trigger is preserved and indexed. The migration logic is unit-tested in
// x/clob/keeper (TestIndexTwapTriggers); this runs it against the full app's state stores, the same
// surface the upgrade handler operates on.
func TestUpgradeMigration_IndexesTwapTriggers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain().WithBlockTime(time.Unix(1000, 0))
	k := tApp.App.ClobKeeper
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)

	// An orphan trigger key whose parent placement is missing.
	orphanTriggerSuborderId := clobtypes.OrderId{
		SubaccountId: constants.Carl_Num0,
		ClientId:     7,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	triggerStore.Set(clobtypes.GetTWAPTriggerKey(1000, orphanTriggerSuborderId), []byte{})

	// A live TWAP parent. SetTWAPOrderPlacement creates its first-due trigger; leave it unindexed as
	// pre-upgrade state and add a later duplicate trigger for the same suborder id.
	twapOrder := constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO
	k.SetTWAPOrderPlacement(ctx, twapOrder, 5)
	liveSuborderId := clobtypes.OrderId{
		SubaccountId: twapOrder.OrderId.SubaccountId,
		ClientId:     twapOrder.OrderId.ClientId,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrder.OrderId.ClobPairId,
	}
	k.GetTWAPTriggerIndexStore(ctx).Delete(liveSuborderId.ToStateKey())
	triggerStore.Set(clobtypes.GetTWAPTriggerKey(1030, liveSuborderId), []byte{})

	indexed, deleted := k.IndexTwapTriggers(ctx)
	require.Equal(t, 1, indexed)
	require.Equal(t, 2, deleted)

	// The orphan trigger is gone.
	_, _, found := k.GetTwapTriggerPlacement(ctx, orphanTriggerSuborderId)
	require.False(t, found)

	// The live parent placement is untouched; only its first-due trigger remains, and it is indexed.
	_, found = k.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.True(t, found)
	_, scannedTime, found := k.GetTwapTriggerPlacement(ctx, liveSuborderId)
	require.True(t, found)
	require.Equal(t, int64(1000), scannedTime)
	indexedTime, found := k.GetPendingTwapTrigger(ctx, liveSuborderId)
	require.True(t, found)
	require.Equal(t, int64(1000), indexedTime)
}
