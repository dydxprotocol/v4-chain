//go:build conditional_order_e2e

// These full-app e2e tests exercise the index-backed conditional-order triggering path end to end
// (signed-tx admission, bounded triggering, expiry/backfill). They run full CheckTx→DeliverTx→
// AdvanceToBlock cycles which are expensive under -race, and the x/clob/e2e package already sits
// near the 20m -race package budget, so they are excluded from the default `go test ./...` sweep
// and gated behind the `conditional_order_e2e` build tag. Run them during validation with:
//
//	go test -tags conditional_order_e2e ./x/clob/e2e/ -run TestConditionalTrigger_
//
// The throughput-measurement tests additionally require BENCH_DYDX_UNTRIG_ORDERS=<N> to run.
// Correctness of the underlying logic is covered by the (always-on) keeper unit tests.
package clob_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// condE2EBenchN returns the order count for the opt-in throughput-measurement tests. These tests
// seed thousands of orders and run full AdvanceToBlock cycles, which is far too slow for the
// standard unit / -race / container suites, so they run ONLY when BENCH_DYDX_UNTRIG_ORDERS is set
// (they are timing-report tools, not pass/fail assertions). When it is unset, the caller skips.
func condE2EBenchN(t *testing.T) int {
	t.Helper()
	s := strings.TrimSpace(os.Getenv("BENCH_DYDX_UNTRIG_ORDERS"))
	if s == "" {
		t.Skip("measurement test; set BENCH_DYDX_UNTRIG_ORDERS=<N> to run")
	}
	v, err := strconv.Atoi(s)
	require.NoError(t, err)
	require.Positive(t, v)
	return v
}

func condE2EGenesis() types.GenesisDoc {
	genesis := testapp.DefaultGenesis()
	testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(g *satypes.GenesisState) {
		g.Subaccounts = []satypes.Subaccount{
			constants.Alice_Num0_100_000USD,
		}
	})
	testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(g *pricestypes.GenesisState) {
		*g = constants.TestPricesGenesisState
	})
	testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(g *perptypes.GenesisState) {
		g.Params = constants.PerpetualsGenesisParams
		g.LiquidityTiers = constants.LiquidityTiers
		g.Perpetuals = []perptypes.Perpetual{constants.BtcUsd_20PercentInitial_10PercentMaintenance}
	})
	testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(g *clobtypes.GenesisState) {
		g.ClobPairs = []clobtypes.ClobPair{constants.ClobPair_Btc}
		g.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
		g.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{
			MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
				{NumBlocks: 100, Limit: 50},
			},
		}
	})
	testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(g *feetiertypes.GenesisState) {
		g.Params = constants.PerpetualFeeParamsNoFee
	})
	return genesis
}

func e2eSeedConditionalOrder(clientID uint32) clobtypes.Order {
	return clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     clientID,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   constants.ClobPair_Btc.Id,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        1_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 1_900_000_000},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 1,
	}
}

func e2eSignedConditionalOrder(clientID uint32, gtbt uint32) clobtypes.Order {
	return clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     clientID,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   constants.ClobPair_Btc.Id,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: gtbt},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 49_999_000_000,
	}
}

func TestConditionalTrigger_AppAdvanceMeasure(t *testing.T) {
	n := condE2EBenchN(t)
	tApp := testapp.NewTestAppBuilder(t).
		WithGenesisDocFn(condE2EGenesis).
		WithNonDeterminismChecksEnabled(false).
		Build()

	seedCtx := tApp.App.NewUncachedContext(false, tmproto.Header{})
	for i := 0; i < n; i++ {
		order := e2eSeedConditionalOrder(uint32(i))
		tApp.App.ClobKeeper.SetLongTermOrderPlacement(seedCtx, order, 1)
		tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(seedCtx, time.Unix(1_900_000_000, 0), order.OrderId)
	}

	_ = tApp.InitChain()
	start := time.Now()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	elapsed := time.Since(start)

	untriggered := tApp.App.ClobKeeper.GetAllUntriggeredConditionalOrders(ctx)
	require.Len(t, untriggered, n)
	fmt.Printf(
		"app-block: N=%d AdvanceToBlock(EndBlocker+Commit) elapsed=%.4fs (%.2fus/order)\n",
		n,
		elapsed.Seconds(),
		float64(elapsed.Microseconds())/float64(n),
	)
}

func e2eCheckAndCollectSignedOrders(
	t *testing.T,
	tApp *testapp.TestApp,
	ctx sdk.Context,
	orders []clobtypes.Order,
) [][]byte {
	t.Helper()

	deliverTxs := make([][]byte, 0, len(orders)+2)
	deliverTxs = append(deliverTxs, constants.ValidEmptyMsgProposedOperationsTxBytes)
	for _, order := range orders {
		for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
			ctx,
			tApp.App,
			*clobtypes.NewMsgPlaceOrder(order),
		) {
			resp := tApp.CheckTx(checkTx)
			require.Conditionf(t, resp.IsOK, "signed conditional order CheckTx must pass. resp=%+v", resp)
			deliverTxs = append(deliverTxs, checkTx.Tx)
		}
	}
	deliverTxs = append(deliverTxs, constants.EmptyMsgAddPremiumVotesTxBytes)
	return deliverTxs
}

