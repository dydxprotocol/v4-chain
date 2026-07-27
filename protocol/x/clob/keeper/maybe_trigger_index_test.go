package keeper_test

// Packet 2 acceptance tests for the rewritten MaybeTriggerConditionalOrders.
//
// Test groups:
//  1. TestMaybeTrigger_EquivalenceVsFullScan — for randomized sets of untriggered
//     conditionals, the indexed MaybeTriggerConditionalOrders triggers the EXACT SAME
//     set of order IDs (set equality) as the reference full-scan oracle embedded here.
//
//  2. TestMaybeTrigger_BudgetAndDeferral — when crossed orders exceed
//     MaxConditionalTriggersPerBlock, exactly the budget triggers on the first block,
//     remainder trigger on subsequent blocks, total equals unbounded result.
//
//  3. TestMaybeTrigger_Determinism — repeated runs (same state, same prices) produce
//     identical triggered-id lists.
//
//  4. TestMaybeTrigger_NoPriceCrossingIsConstantWork — no orders are triggered when the
//     oracle price doesn't cross any trigger (no-price-crossing regression).
//
//  5. TestMaybeTrigger_OneCrossedOrderNoGlobalScan — a single crossed order triggers;
//     single-crossing regression.

import (
	"sort"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------
// Helpers shared by all Packet-2 tests
// ----------------------------------------------------------------
// trigIndexTestKeeper builds a context with BTC (clobPairId=0) and ETH (clobPairId=1).
// Reuses indexTestKeeperWithEth from conditional_order_index_test.go (same package).
func trigIndexTestKeeper(t *testing.T) keepertest.ClobKeepersTestContext {
	t.Helper()
	return indexTestKeeperWithEth(t)
}

// makeTrigTestConditionalOrder builds a conditional order for testing.
// dir selects LTE (take-profit BUY) or GTE (stop-loss BUY) direction.
// clientId must be unique per (subaccountId, clobPairId).
func makeTrigTestOrder(
	subaccountId satypes.SubaccountId,
	clientId uint32,
	clobPairId uint32,
	dir byte,
	triggerSubticks uint64,
) clobtypes.Order {
	var side clobtypes.Order_Side
	var condType clobtypes.Order_ConditionType
	if dir == clobkeeper.TriggerDirectionLTE {
		// take-profit BUY → triggers when oracle ≤ triggerSubticks
		side = clobtypes.Order_SIDE_BUY
		condType = clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT
	} else {
		// stop-loss BUY → triggers when oracle ≥ triggerSubticks
		side = clobtypes.Order_SIDE_BUY
		condType = clobtypes.Order_CONDITION_TYPE_STOP_LOSS
	}
	return clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: subaccountId,
			ClientId:     clientId,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   clobPairId,
		},
		Side:                            side,
		Quantums:                        1_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: condTestGTBT},
		ConditionType:                   condType,
		ConditionalOrderTriggerSubticks: triggerSubticks,
	}
}

// placeOrder places an order into the keeper and registers an expiry.
func placeOrder(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context, order clobtypes.Order) {
	t.Helper()
	expiry := ctx.BlockTime().Add(clobtypes.StatefulOrderTimeWindow - time.Hour)
	k.SetLongTermOrderPlacement(ctx, order, 1)
	k.AddStatefulOrderIdExpiration(ctx, expiry, order.OrderId)
}

// orderIdSet converts a slice of OrderIds to a set of state-key strings.
func orderIdSet(ids []clobtypes.OrderId) map[string]struct{} {
	s := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		s[string(id.ToStateKey())] = struct{}{}
	}
	return s
}

// ----------------------------------------------------------------
// Reference oracle: full-scan implementation (the old MaybeTriggerConditionalOrders logic)
// ----------------------------------------------------------------
// referenceTriggerFullScan is a faithful re-implementation of the OLD full-scan
// MaybeTriggerConditionalOrders.  It is used as the oracle in the equivalence test.
// It does NOT trigger orders in state — it only computes the set of order IDs that
// would be triggered, by reading the full untriggered set and applying PollTriggeredConditionalOrders.
// This mirrors the old code exactly:
//  1. GetAllUntriggeredConditionalOrders → OrganizeUntriggeredConditionalOrdersFromState
//  2. Sort keys.
//  3. For each clob pair: oracle, clamped-min, clamped-max via PollTriggeredConditionalOrders.
//
// NOTE: the clamped-trade-price section requires a private keeper method. In test contexts where
// trade prices are not seeded (GetTradePricesForPerpetual returns found=false), the section is
// naturally skipped. The test cases in this file do not seed trade prices, so the oracle faithfully
// matches the new implementation for oracle-price crossing equivalence.
//
// This function must be called BEFORE the new MaybeTriggerConditionalOrders since both
// consume the same state (the new code writes state; this function does not).
func referenceTriggerFullScanSet(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) map[string]struct{} {
	t.Helper()

	allOrders := k.GetAllUntriggeredConditionalOrders(ctx)
	clobPairToUntriggered := clobkeeper.OrganizeUntriggeredConditionalOrdersFromState(allOrders)

	// Sort keys (same deterministic order as the old implementation).
	sortedKeys := lib.GetSortedKeys[clobtypes.SortedClobPairId](clobPairToUntriggered)

	triggered := make(map[string]struct{})

	for _, clobPairId := range sortedKeys {
		untriggered := clobPairToUntriggered[clobPairId]
		clobPair, found := k.GetClobPair(ctx, clobPairId)
		if !found {
			continue
		}

		oraclePrice := k.GetOraclePriceSubticksRat(ctx, clobPair)

		// oracle price — this is always evaluated.
		for _, id := range untriggered.PollTriggeredConditionalOrders(oraclePrice) {
			triggered[string(id.ToStateKey())] = struct{}{}
		}

		// clamped trade prices — skipped when no trade prices have been recorded this block
		// (GetTradePricesForPerpetual returns found=false). All test cases here fall in this
		// category because they use simple keeper setups that don't execute matches.
		//
		// If a future test seeds trade prices, add a call to GetTradePricesForPerpetual here
		// and construct the clamped prices using the same logic as getClampedTradePricesForTriggering.
	}

	return triggered
}

