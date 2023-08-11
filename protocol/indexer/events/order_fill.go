package events

import (
	"fmt"

	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// NewOrderFillEvent creates a new OrderFillEvent proto message given the maker and taker orders along
// with the fill amount. Note: This function does no validation of the input maker/taker orders or
// the fill amount and assumes all such validation has been done before constructing the event.
func NewOrderFillEvent(
	makerOrder clobtypes.Order,
	takerOrder clobtypes.Order,
	fillAmount satypes.BaseQuantums,
) *OrderFillEvent {
	return &OrderFillEvent{
		MakerOrder: makerOrder,
		TakerOrder: &OrderFillEvent_Order{
			Order: &takerOrder,
		},
		FillAmount: fillAmount.ToUint64(),
	}
}

// NewLiquidationOrderFillEvent creates a new OrderFillEvent proto message given the maker and liquidation
// taker orders along with the fill amount. Panics if the taker order is not a liquidation order.
func NewLiquidationOrderFillEvent(
	makerOrder clobtypes.Order,
	liquidationTakerOrder clobtypes.MatchableOrder,
	fillAmount satypes.BaseQuantums,
) *OrderFillEvent {
	if !liquidationTakerOrder.IsLiquidation() {
		panic(fmt.Sprintf("liquidationTakerOrder is not a liquidation order: %v", liquidationTakerOrder))
	}
	liquidationOrder := LiquidationOrder{
		Liquidated:  liquidationTakerOrder.GetSubaccountId(),
		ClobPairId:  liquidationTakerOrder.GetClobPairId().ToUint32(),
		PerpetualId: liquidationTakerOrder.MustGetLiquidatedPerpetualId(),
		TotalSize:   uint64(liquidationTakerOrder.GetBaseQuantums()),
		IsBuy:       liquidationTakerOrder.IsBuy(),
		Subticks:    uint64(liquidationTakerOrder.GetOrderSubticks()),
	}
	return &OrderFillEvent{
		MakerOrder: makerOrder,
		TakerOrder: &OrderFillEvent_LiquidationOrder{LiquidationOrder: &liquidationOrder},
		FillAmount: fillAmount.ToUint64(),
	}
}