func TestConditionalTrigger_AdmissionE2E(t *testing.T) {
	const k = 40
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(condE2EGenesis).Build()
	ctx := tApp.InitChain()
	gtbt := uint32(ctx.BlockTime().Add(10 * time.Minute).Unix())

	orders := make([]clobtypes.Order, 0, k)
	for i := 0; i < k; i++ {
		orders = append(orders, e2eSignedConditionalOrder(uint32(i), gtbt))
	}

	deliverTxs := e2eCheckAndCollectSignedOrders(t, tApp, ctx, orders)
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{DeliverTxsOverride: deliverTxs})

	untriggered := tApp.App.ClobKeeper.GetAllUntriggeredConditionalOrders(ctx)
	require.Len(t, untriggered, k)
	for _, order := range orders {
		require.False(t, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, order.OrderId))
	}
	fmt.Printf("admit-e2e: %d/%d signed conditionals admitted and resting untriggered at block 2\n", k, k)

	for blk := uint32(3); blk <= 6; blk++ {
		ctx = tApp.AdvanceToBlock(blk, testapp.AdvanceToBlockOptions{})
		still := tApp.App.ClobKeeper.GetAllUntriggeredConditionalOrders(ctx)
		require.Len(t, still, k)
		fmt.Printf("admit-e2e: block %d: untriggered resting set still=%d\n", blk, len(still))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// conditional-order EndBlocker work Phase B validation tests (committed-state; governance flag ON)
// ─────────────────────────────────────────────────────────────────────────────

// TestConditionalTrigger_AppAdvanceBoundedMeasure measures the steady-state bounded trigger path.
//
// The bounded index path is always active, so per-block EndBlocker cost is O(0 crossed entries),
// independent of N. The timing must be flat regardless of how large N is.
//
// Committed-state note: seeding happens via NewUncachedContext (same pattern as Phase A) so the
// orders are written directly to the commit store (IAVL-backed), making the iterator cost
// realistic. The trigger config is set and its incremental activation passes are completed on the
// uncached context BEFORE InitChain. After InitChain + Commit, both stores are committed and the
// timed block measures the steady-state bounded path, excluding activation work.
func TestConditionalTrigger_AppAdvanceBoundedMeasure(t *testing.T) {
	n := condE2EBenchN(t)
	tApp := testapp.NewTestAppBuilder(t).
		WithGenesisDocFn(condE2EGenesis).
		WithNonDeterminismChecksEnabled(false).
		Build()

	seedCtx := tApp.App.NewUncachedContext(false, tmproto.Header{})
	for i := 0; i < n; i++ {
		order := e2eSeedConditionalOrder(uint32(i))
		tApp.App.ClobKeeper.SetLongTermOrderPlacement(seedCtx, order, 1)
		tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(seedCtx, time.Unix(1_900_000_000, 0), order.OrderId)
	}

	// The bounded trigger path is always active; the trigger-price index was maintained
	// incrementally as the orders were seeded via SetLongTermOrderPlacement. Set an explicit budget.
	tApp.App.ClobKeeper.SetConditionalOrderTriggerConfig(seedCtx, clobkeeper.ConditionalOrderTriggerConfig{
		MaxTriggersPerBlock: clobkeeper.MaxConditionalTriggersPerBlock,
	})

	_ = tApp.InitChain()
	start := time.Now()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	elapsed := time.Since(start)

	untriggered := tApp.App.ClobKeeper.GetAllUntriggeredConditionalOrders(ctx)
	// All N orders remain untriggered: their trigger price is 1 subtick (far below oracle ~50B),
	// so the bounded LTE scan finds zero crossings and the untriggered set is unchanged.
	require.Len(t, untriggered, n)

	fmt.Printf(
		"bounded: flag=ON N=%d AdvanceToBlock(EndBlocker+Commit) elapsed=%.4fs "+
			"(%.2fus/order) - per-block cost bounded by crossings, not N\n",
		n,
		elapsed.Seconds(),
		float64(elapsed.Microseconds())/float64(n),
	)
}

// TestConditionalTrigger_BoundedLegitTrigger is the Phase B negative test (legitimate-behavior
// check): with the governance flag ON, a conditional order whose trigger price is crossed by
// the oracle price MUST still trigger.
//
// Setup: N-1 non-crossing orders with trigger=1 (never cross), plus 1 legitimate TAKE_PROFIT BUY with
// trigger=51_000_000_000 subticks.  Oracle price at genesis is ~50B subticks (BTC $50k), so
// oracle (50B) <= trigger (51B) → the order is in the LTE bucket and IS crossed.
//
// Flag ON starts incremental activation, which this test completes before AdvanceToBlock(2).
// Exactly 1 order must trigger; the remaining N-1 non-crossing orders must stay untriggered.
func TestConditionalTrigger_BoundedLegitTrigger(t *testing.T) {
	const nNonCrossing = 200 // small enough for a fast deterministic test
	tApp := testapp.NewTestAppBuilder(t).
		WithGenesisDocFn(condE2EGenesis).
		WithNonDeterminismChecksEnabled(false).
		Build()

	seedCtx := tApp.App.NewUncachedContext(false, tmproto.Header{})

	// Seed nNonCrossing non-crossing orders with trigger=1 (never cross for LTE with oracle ~50B).
	for i := 0; i < nNonCrossing; i++ {
		order := e2eSeedConditionalOrder(uint32(i))
		tApp.App.ClobKeeper.SetLongTermOrderPlacement(seedCtx, order, 1)
		tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(seedCtx, time.Unix(1_900_000_000, 0), order.OrderId)
	}

	// Legitimate order: TAKE_PROFIT BUY with trigger=51_000_000_000 subticks.
	// Oracle (~50B) <= trigger (51B) → CROSSED → should trigger at block 2.
	legitClientID := uint32(nNonCrossing) // distinct from non-crossing IDs
	legitOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     legitClientID,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   constants.ClobPair_Btc.Id,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        1_000_000,
		Subticks:                        55_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 1_900_000_000},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 51_000_000_000, // > oracle ~50B → LTE → crosses
	}
	tApp.App.ClobKeeper.SetLongTermOrderPlacement(seedCtx, legitOrder, 1)
	tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(seedCtx, time.Unix(1_900_000_000, 0), legitOrder.OrderId)

	// The bounded trigger path is always active; the trigger-price index was maintained
	// incrementally as the orders were seeded via SetLongTermOrderPlacement. Set an explicit budget.
	tApp.App.ClobKeeper.SetConditionalOrderTriggerConfig(seedCtx, clobkeeper.ConditionalOrderTriggerConfig{
		MaxTriggersPerBlock: clobkeeper.MaxConditionalTriggersPerBlock,
	})

	_ = tApp.InitChain()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Legitimate order must have triggered.
	require.True(t,
		tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, legitOrder.OrderId),
		"legitimate order (trigger=51B, oracle~50B, LTE-crossed) must be triggered with flag ON",
	)

	// Non-crossing orders must remain untriggered.
	untriggered := tApp.App.ClobKeeper.GetAllUntriggeredConditionalOrders(ctx)
	require.Len(t, untriggered, nNonCrossing,
		"all %d non-crossing orders (trigger=1) must remain untriggered; "+
			"the bounded scan must not spuriously trigger them", nNonCrossing)

	fmt.Printf(
		"legit-trigger: flag=ON non_crossing=%d legit_triggered=true remaining=%d\n",
		nNonCrossing,
		len(untriggered),
	)
}

