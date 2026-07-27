package keeper_test

// Packet 1 acceptance tests: trigger-price secondary index consistency and range-scan correctness.
//
// Two test groups:
//
//  1. TestConditionalOrderIndex_InvariantPlaceCancelTriggerExpire – invariant: after arbitrary
//     place/cancel/trigger/expire sequences the set of orderIds reconstructed from the index
//     (per clobPairId + direction) EXACTLY equals GetAllUntriggeredConditionalOrders.
//
//  2. TestConditionalOrderIndex_CrossedRangeScan – ordering/crossed test: big-endian range scan
//     yields ascending subticks; IterateCrossedConditionalOrders returns exactly the crossed set
//     for representative prices, boundary-inclusive.

import (
	"sort"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// makeConditionalOrder constructs a minimal conditional order for tests with unique OrderIds.
func makeConditionalOrder(
	subaccountId satypes.SubaccountId,
	clientId uint32,
	clobPairId uint32,
	side clobtypes.Order_Side,
	conditionType clobtypes.Order_ConditionType,
	triggerSubticks uint64,
) clobtypes.Order {
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
		ConditionType:                   conditionType,
		ConditionalOrderTriggerSubticks: triggerSubticks,
	}
}

// indexTestKeeperWithEth builds a test context with BTC (clobPairId=0) and ETH (clobPairId=1).
func indexTestKeeperWithEth(t *testing.T) keepertest.ClobKeepersTestContext {
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

	// BTC perpetual + clob pair (clobPairId=0, perpetualId=0)
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

	// ETH perpetual + clob pair (clobPairId=1, perpetualId=1)
	_, err = ks.PerpetualsKeeper.CreatePerpetual(
		ctx,
		constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.Id,
		constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.Ticker,
		constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.MarketId,
		constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.AtomicResolution,
		constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.DefaultFundingPpm,
		constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.LiquidityTier,
		constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.MarketType,
	)
	require.NoError(t, err)
	_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
		ctx,
		constants.ClobPair_Eth.Id,
		constants.ClobPair_Eth.MustGetPerpetualId(),
		satypes.BaseQuantums(constants.ClobPair_Eth.StepBaseQuantums),
		constants.ClobPair_Eth.QuantumConversionExponent,
		constants.ClobPair_Eth.SubticksPerTick,
		constants.ClobPair_Eth.Status,
	)
	require.NoError(t, err)

	return ks
}

// collectIndexOrderIds reads all order IDs from the trigger-price index store and returns them
// as a map: (clobPairId, direction) → set of orderId state keys.
// It decodes the composite key layout: <clobPairId:4><dir:1><subticks:8><orderId:N>.
func collectIndexOrderIds(
	t *testing.T,
	k *clobkeeper.Keeper,
	ctx sdk.Context,
) map[uint32]map[byte]map[string]struct{} {
	t.Helper()

	// result[clobPairId][direction] = set of orderId.ToStateKey() strings
	result := make(map[uint32]map[byte]map[string]struct{})

	store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
	it := store.Iterator(nil, nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		rawKey := it.Key()
		// Fixed prefix: <clobPairId:4><dir:1><subticks:8><sequenceKey:8> = 21 bytes.
		require.Greater(t, len(rawKey), 21, "malformed index key: %x", rawKey)

		clobPairIdBytes := rawKey[0:4]
		dirByte := rawKey[4]
		orderKeyBytes := rawKey[21:]

		clobPairId := uint32(clobPairIdBytes[0])<<24 |
			uint32(clobPairIdBytes[1])<<16 |
			uint32(clobPairIdBytes[2])<<8 |
			uint32(clobPairIdBytes[3])

		if _, ok := result[clobPairId]; !ok {
			result[clobPairId] = make(map[byte]map[string]struct{})
		}
		if _, ok := result[clobPairId][dirByte]; !ok {
			result[clobPairId][dirByte] = make(map[string]struct{})
		}
		result[clobPairId][dirByte][string(orderKeyBytes)] = struct{}{}
	}
	return result
}

