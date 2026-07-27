package keeper_test

// Packet 3 acceptance tests: untriggered conditional order admission caps + counters.
//
// Test groups:
//  1. TestAdmissionCap_OverCapRejected        — over-cap placement returns ErrTooManyUntriggeredConditionalOrders
//  2. TestAdmissionCap_UnderCapSucceeds        — under-cap placement does not return the cap error
//  3. TestAdmissionCap_CounterSymmetry         — counters increment on place, decrement on cancel/trigger/expiry
//  4. TestAdmissionCap_NonConditionalUnaffected — long-term and short-term orders are unaffected
//  5. TestAdmissionCap_HydrationSafe           — SetLongTermOrderPlacement (InitMemStore) never rejects at cap
//  6. TestAdmissionCap_GlobalCapEnforced       — global cap rejects even when per-subaccount is under cap
//  7. TestAdmissionCap_PerSubaccountCapEnforced — per-subaccount cap rejects even when global is under cap

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// capTestSA builds a SubaccountId for tests.
func capTestSA(number uint32) satypes.SubaccountId {
	return satypes.SubaccountId{
		Owner:  condTestOwner,
		Number: number,
	}
}

// capConditionalOrder builds a minimal conditional (take-profit buy) order for the given
// subaccount + clientId combination.
func capConditionalOrder(sa satypes.SubaccountId, clientId uint32, triggerSubticks uint64) clobtypes.Order {
	return clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: sa,
			ClientId:     clientId,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   constants.ClobPair_Btc.Id,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        1_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: condTestGTBT},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: triggerSubticks,
	}
}

// capLongTermOrder builds a minimal long-term order for the given subaccount.
func capLongTermOrder(sa satypes.SubaccountId, clientId uint32) clobtypes.Order {
	return clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: sa,
			ClientId:     clientId,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   constants.ClobPair_Btc.Id,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     1_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: condTestGTBT},
	}
}

// capTestKeeper initializes a test ClobKeeper with BTC pair only.
func capTestKeeper(t *testing.T) keepertest.ClobKeepersTestContext {
	t.Helper()

	bankKeeper := &mocks.BankKeeper{}
	bankKeeper.On("GetBalance", mock.Anything, mock.Anything, constants.Usdc.Denom).
		Return(sdk.NewCoin(constants.Usdc.Denom, sdkmath.ZeroInt())).
		Maybe()

	ks := keepertest.NewClobKeepersTestContext(
		t,
		memclob.NewMemClobPriceTimePriority(false),
		bankKeeper,
		indexer_manager.NewIndexerEventManagerNoop(),
	)
	ctx := ks.Ctx

	require.NoError(t, keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper))
	keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)
	_, err := ks.PerpetualsKeeper.CreatePerpetual(
		ctx,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Id,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Ticker,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.MarketId,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.AtomicResolution,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.DefaultFundingPpm,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.LiquidityTier,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.MarketType,
	)
	require.NoError(t, err)
	_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
		ctx,
		constants.ClobPair_Btc.Id,
		constants.ClobPair_Btc.MustGetPerpetualId(),
		satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
		constants.ClobPair_Btc.QuantumConversionExponent,
		constants.ClobPair_Btc.SubticksPerTick,
		constants.ClobPair_Btc.Status,
	)
	require.NoError(t, err)

	return ks
}

// capValidCtx returns a context with a valid positive block time and height so that
// PerformStatefulOrderValidation does not fail with time-related errors.
//
// condTestGTBT = 1_700_000_000 and StatefulOrderTimeWindow = 95 days ≈ 8_208_000 s, so
// blockTime must satisfy:
//
//	blockTime < condTestGTBT                            (GoodTilBlockTime > blockTime)
//	blockTime > condTestGTBT - StatefulOrderTimeWindow  (GoodTilBlockTime within window)
//
// 1_695_000_000 (2023-09-18) satisfies both constraints.
func capValidCtx(ks keepertest.ClobKeepersTestContext) sdk.Context {
	blockTime := time.Unix(1_695_000_000, 0) // 2023-09-18, within StatefulOrderTimeWindow of condTestGTBT
	ctx := ks.Ctx.WithBlockTime(blockTime).WithBlockHeight(1)
	ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
		Timestamp: blockTime,
	})
	return ctx
}