// ----------------------------------------------------------------
// Test 1: Equivalence vs full-scan oracle
// ----------------------------------------------------------------
// TestMaybeTrigger_EquivalenceVsFullScan verifies that the indexed
// MaybeTriggerConditionalOrders triggers the same set of order IDs (set equality)
// as the reference full-scan oracle, across multiple clob pairs, both directions,
// and a range of oracle/trade prices — all within the budget.
func TestMaybeTrigger_EquivalenceVsFullScan(t *testing.T) {
	// Table-driven: each case describes a population of untriggered orders and
	// the oracle price used to trigger them.  We verify that the indexed path
	// produces the same triggered SET as the full-scan oracle.
	// BTC oracle price in the test setup is 50_000_000_000 subticks (5B price * 10^exponents).
	// LTE orders cross when oracle ≤ triggerSubticks → need triggerSubticks ≥ 50_000_000_000.
	// GTE orders cross when oracle ≥ triggerSubticks → need triggerSubticks ≤ 50_000_000_000.
	const btcOracleSubticks = uint64(50_000_000_000)
	// Above oracle: triggerSubticks > btcOracleSubticks → LTE crosses; GTE does not.
	const aboveOracle = uint64(60_000_000_000)
	// At oracle: triggerSubticks == btcOracleSubticks → both LTE and GTE cross (boundary, inclusive).
	const atOracle = btcOracleSubticks
	// Below oracle: triggerSubticks < btcOracleSubticks → GTE crosses; LTE does not.
	const belowOracle = uint64(40_000_000_000)

	cases := []struct {
		name  string
		setup func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context)
	}{
		{
			name: "empty_set_no_triggers",
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				// no orders — both oracle and new return empty
			},
		},
		{
			name: "all_lte_crossed_above_oracle",
			// LTE orders with triggerSubticks > oracle → all triggered.
			// oracle (50B) ≤ triggerSubticks (60B) → crossed.
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				for i := 0; i < 10; i++ {
					order := makeTrigTestOrder(
						constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionLTE, aboveOracle+uint64(i),
					)
					placeOrder(t, k, ctx, order)
				}
			},
		},
		{
			name: "all_lte_uncrossed_below_oracle",
			// LTE orders with triggerSubticks < oracle → none triggered.
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				for i := 0; i < 10; i++ {
					order := makeTrigTestOrder(
						constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionLTE, belowOracle-uint64(i+1),
					)
					placeOrder(t, k, ctx, order)
				}
			},
		},
		{
			name: "all_gte_crossed_below_oracle",
			// GTE orders with triggerSubticks < oracle → all triggered.
			// oracle (50B) ≥ triggerSubticks (40B) → crossed.
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				for i := 0; i < 10; i++ {
					order := makeTrigTestOrder(
						constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, belowOracle-uint64(i+1),
					)
					placeOrder(t, k, ctx, order)
				}
			},
		},
		{
			name: "mixed_partial_crossing",
			// Some LTE crossed (above oracle), some GTE crossed (below oracle).
			// LTE: at 60B (crossed) and 40B (not crossed, 40B < oracle 50B → LTE requires triggerSubticks ≥ oracle).
			// GTE: at 40B (crossed) and 60B (not crossed, 60B > oracle 50B → GTE requires triggerSubticks ≤ oracle).
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				// LTE: above oracle → crossed; below oracle → not crossed
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 0, 0, clobkeeper.TriggerDirectionLTE, aboveOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 1, 0, clobkeeper.TriggerDirectionLTE, belowOracle))
				// GTE: below oracle → crossed; above oracle → not crossed
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num1, 0, 0, clobkeeper.TriggerDirectionGTE, belowOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num1, 1, 0, clobkeeper.TriggerDirectionGTE, aboveOracle))
			},
		},
		{
			name: "multi_pair_mixed",
			// Orders on both BTC (pair 0) and ETH (pair 1).
			// ETH oracle is also set by CreateTestMarkets (market 1 price = ThreeBillion = 3B).
			// We use above/below oracle relative to BTC for pair 0.
			// For pair 1, use values around the ETH oracle — we only need equivalence, not specific triggering.
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				// BTC (pair 0): LTE above oracle → crossed; GTE below oracle → crossed
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 0, 0, clobkeeper.TriggerDirectionLTE, aboveOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 1, 0, clobkeeper.TriggerDirectionGTE, belowOracle))
				// BTC (pair 0): uncrossed orders
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 2, 0, clobkeeper.TriggerDirectionLTE, belowOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 3, 0, clobkeeper.TriggerDirectionGTE, aboveOracle))
				// ETH (pair 1): add some orders; equivalence test only cares about SET equality.
				placeOrder(t, k, ctx, makeTrigTestOrder(
					constants.Alice_Num1, 0, 1, clobkeeper.TriggerDirectionLTE, 1_000_000_000_000,
				))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num1, 1, 1, clobkeeper.TriggerDirectionGTE, 1))
			},
		},
		{
			name: "nothing_crossed_lte_far_below_oracle",
			// All LTE orders have triggerSubticks far below oracle → none triggered.
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				for i := 0; i < 20; i++ {
					order := makeTrigTestOrder(constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionLTE, 1)
					placeOrder(t, k, ctx, order)
				}
			},
		},
		{
			name: "boundary_price_exact_at_oracle",
			// Boundary: triggerSubticks exactly equal to oracle price.
			// LTE: oracle ≤ triggerSubticks → oracle == 50B ≤ 50B → triggered.
			// GTE: oracle ≥ triggerSubticks → oracle == 50B ≥ 50B → triggered.
			// LTE with triggerSubticks < oracle (1 below) → not triggered.
			// GTE with triggerSubticks > oracle (1 above) → not triggered.
			setup: func(t *testing.T, k *clobkeeper.Keeper, ctx sdk.Context) {
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 0, 0, clobkeeper.TriggerDirectionLTE, atOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 1, 0, clobkeeper.TriggerDirectionGTE, atOracle))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 2, 0, clobkeeper.TriggerDirectionLTE, atOracle-1))
				placeOrder(t, k, ctx, makeTrigTestOrder(constants.Alice_Num0, 3, 0, clobkeeper.TriggerDirectionGTE, atOracle+1))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ks := trigIndexTestKeeper(t)
			ctx := ks.Ctx
			k := ks.ClobKeeper

			tc.setup(t, k, ctx)

			// Compute the reference triggered SET using the old full-scan oracle.
			// Must be called before the new code runs (it reads state).
			refSet := referenceTriggerFullScanSet(t, k, ctx)

			// Enable the bounded (indexed) path so this equivalence test proves the FIXED path
			// triggers the same set as legacy (high budget → not budget-limited here).
			enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))

			// NOTE: The reference oracle consumed the in-memory UntriggeredConditionalOrders
			// struct via PollTriggeredConditionalOrders but did NOT write any state.
			// The new implementation DOES write state. So we run the new code now.
			triggered := k.MaybeTriggerConditionalOrders(ctx)
			newSet := orderIdSet(triggered)

			// Set equality: same order IDs triggered.
			require.Equal(t, refSet, newSet,
				"indexed MaybeTriggerConditionalOrders must trigger the same set of orders as the full-scan oracle")

			// All triggered orders must now be in the triggered state store.
			for _, id := range triggered {
				_, found := k.GetTriggeredConditionalOrderPlacement(ctx, id)
				require.True(t, found, "triggered order %v must be in triggered state", id)
			}
		})
	}
}

