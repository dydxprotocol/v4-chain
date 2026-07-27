package keeper_test

// Tests for the governance-tunable conditional-order trigger config and the always-on bounded
// trigger / expiry paths.
//
// The mitigation is activated by the state-breaking v9.6 upgrade and is always on; the config only
// tunes the per-block trigger budget, the per-block expiry-removal budget, and the admission caps.
//
// Test groups:
//  1. TestTriggerConfig_DefaultAndRoundTrip — absent config decodes to the package-constant
//     defaults; set/get round-trips; zero fields normalize to defaults.
//  2. TestTriggerConfig_AppliesBudget — the per-block trigger budget bounds triggers per block and
//     defers the remainder deterministically to subsequent blocks.
//  3. TestTriggerConfig_CustomBudgetHonored — a governance-set custom budget is honored.
//  4. TestTriggerConfig_Deterministic — the bounded path is deterministic across independent runs.
//  5. TestRemoveExpiredStatefulOrders_Bounded — expiry removal is bounded by MaxRemovalsPerBlock per
//     block, draining the remainder over subsequent blocks with no loss.

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// enableTriggerConfig sets the per-block trigger budget. The bounded path is always active (the
// state-breaking upgrade builds the index), so this only tunes the budget; the trigger-price index
// is maintained incrementally by the placement hooks that seeded the orders.
func enableTriggerConfig(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context, budget uint32) {
	t.Helper()
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{
		MaxTriggersPerBlock: budget,
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// 1. Default + round-trip
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_DefaultAndRoundTrip(t *testing.T) {
	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// Absent → default → package-constant budgets/caps.
	got := k.GetConditionalOrderTriggerConfig(ctx)
	require.Equal(t, uint32(clobkeeper.MaxConditionalTriggersPerBlock), got.MaxTriggersPerBlock)
	require.Equal(t, uint32(clobkeeper.MaxConditionalRemovalsPerBlock), got.MaxRemovalsPerBlock)

	// Round-trip a custom config.
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{
		MaxTriggersPerBlock: 250,
		MaxRemovalsPerBlock: 400,
	})
	got = k.GetConditionalOrderTriggerConfig(ctx)
	require.Equal(t, uint32(250), got.MaxTriggersPerBlock)
	require.Equal(t, uint32(400), got.MaxRemovalsPerBlock)

	// Normalization: zero fields fall back to the package constants.
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{})
	got = k.GetConditionalOrderTriggerConfig(ctx)
	require.Equal(t, uint32(clobkeeper.MaxConditionalTriggersPerBlock), got.MaxTriggersPerBlock)
	require.Equal(t, uint32(clobkeeper.MaxConditionalRemovalsPerBlock), got.MaxRemovalsPerBlock)
}

// ──────────────────────────────────────────────────────────────────────────────
// 2. Per-block budget bounds triggers and defers deterministically
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_AppliesBudget(t *testing.T) {
	const budget = uint32(40)
	const extra = 15
	total := int(budget) + extra

	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	for i := 0; i < total; i++ {
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, uint64(1+i),
		))
	}

	enableTriggerConfig(t, k, ctx, budget)

	// Block 1: exactly `budget` trigger.
	block1 := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, block1, int(budget), "exactly the budget should trigger on block 1")

	// `extra` remain.
	require.Len(t, k.GetAllUntriggeredConditionalOrders(ctx), extra)

	// Block 2: the remaining `extra` trigger.
	block2 := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, block2, extra, "remaining orders trigger on block 2")

	// Total equals the unbounded result; no duplicates.
	all := orderIdSet(append(block1, block2...))
	require.Len(t, all, total)
	require.Empty(t, k.GetAllUntriggeredConditionalOrders(ctx))
}

// ──────────────────────────────────────────────────────────────────────────────
// 3. Custom budget from SetConditionalOrderTriggerConfig
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_CustomBudgetHonored(t *testing.T) {
	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// 10 crossed orders, custom budget of 3.
	for i := 0; i < 10; i++ {
		placeOrder(t, k, ctx, makeTrigTestOrder(
			satypes.SubaccountId{Owner: condTestOwner, Number: uint32(i)},
			uint32(i), 0, clobkeeper.TriggerDirectionGTE, uint64(1+i)))
	}

	enableTriggerConfig(t, k, ctx, 3)

	triggered := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, triggered, 3, "custom budget of 3 must be honored")
	require.Len(t, k.GetAllUntriggeredConditionalOrders(ctx), 7)
}

// ──────────────────────────────────────────────────────────────────────────────
// 4. Determinism of the bounded path across independent runs
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_Deterministic(t *testing.T) {
	run := func() []clobtypes.OrderId {
		ks := newCondTestKeepers(t)
		ctx := ks.Ctx
		k := ks.ClobKeeper
		expiry := ctx.BlockTime().Add(clobtypes.StatefulOrderTimeWindow - time.Hour)
		for i := 0; i < 20; i++ {
			o := makeTrigTestOrder(constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, uint64(1+i))
			k.SetLongTermOrderPlacement(ctx, o, 1)
			k.AddStatefulOrderIdExpiration(ctx, expiry, o.OrderId)
		}
		enableTriggerConfig(t, k, ctx, 7)
		return k.MaybeTriggerConditionalOrders(ctx)
	}
	require.Equal(t, run(), run(), "bounded path must be deterministic across runs")
}

// TestRemoveExpiredStatefulOrders_Bounded proves the per-block expiry budget:
// RemoveExpiredStatefulOrders removes at most MaxRemovalsPerBlock expired entries per block,
// draining the remainder over subsequent blocks, keeping mass-simultaneous expiry a bounded
// per-block workload with no loss.
func TestRemoveExpiredStatefulOrders_Bounded(t *testing.T) {
	const n = 10
	const budget = 4
	expiry := time.Unix(1_700_000_000, 0)

	seed := func(k *clobkeeper.Keeper, ctx sdk.Context) {
		for i := 0; i < n; i++ {
			o := newCondTestOrder(i, condTestTriggerFarBelow)
			k.SetLongTermOrderPlacement(ctx, o, 1)
			k.AddStatefulOrderIdExpiration(ctx, expiry, o.OrderId)
		}
	}

	// With MaxRemovalsPerBlock=budget: each block removes at most `budget`; the remainder drains
	// over subsequent blocks; total removed equals n with no loss.
	ks := newCondTestKeepers(t)
	ks.ClobKeeper.SetConditionalOrderTriggerConfig(ks.Ctx, clobkeeper.ConditionalOrderTriggerConfig{
		MaxTriggersPerBlock: clobkeeper.MaxConditionalTriggersPerBlock,
		MaxRemovalsPerBlock: budget,
	})
	seed(ks.ClobKeeper, ks.Ctx)

	first := ks.ClobKeeper.RemoveExpiredStatefulOrders(ks.Ctx, expiry)
	require.Len(t, first, budget, "first block removes exactly the per-block removal budget")

	total := len(first)
	for {
		r := ks.ClobKeeper.RemoveExpiredStatefulOrders(ks.Ctx, expiry)
		if len(r) == 0 {
			break
		}
		require.LessOrEqual(t, len(r), budget, "no block exceeds the removal budget")
		total += len(r)
	}
	require.Equal(t, n, total, "all expired orders are eventually drained, none lost")
}
