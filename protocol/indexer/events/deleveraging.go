package events

import (
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// NewDeleveragingEvent creates a DeleveragingEvent representing a deleveraging
// where a liquidated subaccount's position is offset by another subaccount.
func NewDeleveragingEvent(
	liquidatedSubaccountId satypes.SubaccountId,
	offsettingSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	fillAmount satypes.BaseQuantums,
	price satypes.BaseQuantums,
	isBuy bool,
) *DeleveragingEventV1 {
	indexerLiquidatedSubaccountId := v1.SubaccountIdToIndexerSubaccountId(liquidatedSubaccountId)
	indexerOffsettingSubaccountId := v1.SubaccountIdToIndexerSubaccountId(offsettingSubaccountId)
	return &DeleveragingEventV1{
		Liquidated:  indexerLiquidatedSubaccountId,
		Offsetting:  indexerOffsettingSubaccountId,
		PerpetualId: perpetualId,
		FillAmount:  fillAmount.ToUint64(),
		Price:       price.ToUint64(),
		IsBuy:       isBuy,
	}
}
