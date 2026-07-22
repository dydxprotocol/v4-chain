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
// (lowest triggerSubticks for LTE, lowest triggerSubticks for GTE) are triggered this block;
// the remainder are deferred to subsequent blocks where the budget refills, making deferral
// deterministic and consistent across all nodes.
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
	if !cfg.Enabled {
		// Legacy path — identical to the pre-fix behavior. Rolling-deploy safe.
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
//   - The per-block budget is cfg.MaxTriggersPerBlock, drained nearest-crossing-first across clob
//     pairs in ascending id.
//   - Any orders beyond the budget are deterministically deferred to subsequent blocks (the index
//     is ordered, so the same nearest-crossing prefix is chosen on every node every block).
//
// Ordering: (clobPairId ascending, oracle→min-trade→max-trade, LTE→GTE, ascending
// triggerSubticks/orderId). Deterministic across nodes. Differs from legacy time-priority ordering,
// which is why enabling the flag must happen at a coordinated height.
func (k Keeper) maybeTriggerConditionalOrdersBounded(
	ctx sdk.Context,
	cfg ConditionalOrderTriggerConfig,
) (allTriggeredOrderIds []types.OrderId) {
	allTriggeredOrderIds = make([]types.OrderId, 0)
	remaining := int(cfg.MaxTriggersPerBlock)

	// GetAllClobPairs already returns clob pairs sorted in ascending ClobPair.Id order,
	// preserving the deterministic per-pair iteration order of the prior implementation.
	clobPairs := k.GetAllClobPairs(ctx)

	for _, clobPair := range clobPairs {
		if remaining <= 0 {
			break
		}

		clobPairId := types.ClobPairId(clobPair.Id)
		perpetualId := clobPair.MustGetPerpetualId()
		oraclePrice := k.GetOraclePriceSubticksRat(ctx, clobPair)

		// Trigger conditional orders using the oracle price.
		triggered := k.triggerCrossedOrdersFromIndex(
			ctx, clobPairId, oraclePrice, perpetualId, metrics.OraclePrice, remaining,
		)
		allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)
		remaining -= len(triggered)

		if remaining <= 0 {
			continue
		}

		// Trigger conditional orders using the clamped trade prices.
		clampedMinTradePrice,
			clampedMaxTradePrice,
			found := k.getClampedTradePricesForTriggering(
			ctx,
			perpetualId,
			oraclePrice,
		)

		if found {
			triggered = k.triggerCrossedOrdersFromIndex(
				ctx, clobPairId, clampedMinTradePrice, perpetualId, metrics.MinTradePrice, remaining,
			)
			allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)
			remaining -= len(triggered)

			if remaining > 0 {
				triggered = k.triggerCrossedOrdersFromIndex(
					ctx, clobPairId, clampedMaxTradePrice, perpetualId, metrics.MaxTradePrice, remaining,
				)
				allTriggeredOrderIds = append(allTriggeredOrderIds, triggered...)
				remaining -= len(triggered)
			}
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

// triggerCrossedOrdersFromIndex performs a crossed-order range scan on the trigger-price index
// for a single clob pair and price, triggering up to `budget` orders. It returns the list of
// triggered order ids in the order they were visited (ascending triggerSubticks, then orderId
// bytes). This is O(crossed) — it reads only the orders whose trigger boundary the given price
// crosses, not the full untriggered set.
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
) (triggeredOrderIds []types.OrderId) {
	triggeredOrderIds = make([]types.OrderId, 0)

	// Emit the price gauge (mirrors TriggerOrdersWithPrice metric).
	priceFloat, _ := price.Float32()
	labels := []gometrics.Label{
		metrics.GetLabelForStringValue(metrics.Type, priceType),
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
	}
	metrics.SetGaugeWithLabels(metrics.ClobConditionalOrderTriggerPrice, priceFloat, labels...)

	// LTE-direction: orders trigger when oracle_price ≤ triggerSubticks.
	// Pessimistic rounding: ceil(price) for the subticks threshold.
	if budget > 0 {
		ltePriceSubticks := lib.BigRatRound(price, true).Uint64() // ceil
		k.IterateCrossedConditionalOrders(
			ctx,
			uint32(clobPairId),
			TriggerDirectionLTE,
			ltePriceSubticks,
			func(orderId types.OrderId) bool {
				if budget <= 0 {
					return false
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
				budget--
				return true
			},
		)
	}

	// GTE-direction: orders trigger when oracle_price ≥ triggerSubticks.
	// Pessimistic rounding: floor(price) for the subticks threshold.
	if budget > 0 {
		gtePriceSubticks := lib.BigRatRound(price, false).Uint64() // floor
		k.IterateCrossedConditionalOrders(
			ctx,
			uint32(clobPairId),
			TriggerDirectionGTE,
			gtePriceSubticks,
			func(orderId types.OrderId) bool {
				if budget <= 0 {
					return false
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
				budget--
				return true
			},
		)
	}

	return triggeredOrderIds
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