// ----------------------------------------------------------------
// Test 2: Budget and deferral
// ----------------------------------------------------------------
// TestMaybeTrigger_BudgetAndDeferral verifies that when more orders are crossed than
// MaxConditionalTriggersPerBlock, exactly the budget fires this block (nearest-crossing
// first by price/orderId order) and the remainder fire in subsequent blocks.
func TestMaybeTrigger_BudgetAndDeferral(t *testing.T) {
	// We need budget+extra orders all crossed at once.
	// To control the budget for the test without changing the package constant, we use
	// a large number of orders that exceeds MaxConditionalTriggersPerBlock.
	const budget = clobkeeper.MaxConditionalTriggersPerBlock
	// Place budget+50 LTE orders all with triggerSubticks well above the oracle price.
	const extra = 50
	total := budget + extra

	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// Place `total` GTE orders (stop-loss buys), all with triggerSubticks well below the oracle price.
	// GTE orders trigger when oracle ≥ triggerSubticks.  The BTC oracle in the test setup is
	// 50_000_000_000 subticks.  We use triggerSubticks = 1+i so all are crossed immediately.
	// Index key order: ascending by (triggerSubticks, orderId) for GTE direction.
	orders := make([]clobtypes.Order, total)
	for i := 0; i < total; i++ {
		orders[i] = makeTrigTestOrder(
			constants.Alice_Num0,
			uint32(i),
			0,
			clobkeeper.TriggerDirectionGTE,
			uint64(1+i), // unique, ascending; all far below oracle price
		)
		placeOrder(t, k, ctx, orders[i])
	}

	// Enable the bounded trigger path (the budget only applies when the config flag is on).
	enableTriggerConfig(t, k, ctx, uint32(budget))

	// Oracle price = 50_000_000_000 → all GTE orders are crossed (triggerSubticks ≤ oracle for all).
	// Pessimistic GTE rounding: floor(oracle) = 50_000_000_000, all orders with triggerSubticks ≤ 50B fire.
	// With the budget, exactly `budget` should fire on the first block.
	triggered1 := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, triggered1, budget,
		"first block: exactly MaxConditionalTriggersPerBlock orders should trigger")

	// Exactly `budget` orders fire this block and the remaining `extra` defer to the next block;
	// this test asserts the COUNTS and cross-block completeness (which specific orders fire first —
	// nearest-crossing-first — is covered by TestIterateCrossedConditionalOrders_NearestFirstOrdering).
	triggeredSet1 := orderIdSet(triggered1)
	for _, id := range triggered1 {
		_, found := k.GetTriggeredConditionalOrderPlacement(ctx, id)
		require.True(t, found, "triggered order %v must be in triggered store", id)
	}

	// Remainder: `extra` orders still untriggered.
	still := k.GetAllUntriggeredConditionalOrders(ctx)
	require.Len(t, still, extra,
		"after first block, exactly extra=%d orders should remain untriggered", extra)

	// Verify triggered orders are NOT in the untriggered store.
	for _, order := range still {
		_, inTriggered := triggeredSet1[string(order.OrderId.ToStateKey())]
		require.False(t, inTriggered, "order still untriggered should not be in triggered set")
	}

	// Second block: remaining `extra` orders should now trigger (all are crossed).
	triggered2 := k.MaybeTriggerConditionalOrders(ctx)
	require.Len(t, triggered2, extra,
		"second block: remaining extra=%d orders should trigger", extra)

	// Total triggered = budget + extra = total.
	allTriggered := append(triggered1, triggered2...)
	allSet := orderIdSet(allTriggered)
	require.Len(t, allSet, total,
		"total triggered across both blocks must equal all placed orders (no duplicates)")

	// No untriggered orders remain.
	finalRemaining := k.GetAllUntriggeredConditionalOrders(ctx)
	require.Empty(t, finalRemaining, "no untriggered orders should remain after both blocks")
}

