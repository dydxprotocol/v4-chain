package events

import (
	"fmt"
	"math/big"

	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// NewOrderFillEvent creates a new OrderFillEvent proto message given the maker and taker orders along
// with the fill and fee amounts. Note: This function does no validation of the input maker/taker orders
// or the fill amount and assumes all such validation has been done before constructing the event.
func NewOrderFillEvent(
	makerOrder clobtypes.Order,
	takerOrder clobtypes.Order,
	fillAmount satypes.BaseQuantums,
	makerFee int64,
	takerFee int64,
	makerBuilderFee uint64,
	takerBuilderFee uint64,
	totalFilledMaker satypes.BaseQuantums,
	totalFilledTaker satypes.BaseQuantums,
	affiliateRevShareQuoteQuantums *big.Int,
	makerOrderRouterFee uint64,
	takerOrderRouterFee uint64,
) *OrderFillEventV1 {
	indexerTakerOrder := v1.OrderToIndexerOrder(takerOrder)
	makerBuilderAddress := getBuilderAddress(makerOrder)
	takerBuilderAddress := getBuilderAddress(takerOrder)
	return &OrderFillEventV1{
		MakerOrder: v1.OrderToIndexerOrder(makerOrder),
		TakerOrder: &OrderFillEventV1_Order{
			Order: &indexerTakerOrder,
		},
		FillAmount:          fillAmount.ToUint64(),
		MakerFee:            makerFee,
		TakerFee:            takerFee,
		MakerBuilderFee:     makerBuilderFee,
		TakerBuilderFee:     takerBuilderFee,
		MakerBuilderAddress: makerBuilderAddress,
		TakerBuilderAddress: takerBuilderAddress,
		TotalFilledMaker:    totalFilledMaker.ToUint64(),
		TotalFilledTaker:    totalFilledTaker.ToUint64(),
		// Since revshare is always less than taker fee, this will not overflow.
		AffiliateRevShare:       affiliateRevShareQuoteQuantums.Uint64(),
		MakerOrderRouterFee:     makerOrderRouterFee,
		TakerOrderRouterFee:     takerOrderRouterFee,
		MakerOrderRouterAddress: makerOrder.GetOrderRouterAddress(),
		TakerOrderRouterAddress: takerOrder.GetOrderRouterAddress(),
	}
}

// NewLiquidationOrderFillEvent creates a new OrderFillEvent proto message given the maker and liquidation
// taker orders along with the fill and fee amounts. Panics if the taker order is not a liquidation order.
// The taker fee here refers to the special liquidation fee, not the standard taker fee.
func NewLiquidationOrderFillEvent(
	makerOrder clobtypes.Order,
	liquidationTakerOrder clobtypes.MatchableOrder,
	fillAmount satypes.BaseQuantums,
	makerFee int64,
	takerFee int64,
	makerBuilderFee uint64,
	totalFilledMaker satypes.BaseQuantums,
	affiliateRevShareQuoteQuantums *big.Int,
	makerOrderRouterFee uint64,
) *OrderFillEventV1 {
	if !liquidationTakerOrder.IsLiquidation() {
		panic(fmt.Sprintf("liquidationTakerOrder is not a liquidation order: %v", liquidationTakerOrder))
	}

	liquidationOrder := LiquidationOrderV1{
		Liquidated:  v1.SubaccountIdToIndexerSubaccountId(liquidationTakerOrder.GetSubaccountId()),
		ClobPairId:  liquidationTakerOrder.GetClobPairId().ToUint32(),
		PerpetualId: liquidationTakerOrder.MustGetLiquidatedPerpetualId(),
		TotalSize:   uint64(liquidationTakerOrder.GetBaseQuantums()),
		IsBuy:       liquidationTakerOrder.IsBuy(),
		Subticks:    uint64(liquidationTakerOrder.GetOrderSubticks()),
	}

	makerBuilderAddress := getBuilderAddress(makerOrder)
	return &OrderFillEventV1{
		MakerOrder:       v1.OrderToIndexerOrder(makerOrder),
		TakerOrder:       &OrderFillEventV1_LiquidationOrder{LiquidationOrder: &liquidationOrder},
		FillAmount:       fillAmount.ToUint64(),
		MakerFee:         makerFee,
		TakerFee:         takerFee,
		TotalFilledMaker: totalFilledMaker.ToUint64(),
		TotalFilledTaker: fillAmount.ToUint64(),
		// Since revshare is always less than taker fee, this will not overflow.
		AffiliateRevShare:       affiliateRevShareQuoteQuantums.Uint64(),
		MakerBuilderAddress:     makerBuilderAddress,
		MakerBuilderFee:         makerBuilderFee,
		TakerBuilderFee:         0, // protocol generated liquidation orders have no builder fee
		MakerOrderRouterAddress: makerOrder.GetOrderRouterAddress(),
		TakerOrderRouterFee:     0, // protocol generated liquidation orders have no order router fee
	}
}

func getBuilderAddress(order clobtypes.Order) string {
	builderAddress := ""
	if order.GetBuilderCodeParameters() != nil {
		builderAddress = order.GetBuilderCodeParameters().GetBuilderAddress()
	}
	return builderAddress
}
