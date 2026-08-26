package keeper

import (
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	gometrics "github.com/hashicorp/go-metrics"
)

// MaxConditionalTriggersPerBlock is the maximum number of conditional orders that may be
// triggered in a single block by MaybeTriggerConditionalOrders. This acts as a per-block
// budget that bounds EndBlocker work to O(budget) rather than O(global N). A conservative
// high default is chosen so that normal operation and equivalence tests are unaffected; the
// value is governance-tunable via MsgUpdateConditionalOrderTriggerConfig.
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
// It runs the bounded, crossing-priority path (maybeTriggerConditionalOrdersBounded), which uses
// the trigger-price secondary index to do O(crossed + budget) work per block instead of
// O(global N). The trigger-price index is built by the state-breaking upgrade handler and
// maintained incrementally thereafter, so the bounded path is always authoritative.
func (k Keeper) MaybeTriggerConditionalOrders(ctx sdk.Context) (allTriggeredOrderIds []types.OrderId) {
	defer metrics.ModuleMeasureSince(
		types.ModuleName,
		metrics.ClobMaybeTriggerConditionalOrders,
		time.Now(),
	)

	return k.maybeTriggerConditionalOrdersBounded(ctx, k.GetConditionalOrderTriggerConfig(ctx))
}

// maybeTriggerConditionalOrdersBounded is the bounded trigger path.
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