// ----------------------------------------------------------------
// Test 3: Determinism
// ----------------------------------------------------------------
// TestMaybeTrigger_Determinism verifies that repeated calls to MaybeTriggerConditionalOrders
// with the same state produce identical triggered-id lists (no map-iteration nondeterminism).
func TestMaybeTrigger_Determinism(t *testing.T) {
	const orderCount = 50

	// Build fresh state twice and verify the triggered-id sequences are identical.
	// Because we can't "reset" a keeper between runs, we use two independent keepers
	// and place the same orders in each.
	setupAndTrigger := func() []clobtypes.OrderId {
		ks := newCondTestKeepers(t)
		ctx := ks.Ctx
		k := ks.ClobKeeper

		for i := 0; i < orderCount; i++ {
			dir := clobkeeper.TriggerDirectionLTE
			if i%2 == 0 {
				dir = clobkeeper.TriggerDirectionGTE
			}
			order := makeTrigTestOrder(
				constants.Alice_Num0,
				uint32(i),
				0,
				dir,
				uint64(100+i*10),
			)
			placeOrder(t, k, ctx, order)
		}

		// Enable the bounded (indexed) path.
		enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))
		// Oracle price = 300 → half of LTE and half of GTE are crossed.
		return k.MaybeTriggerConditionalOrders(ctx)
	}

	run1 := setupAndTrigger()
	run2 := setupAndTrigger()

	require.Equal(t, run1, run2, "two identical runs must produce identical triggered-id sequences")
}

// ----------------------------------------------------------------
// Test 4: No price crossing = O(1) work (no-price-crossing regression)
// ----------------------------------------------------------------
// TestMaybeTrigger_NoPriceCrossingIsConstantWork verifies that when no conditional orders
// are crossed by the oracle price, MaybeTriggerConditionalOrders returns empty and no
// orders are triggered — regardless of how many untriggered orders exist.
// This is the no-price-crossing regression guard: after the fix, per-block cost does NOT
// scale with the global untriggered count when no orders are crossed.
func TestMaybeTrigger_NoPriceCrossingIsConstantWork(t *testing.T) {
	const n = 5000
	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// All LTE orders with triggerSubticks = 1 (far below any oracle price).
	for i := 0; i < n; i++ {
		order := makeTrigTestOrder(constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionLTE, 1)
		placeOrder(t, k, ctx, order)
	}

	// Enable the bounded (indexed) path so this exercises the fixed O(crossed) scan.
	enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))

	// A malformed sentinel in the authoritative SO/U store makes any broad/full-store scan panic
	// during unmarshal. The indexed no-crossing path must never read it.
	sourceStore := k.GetUntriggeredConditionalOrderPlacementStore(ctx)
	poisonKey := []byte("no-crossing-global-scan-sentinel")
	sourceStore.Set(poisonKey, []byte{0xff})

	for block := 1; block <= 3; block++ {
		var triggered []clobtypes.OrderId
		require.NotPanics(t, func() {
			triggered = k.MaybeTriggerConditionalOrders(ctx)
		}, "bounded no-crossing evaluation must not scan the authoritative order store")
		require.Empty(t, triggered,
			"block=%d: no orders should be triggered when price doesn't cross any trigger", block)
	}
	sourceStore.Delete(poisonKey)
	require.Len(t, k.GetAllUntriggeredConditionalOrders(ctx), n)
}

