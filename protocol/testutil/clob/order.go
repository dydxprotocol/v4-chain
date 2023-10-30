package clob

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type OrderModifierOption func(cp *clobtypes.Order)

func WithGTB(gtb uint32) OrderModifierOption {
	return func(o *clobtypes.Order) {
		o.GoodTilOneof = &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: gtb,
		}
	}
}

func WithSide(side clobtypes.Order_Side) OrderModifierOption {
	return func(o *clobtypes.Order) {
		o.Side = side
	}
}

func WithClobPairid(id uint32) OrderModifierOption {
	return func(o *clobtypes.Order) {
		o.OrderId.ClobPairId = id
	}
}

func WithSubaccountId(subaccountId satypes.SubaccountId) OrderModifierOption {
	return func(o *clobtypes.Order) {
		o.OrderId.SubaccountId = subaccountId
	}
}

func WithClientId(clientId uint32) OrderModifierOption {
	return func(o *clobtypes.Order) {
		o.OrderId.ClientId = clientId
	}
}

// GenerateOrderUsingTemplate is a helper function to generate an test order with a template and
// opitonal modifier options.
// Example usage:
//
//	  clobtest.GenerateOrderUsingTemplate(
//	    OrderTemplate_ShortTerm_Btc,
//	    clobtest.WithSide(clobtypes.Order_SIDE_SELL),
//		clobtest.WithSubaccountId(Alice_Num0),
//		clobtest.WithClobPairid(TestEthMarketId),
//		clobtest.WithGTB(TestGTB),
//	  )
func GenerateOrderUsingTemplate(order clobtypes.Order, optionalModifications ...OrderModifierOption) clobtypes.Order {
	for _, opt := range optionalModifications {
		opt(&order)
	}

	return order
}