// TestConditionalOrderIndex_InvariantPlaceCancelTriggerExpire verifies that after arbitrary
// place/cancel/trigger/expire sequences the trigger-price index exactly mirrors the set of
// untriggered conditional orders (per clobPairId + direction), across both BTC and ETH pairs.
func TestConditionalOrderIndex_InvariantPlaceCancelTriggerExpire(t *testing.T) {
	ks := indexTestKeeperWithEth(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// The trigger-price index is consensus state maintained by the placement/cancel/trigger/expiry
	// hooks ONLY when the mitigation flag is enabled (rolling-deploy gating). Enable it so the
	// hooks populate the index for this invariant test.
	enableTriggerConfig(t, k, ctx, clobkeeper.MaxConditionalTriggersPerBlock)

	expiry := ctx.BlockTime().Add(clobtypes.StatefulOrderTimeWindow - time.Hour)

	// --- Helper: build expected index from GetAllUntriggeredConditionalOrders ---
	expectedIndex := func() map[uint32]map[byte]map[string]struct{} {
		all := k.GetAllUntriggeredConditionalOrders(ctx)
		result := make(map[uint32]map[byte]map[string]struct{})
		for _, order := range all {
			cp := uint32(order.GetClobPairId())
			var dir byte
			if order.IsTakeProfitOrder() && order.IsBuy() ||
				order.IsStopLossOrder() && !order.IsBuy() {
				dir = clobkeeper.TriggerDirectionLTE
			} else {
				dir = clobkeeper.TriggerDirectionGTE
			}
			if _, ok := result[cp]; !ok {
				result[cp] = make(map[byte]map[string]struct{})
			}
			if _, ok := result[cp][dir]; !ok {
				result[cp][dir] = make(map[string]struct{})
			}
			result[cp][dir][string(order.OrderId.ToStateKey())] = struct{}{}
		}
		return result
	}

	checkInvariant := func(label string) {
		t.Helper()
		expected := expectedIndex()
		actual := collectIndexOrderIds(t, k, ctx)
		require.Equal(t, expected, actual, "index != full-set at step: %s", label)
	}

	// Build orders manually to ensure unique OrderIds (ClientIds must be distinct per owner+pair).
	// Clob0 (BTC), all Alice_Num0 — ClientIds 0–3.
	// Take-profit BUY → LTE direction (trigger when oracle ≤ trigger price)
	tpBuyClob0 := makeConditionalOrder(constants.Alice_Num0, 0, 0, condSideBuy, condTypeTP, 20)
	// Stop-loss BUY → GTE direction (trigger when oracle ≥ trigger price)
	slBuyClob0 := makeConditionalOrder(constants.Alice_Num0, 1, 0, condSideBuy, condTypeSL, 20)
	// Stop-loss SELL → LTE direction
	slSellClob0 := makeConditionalOrder(constants.Alice_Num0, 2, 0, condSideSell, condTypeSL, 20)
	// Take-profit SELL → GTE direction
	tpSellClob0 := makeConditionalOrder(constants.Alice_Num0, 3, 0, condSideSell, condTypeTP, 20)
	// Second GTE order on clob0 (different client) for non-trivial set size
	slBuyClob0B := makeConditionalOrder(constants.Alice_Num0, 4, 0, condSideBuy, condTypeSL, 25)

	// Clob1 (ETH), Alice_Num1 to avoid OrderId collision with clob0 orders.
	// ClientId 0-3 but different SubaccountId.
	tpBuyClob1 := makeConditionalOrder(constants.Alice_Num1, 0, 1, condSideBuy, condTypeTP, 30)
	slBuyClob1 := makeConditionalOrder(constants.Alice_Num1, 1, 1, condSideBuy, condTypeSL, 30)
	slSellClob1 := makeConditionalOrder(constants.Alice_Num1, 2, 1, condSideSell, condTypeSL, 30)
	tpSellClob1 := makeConditionalOrder(constants.Alice_Num1, 3, 1, condSideSell, condTypeTP, 30)

	// Step 1: place all orders across both clob pairs
	for _, order := range []clobtypes.Order{
		tpBuyClob0, slBuyClob0, slSellClob0, tpSellClob0, slBuyClob0B,
		tpBuyClob1, slBuyClob1, slSellClob1, tpSellClob1,
	} {
		o := order // capture
		k.SetLongTermOrderPlacement(ctx, o, 1)
		k.AddStatefulOrderIdExpiration(ctx, expiry, o.OrderId)
	}
	checkInvariant("after placing all orders")

	// Step 2: cancel one LTE and one GTE from clob0 via DeleteLongTermOrderPlacement
	k.DeleteLongTermOrderPlacement(ctx, tpBuyClob0.OrderId)  // removes LTE on clob0
	k.DeleteLongTermOrderPlacement(ctx, slBuyClob0B.OrderId) // removes GTE on clob0
	checkInvariant("after cancelling tpBuyClob0 and slBuyClob0B")

	// Step 3: trigger one order on clob0 (slBuyClob0 is GTE direction)
	k.MustTriggerConditionalOrder(ctx, slBuyClob0.OrderId)
	checkInvariant("after triggering slBuyClob0")

	// Step 4: expire all remaining orders via DeleteLongTermOrderPlacement
	// (abci.go calls DeleteLongTermOrderPlacement for each expired id)
	remaining := k.GetAllUntriggeredConditionalOrders(ctx)
	for _, order := range remaining {
		k.DeleteLongTermOrderPlacement(ctx, order.OrderId)
	}
	checkInvariant("after expiring all remaining orders")
	require.Empty(t, k.GetAllUntriggeredConditionalOrders(ctx), "no untriggered orders should remain")

	// Step 5: re-place a subset and verify invariant again
	for _, order := range []clobtypes.Order{slSellClob0, tpSellClob1} {
		o := order
		k.SetLongTermOrderPlacement(ctx, o, 2)
		k.AddStatefulOrderIdExpiration(ctx, expiry, o.OrderId)
	}
	checkInvariant("after re-placing two orders")

	// Step 6: cancel remaining; index must be empty
	for _, order := range []clobtypes.Order{slSellClob0, tpSellClob1} {
		k.DeleteLongTermOrderPlacement(ctx, order.OrderId)
	}
	checkInvariant("after final cancel, index must be empty")
	actual := collectIndexOrderIds(t, k, ctx)
	require.Empty(t, actual, "index store must be empty after all removals")
}

// triggerSubticksFromIndexKey extracts the triggerSubticks field from a raw index key
// (relative to the prefix store). Layout: <clobPairId:4><dir:1><subticks:8><orderId:N>.
func triggerSubticksFromIndexKey(t *testing.T, rawKey []byte) uint64 {
	t.Helper()
	require.GreaterOrEqual(t, len(rawKey), 13, "index key too short")
	hi := uint64(rawKey[5])<<56 | uint64(rawKey[6])<<48 | uint64(rawKey[7])<<40 |
		uint64(rawKey[8])<<32 | uint64(rawKey[9])<<24 | uint64(rawKey[10])<<16 |
		uint64(rawKey[11])<<8 | uint64(rawKey[12])
	return hi
}

// TestConditionalOrderIndex_CrossedRangeScan verifies:
//  1. Forward iteration order is ascending in subticks (big-endian encoding is correct).
//  2. IterateCrossedConditionalOrders for LTE direction returns exactly orders with
//     triggerSubticks >= priceSubticks (boundary-inclusive).
//  3. IterateCrossedConditionalOrders for GTE direction returns exactly orders with
//     triggerSubticks <= priceSubticks (boundary-inclusive).
//  4. Both directions return empty sets for prices that cross nothing.
func TestConditionalOrderIndex_CrossedRangeScan(t *testing.T) {
	ks := newCondTestKeepers(t) // BTC only
	ctx := ks.Ctx
	k := ks.ClobKeeper

	// Enable the mitigation so the placement hooks maintain the (consensus-gated) index.
	enableTriggerConfig(t, k, ctx, clobkeeper.MaxConditionalTriggersPerBlock)

	// Build a set of LTE-direction orders (take-profit buy) with varying triggerSubticks:
	//   subticks values: 10, 20, 25, 30
	lteMakeOrder := func(clientId uint32, triggerSubticks uint64) clobtypes.Order {
		return clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     clientId,
				OrderFlags:   clobtypes.OrderIdFlags_Conditional,
				ClobPairId:   0,
			},
			Side:                            condSideBuy,
			Quantums:                        1_000_000,
			Subticks:                        50_000_000_000,
			GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: condTestGTBT},
			ConditionType:                   condTypeTP,
			ConditionalOrderTriggerSubticks: triggerSubticks,
		}
	}

	// Build a set of GTE-direction orders (stop-loss buy) with varying triggerSubticks:
	//   subticks values: 50, 100, 200, 300
	gteMakeOrder := func(clientId uint32, triggerSubticks uint64) clobtypes.Order {
		return clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num1,
				ClientId:     clientId,
				OrderFlags:   clobtypes.OrderIdFlags_Conditional,
				ClobPairId:   0,
			},
			Side:                            condSideBuy,
			Quantums:                        1_000_000,
			Subticks:                        50_000_000_000,
			GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: condTestGTBT},
			ConditionType:                   condTypeSL,
			ConditionalOrderTriggerSubticks: triggerSubticks,
		}
	}

	// LTE orders (triggerSubticks ascending): 10, 20, 25, 30
	lteSubticks := []uint64{10, 20, 25, 30}
	lteOrders := make([]clobtypes.Order, 0, len(lteSubticks))
	for i, sub := range lteSubticks {
		order := lteMakeOrder(uint32(i), sub)
		k.SetLongTermOrderPlacement(ctx, order, 1)
		lteOrders = append(lteOrders, order)
	}

	// GTE orders (triggerSubticks ascending): 50, 100, 200, 300
	gteSubticks := []uint64{50, 100, 200, 300}
	gteOrders := make([]clobtypes.Order, 0, len(gteSubticks))
	for i, sub := range gteSubticks {
		order := gteMakeOrder(uint32(i), sub)
		k.SetLongTermOrderPlacement(ctx, order, 1)
		gteOrders = append(gteOrders, order)
	}

	// --- Sub-test 1: big-endian ascending ordering within each direction bucket ---
	t.Run("ascending_subticks_order_in_store", func(t *testing.T) {
		store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)

		// LTE bucket: scan all keys for clobPairId=0, dir=LTE
		lteStart := make([]byte, 4+1)
		lteStart[4] = clobkeeper.TriggerDirectionLTE
		lteEnd := make([]byte, 4+1)
		lteEnd[4] = clobkeeper.TriggerDirectionGTE
		lteIt := store.Iterator(lteStart, lteEnd)
		defer lteIt.Close()

		var observedLTE []uint64
		for ; lteIt.Valid(); lteIt.Next() {
			observedLTE = append(observedLTE, triggerSubticksFromIndexKey(t, lteIt.Key()))
		}

		// Must be non-decreasing (ascending).
		sorted := append([]uint64{}, observedLTE...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
		require.Equal(t, sorted, observedLTE, "LTE bucket keys must be ascending in subticks")
		require.Equal(t, lteSubticks, observedLTE, "LTE bucket must contain exactly the placed subticks")

		// GTE bucket: scan all keys for clobPairId=0, dir=GTE
		gteStart := make([]byte, 4+1)
		gteStart[4] = clobkeeper.TriggerDirectionGTE
		// End: exclusive upper bound = next clobPairId (0+1=1, big-endian)
		gteEnd := make([]byte, 4)
		gteEnd[3] = 1 // clobPairId=1 big-endian
		gteIt := store.Iterator(gteStart, gteEnd)
		defer gteIt.Close()

		var observedGTE []uint64
		for ; gteIt.Valid(); gteIt.Next() {
			observedGTE = append(observedGTE, triggerSubticksFromIndexKey(t, gteIt.Key()))
		}
		sortedGTE := append([]uint64{}, observedGTE...)
		sort.Slice(sortedGTE, func(i, j int) bool { return sortedGTE[i] < sortedGTE[j] })
		require.Equal(t, sortedGTE, observedGTE, "GTE bucket keys must be ascending in subticks")
		require.Equal(t, gteSubticks, observedGTE, "GTE bucket must contain exactly the placed subticks")
	})

	// --- Sub-test 2: LTE crossed range scan ---
	// LTE condition: order triggers when oraclePrice ≤ triggerSubticks.
	// So a scan at priceSubticks P should return all orders with triggerSubticks ≥ P.
	t.Run("lte_crossed_scan", func(t *testing.T) {
		collectCrossed := func(price uint64) map[string]struct{} {
			got := make(map[string]struct{})
			k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionLTE, price, func(id clobtypes.OrderId) bool {
				got[string(id.ToStateKey())] = struct{}{}
				return true
			})
			return got
		}

		// Price = 5: all LTE orders are crossed (all have triggerSubticks ≥ 5)
		crossed5 := collectCrossed(5)
		require.Len(t, crossed5, 4)

		// Price = 10: boundary — order with triggerSubticks=10 IS crossed (inclusive)
		crossed10 := collectCrossed(10)
		require.Len(t, crossed10, 4, "price=10 is exactly on boundary, must be included")

		// Price = 11: order with triggerSubticks=10 is NOT crossed; 20,25,30 are
		crossed11 := collectCrossed(11)
		require.Len(t, crossed11, 3)
		require.NotContains(t, crossed11, string(lteOrders[0].OrderId.ToStateKey())) // subticks=10

		// Price = 25: orders 25 and 30 are crossed
		crossed25 := collectCrossed(25)
		require.Len(t, crossed25, 2)
		require.Contains(t, crossed25, string(lteOrders[2].OrderId.ToStateKey())) // subticks=25
		require.Contains(t, crossed25, string(lteOrders[3].OrderId.ToStateKey())) // subticks=30

		// Price = 30: boundary for highest order — only it is crossed
		crossed30 := collectCrossed(30)
		require.Len(t, crossed30, 1)
		require.Contains(t, crossed30, string(lteOrders[3].OrderId.ToStateKey())) // subticks=30

		// Price = 31: nothing crossed
		crossed31 := collectCrossed(31)
		require.Empty(t, crossed31, "price above all triggerSubticks: nothing should be crossed")

		// Verify GTE orders are never included in LTE scan
		for _, gteOrder := range gteOrders {
			for _, price := range []uint64{0, 50, 100, 300} {
				got := collectCrossed(price)
				require.NotContains(t, got, string(gteOrder.OrderId.ToStateKey()),
					"GTE order must never appear in LTE scan")
			}
		}
	})

	// --- Sub-test 3: GTE crossed range scan ---
	// GTE condition: order triggers when oraclePrice ≥ triggerSubticks.
	// So a scan at priceSubticks P should return all orders with triggerSubticks ≤ P.
	t.Run("gte_crossed_scan", func(t *testing.T) {
		collectCrossed := func(price uint64) map[string]struct{} {
			got := make(map[string]struct{})
			k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionGTE, price, func(id clobtypes.OrderId) bool {
				got[string(id.ToStateKey())] = struct{}{}
				return true
			})
			return got
		}

		// Price = 49: nothing crossed (all have triggerSubticks ≥ 50)
		crossed49 := collectCrossed(49)
		require.Empty(t, crossed49)

		// Price = 50: boundary — order with triggerSubticks=50 IS crossed (inclusive)
		crossed50 := collectCrossed(50)
		require.Len(t, crossed50, 1)
		require.Contains(t, crossed50, string(gteOrders[0].OrderId.ToStateKey())) // subticks=50

		// Price = 100: boundary — orders 50 and 100 are crossed
		crossed100 := collectCrossed(100)
		require.Len(t, crossed100, 2)
		require.Contains(t, crossed100, string(gteOrders[0].OrderId.ToStateKey())) // subticks=50
		require.Contains(t, crossed100, string(gteOrders[1].OrderId.ToStateKey())) // subticks=100

		// Price = 201: orders 50, 100, 200 are crossed; 300 is not
		crossed201 := collectCrossed(201)
		require.Len(t, crossed201, 3)
		require.NotContains(t, crossed201, string(gteOrders[3].OrderId.ToStateKey())) // subticks=300

		// Price = 300: all crossed
		crossed300 := collectCrossed(300)
		require.Len(t, crossed300, 4)

		// Price = 1_000_000: all crossed (price above all trigger prices)
		crossedHigh := collectCrossed(1_000_000)
		require.Len(t, crossedHigh, 4)

		// Verify LTE orders are never included in GTE scan
		for _, lteOrder := range lteOrders {
			for _, price := range []uint64{0, 10, 25, 100} {
				got := collectCrossed(price)
				require.NotContains(t, got, string(lteOrder.OrderId.ToStateKey()),
					"LTE order must never appear in GTE scan")
			}
		}
	})

	// --- Sub-test 4: empty scan when nothing is crossed ---
	t.Run("empty_scan_when_nothing_crossed", func(t *testing.T) {
		var ids []clobtypes.OrderId
		// LTE: price above all trigger prices → nothing crossed
		k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionLTE, 10_000, func(id clobtypes.OrderId) bool {
			ids = append(ids, id)
			return true
		})
		require.Empty(t, ids)

		// GTE: price below all trigger prices → nothing crossed
		k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionGTE, 1, func(id clobtypes.OrderId) bool {
			ids = append(ids, id)
			return true
		})
		require.Empty(t, ids)
	})
}