// ----------------------------------------------------------------
// Test 4e: Per-block budget is fair-shared across markets
// ----------------------------------------------------------------
// TestMaybeTrigger_FairShareAcrossMarkets verifies that a low-id market cannot consume the whole
// chain-wide trigger budget and starve higher-id markets. Before the fix, a single shared budget was
// drained in ascending pair-id order, so pair 0 (with many crossed orders) could take the entire
// budget every block while pair 1 stayed deferred.
func TestMaybeTrigger_FairShareAcrossMarkets(t *testing.T) {
	ks := trigIndexTestKeeper(t) // BTC = pair 0, ETH = pair 1
	k := ks.ClobKeeper
	ctx := ks.Ctx

	// Many crossed GTE orders on BOTH pairs (triggerSubticks well below each oracle price).
	const perPair = 10
	for i := 0; i < perPair; i++ {
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, uint64(1+i),
		))
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num1, uint32(i), 1, clobkeeper.TriggerDirectionGTE, uint64(1+i),
		))
	}

	// Budget of 4 across 2 active pairs → each pair gets a fair share of 2 this block.
	enableTriggerConfig(t, k, ctx, 4)
	triggered := k.MaybeTriggerConditionalOrders(ctx)

	perPairCount := map[uint32]int{}
	for _, id := range triggered {
		perPairCount[id.ClobPairId]++
	}
	require.Len(t, triggered, 4, "exactly the budget (4) should trigger this block")
	require.Equal(t, 2, perPairCount[0], "BTC (pair 0) must not consume the whole budget")
	require.Equal(t, 2, perPairCount[1], "ETH (pair 1) must get its fair share, not be starved")
}

// TestMaybeTrigger_FairShareRotatesWhenBudgetIsSmallerThanMarkets verifies cross-block fairness.
// With a budget of one and two continuously backlogged markets, a fixed ascending-id scheduler
// would serve pair 0 forever. The persisted cursor must alternate which pair starts each block.
func TestMaybeTrigger_FairShareRotatesWhenBudgetIsSmallerThanMarkets(t *testing.T) {
	ks := trigIndexTestKeeper(t)
	k := ks.ClobKeeper
	ctx := ks.Ctx

	const perPair = 3
	for i := 0; i < perPair; i++ {
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, uint64(1+i),
		))
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num1, uint32(i), 1, clobkeeper.TriggerDirectionGTE, uint64(1+i),
		))
	}

	enableTriggerConfig(t, k, ctx, 1)
	expectedPairOrder := []uint32{0, 1, 0, 1}
	for block, expectedPair := range expectedPairOrder {
		triggered := k.MaybeTriggerConditionalOrders(ctx)
		require.Len(t, triggered, 1, "block %d must consume exactly one trigger slot", block+1)
		require.Equal(t, expectedPair, triggered[0].ClobPairId,
			"the scheduling cursor must rotate across backlogged markets")
	}
}

// ----------------------------------------------------------------
// Test 4a: GTE crossed orders are visited nearest-first under budget
// ----------------------------------------------------------------
// TestIterateCrossedConditionalOrders_NearestFirstOrdering verifies that both trigger directions
// visit crossed orders NEAREST-CROSSING first, so a per-block budget triggers the orders closest to
// the current price before farther ones. The pre-fix GTE scan iterated ascending over [0, price],
// which is FARTHEST-first (deep out-of-the-money stops fired before near-price ones); the fix
// iterates GTE descending. LTE was and remains ascending (nearest-first).
func TestIterateCrossedConditionalOrders_NearestFirstOrdering(t *testing.T) {
	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper
	// Enable so the placement hooks maintain the consensus-gated index.
	enableTriggerConfig(t, k, ctx, clobkeeper.MaxConditionalTriggersPerBlock)

	// GTE (stop-loss buy) orders at ascending trigger subticks; clientId i -> subticks.
	gteSubticks := []uint64{50, 100, 200, 300}
	gteByClient := make(map[uint32]uint64, len(gteSubticks))
	for i, sub := range gteSubticks {
		o := makeTrigTestOrder(constants.Alice_Num0, uint32(i), 0, clobkeeper.TriggerDirectionGTE, sub)
		k.SetLongTermOrderPlacement(ctx, o, 1)
		gteByClient[uint32(i)] = sub
	}

	// Scan GTE at price 300: all crossed. Visitation order must be DESCENDING (nearest-first).
	var gteVisited []uint64
	k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionGTE, 300,
		func(id clobtypes.OrderId) bool {
			gteVisited = append(gteVisited, gteByClient[id.ClientId])
			return true
		})
	require.Equal(t, []uint64{300, 200, 100, 50}, gteVisited,
		"GTE crossed orders must be visited nearest-first (highest subticks ≤ price first)")

	// LTE (take-profit buy) orders at ascending trigger subticks on a different subaccount.
	lteSubticks := []uint64{10, 20, 25, 30}
	lteByClient := make(map[uint32]uint64, len(lteSubticks))
	for i, sub := range lteSubticks {
		o := makeTrigTestOrder(constants.Alice_Num1, uint32(i), 0, clobkeeper.TriggerDirectionLTE, sub)
		k.SetLongTermOrderPlacement(ctx, o, 1)
		lteByClient[uint32(i)] = sub
	}

	// Scan LTE at price 5: all crossed. Visitation order must be ASCENDING (nearest-first).
	var lteVisited []uint64
	k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionLTE, 5,
		func(id clobtypes.OrderId) bool {
			lteVisited = append(lteVisited, lteByClient[id.ClientId])
			return true
		})
	require.Equal(t, []uint64{10, 20, 25, 30}, lteVisited,
		"LTE crossed orders must be visited nearest-first (lowest subticks ≥ price first)")
}

