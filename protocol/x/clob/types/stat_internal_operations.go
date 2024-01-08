package types

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	gometrics "github.com/hashicorp/go-metrics"
)

// OperationsStats is a struct that holds stats about a list of
// operations that were proposed in a block.
type OperationsStats struct {
	// Stat fields.
	MatchedShortTermOrdersCount            uint
	MatchedLongTermOrdersCount             uint
	MatchedConditionalOrdersCount          uint
	TakerOrdersCount                       uint
	LiquidationOrdersCount                 uint
	DeleveragingOperationsCount            uint
	TotalFillsCount                        uint
	LongTermOrderRemovalsCount             uint
	ConditionalOrderRemovalsCount          uint
	UniqueSubaccountsLiquidated            uint
	UniqueSubaccountsDeleveraged           uint
	UniqueSubaccountsOffsettingDeleveraged uint

	// Internal fields for calculating stats.
	uniqueSubaccountsLiquidated            map[satypes.SubaccountId]bool
	uniqueSubaccountsDeleveraged           map[satypes.SubaccountId]bool
	uniqueSubaccountsOffsettingDeleveraged map[satypes.SubaccountId]bool
	uniqueMatchedOrderIds                  map[OrderId]bool
}

// StatMsgProposedOperations generates statistics from a list of internal operations
// and returns them in an `OperationsStats` struct.
func StatMsgProposedOperations(
	rawOperations []OperationRaw,
) OperationsStats {
	stats := OperationsStats{
		uniqueSubaccountsLiquidated:            make(map[satypes.SubaccountId]bool),
		uniqueSubaccountsDeleveraged:           make(map[satypes.SubaccountId]bool),
		uniqueSubaccountsOffsettingDeleveraged: make(map[satypes.SubaccountId]bool),
		uniqueMatchedOrderIds:                  make(map[OrderId]bool),
	}

	for _, operation := range rawOperations {
		switch castedOperation := operation.Operation.(type) {
		case *OperationRaw_Match:
			switch castedOperation.Match.Match.(type) {
			case *ClobMatch_MatchOrders:
				stats.TakerOrdersCount++

				stats.statMatchedOrderId(castedOperation.Match.GetMatchOrders().GetTakerOrderId())
				for _, makerOrderId := range castedOperation.Match.GetMatchOrders().GetFills() {
					stats.TotalFillsCount++
					stats.statMatchedOrderId(makerOrderId.GetMakerOrderId())
				}
			case *ClobMatch_MatchPerpetualLiquidation:
				stats.LiquidationOrdersCount++

				for _, makerOrderId := range castedOperation.Match.GetMatchPerpetualLiquidation().GetFills() {
					stats.TotalFillsCount++
					stats.statMatchedOrderId(makerOrderId.GetMakerOrderId())
				}

				liquidated := castedOperation.Match.GetMatchPerpetualLiquidation().GetLiquidated()
				if _, exists := stats.uniqueSubaccountsLiquidated[liquidated]; !exists {
					stats.UniqueSubaccountsLiquidated++
					stats.uniqueSubaccountsLiquidated[liquidated] = true
				}
			case *ClobMatch_MatchPerpetualDeleveraging:
				stats.DeleveragingOperationsCount++

				deleveraged := castedOperation.Match.GetMatchPerpetualDeleveraging().GetLiquidated()
				if _, exists := stats.uniqueSubaccountsDeleveraged[deleveraged]; !exists {
					stats.UniqueSubaccountsDeleveraged++
					stats.uniqueSubaccountsDeleveraged[deleveraged] = true
				}

				for _, makerFill := range castedOperation.Match.GetMatchPerpetualDeleveraging().GetFills() {
					offsetting := makerFill.OffsettingSubaccountId
					if _, exists := stats.
						uniqueSubaccountsOffsettingDeleveraged[offsetting]; !exists {
						stats.UniqueSubaccountsOffsettingDeleveraged++
						stats.uniqueSubaccountsOffsettingDeleveraged[offsetting] = true
					}
				}
			}
		case *OperationRaw_OrderRemoval:
			orderId := castedOperation.OrderRemoval.GetOrderId()
			if orderId.IsConditionalOrder() {
				stats.ConditionalOrderRemovalsCount++
			}

			if orderId.IsLongTermOrder() {
				stats.LongTermOrderRemovalsCount++
			}
		}
	}

	return stats
}

// EmitStats emits stats about the internal operations. It includes the provided ABCI callback as a
// label in the stat.
func (stats *OperationsStats) EmitStats(abciCallback string) {
	type Stat struct {
		keys  []string
		value float32
	}

	statsList := []Stat{
		{
			keys:  []string{ModuleName, metrics.NumMatchedShortTermOrders},
			value: float32(stats.MatchedShortTermOrdersCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumMatchedLongTermOrders},
			value: float32(stats.MatchedLongTermOrdersCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumMatchedConditionalOrders},
			value: float32(stats.MatchedConditionalOrdersCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumMatchTakerOrders},
			value: float32(stats.TakerOrdersCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumMatchedLiquidationOrders},
			value: float32(stats.LiquidationOrdersCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumMatchPerpDeleveragingOperations},
			value: float32(stats.DeleveragingOperationsCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumFills},
			value: float32(stats.TotalFillsCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumLongTermOrderRemovals},
			value: float32(stats.LongTermOrderRemovalsCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumConditionalOrderRemovals},
			value: float32(stats.ConditionalOrderRemovalsCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumUniqueSubaccountsLiquidated},
			value: float32(stats.UniqueSubaccountsLiquidated),
		},
		{
			keys:  []string{ModuleName, metrics.NumUniqueSubaccountsDeleveraged},
			value: float32(stats.UniqueSubaccountsDeleveraged),
		},
		{
			keys:  []string{ModuleName, metrics.NumUniqueSubaccountsOffsettingDeleveraged},
			value: float32(stats.UniqueSubaccountsOffsettingDeleveraged),
		},
	}

	labels := []gometrics.Label{
		metrics.GetLabelForStringValue(
			metrics.Callback,
			abciCallback,
		),
	}
	for _, stat := range statsList {
		gometrics.AddSampleWithLabels(
			stat.keys,
			stat.value,
			labels,
		)
	}
}

// statMatchedOrderId updates the number of unique matched order IDs. If the order ID was already
// matched it will early return. Note this function should only be called with matched order IDs.
func (stats *OperationsStats) statMatchedOrderId(orderId OrderId) {
	if _, exists := stats.uniqueMatchedOrderIds[orderId]; exists {
		return
	}

	if orderId.IsShortTermOrder() {
		stats.MatchedShortTermOrdersCount++
	}

	if orderId.IsLongTermOrder() {
		stats.MatchedLongTermOrdersCount++
	}

	if orderId.IsConditionalOrder() {
		stats.MatchedConditionalOrdersCount++
	}

	stats.uniqueMatchedOrderIds[orderId] = true
}
