package types

import (
	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// OperationsStats is a struct that holds stats about a list of
// operations that were proposed in a block.
type OperationsStats struct {
	// Stat fields.
	MatchedShortTermOrdersCount            uint
	MatchedLongTermOrdersCount             uint
	MatchedConditionalOrdersCount          uint
	RegularMatchesCount                    uint
	LiquidationMatchesCount                uint
	DeleveragingMatchesCount               uint
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
				stats.RegularMatchesCount++

				stats.StatMatchedOrderId(castedOperation.Match.GetMatchOrders().GetTakerOrderId())
				for _, makerOrderId := range castedOperation.Match.GetMatchOrders().GetFills() {
					stats.StatMatchedOrderId(makerOrderId.GetMakerOrderId())
				}
			case *ClobMatch_MatchPerpetualLiquidation:
				stats.LiquidationMatchesCount++

				for _, makerOrderId := range castedOperation.Match.GetMatchPerpetualLiquidation().GetFills() {
					stats.StatMatchedOrderId(makerOrderId.GetMakerOrderId())
				}

				liquidated := castedOperation.Match.GetMatchPerpetualLiquidation().GetLiquidated()
				if _, exists := stats.uniqueSubaccountsLiquidated[liquidated]; !exists {
					stats.UniqueSubaccountsLiquidated++
					stats.uniqueSubaccountsLiquidated[liquidated] = true
				}
			case *ClobMatch_MatchPerpetualDeleveraging:
				stats.DeleveragingMatchesCount++

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
			keys:  []string{ModuleName, metrics.NumFills},
			value: float32(stats.RegularMatchesCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumMatchedLiquidationOrders},
			value: float32(stats.LiquidationMatchesCount),
		},
		{
			keys:  []string{ModuleName, metrics.NumMatchPerpDeleveragingOperations},
			value: float32(stats.DeleveragingMatchesCount),
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

	for _, stat := range statsList {
		telemetry.SetGaugeWithLabels(
			stat.keys,
			stat.value,
			[]gometrics.Label{
				metrics.GetLabelForStringValue(
					metrics.Callback,
					abciCallback,
				),
			},
		)
	}
}

// StatMatchedOrderId updates the number of unique matched order IDs. If the order ID was already
// matched it will early return. Note this function should only be called with matched order IDs.
func (stats *OperationsStats) StatMatchedOrderId(orderId OrderId) {
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