// ----------------------------------------------------------------
// Test 4d: Equal-price ties resolve by placement time, not client id
// ----------------------------------------------------------------
// TestMaybeTrigger_EqualPriceTieBreaksByPlacementTime verifies that when two orders share the same
// trigger subticks, the OLDER order (earlier placement) triggers first under a per-block budget,
// regardless of client id. Before the fix, the index tie-break was the raw orderId (which embeds
// the client-chosen ClientId), so a newer order with a lower client id could jump an older one.
func TestMaybeTrigger_EqualPriceTieBreaksByPlacementTime(t *testing.T) {
	const sameSubticksLTE = uint64(60_000_000_000) // above oracle (50B) -> LTE crosses
	const sameSubticksGTE = uint64(40_000_000_000) // below oracle (50B) -> GTE crosses

	// Each case places an OLDER order with a HIGHER client id first, then a NEWER order with a
	// LOWER client id. Time priority must trigger the older (high-client-id) order first; a
	// client-id tie-break would (incorrectly) pick the newer low-client-id order.
	cases := []struct {
		name string
		dir  byte
		sub  uint64
	}{
		{"lte", clobkeeper.TriggerDirectionLTE, sameSubticksLTE},
		{"gte", clobkeeper.TriggerDirectionGTE, sameSubticksGTE},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ks := trigIndexTestKeeper(t)
			k := ks.ClobKeeper
			ctx := ks.Ctx

			older := makeTrigTestOrder(constants.Alice_Num0, 99 /* high client id */, 0, tc.dir, tc.sub)
			newer := makeTrigTestOrder(constants.Alice_Num0, 1 /* low client id */, 0, tc.dir, tc.sub)
			placeOrder(t, k, ctx, older) // placed first -> lower placement sequence
			placeOrder(t, k, ctx, newer) // placed second -> higher placement sequence

			// Budget of 1 so only the highest-priority (oldest) order triggers this block.
			enableTriggerConfig(t, k, ctx, 1)

			triggered := k.MaybeTriggerConditionalOrders(ctx)
			require.Equal(t, []clobtypes.OrderId{older.OrderId}, triggered,
				"the older order must trigger first even though it has the higher client id")
		})
	}
}

// ----------------------------------------------------------------
// Test 4b: Zero oracle price does not panic the EndBlocker
// ----------------------------------------------------------------
// TestMaybeTrigger_ZeroOraclePriceDoesNotPanic verifies that the enabled bounded path does not
// panic when a clob pair's market oracle price is still 0 (e.g. a freshly listed market that has
// not yet received a valid non-zero price update). The pre-fix bounded path fetched the oracle
// price for every clob pair unconditionally and panicked (ErrZeroPriceForOracle); the legacy path
// only ever visited pairs that had resting orders, so this all-pairs visit was a new exposure.
func TestMaybeTrigger_ZeroOraclePriceDoesNotPanic(t *testing.T) {
	const btcMarketId = uint32(0)

	setBtcPrice := func(t *testing.T, ks keepertest.ClobKeepersTestContext, price uint64) {
		t.Helper()
		require.NoError(t, ks.PricesKeeper.UpdateMarketPrices(ks.Ctx,
			[]*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				{MarketId: btcMarketId, Price: price},
			},
		))
	}

	t.Run("empty_price_zero_pair_is_skipped", func(t *testing.T) {
		ks := trigIndexTestKeeper(t)
		k := ks.ClobKeeper
		ctx := ks.Ctx

		// BTC market drops to price 0; no conditional orders rest on any pair.
		setBtcPrice(t, ks, 0)
		enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))

		require.NotPanics(t, func() {
			require.Empty(t, k.MaybeTriggerConditionalOrders(ctx))
		}, "a price-0 pair with no resting conditionals must be skipped, not panic")
	})

	t.Run("price_zero_pair_with_resting_conditional_defers_then_triggers", func(t *testing.T) {
		ks := trigIndexTestKeeper(t)
		k := ks.ClobKeeper
		ctx := ks.Ctx

		// A GTE (stop-loss buy) order that will cross once BTC is at its normal 50B oracle price.
		order := makeTrigTestOrder(constants.Alice_Num0, 0, 0, clobkeeper.TriggerDirectionGTE, 40_000_000_000)
		placeOrder(t, k, ctx, order)

		// Enable while BTC is at price 0. The pair now HAS an index entry, so it is not skipped by
		// the presence check — the non-panicking oracle guard must handle the zero price.
		setBtcPrice(t, ks, 0)
		enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))

		require.NotPanics(t, func() {
			require.Empty(t, k.MaybeTriggerConditionalOrders(ctx),
				"at price 0 the crossing cannot be evaluated, so nothing triggers")
		}, "a price-0 pair WITH a resting conditional must not panic")

		// The order must remain untriggered (skipped, not lost).
		require.Len(t, k.GetAllUntriggeredConditionalOrders(ctx), 1)

		// Once a valid non-zero price arrives, the deferred order triggers normally — proving the
		// zero-price skip is temporary, not a black hole.
		setBtcPrice(t, ks, 50_000_000_000)
		triggered := k.MaybeTriggerConditionalOrders(ctx)
		require.Len(t, triggered, 1, "order must trigger once a non-zero price is available")
		require.Equal(t, order.OrderId, triggered[0])
	})
}