func TestConditionalTrigger_SignedSingleOwnerAccumulation(t *testing.T) {
	const perWindow = 50
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(condE2EGenesis).Build()
	ctx := tApp.InitChain()
	gtbt := uint32(ctx.BlockTime().Add(30 * time.Minute).Unix())

	firstWindow := make([]clobtypes.Order, 0, perWindow)
	for i := 0; i < perWindow; i++ {
		firstWindow = append(firstWindow, e2eSignedConditionalOrder(uint32(i), gtbt))
	}
	deliverTxs := e2eCheckAndCollectSignedOrders(t, tApp, ctx, firstWindow)
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{DeliverTxsOverride: deliverTxs})
	require.Len(t, tApp.App.ClobKeeper.GetAllUntriggeredConditionalOrders(ctx), perWindow)
	fmt.Printf("accum-e2e: window=1 admitted=%d resting=%d\n", perWindow, perWindow)

	ctx = tApp.AdvanceToBlock(103, testapp.AdvanceToBlockOptions{})
	secondWindow := make([]clobtypes.Order, 0, perWindow)
	for i := 0; i < perWindow; i++ {
		secondWindow = append(secondWindow, e2eSignedConditionalOrder(uint32(perWindow+i), gtbt))
	}
	deliverTxs = e2eCheckAndCollectSignedOrders(t, tApp, ctx, secondWindow)
	ctx = tApp.AdvanceToBlock(104, testapp.AdvanceToBlockOptions{DeliverTxsOverride: deliverTxs})

	untriggered := tApp.App.ClobKeeper.GetAllUntriggeredConditionalOrders(ctx)
	require.Len(t, untriggered, 2*perWindow)
	require.Equal(t, uint32(2*perWindow), tApp.App.ClobKeeper.GetStatefulOrderCount(ctx, constants.Alice_Num0))
	fmt.Printf("accum-e2e: same owner admitted second rate window; resting=%d stateful_count=%d\n",
		len(untriggered),
		tApp.App.ClobKeeper.GetStatefulOrderCount(ctx, constants.Alice_Num0),
	)
}
