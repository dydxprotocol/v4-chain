package keeper

import (
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	gometrics "github.com/hashicorp/go-metrics"
)

// UntriggeredConditionalOrders is an in-memory struct stored on the clob Keeper.
// It is intended to efficiently store placed conditional orders and poll out triggered
// conditional orders on oracle price changes for a given ClobPairId.
// All orders contained in this data structure are placed conditional orders with the same
// ClobPairId and are untriggered, unexpired, and uncancelled.
// Note that we are using a Order list for the initial implementation, but for
// optimal runtime a an AVL-tree backed priority queue would work.
// TODO(CLOB-717) Change list to use priority queue.
type UntriggeredConditionalOrders struct {
	// All untriggered take profit buy orders and stop loss sell orders sorted by time priority.
	// These orders will be triggered when the oracle price goes lower than or equal to the trigger price.
	// This array functions like a max heap.
	OrdersToTriggerWhenOraclePriceLTETriggerPrice []types.Order

	// All untriggered take profit sell orders and stop loss buy orders sorted by time priority.
	// These orders will be triggered when the oracle price goes greater than or equal to the trigger price.
	// This array functions like a min heap.
	OrdersToTriggerWhenOraclePriceGTETriggerPrice []types.Order
}

func (k Keeper) NewUntriggeredConditionalOrders() *UntriggeredConditionalOrders {
	return NewUntriggeredConditionalOrders()
}

func NewUntriggeredConditionalOrders() *UntriggeredConditionalOrders {
	return &UntriggeredConditionalOrders{
		OrdersToTriggerWhenOraclePriceLTETriggerPrice: make([]types.Order, 0),
		OrdersToTriggerWhenOraclePriceGTETriggerPrice: make([]types.Order, 0),
	}
}

// IsEmpty returns true if the UntriggeredConditionalOrders' order slices are both empty.
func (untriggeredOrders *UntriggeredConditionalOrders) IsEmpty() bool {
	return len(untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice) == 0 &&
		len(untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice) == 0
}

// AddUntriggeredConditionalOrder adds an untriggered conditional order to the UntriggeredConditionalOrders
// data structure. It will panic if the order is not a conditional order.
func (untriggeredOrders *UntriggeredConditionalOrders) AddUntriggeredConditionalOrder(order types.Order) {
	order.MustBeConditionalOrder()

	if order.IsTakeProfitOrder() {
		if order.IsBuy() {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice,
				order,
			)
		} else {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice,
				order,
			)
		}
	}

	if order.IsStopLossOrder() {
		if order.IsBuy() {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice,
				order,
			)
		} else {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice,
				order,
			)
		}
	}
}