// ----------------------------------------------------------------
// Test 4c: Expired-but-un-drained orders are not triggered
// ----------------------------------------------------------------
// TestMaybeTrigger_ExpiredOrderNotTriggered verifies that a conditional order past its
// GoodTilBlockTime is not triggered even though it is still resting in state. During traversal its
// stale trigger-price index entry is removed, while the budgeted expiry prune retains ownership of
// deleting the authoritative placement.
func TestMaybeTrigger_ExpiredOrderNotTriggered(t *testing.T) {
	// GTE (stop-loss buy) crossing at the BTC oracle price (50B); triggerSubticks 40B < oracle.
	newSetup := func(t *testing.T) (*clobkeeper.Keeper, sdk.Context, clobtypes.Order) {
		ks := trigIndexTestKeeper(t)
		k := ks.ClobKeeper
		ctx := ks.Ctx
		order := makeTrigTestOrder(constants.Alice_Num0, 0, 0, clobkeeper.TriggerDirectionGTE, 40_000_000_000)
		placeOrder(t, k, ctx, order)
		enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))
		return k, ctx, order
	}

	t.Run("expired_order_is_not_triggered", func(t *testing.T) {
		k, ctx, order := newSetup(t)
		// Advance block time PAST the order's GoodTilBlockTime without running the (budgeted) expiry
		// prune — the order is still resting in state and in the index.
		expiredCtx := ctx.WithBlockTime(time.Unix(int64(condTestGTBT)+1, 0))
		require.Empty(t, k.MaybeTriggerConditionalOrders(expiredCtx),
			"an order past its GoodTilBlockTime must not trigger even though it is still resting")
		// It remains untriggered (left for the expiry prune), not moved to triggered state.
		require.Len(t, k.GetAllUntriggeredConditionalOrders(expiredCtx), 1)
		_, foundTriggered := k.GetTriggeredConditionalOrderPlacement(expiredCtx, order.OrderId)
		require.False(t, foundTriggered, "expired order must not be moved to triggered state")
		remainingIndexEntries := 0
		k.IterateCrossedConditionalOrders(
			expiredCtx,
			0,
			clobkeeper.TriggerDirectionGTE,
			50_000_000_000,
			func(clobtypes.OrderId) bool {
				remainingIndexEntries++
				return true
			},
		)
		require.Zero(t, remainingIndexEntries, "expired index entry must be removed during traversal")
	})

	t.Run("unexpired_order_triggers", func(t *testing.T) {
		k, ctx, order := newSetup(t)
		// Block time before the GoodTilBlockTime → the same crossed order triggers normally,
		// proving the GTBT guard blocks only expired orders.
		validCtx := ctx.WithBlockTime(time.Unix(int64(condTestGTBT)-1, 0))
		require.Equal(t, []clobtypes.OrderId{order.OrderId}, k.MaybeTriggerConditionalOrders(validCtx))
	})
}

// TestMaybeTrigger_ExpiredCleanupConsumesBudget verifies an expired backlog cannot turn the
// bounded trigger path back into an unbounded scan. Each stale entry removed consumes one unit of
// the same per-block work budget as a trigger.
func TestMaybeTrigger_ExpiredCleanupConsumesBudget(t *testing.T) {
	const (
		budget = 3
		total  = 7
	)

	ks := trigIndexTestKeeper(t)
	k := ks.ClobKeeper
	ctx := ks.Ctx
	for i := 0; i < total; i++ {
		placeOrder(t, k, ctx, makeTrigTestOrder(
			constants.Alice_Num0,
			uint32(i),
			0,
			clobkeeper.TriggerDirectionGTE,
			40_000_000_000+uint64(i),
		))
	}
	enableTriggerConfig(t, k, ctx, budget)
	expiredCtx := ctx.WithBlockTime(time.Unix(int64(condTestGTBT)+1, 0))

	countIndexEntries := func() int {
		count := 0
		k.IterateCrossedConditionalOrders(
			expiredCtx,
			0,
			clobkeeper.TriggerDirectionGTE,
			50_000_000_000,
			func(clobtypes.OrderId) bool {
				count++
				return true
			},
		)
		return count
	}

	require.Empty(t, k.MaybeTriggerConditionalOrders(expiredCtx))
	require.Equal(t, total-budget, countIndexEntries())
	require.Len(t, k.GetAllUntriggeredConditionalOrders(expiredCtx), total,
		"trigger traversal must not delete authoritative expiry state")

	require.Empty(t, k.MaybeTriggerConditionalOrders(expiredCtx))
	require.Equal(t, total-(2*budget), countIndexEntries())

	require.Empty(t, k.MaybeTriggerConditionalOrders(expiredCtx))
	require.Zero(t, countIndexEntries())
}