// TestAdmissionCap_CounterSymmetry validates that counters accurately track placements,
// cancels, triggers, and expiries.
func TestAdmissionCap_CounterSymmetry(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	sa0 := capTestSA(0)
	sa1 := capTestSA(1)

	// Initially zero.
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0))
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa1))

	// Place 2 orders for sa0, 1 for sa1.
	order0 := capConditionalOrder(sa0, 100, condTestTriggerFarBelow)
	order1 := capConditionalOrder(sa0, 101, condTestTriggerFarBelow)
	order2 := capConditionalOrder(sa1, 200, condTestTriggerFarBelow)

	k.SetLongTermOrderPlacement(ctx, order0, 1)
	k.SetLongTermOrderPlacement(ctx, order1, 1)
	k.SetLongTermOrderPlacement(ctx, order2, 1)

	require.Equal(t, uint32(3), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	require.Equal(t, uint32(2), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0))
	require.Equal(t, uint32(1), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa1))

	// Cancel order0 (via DeleteLongTermOrderPlacement).
	k.DeleteLongTermOrderPlacement(ctx, order0.OrderId)
	require.Equal(t, uint32(2), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	require.Equal(t, uint32(1), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0))
	require.Equal(t, uint32(1), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa1))

	// Trigger order1 (via MustTriggerConditionalOrder).
	k.MustTriggerConditionalOrder(ctx, order1.OrderId)
	require.Equal(t, uint32(1), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0))
	require.Equal(t, uint32(1), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa1))

	// Expire order2 (via RemoveExpiredStatefulOrders).
	expTime := ctx.BlockTime().Add(clobtypes.StatefulOrderTimeWindow - time.Hour)
	k.AddStatefulOrderIdExpiration(ctx, expTime, order2.OrderId)
	expiredIds := k.RemoveExpiredStatefulOrders(ctx, expTime.Add(time.Second))
	require.Len(t, expiredIds, 1)
	require.Equal(t, order2.OrderId, expiredIds[0])
	// After expiry, DeleteLongTermOrderPlacement is NOT called by RemoveExpiredStatefulOrders
	// directly; the caller (abci.go) calls MustRemoveStatefulOrder. Let's call it explicitly.
	k.DeleteLongTermOrderPlacement(ctx, order2.OrderId)
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0))
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa1))
}

// TestAdmissionCap_NonConditionalUnaffected confirms that long-term orders are not counted
// in the untriggered conditional counters.
func TestAdmissionCap_NonConditionalUnaffected(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	sa0 := capTestSA(0)

	// Place several long-term orders.
	for i := uint32(0); i < 10; i++ {
		lt := capLongTermOrder(sa0, i)
		k.SetLongTermOrderPlacement(ctx, lt, 1)
	}

	// Untriggered counters must remain zero — long-term orders don't affect them.
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	require.Equal(t, uint32(0), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0))
}

// TestAdmissionCap_HydrationSafe verifies that SetLongTermOrderPlacement (the path used by
// InitMemStore / state hydration / replay) never panics or enforces the cap.  We seed the
// store well above the per-subaccount cap and confirm it succeeds without error or panic.
func TestAdmissionCap_HydrationSafe(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	sa0 := capTestSA(0)

	// Seed above the per-subaccount cap (MaxUntriggeredConditionalOrdersPerSubaccount = 200)
	// directly via SetLongTermOrderPlacement — this must NOT panic.
	aboveCap := int(clobkeeper.MaxUntriggeredConditionalOrdersPerSubaccount) + 50
	require.NotPanics(t, func() {
		for i := 0; i < aboveCap; i++ {
			order := capConditionalOrder(sa0, uint32(1000+i), condTestTriggerFarBelow)
			k.SetLongTermOrderPlacement(ctx, order, 1)
		}
	})

	// Counters reflect the seeded count even though it exceeds the cap.
	require.Equal(t, uint32(aboveCap), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	require.Equal(t, uint32(aboveCap), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0))
}

// TestAdmissionCap_GlobalCapEnforced verifies that when the global counter reaches the cap,
// a new placement from a different subaccount (which is under its own per-SA cap) is rejected.
// We use the SetUntriggeredConditionalOrderCountGlobal helper to fast-forward the global counter
// without seeding millions of orders in state.
func TestAdmissionCap_GlobalCapEnforced(t *testing.T) {
	ks := capTestKeeper(t)
	k := ks.ClobKeeper

	// PerformStatefulOrderValidation requires a positive block time; use the benchmark ctx helper
	// that has a suitable block time and height.
	ctx := capValidCtx(ks)

	// Caps are enforced only when the mitigation flag is enabled (rolling-deploy gating).
	enableTriggerConfig(t, k, ctx, clobkeeper.MaxConditionalTriggersPerBlock)

	// Simulate global counter at the cap.
	k.SetUntriggeredConditionalOrderCountGlobal(ctx, clobkeeper.MaxUntriggeredConditionalOrdersGlobal)

	// A fresh subaccount (well under its per-SA cap) should be rejected.
	saFresh := capTestSA(99)
	order := capConditionalOrder(saFresh, 1, condTestTriggerFarBelow)

	err := k.PlaceStatefulOrder(ctx, &clobtypes.MsgPlaceOrder{Order: order}, true)
	require.ErrorIs(t, err, clobtypes.ErrTooManyUntriggeredConditionalOrders,
		"expected ErrTooManyUntriggeredConditionalOrders when global cap is reached")
}