// RemoveUntriggeredConditionalOrders removes a list of order ids from the `UntriggeredConditionalOrders`
// data structure. This function will panic if the order ids contained involve more than one ClobPairId.
func (untriggeredOrders *UntriggeredConditionalOrders) RemoveUntriggeredConditionalOrders(
	orderIdsToRemove []types.OrderId,
) {
	if len(orderIdsToRemove) == 0 {
		return
	}

	// all orders should have the same ClobPairId
	clobPairId := types.ClobPairId(orderIdsToRemove[0].GetClobPairId())
	for _, orderId := range orderIdsToRemove {
		orderClobPairId := types.ClobPairId(orderId.GetClobPairId())
		if types.ClobPairId(orderId.GetClobPairId()) != clobPairId {
			panic(
				fmt.Sprintf(
					"RemoveExpiredUntriggeredConditionalOrders: expected all orders to have the same ClobPairId. "+
						"Got %v and %v.",
					clobPairId,
					orderClobPairId,
				),
			)
		}
	}

	orderIdsToRemoveSet := lib.UniqueSliceToSet(orderIdsToRemove)

	newOrdersToTriggerWhenOraclePriceLTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice {
		if _, exists := orderIdsToRemoveSet[order.OrderId]; !exists {
			newOrdersToTriggerWhenOraclePriceLTETriggerPrice = append(newOrdersToTriggerWhenOraclePriceLTETriggerPrice, order)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = newOrdersToTriggerWhenOraclePriceLTETriggerPrice

	newOrdersToTriggerWhenOraclePriceGTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice {
		if _, exists := orderIdsToRemoveSet[order.OrderId]; !exists {
			newOrdersToTriggerWhenOraclePriceGTETriggerPrice = append(newOrdersToTriggerWhenOraclePriceGTETriggerPrice, order)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = newOrdersToTriggerWhenOraclePriceGTETriggerPrice
}

// PollTriggeredConditionalOrders removes all triggered conditional orders from the
// `UntriggeredConditionalOrders` struct given a new oracle price for a clobPairId. It returns
// a list of order ids that were triggered. This is only called in EndBlocker. We round up to the nearest
// subtick int for LTE and down to the nearest subtick int for GTE conditions. This is pessimistic rounding,
// as we want to trigger orders only when we are sure they are triggerable.
func (untriggeredOrders *UntriggeredConditionalOrders) PollTriggeredConditionalOrders(
	oraclePriceSubticksRat *big.Rat,
) []types.OrderId {
	triggeredOrderIds := make([]types.OrderId, 0)
	pessimisticLTESubticks := types.Subticks(lib.BigRatRound(oraclePriceSubticksRat, true).Uint64())
	// For the lte array, find all orders that are triggered when oracle price goes lower
	// than or equal to the trigger price.

	newOrdersToTriggerWhenOraclePriceLTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice {
		if order.CanTrigger(pessimisticLTESubticks) {
			triggeredOrderIds = append(triggeredOrderIds, order.OrderId)
		} else {
			newOrdersToTriggerWhenOraclePriceLTETriggerPrice = append(
				newOrdersToTriggerWhenOraclePriceLTETriggerPrice,
				order,
			)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = newOrdersToTriggerWhenOraclePriceLTETriggerPrice

	pessimisticGTESubticks := types.Subticks(lib.BigRatRound(oraclePriceSubticksRat, false).Uint64())
	// For the gte array, find all orders that are triggered when oracle price goes greater
	// than or equal to the trigger price.
	newOrdersToTriggerWhenOraclePriceGTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice {
		if order.CanTrigger(pessimisticGTESubticks) {
			triggeredOrderIds = append(triggeredOrderIds, order.OrderId)
		} else {
			newOrdersToTriggerWhenOraclePriceGTETriggerPrice = append(
				newOrdersToTriggerWhenOraclePriceGTETriggerPrice,
				order,
			)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = newOrdersToTriggerWhenOraclePriceGTETriggerPrice

	return triggeredOrderIds
}

// OrganizeUntriggeredConditionalOrdersFromState takes in a list of conditional orders read from
// state, organize them and return in form of `UntriggeredConditionalOrders` struct.
func OrganizeUntriggeredConditionalOrdersFromState(
	conditonalOrdersFromState []types.Order,
) map[types.ClobPairId]*UntriggeredConditionalOrders {
	ret := make(map[types.ClobPairId]*UntriggeredConditionalOrders)

	for _, order := range conditonalOrdersFromState {
		clobPairId := types.ClobPairId(order.GetClobPairId())
		untriggeredConditionalOrders, exists := ret[clobPairId]
		if !exists {
			untriggeredConditionalOrders = NewUntriggeredConditionalOrders()
			ret[clobPairId] = untriggeredConditionalOrders
		}
		untriggeredConditionalOrders.AddUntriggeredConditionalOrder(order)
	}

	return ret
}

// MaxConditionalTriggersPerBlock is the maximum number of conditional orders that may be
// triggered in a single block by MaybeTriggerConditionalOrders. This acts as a per-block
// budget that bounds EndBlocker work to O(budget) rather than O(global N). A conservative
// high default is chosen so that normal operation and equivalence tests are unaffected; the
// value will be made governance-tunable in Packet 3.
//
// Rationale for 1000: the observed legitimate peak is on the order of tens of triggers per
// block; 1000 provides a 10–100× safety margin while keeping worst-case EndBlocker work
// firmly bounded. When the crossed set exceeds the budget, only the nearest-crossing orders
// (for LTE, the lowest triggerSubticks ≥ price; for GTE, the highest triggerSubticks ≤ price)
// are triggered this block; the remainder are deferred to subsequent blocks where the budget
// refills, making deferral deterministic and consistent across all nodes.
const MaxConditionalTriggersPerBlock = 1000

// MaybeTriggerConditionalOrders queries the prices module for price updates and triggers
// any conditional orders that can be triggered. It is called in EndBlocker and returns the list
// of triggered conditional order ids, written to
// `ProcessProposerMatchesEvents.ConditionalOrderIdsTriggeredInLastBlock`.
//
// FLAG-GATED BEHAVIOR (consensus-safe rolling deploy):
//
//   - When the ConditionalOrderTriggerConfig is DISABLED (the default when unset), this function
//     runs maybeTriggerConditionalOrdersLegacy, which is byte-for-byte identical to the pre-fix
//     implementation: it reads the full untriggered set, organizes and time-sorts it, and triggers
//     against oracle + clamped trade prices with the original ordering and state writes. Because
//     the legacy path is bit-identical, the new binary can be rolled out on all nodes without a
//     version split — every node produces the same app hash until the flag is flipped.
//
//   - When the config is ENABLED, this function runs the bounded, crossing-priority path
//     (maybeTriggerConditionalOrdersBounded), which uses the trigger-price secondary index to do
//     O(crossed + budget) work per block instead of O(global N). The flag MUST be flipped at a
//     coordinated (governed) height, since enabling it changes the triggered set (under budget
//     pressure) and ordering (price-priority instead of time-priority) — both observable in block
//     state.
func (k Keeper) MaybeTriggerConditionalOrders(ctx sdk.Context) (allTriggeredOrderIds []types.OrderId) {
	defer metrics.ModuleMeasureSince(
		types.ModuleName,
		metrics.ClobMaybeTriggerConditionalOrders,
		time.Now(),
	)

	cfg := k.GetConditionalOrderTriggerConfig(ctx)
	if !cfg.Enabled || !k.IsConditionalOrderTriggerIndexReady(ctx) {
		// Legacy remains authoritative while disabled and throughout incremental activation.
		return k.maybeTriggerConditionalOrdersLegacy(ctx)
	}
	return k.maybeTriggerConditionalOrdersBounded(ctx, cfg)
}

// maybeTriggerConditionalOrdersLegacy is the ORIGINAL, pre-fix implementation preserved verbatim.
// It is executed when the ConditionalOrderTriggerConfig is disabled so that a chain running the
// new binary behaves bit-identically to the old binary (same triggered set, same ordering, same
// state writes, same metrics) — enabling a rolling deploy with delayed, coordinated activation.
//
// DO NOT change the behavior of this function; any change here would break the rolling-deploy
// equivalence guarantee. It intentionally performs the full-set read + organize + sort.
func (k Keeper) maybeTriggerConditionalOrdersLegacy(ctx sdk.Context) (allTriggeredOrderIds []types.OrderId) {
	clobPairToUntriggeredConditionals := OrganizeUntriggeredConditionalOrdersFromState(
		k.GetAllUntriggeredConditionalOrders(ctx),
	)

	// Sort the keys for the untriggered conditional orders struct. We need to trigger
	// the conditional orders in an ordered way to have deterministic state writes.
	sortedKeys := lib.GetSortedKeys[types.SortedClobPairId](clobPairToUntriggeredConditionals)

	allTriggeredOrderIds = make([]types.OrderId, 0)
	// For all clob pair ids in UntriggeredConditionalOrders, fetch the updated
	// oracle price and poll out triggered conditional orders.
	for _, clobPairId := range sortedKeys {
		untriggered := clobPairToUntriggeredConditionals[clobPairId]
		clobPair, found := k.GetClobPair(ctx, clobPairId)

		// Error log and skip to next clob pair id if invalid clob pair id found.
		if !found {
			log.ErrorLogWithError(
				ctx,
				"MaybeTriggerConditionalOrders: Failed to fetch Clob Pair for Clob Pair Id",
				types.ErrInvalidClob,
				log.ClobPairId, clobPairId,
			)
			continue
		}

		// Trigger conditional orders using the oracle price.
		perpetualId := clobPair.MustGetPerpetualId()
		oraclePrice := k.GetOraclePriceSubticksRat(ctx, clobPair)
		triggered := k.TriggerOrdersWithPrice(ctx, untriggered, oraclePrice, perpetualId, metrics.OraclePrice)
		allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)

		// Trigger conditional orders using the last traded price.
		clampedMinTradePrice,
			clampedMaxTradePrice,
			found := k.getClampedTradePricesForTriggering(
			ctx,
			perpetualId,
			oraclePrice,
		)

		if found {
			triggered = k.TriggerOrdersWithPrice(ctx, untriggered, clampedMinTradePrice, perpetualId, metrics.MinTradePrice)
			allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)

			triggered = k.TriggerOrdersWithPrice(ctx, untriggered, clampedMaxTradePrice, perpetualId, metrics.MaxTradePrice)
			allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)
		}

		// Gauge the number of untriggered orders.
		metrics.SetGaugeWithLabels(
			metrics.ClobNumUntriggeredOrders,
			float32(
				len(untriggered.OrdersToTriggerWhenOraclePriceGTETriggerPrice)+
					len(untriggered.OrdersToTriggerWhenOraclePriceLTETriggerPrice),
			),
			metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
		)
	}

	return allTriggeredOrderIds
}

// maybeTriggerConditionalOrdersBounded is the fixed, bounded trigger path. It runs only when the
// ConditionalOrderTriggerConfig is enabled.
//
// Work is O(crossed + budget), independent of the global untriggered count N, because it uses the
// trigger-price secondary index to visit ONLY the orders whose trigger boundary the current price
// crosses (nearest-crossing first). Far-from-market orders that cannot cross are never visited, so
// they do not contribute to per-block scan cost; genuinely crossing orders (real liquidity /
// taking) are processed nearest-first.
//
// Prioritization under budget pressure (prioritize real crossing liquidity):
//   - The per-block budget is cfg.MaxTriggersPerBlock, drained nearest-crossing-first within each
//     clob pair. The first pair rotates across blocks so a budget smaller than the active-pair count
//     cannot permanently starve later pair ids.
//   - Any orders beyond the budget are deterministically deferred to subsequent blocks (the index
//     is ordered, so the same nearest-crossing prefix is chosen on every node every block).
//
// Ordering: (rotating clobPairId start, oracle→min-trade→max-trade, LTE→GTE, nearest trigger
// subticks, placement time). Deterministic across nodes. Differs from legacy time-priority ordering,
// which is why enabling the flag must happen at a coordinated height.
func (k Keeper) maybeTriggerConditionalOrdersBounded(
	ctx sdk.Context,
	cfg ConditionalOrderTriggerConfig,
) (allTriggeredOrderIds []types.OrderId) {
	allTriggeredOrderIds = make([]types.OrderId, 0)
	totalBudget := int(cfg.MaxTriggersPerBlock)

	// Determine the ACTIVE pairs (those with resting untriggered conditionals) in ascending id
	// order. GetAllClobPairs is already sorted ascending, so this is deterministic across nodes.
	// (This also implements the empty-pair skip: pairs with no index entries are never
	// visited, so their oracle price is never fetched and cannot panic.)
	activePairs := make([]types.ClobPair, 0)
	for _, clobPair := range k.GetAllClobPairs(ctx) {
		if k.clobPairHasTriggerIndexEntries(ctx, uint32(clobPair.Id)) {
			activePairs = append(activePairs, clobPair)
		}
	}

	// Rotate the deterministic pair order at a persisted cursor. Advancing by one active pair per
	// block guarantees eventual service even when totalBudget < len(activePairs), while also rotating
	// which market receives Phase-2 leftover capacity. If the cursor's pair is no longer active,
	// resume at the next greater active pair id, wrapping to the first pair.
	scheduledPairs := activePairs
	if len(activePairs) > 0 {
		start := 0
		if nextClobPairId, found := k.getConditionalOrderTriggerNextClobPairId(ctx); found {
			start = len(activePairs) // sentinel: wrap if every active id is below the cursor
			for i, clobPair := range activePairs {
				if uint32(clobPair.Id) >= nextClobPairId {
					start = i
					break
				}
			}
			if start == len(activePairs) {
				start = 0
			}
		}
		scheduledPairs = make([]types.ClobPair, 0, len(activePairs))
		scheduledPairs = append(scheduledPairs, activePairs[start:]...)
		scheduledPairs = append(scheduledPairs, activePairs[:start]...)

		nextStart := 0
		if len(scheduledPairs) > 1 {
			nextStart = 1
		}
		k.setConditionalOrderTriggerNextClobPairId(ctx, uint32(scheduledPairs[nextStart].Id))
	}

	// Fair-share the chain-wide per-block trigger budget across active pairs so a low-id pair
	// cannot consume the whole budget and starve higher-id pairs.
	if totalBudget > 0 && len(scheduledPairs) > 0 {
		// Phase 1: each active pair gets an equal share (at least 1) of the budget, in rotating
		// scheduling order. When there are more pairs than budget, advancing the starting pair every
		// block ensures every active pair eventually enters the serviced prefix.
		perPairBudget := totalBudget / len(activePairs)
		if perPairBudget < 1 {
			perPairBudget = 1
		}

		spent := 0
		for _, clobPair := range scheduledPairs {
			if spent >= totalBudget {
				break
			}
			pairBudget := perPairBudget
			if rem := totalBudget - spent; pairBudget > rem {
				pairBudget = rem
			}
			triggered, processed := k.triggerPairWithinBudget(ctx, clobPair, pairBudget)
			allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)
			spent += processed
		}

		// Phase 2: redistribute any leftover budget (from pairs that had fewer crossings than their
		// share) to the remaining crossed orders, in id order. Orders triggered in Phase 1 were
		// removed from the index, so re-scanning a pair here continues past them.
		for _, clobPair := range scheduledPairs {
			if spent >= totalBudget {
				break
			}
			triggered, processed := k.triggerPairWithinBudget(ctx, clobPair, totalBudget-spent)
			allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)
			spent += processed
		}
	}

	// Emit the global count of resting untriggered conditional orders. Sourced O(1) from the
	// admission-cap counter rather than the prior O(N) full-set read, so the monitoring signal is
	// preserved without reintroducing the vulnerable per-block scan. Global (not per-perpetual).
	metrics.SetGaugeWithLabels(
		metrics.ClobNumUntriggeredOrders,
		float32(k.GetUntriggeredConditionalOrderCountGlobal(ctx)),
	)

	return allTriggeredOrderIds
}

// triggerPairWithinBudget processes up to `budget` crossed conditional-order index entries for a
// single clob pair — first against the oracle price, then the clamped min/max trade prices. Each
// triggered order and each expired/orphaned entry removed during traversal consumes one unit of
// work. It returns both the triggered order ids and the number of index entries processed.
func (k Keeper) triggerPairWithinBudget(
	ctx sdk.Context,
	clobPair types.ClobPair,
	budget int,
) (triggeredOrderIds []types.OrderId, processed int) {
	triggeredOrderIds = make([]types.OrderId, 0)
	if budget <= 0 {
		return triggeredOrderIds, 0
	}

	clobPairId := types.ClobPairId(clobPair.Id)

	// Skip pairs with nothing to trigger, and read the oracle price without panicking on price 0.
	if !k.clobPairHasTriggerIndexEntries(ctx, uint32(clobPairId)) {
		return triggeredOrderIds, 0
	}
	perpetualId, oraclePrice, ok := k.getBoundedTriggerOraclePrice(ctx, clobPair)
	if !ok {
		return triggeredOrderIds, 0
	}

	remaining := budget

	// Trigger conditional orders using the oracle price.
	triggered, priceProcessed := k.triggerCrossedOrdersFromIndex(
		ctx, clobPairId, oraclePrice, perpetualId, metrics.OraclePrice, remaining,
	)
	triggeredOrderIds = append(triggeredOrderIds, triggered...)
	processed += priceProcessed
	remaining -= priceProcessed
	if remaining <= 0 {
		return triggeredOrderIds, processed
	}

	// Trigger conditional orders using the clamped trade prices.
	clampedMinTradePrice, clampedMaxTradePrice, found := k.getClampedTradePricesForTriggering(
		ctx,
		perpetualId,
		oraclePrice,
	)
	if found {
		triggered, priceProcessed = k.triggerCrossedOrdersFromIndex(
			ctx, clobPairId, clampedMinTradePrice, perpetualId, metrics.MinTradePrice, remaining,
		)
		triggeredOrderIds = append(triggeredOrderIds, triggered...)
		processed += priceProcessed
		remaining -= priceProcessed

		if remaining > 0 {
			triggered, priceProcessed = k.triggerCrossedOrdersFromIndex(
				ctx, clobPairId, clampedMaxTradePrice, perpetualId, metrics.MaxTradePrice, remaining,
			)
			triggeredOrderIds = append(triggeredOrderIds, triggered...)
			processed += priceProcessed
		}
	}

	return triggeredOrderIds, processed
}

// triggerCrossedOrdersFromIndex performs a crossed-order range scan on the trigger-price index
// for a single clob pair and price, processing up to `budget` entries. Each triggered order and
// each expired/orphaned index entry removed during traversal consumes one unit of work. It returns
// the triggered order ids plus the number of entries processed.
//
// Pessimistic rounding (matching PollTriggeredConditionalOrders):
//   - LTE direction: ceil(price) — only trigger if we are sure the price is truly ≤ triggerSubticks.
//   - GTE direction: floor(price) — only trigger if we are sure the price is truly ≥ triggerSubticks.
//
// The function scans LTE-crossed orders first, then GTE-crossed orders, matching the append order
// of the prior PollTriggeredConditionalOrders implementation.
func (k Keeper) triggerCrossedOrdersFromIndex(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	price *big.Rat,
	perpetualId uint32,
	priceType string,
	budget int,
) (triggeredOrderIds []types.OrderId, processed int) {
	triggeredOrderIds = make([]types.OrderId, 0)

	// Block time is the logical validity boundary even if the bounded expiry drain has not yet
	// physically removed the order from consensus state.
	blockTime := ctx.BlockTime()

	// Emit the price gauge (mirrors TriggerOrdersWithPrice metric).
	priceFloat, _ := price.Float32()
	labels := []gometrics.Label{
		metrics.GetLabelForStringValue(metrics.Type, priceType),
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
	}
	metrics.SetGaugeWithLabels(metrics.ClobConditionalOrderTriggerPrice, priceFloat, labels...)

	scanDirection := func(direction byte, priceSubticks uint64) {
		if processed >= budget {
			return
		}
		k.iterateCrossedConditionalOrders(
			ctx,
			uint32(clobPairId),
			direction,
			priceSubticks,
			func(orderId types.OrderId, indexKey []byte) bool {
				if processed >= budget {
					return false
				}
				processed++

				placement, found := k.GetUntriggeredConditionalOrderPlacement(ctx, orderId)
				if !found || placement.Order.IsStatefulOrderExpired(blockTime) {
					// Remove the stale secondary-index entry immediately. The bounded expiry queue
					// remains responsible for deleting an expired placement from consensus state.
					k.removeConditionalOrderTriggerPriceIndexEntry(ctx, indexKey)
					return processed < budget
				}

				k.MustTriggerConditionalOrder(ctx, orderId)
				k.GetIndexerEventManager().AddTxnEvent(
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewConditionalOrderTriggeredEvent(orderId),
					),
				)
				metrics.IncrCountMetricWithLabels(
					types.ModuleName,
					metrics.ClobConditionalOrderTriggered,
					append(orderId.GetOrderIdLabels(), labels...)...,
				)
				triggeredOrderIds = append(triggeredOrderIds, orderId)
				return processed < budget
			},
		)
	}

	// LTE-direction: oracle_price <= triggerSubticks. Pessimistic rounding uses ceil(price).
	scanDirection(TriggerDirectionLTE, lib.BigRatRound(price, true).Uint64())
	// GTE-direction: oracle_price >= triggerSubticks. Pessimistic rounding uses floor(price).
	scanDirection(TriggerDirectionGTE, lib.BigRatRound(price, false).Uint64())

	return triggeredOrderIds, processed
}

// TriggerOrdersWithPrice triggers all untriggered conditional orders using the given price. It returns
// a list of order ids that were triggered. This function is called in EndBlocker.
// It removes all triggered conditional orders from the `UntriggeredConditionalOrders ` struct.
func (k Keeper) TriggerOrdersWithPrice(
	ctx sdk.Context,
	untriggered *UntriggeredConditionalOrders,
	price *big.Rat,
	perpetualId uint32,
	priceType string,
) (triggeredOrderIds []types.OrderId) {
	triggeredOrderIds = untriggered.PollTriggeredConditionalOrders(price)

	// Emit metrics.
	priceFloat, _ := price.Float32()
	labels := []gometrics.Label{
		metrics.GetLabelForStringValue(metrics.Type, priceType),
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
	}
	metrics.SetGaugeWithLabels(metrics.ClobConditionalOrderTriggerPrice, priceFloat, labels...)

	// State write - move the conditional order placement in state from untriggered to triggered state.
	// Emit an event for each triggered conditional order.
	for _, orderId := range triggeredOrderIds {
		k.MustTriggerConditionalOrder(
			ctx,
			orderId,
		)
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexerevents.StatefulOrderEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewConditionalOrderTriggeredEvent(
					orderId,
				),
			),
		)

		metrics.IncrCountMetricWithLabels(
			types.ModuleName,
			metrics.ClobConditionalOrderTriggered,
			append(orderId.GetOrderIdLabels(), labels...)...,
		)
	}
	return triggeredOrderIds
}