// ----------------------------------------------------------------
// Test 5: One crossed order triggers without a global scan (single-crossing regression)
// ----------------------------------------------------------------
// TestMaybeTrigger_OneCrossedOrderNoBroadScan verifies that exactly one order triggers
// when only one order is crossed, even when a large number of uncrossed orders exist.
// This is the single-crossing regression guard: after the fix, triggering 1 order
// does NOT require reading the entire global untriggered set.
func TestMaybeTrigger_OneCrossedOrderNoBroadScan(t *testing.T) {
	const n = 5000
	require.GreaterOrEqual(t, n, 2)

	ks := newCondTestKeepers(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// n-1 LTE orders with triggerSubticks = 1 (far below oracle price → not crossed).
	// 1 GTE order with triggerSubticks = 1 (well below oracle price → crossed).
	var crossedOrder clobtypes.Order
	for i := 0; i < n; i++ {
		if i == n-1 {
			// This one will be crossed: GTE with triggerSubticks = 1 < oracle price.
			crossedOrder = makeTrigTestOrder(constants.Alice_Num0, uint32(n), 0, clobkeeper.TriggerDirectionGTE, 1)
			placeOrder(t, k, ctx, crossedOrder)
		} else {
			// LTE with triggerSubticks = 1; LTE oracle price = ceil(oracle) so at oracle
			// price >> 1, LTE orders with triggerSubticks = 1 are NOT crossed.
			order := makeTrigTestOrder(
				satypes.SubaccountId{Owner: condTestOwner, Number: uint32(i / 10)},
				uint32(i),
				0,
				clobkeeper.TriggerDirectionLTE,
				1,
			)
			placeOrder(t, k, ctx, order)
		}
	}

	// Enable the bounded (indexed) path so this exercises the fixed O(crossed) scan.
	enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))
	sourceStore := k.GetUntriggeredConditionalOrderPlacementStore(ctx)
	poisonKey := []byte("single-crossing-global-scan-sentinel")
	sourceStore.Set(poisonKey, []byte{0xff})

	var triggered []clobtypes.OrderId
	require.NotPanics(t, func() {
		triggered = k.MaybeTriggerConditionalOrders(ctx)
	}, "one indexed crossing must not scan the authoritative order store")
	require.Len(t, triggered, 1, "exactly one order should trigger")
	require.Equal(t, crossedOrder.OrderId, triggered[0], "the GTE order at triggerSubticks=1 should be triggered")

	_, found := k.GetTriggeredConditionalOrderPlacement(ctx, crossedOrder.OrderId)
	require.True(t, found, "triggered order must be in triggered state store")

	sourceStore.Delete(poisonKey)
	still := k.GetAllUntriggeredConditionalOrders(ctx)
	require.Len(t, still, n-1, "n-1 untriggered orders should remain")
}

// ----------------------------------------------------------------
// Helpers needed from outside the keeper package
// ----------------------------------------------------------------
// trigEquivTestKeeper is a simpler single-BTC-pair keeper (same as newCondTestKeepers).
func trigEquivTestKeeperBtcOnly(t *testing.T) keepertest.ClobKeepersTestContext {
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

// ----------------------------------------------------------------
// Test 6: Randomized large equivalence (property test)
// ----------------------------------------------------------------
// TestMaybeTrigger_RandomizedEquivalence runs a larger equivalence check with a more
// varied spread of triggerSubticks and oracle prices.
func TestMaybeTrigger_RandomizedEquivalence(t *testing.T) {
	ks := trigEquivTestKeeperBtcOnly(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// Place orders with a spread of triggerSubticks covering both directions.
	// We stay well within MaxConditionalTriggersPerBlock for equivalence comparison.
	const numOrders = 200
	for i := 0; i < numOrders; i++ {
		var dir byte
		var sub uint64
		if i%2 == 0 {
			dir = clobkeeper.TriggerDirectionLTE
			sub = uint64(50 + i*3) // spread: 50, 56, 62, ...
		} else {
			dir = clobkeeper.TriggerDirectionGTE
			sub = uint64(50 + i*3)
		}

		var sa satypes.SubaccountId
		if i < 100 {
			sa = constants.Alice_Num0
		} else {
			sa = constants.Alice_Num1
		}

		order := makeTrigTestOrder(sa, uint32(i), 0, dir, sub)
		placeOrder(t, k, ctx, order)
	}

	// Reference full-scan set (read state, no writes).
	refSet := referenceTriggerFullScanSet(t, k, ctx)

	// Enable the bounded (indexed) path.
	enableTriggerConfig(t, k, ctx, uint32(clobkeeper.MaxConditionalTriggersPerBlock))

	// New indexed trigger (writes state).
	triggered := k.MaybeTriggerConditionalOrders(ctx)
	newSet := orderIdSet(triggered)

	require.Equal(t, refSet, newSet,
		"randomized equivalence: indexed and full-scan must trigger identical order sets")

	// Verify none of the triggered ids remain untriggered.
	still := k.GetAllUntriggeredConditionalOrders(ctx)
	stillSet := make(map[string]struct{}, len(still))
	for _, o := range still {
		stillSet[string(o.OrderId.ToStateKey())] = struct{}{}
	}
	for _, id := range triggered {
		_, inStill := stillSet[string(id.ToStateKey())]
		require.False(t, inStill, "triggered order %v must not remain in untriggered store", id)
	}
}

// ----------------------------------------------------------------
// Compile-time check: MaxConditionalTriggersPerBlock is exported and positive.
// ----------------------------------------------------------------
func TestMaxConditionalTriggersPerBlock_IsPositive(t *testing.T) {
	require.Positive(t, clobkeeper.MaxConditionalTriggersPerBlock,
		"MaxConditionalTriggersPerBlock must be a positive constant")
}

// ----------------------------------------------------------------
// Verify that the in-package sort/dedup helpers still work.
// ----------------------------------------------------------------
func TestOrderIdSet_NoDuplicates(t *testing.T) {
	ids := []clobtypes.OrderId{
		{SubaccountId: constants.Alice_Num0, ClientId: 0, OrderFlags: clobtypes.OrderIdFlags_Conditional, ClobPairId: 0},
		{SubaccountId: constants.Alice_Num0, ClientId: 1, OrderFlags: clobtypes.OrderIdFlags_Conditional, ClobPairId: 0},
		{SubaccountId: constants.Alice_Num0, ClientId: 0, OrderFlags: clobtypes.OrderIdFlags_Conditional, ClobPairId: 0},
	}
	s := orderIdSet(ids)
	require.Len(t, s, 2, "set must de-duplicate")

	// Also verify SortedOrders helper works for downstream consumer compatibility.
	sort.Sort(clobtypes.SortedOrders([]clobtypes.OrderId{})) // compile-check
}