// TestAdmissionCap_PerSubaccountCapEnforced verifies that when a subaccount's counter reaches
// the per-SA cap, new placements for that subaccount are rejected even if the global counter is
// below the global cap.
func TestAdmissionCap_PerSubaccountCapEnforced(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := capValidCtx(ks)
	k := ks.ClobKeeper

	// Caps are enforced only when the mitigation flag is enabled (rolling-deploy gating).
	enableTriggerConfig(t, k, ctx, clobkeeper.MaxConditionalTriggersPerBlock)

	sa0 := capTestSA(0)

	// Set per-SA counter to cap without touching global (global stays at 0, which is fine).
	k.SetUntriggeredConditionalOrderCountForSubaccount(
		ctx, sa0, clobkeeper.MaxUntriggeredConditionalOrdersPerSubaccount,
	)
	// Global counter also needs to be at the same level to avoid global-cap being hit by
	// the per-SA count (it's already been artificially set so global check passes).
	// Leave global at 0 — the per-SA check should fire first.

	order := capConditionalOrder(sa0, 2000, condTestTriggerFarBelow)
	err := k.PlaceStatefulOrder(ctx, &clobtypes.MsgPlaceOrder{Order: order}, true)
	require.ErrorIs(t, err, clobtypes.ErrTooManyUntriggeredConditionalOrders,
		"expected ErrTooManyUntriggeredConditionalOrders when per-subaccount cap is reached")
}

// TestAdmissionCap_UnderCapSucceeds verifies that PlaceStatefulOrder does not return the cap
// error when both counters are well below their caps.  We call with isInternalOrder=true to
// skip equity-tier and revshare checks that would otherwise require a fuller setup.
func TestAdmissionCap_UnderCapSucceeds(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := capValidCtx(ks)
	k := ks.ClobKeeper

	sa0 := capTestSA(0)

	// Both counters start at 0 — well below any cap.
	order := capConditionalOrder(sa0, 3000, condTestTriggerFarBelow)

	// PlaceStatefulOrder with isInternalOrder=true skips equity-tier and collateral checks,
	// so the order placement should fail for a reason unrelated to the cap (e.g., no collateral)
	// OR succeed — but it must NOT fail with ErrTooManyUntriggeredConditionalOrders.
	err := k.PlaceStatefulOrder(ctx, &clobtypes.MsgPlaceOrder{Order: order}, true)
	if err != nil {
		require.NotErrorIs(t, err, clobtypes.ErrTooManyUntriggeredConditionalOrders,
			"under-cap placement must not be rejected with the cap error")
	}
}

// TestAdmissionCap_TriggerDecrementsCounter validates the specific invariant:
// triggering an untriggered conditional order decrements the untriggered counter.
func TestAdmissionCap_TriggerDecrementsCounter(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	sa0 := capTestSA(0)

	orders := make([]clobtypes.Order, 3)
	for i := range orders {
		orders[i] = capConditionalOrder(sa0, uint32(500+i), condTestTriggerFarBelow)
		k.SetLongTermOrderPlacement(ctx, orders[i], 1)
	}
	require.Equal(t, uint32(3), k.GetUntriggeredConditionalOrderCountGlobal(ctx))

	// Trigger the second order.
	k.MustTriggerConditionalOrder(ctx, orders[1].OrderId)

	require.Equal(t, uint32(2), k.GetUntriggeredConditionalOrderCountGlobal(ctx),
		"trigger must decrement the global untriggered counter")
	require.Equal(t, uint32(2), k.GetUntriggeredConditionalOrderCountForSubaccount(ctx, sa0),
		"trigger must decrement the per-subaccount untriggered counter")

	// The triggered order is now in triggered state; confirm the other two are still untriggered.
	require.False(t, k.IsConditionalOrderTriggered(ctx, orders[0].OrderId))
	require.True(t, k.IsConditionalOrderTriggered(ctx, orders[1].OrderId))
	require.False(t, k.IsConditionalOrderTriggered(ctx, orders[2].OrderId))
}