// getBoundedTriggerOraclePrice returns the perpetual id and oracle price (in subticks) for the
// given clob pair for use by the bounded trigger path, or ok=false when the price cannot be used
// this block — the perpetual/market cannot be resolved, or the oracle price is 0 (a market that has
// never received a valid price update). Unlike the shared GetOraclePriceSubticksRat, this never
// panics: a price-0 pair simply has no meaningful crossing and is skipped, so a freshly listed pair
// at price 0 cannot halt the EndBlocker.
func (k Keeper) getBoundedTriggerOraclePrice(
	ctx sdk.Context,
	clobPair types.ClobPair,
) (perpetualId uint32, oraclePrice *big.Rat, ok bool) {
	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		return 0, nil, false
	}
	perpetual, marketPrice, err := k.perpetualsKeeper.GetPerpetualAndMarketPrice(ctx, perpetualId)
	if err != nil {
		return 0, nil, false
	}
	oraclePrice = types.PriceToSubticks(
		marketPrice,
		clobPair,
		perpetual.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	if oraclePrice.Sign() == 0 {
		return perpetualId, nil, false
	}
	return perpetualId, oraclePrice, true
}

func (k Keeper) getClampedTradePricesForTriggering(
	ctx sdk.Context,
	perpetualId uint32,
	oraclePrice *big.Rat,
) (
	clampedMinTradePrice *big.Rat,
	clampedMaxTradePrice *big.Rat,
	found bool,
) {
	minTradePriceSubticks, maxTradePriceSubticks, found := k.GetTradePricesForPerpetual(ctx, perpetualId)
	if found {
		// Get the perpetual.
		perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId)
		if err != nil {
			panic(
				fmt.Errorf(
					"EndBlocker: untriggeredConditionalOrders failed to find perpetualId %+v",
					perpetualId,
				),
			)
		}

		// Get the market param.
		marketParam, exists := k.pricesKeeper.GetMarketParam(ctx, perpetual.Params.MarketId)
		if !exists {
			panic(
				fmt.Errorf(
					"EndBlocker: untriggeredConditionalOrders failed to find marketParam %+v",
					perpetual.Params.MarketId,
				),
			)
		}

		// Calculate the max allowed range.
		maxAllowedRange := lib.BigRatMulPpm(oraclePrice, marketParam.MinPriceChangePpm)
		maxAllowedRange.Mul(maxAllowedRange, new(big.Rat).SetUint64(types.ConditionalOrderTriggerMultiplier))

		upperBound := new(big.Rat).Add(oraclePrice, maxAllowedRange)
		lowerBound := new(big.Rat).Sub(oraclePrice, maxAllowedRange)

		// Clamp the min and max trade prices to the upper and lower bounds.
		clampedMinTradePrice = lib.BigRatClamp(
			new(big.Rat).SetUint64(minTradePriceSubticks.ToUint64()),
			lowerBound,
			upperBound,
		)
		clampedMaxTradePrice = lib.BigRatClamp(
			new(big.Rat).SetUint64(maxTradePriceSubticks.ToUint64()),
			lowerBound,
			upperBound,
		)
	}
	return clampedMinTradePrice, clampedMaxTradePrice, found
}