// TestBuildConditionalOrderTriggerPriceIndex verifies the one-shot state-breaking migration that
// the v9.6 upgrade handler runs: given resting untriggered conditional orders that have NO
// trigger-price index entries (the pre-upgrade state), BuildConditionalOrderTriggerPriceIndex
// reconstructs the index so that every order is present, the build is idempotent, and orders are
// subsequently triggerable via the index.
func TestBuildConditionalOrderTriggerPriceIndex(t *testing.T) {
	ks := indexTestKeeperWithEth(t)
	ctx := ks.Ctx
	k := ks.ClobKeeper

	expiry := ctx.BlockTime().Add(clobtypes.StatefulOrderTimeWindow - 1)

	// Build a representative set of conditional orders across both directions and clob pairs.
	tpBuyBtc := makeConditionalOrder(constants.Alice_Num0, 10, 0, condSideBuy, condTypeTP, 100)   // LTE
	slBuyBtc := makeConditionalOrder(constants.Alice_Num0, 11, 0, condSideBuy, condTypeSL, 200)   // GTE
	tpSellEth := makeConditionalOrder(constants.Alice_Num1, 12, 1, condSideSell, condTypeTP, 300) // GTE
	slSellEth := makeConditionalOrder(constants.Alice_Num1, 13, 1, condSideSell, condTypeSL, 50)  // LTE

	orders := []clobtypes.Order{tpBuyBtc, slBuyBtc, tpSellEth, slSellEth}

	// Step 1: place all orders normally (the placement hooks write index entries).
	for _, o := range orders {
		order := o
		k.SetLongTermOrderPlacement(ctx, order, 1)
		k.AddStatefulOrderIdExpiration(ctx, expiry, order.OrderId)
	}

	// Step 2: manually remove all index entries to simulate the "pre-upgrade" state (orders in the
	// SO/U: store but no TPIdx: entries).
	{
		indexStore := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
		it := indexStore.Iterator(nil, nil)
		var keysToDelete [][]byte
		for ; it.Valid(); it.Next() {
			keyCopy := make([]byte, len(it.Key()))
			copy(keyCopy, it.Key())
			keysToDelete = append(keysToDelete, keyCopy)
		}
		it.Close()
		for _, key := range keysToDelete {
			indexStore.Delete(key)
		}
	}

	require.Empty(t, collectIndexOrderIds(t, k, ctx),
		"index must be empty after manual removal (pre-upgrade simulation)")
	require.Len(t, k.GetAllUntriggeredConditionalOrders(ctx), len(orders),
		"orders must still be in the untriggered store after index removal")

	// Step 3: run the one-shot migration build.
	k.BuildConditionalOrderTriggerPriceIndex(ctx)

	// Step 4: assert the rebuilt index exactly equals what GetAllUntriggeredConditionalOrders implies.
	expectedIndex := func() map[uint32]map[byte]map[string]struct{} {
		result := make(map[uint32]map[byte]map[string]struct{})
		for _, order := range k.GetAllUntriggeredConditionalOrders(ctx) {
			cp := uint32(order.GetClobPairId())
			dir := clobkeeper.TriggerDirectionGTE
			if order.IsTakeProfitOrder() && order.IsBuy() ||
				order.IsStopLossOrder() && !order.IsBuy() {
				dir = clobkeeper.TriggerDirectionLTE
			}
			if _, ok := result[cp]; !ok {
				result[cp] = make(map[byte]map[string]struct{})
			}
			if _, ok := result[cp][dir]; !ok {
				result[cp][dir] = make(map[string]struct{})
			}
			result[cp][dir][string(order.OrderId.ToStateKey())] = struct{}{}
		}
		return result
	}()
	postBuildIndex := collectIndexOrderIds(t, k, ctx)
	require.Equal(t, expectedIndex, postBuildIndex,
		"post-build index must exactly equal GetAllUntriggeredConditionalOrders")

	// Step 5: idempotency — building again must not change the index.
	k.BuildConditionalOrderTriggerPriceIndex(ctx)
	require.Equal(t, postBuildIndex, collectIndexOrderIds(t, k, ctx),
		"build must be idempotent: a second call must not change the index")

	// Step 6: after build, orders are triggerable via IterateCrossedConditionalOrders.
	var triggeredByLTE []clobtypes.OrderId
	k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionLTE, 100, func(id clobtypes.OrderId) bool {
		triggeredByLTE = append(triggeredByLTE, id)
		return true
	})
	require.Len(t, triggeredByLTE, 1, "LTE scan at price=100 must find tpBuyBtc (triggerSubticks=100)")
	require.Equal(t, tpBuyBtc.OrderId, triggeredByLTE[0])

	var triggeredByGTE []clobtypes.OrderId
	k.IterateCrossedConditionalOrders(ctx, 0, clobkeeper.TriggerDirectionGTE, 200, func(id clobtypes.OrderId) bool {
		triggeredByGTE = append(triggeredByGTE, id)
		return true
	})
	require.Len(t, triggeredByGTE, 1, "GTE scan at price=200 must find slBuyBtc (triggerSubticks=200)")
	require.Equal(t, slBuyBtc.OrderId, triggeredByGTE[0])
}