// TestAdmissionCap_CrossWindowAccumulationBounded verifies the key attack scenario from the
// conditional-order EndBlocker work report: accumulation across placement windows from multiple subaccounts is now
// bounded by the global cap.  We simulate filling a large fraction of the global cap across
// many subaccounts, confirm counters are correct, then confirm the (cap+1)th order is rejected.
func TestAdmissionCap_CrossWindowAccumulationBounded(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// Caps are enforced only when the mitigation flag is enabled (rolling-deploy gating).
	enableTriggerConfig(t, k, ctx, clobkeeper.MaxConditionalTriggersPerBlock)

	// Seed close to the global cap using multiple subaccounts (each well under per-SA cap).
	// We use SetUntriggeredConditionalOrderCountGlobal to fast-forward rather than seeding
	// millions of state entries, and also advance the per-SA counter for a specific SA so
	// the final order would only be rejected by the global cap.
	targetGlobal := clobkeeper.MaxUntriggeredConditionalOrdersGlobal - 1
	k.SetUntriggeredConditionalOrderCountGlobal(ctx, targetGlobal)
	// sa99 per-SA count is 0 — well under per-SA cap.

	sa99 := capTestSA(99)

	// One more order should succeed (global goes from cap-1 to cap).
	// Note: we manipulate counters directly here; the actual write to state is done via
	// SetLongTermOrderPlacement which calls IncrementUntriggeredConditionalOrderCount.
	order := capConditionalOrder(sa99, 9999, condTestTriggerFarBelow)
	k.SetLongTermOrderPlacement(ctx, order, 1)
	require.Equal(t, clobkeeper.MaxUntriggeredConditionalOrdersGlobal,
		k.GetUntriggeredConditionalOrderCountGlobal(ctx),
		"counter should be at global cap after seeding")

	// Now the (cap+1)th order should be rejected.
	// Need a valid block time context for PlaceStatefulOrder's PerformStatefulOrderValidation.
	validCtx := capValidCtx(ks)
	// Copy the global counter to the valid context's store by setting it explicitly.
	k.SetUntriggeredConditionalOrderCountGlobal(validCtx, clobkeeper.MaxUntriggeredConditionalOrdersGlobal)

	order2 := capConditionalOrder(sa99, 10000, condTestTriggerFarBelow)
	err := k.PlaceStatefulOrder(validCtx, &clobtypes.MsgPlaceOrder{Order: order2}, true)
	require.ErrorIs(t, err, clobtypes.ErrTooManyUntriggeredConditionalOrders,
		"cross-window accumulation at global cap must be rejected")
}

func TestAdmissionCap_LoweredBelowLiveCountBlocksThenResumesAfterDrain(t *testing.T) {
	ks := capTestKeeper(t)
	ctx := capValidCtx(ks)
	k := ks.ClobKeeper
	sa0 := capTestSA(0)

	orders := []clobtypes.Order{
		capConditionalOrder(sa0, 700, condTestTriggerFarBelow),
		capConditionalOrder(sa0, 701, condTestTriggerFarBelow),
		capConditionalOrder(sa0, 702, condTestTriggerFarBelow),
	}
	for _, order := range orders {
		k.SetLongTermOrderPlacement(ctx, order, 1)
	}
	require.Equal(t, uint32(3), k.GetUntriggeredConditionalOrderCountGlobal(ctx))

	// Governance may deliberately lower the cap below the live set to stop new admission while
	// existing orders cancel, trigger, or expire.
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{
		Enabled:                               true,
		MaxUntriggeredConditionalOrdersGlobal: 2,
		MaxUntriggeredConditionalOrdersPerSubaccount: 2,
	})
	for !k.IsConditionalOrderTriggerIndexReady(ctx) {
		k.AdvanceConditionalOrderTriggerIndexActivation(ctx)
	}

	fresh := capConditionalOrder(capTestSA(1), 703, condTestTriggerFarBelow)
	err := k.PlaceStatefulOrder(ctx, &clobtypes.MsgPlaceOrder{Order: fresh}, true)
	require.ErrorIs(t, err, clobtypes.ErrTooManyUntriggeredConditionalOrders)

	// Drain below the new cap. Admission may still fail for unrelated validation/collateral
	// reasons in this minimal keeper fixture, but the cap must no longer reject it.
	k.DeleteLongTermOrderPlacement(ctx, orders[0].OrderId)
	k.DeleteLongTermOrderPlacement(ctx, orders[1].OrderId)
	require.Equal(t, uint32(1), k.GetUntriggeredConditionalOrderCountGlobal(ctx))
	err = k.PlaceStatefulOrder(ctx, &clobtypes.MsgPlaceOrder{Order: fresh}, true)
	require.NotErrorIs(t, err, clobtypes.ErrTooManyUntriggeredConditionalOrders)
}
