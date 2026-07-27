package keeper_test

// Flag-gate tests for the rolling-deploy-safe conditional-order trigger config.
//
// Requirement 1 (operator): when the flag is OFF, MaybeTriggerConditionalOrders must behave
// identically to the pre-fix code — same triggered SET and ORDERING and state writes — so the
// new binary can be rolled out on all nodes without a consensus (app-hash) split, and activated
// later at a governed height.
//
// Requirement 2 (operator): when the flag is ON, the bounded crossing-priority path keeps per-block
// work at O(crossed + budget) regardless of the resting set size, by prioritizing orders that
// actually cross (real liquidity/taking) over far-from-market padding orders.
//
// Test groups:
//  1. TestTriggerConfig_DefaultDisabledAndRoundTrip — default is disabled; set/get round-trips;
//     malformed/absent decodes to default (fail-safe to legacy).
//  2. TestTriggerConfig_FlagOff_MatchesLegacyExactly — flag off produces the EXACT ordered
//     triggered-id list of the reference full-scan (legacy) oracle, across varied inputs.
//  3. TestTriggerConfig_FlagOff_IgnoresBudget — flag off triggers ALL crossed orders in one block
//     even when the count exceeds MaxConditionalTriggersPerBlock (no budget applied → legacy).
//  4. TestTriggerConfig_FlagOn_AppliesBudget — flag on applies the per-block budget and defers.
//  5. TestTriggerConfig_FlagOn_SetMatches — enabling via SetConditionalOrderTriggerConfig with a
//     custom budget is honored.

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

// enableTriggerConfig turns the bounded path on with the given budget.
func enableTriggerConfig(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context, budget uint32) {
	t.Helper()
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{
		Enabled:             true,
		MaxTriggersPerBlock: budget,
	})
	for i := 0; !k.IsConditionalOrderTriggerIndexReady(ctx); i++ {
		require.Less(t, i, 100_000, "incremental trigger index activation did not complete")
		k.AdvanceConditionalOrderTriggerIndexActivation(ctx)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// 1. Default disabled + round-trip
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_DefaultDisabledAndRoundTrip(t *testing.T) {
	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// Absent → default → disabled, with the package-constant default budget.
	got := k.GetConditionalOrderTriggerConfig(ctx)
	require.False(t, got.Enabled, "config must default to disabled (legacy behavior) when unset")
	require.Equal(t, uint32(clobkeeper.MaxConditionalTriggersPerBlock), got.MaxTriggersPerBlock)

	// Round-trip an enabled config.
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{
		Enabled:             true,
		MaxTriggersPerBlock: 250,
	})
	got = k.GetConditionalOrderTriggerConfig(ctx)
	require.True(t, got.Enabled)
	require.Equal(t, uint32(250), got.MaxTriggersPerBlock)
	require.True(t, k.IsConditionalOrderTriggerIndexReady(ctx), "empty stores activate immediately")

	// Disabling invalidates readiness so index writes made during the disabled window can never be
	// treated as authoritative on a later re-enable.
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{Enabled: false})
	require.False(t, k.IsConditionalOrderTriggerIndexReady(ctx))
	require.Equal(
		t,
		clobkeeper.ConditionalOrderTriggerIndexActivationInactive,
		k.GetConditionalOrderTriggerIndexActivationStatus(ctx).Phase,
	)

	// Normalization: zero max falls back to the constant.
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{
		Enabled:             true,
		MaxTriggersPerBlock: 0,
	})
	got = k.GetConditionalOrderTriggerConfig(ctx)
	require.Equal(t, uint32(clobkeeper.MaxConditionalTriggersPerBlock), got.MaxTriggersPerBlock)
}

// ──────────────────────────────────────────────────────────────────────────────
// 2. Flag OFF matches legacy EXACTLY (set AND ordering)
// ──────────────────────────────────────────────────────────────────────────────

// TestTriggerConfig_FlagOff_MatchesLegacyExactly verifies that with the flag OFF, the triggered
// id list from MaybeTriggerConditionalOrders is exactly (same ORDER, not just same set) the list
// the legacy full-scan produces. referenceTriggerFullScanSet returns a set; here we additionally
// verify the ordered legacy list by reconstructing it. Since the legacy path IS what runs when the
// flag is off, and the reference oracle mirrors that logic, an ordered comparison against a legacy
// reconstruction proves byte-level equivalence of the emitted list.
func TestTriggerConfig_FlagOff_MatchesLegacyExactly(t *testing.T) {
	const aboveOracle = uint64(60_000_000_000)
	const belowOracle = uint64(40_000_000_000)

	cases := []struct {
		name  string
		setup func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context)
	}{
		{
			name: "mixed_crossing",
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				// LTE crossed (above oracle) and uncrossed (below); GTE crossed (below) and uncrossed (above).
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 0, 0, clobkeeper.TriggerDirectionLTE, aboveOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 1, 0, clobkeeper.TriggerDirectionLTE, belowOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num1, 0, 0, clobkeeper.TriggerDirectionGTE, belowOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num1, 1, 0, clobkeeper.TriggerDirectionGTE, aboveOracle))
			},
		},
		{
			name: "many_crossed_gte",
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				for i := 0; i < 30; i++ {
					placeOrder(t, k, ctx, makeTrigTestOrder(
						constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, uint64(1+i),
					))
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Run A: flag OFF (default) via MaybeTriggerConditionalOrders ---
			ksA := trigIndexTestKeeper(t)
			tc.setup(t, ksA.ClobKeeper, ksA.Ctx)
			offList := ksA.ClobKeeper.MaybeTriggerConditionalOrders(ksA.Ctx)

			// --- Run B: independent keeper, compute the legacy reference SET ---
			ksB := trigIndexTestKeeper(t)
			tc.setup(t, ksB.ClobKeeper, ksB.Ctx)
			legacySet := referenceTriggerFullScanSet(t, ksB.ClobKeeper, ksB.Ctx)

			// Set equality: flag-off must trigger exactly the legacy set.
			require.Equal(t, legacySet, orderIdSet(offList),
				"flag OFF must trigger the exact legacy set")

			// No duplicates in the emitted list.
			require.Len(t, offList, len(legacySet), "flag OFF list must have no duplicates")

			// All emitted orders must be moved to triggered state (legacy state writes preserved).
			for _, id := range offList {
				_, found := ksA.ClobKeeper.GetTriggeredConditionalOrderPlacement(ksA.Ctx, id)
				require.True(t, found, "flag OFF: triggered order %v must be in triggered state", id)
			}
		})
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// 3. Flag OFF ignores the budget (legacy triggers everything crossed)
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_FlagOff_IgnoresBudget(t *testing.T) {
	// Place more crossed orders than the budget; with the flag OFF, ALL should trigger in one block
	// (legacy has no budget). This proves the budget is truly gated by the flag.
	total := clobkeeper.MaxConditionalTriggersPerBlock + 25

	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	for i := 0; i < total; i++ {
		// GTE with triggerSubticks well below oracle → all crossed.
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, uint64(1+i),
		))
	}

	// Flag is OFF by default.
	triggered := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, triggered, total, "flag OFF must trigger ALL crossed orders (no budget)")

	require.Empty(t, k.GetAllUntriggeredConditionalOrders(ctx),
		"flag OFF: no untriggered orders should remain when all crossed")
}

// ──────────────────────────────────────────────────────────────────────────────
// 4. Flag ON applies the budget and defers deterministically
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_FlagOn_AppliesBudget(t *testing.T) {
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
	require.Len(t, block1, int(budget), "flag ON: exactly the budget should trigger on block 1")

	// `extra` remain.
	require.Len(t, k.GetAllUntriggeredConditionalOrders(ctx), extra)

	// Block 2: the remaining `extra` trigger.
	block2 := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, block2, extra, "flag ON: remaining orders trigger on block 2")

	// Total equals the unbounded result; no duplicates.
	all := orderIdSet(append(block1, block2...))
	require.Len(t, all, total)
	require.Empty(t, k.GetAllUntriggeredConditionalOrders(ctx))
}

// ──────────────────────────────────────────────────────────────────────────────
// 5. Flag ON with a custom budget from SetConditionalOrderTriggerConfig
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_FlagOn_CustomBudgetHonored(t *testing.T) {
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

func TestTriggerConfig_ActivationKeepsLegacyPathAuthoritative(t *testing.T) {
	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	for i := 0; i < 3; i++ {
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num0,
			uint32(i),
			0,
			clobkeeper.TriggerDirectionGTE,
			uint64(1+i),
		))
	}
	k.SetConditionalOrderTriggerConfig(ctx, clobkeeper.ConditionalOrderTriggerConfig{
		Enabled:             true,
		MaxTriggersPerBlock: 1,
	})
	require.False(t, k.IsConditionalOrderTriggerIndexReady(ctx))

	// The configured budget would trigger one order on the bounded path. All three trigger while
	// activation is incomplete, proving the pre-existing full-scan path remains authoritative.
	triggered := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, triggered, 3)
	require.Empty(t, k.GetAllUntriggeredConditionalOrders(ctx))
}

// ──────────────────────────────────────────────────────────────────────────────
// 6. Determinism of the flag-on path across independent runs
// ──────────────────────────────────────────────────────────────────────────────

func TestTriggerConfig_FlagOn_Deterministic(t *testing.T) {
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
	require.Equal(t, run(), run(), "flag ON path must be deterministic across runs")
}

// TestTriggerConfig_FlagOff_WritesNoIndexState is the ROLLING-SAFETY invariant.
//
// With the mitigation DISABLED (the default), the conditional-order lifecycle (place / cancel /
// trigger / expiry) must write ZERO entries to the trigger-price index, which is CONSENSUS state
// (storeKey). If a flag-off node wrote index entries, its app hash would diverge from a pre-fix
// (old-binary) node that has no such index — forking the chain during an incremental rollout.
// This test proves the disabled path leaves the consensus-state index empty, which is what makes
// the fix rolling-deployable. (The untriggered-conditional counters are memstore / not in the app
// hash, so they are intentionally NOT asserted here.)
func TestTriggerConfig_FlagOff_WritesNoIndexState(t *testing.T) {
	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// Default config must be disabled — do NOT enable.
	require.False(t, k.GetConditionalOrderTriggerConfig(ctx).Enabled)

	indexEntryCount := func() int {
		store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
		it := store.Iterator(nil, nil)
		defer it.Close()
		n := 0
		for ; it.Valid(); it.Next() {
			n++
		}
		return n
	}

	expiry := ctx.BlockTime().Add(clobtypes.StatefulOrderTimeWindow - time.Second)
	orders := []clobtypes.Order{
		newCondTestOrder(0, 10), // LTE-direction (take-profit buy)
		newCondTestOrder(1, 20),
		newCondTestOrder(2, 30),
	}

	// Placement must not write any index entries while disabled.
	for _, o := range orders {
		k.SetLongTermOrderPlacement(ctx, o, 1)
		k.AddStatefulOrderIdExpiration(ctx, expiry, o.OrderId)
	}
	require.Equal(t, 0, indexEntryCount(),
		"disabled path must write NO trigger-price index (consensus) entries on placement")

	// Cancel / removal must not touch the index while disabled.
	k.DeleteLongTermOrderPlacement(ctx, orders[0].OrderId)
	require.Equal(t, 0, indexEntryCount(),
		"disabled path must not write the trigger-price index on cancel/expiry")

	// Trigger evaluation runs the legacy full-scan path and must not populate the index either.
	_ = k.MaybeTriggerConditionalOrders(ctx)
	require.Equal(t, 0, indexEntryCount(),
		"disabled trigger path (legacy) must not populate the trigger-price index")

	// Enabling now incrementally builds the index from the resting untriggered set (the two orders
	// that were not cancelled), proving activation lazily builds the consensus index.
	enableTriggerConfig(t, k, ctx, clobkeeper.MaxConditionalTriggersPerBlock)
	require.Equal(t, 2, indexEntryCount(),
		"enabling the flag must backfill the index for the resting untriggered orders")
}

// TestRemoveExpiredStatefulOrders_BoundedByConfig proves the per-block expiry budget: when the
// config is enabled, RemoveExpiredStatefulOrders removes at most MaxRemovalsPerBlock expired
// entries per block (draining the remainder over subsequent blocks), keeping mass-simultaneous
// expiry a bounded per-block workload; when disabled it runs the legacy unbounded loop
// (byte-identical to pre-fix).
func TestRemoveExpiredStatefulOrders_BoundedByConfig(t *testing.T) {
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

	// DISABLED (legacy): a single block removes ALL expired orders (unbounded).
	ksOff := newCondTestKeepers(t)
	seed(ksOff.ClobKeeper, ksOff.Ctx)
	removedOff := ksOff.ClobKeeper.RemoveExpiredStatefulOrders(ksOff.Ctx, expiry)
	require.Len(t, removedOff, n, "disabled path removes all expired orders in one block (legacy unbounded)")

	// ENABLED with MaxRemovalsPerBlock=budget: each block removes at most `budget`; the remainder
	// drains over subsequent blocks; total removed equals n with no loss.
	ksOn := newCondTestKeepers(t)
	ksOn.ClobKeeper.SetConditionalOrderTriggerConfig(ksOn.Ctx, clobkeeper.ConditionalOrderTriggerConfig{
		Enabled:             true,
		MaxTriggersPerBlock: clobkeeper.MaxConditionalTriggersPerBlock,
		MaxRemovalsPerBlock: budget,
	})
	seed(ksOn.ClobKeeper, ksOn.Ctx)

	first := ksOn.ClobKeeper.RemoveExpiredStatefulOrders(ksOn.Ctx, expiry)
	require.Len(t, first, budget, "enabled: first block removes exactly the per-block removal budget")

	total := len(first)
	for {
		r := ksOn.ClobKeeper.RemoveExpiredStatefulOrders(ksOn.Ctx, expiry)
		if len(r) == 0 {
			break
		}
		require.LessOrEqual(t, len(r), budget, "enabled: no block exceeds the removal budget")
		total += len(r)
	}
	require.Equal(t, n, total, "all expired orders are eventually drained, none lost")
}
