package memclob

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func getOrderbook() *Orderbook {
	memclob := NewMemClobPriceTimePriority(false)
	memclob.CreateOrderbook(constants.ClobPair_Btc)
	return memclob.mustGetOrderbook(0)
}

func TestOrderbook_Remove_SingleCancelInTilBlock(t *testing.T) {
	c := getOrderbook()

	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	c.orderIdToCancelExpiry[order.OrderId] = order.GetGoodTilBlock()
	c.cancelExpiryToOrderIds[order.GetGoodTilBlock()] = map[types.OrderId]bool{
		order.OrderId: true,
	}

	c.mustRemoveCancel(order.OrderId)

	require.Empty(t, c.orderIdToCancelExpiry)
	require.Empty(t, c.cancelExpiryToOrderIds)
}

func TestMemclobCancels_Remove_TwoCancelsInTilBlock(t *testing.T) {
	c := getOrderbook()

	// TODO(DEC-124): replace with `AddCancel`
	order1 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	c.orderIdToCancelExpiry[order1.OrderId] = order1.GetGoodTilBlock()

	order2 := constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15
	c.orderIdToCancelExpiry[order2.OrderId] = order2.GetGoodTilBlock()

	c.cancelExpiryToOrderIds[order2.GetGoodTilBlock()] = map[types.OrderId]bool{
		order1.OrderId: true,
		order2.OrderId: true,
	}

	c.mustRemoveCancel(order1.OrderId)

	require.Len(t, c.orderIdToCancelExpiry, 1)
	require.NotContains(t, c.orderIdToCancelExpiry, order1.OrderId)
	require.Len(t, c.cancelExpiryToOrderIds, 1)
	require.Len(t, c.cancelExpiryToOrderIds[order1.GetGoodTilBlock()], 1)
	require.NotContains(t, c.orderIdToCancelExpiry, order1.OrderId)
}

func TestMemclobCancels_Remove_PanicsIfGoodTilBlockDoesNotExistInOrderIdToExpiry(t *testing.T) {
	c := getOrderbook()

	order1 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15

	require.Empty(t, c.orderIdToCancelExpiry)
	require.Empty(t, c.cancelExpiryToOrderIds)
	require.Panics(t, func() {
		c.mustRemoveCancel(order1.OrderId)
	})
}

func TestMemclobCancels_Remove_PanicsIfGoodTilBlockDoesNotExistInExpiryToOrderIds(t *testing.T) {
	c := getOrderbook()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.orderIdToCancelExpiry[orderId] = goodTilBlock

	// Removing cancel should panic.
	require.Panics(t, func() {
		c.mustRemoveCancel(orderId)
	})
}

func TestMemclobCancels_Remove_PanicsIfOrderIdDoesNotExistInExpiryToOrderIdsSubmap(t *testing.T) {
	c := getOrderbook()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.orderIdToCancelExpiry[orderId] = goodTilBlock
	c.cancelExpiryToOrderIds[goodTilBlock] = make(map[types.OrderId]bool)

	// Removing cancel should panic.
	require.Panics(t, func() {
		c.mustRemoveCancel(orderId)
	})
}

func TestMemclobCancels_Add_PanicsIfAlreadyExistsInOrderIdToExpiry(t *testing.T) {
	c := getOrderbook()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.orderIdToCancelExpiry[orderId] = goodTilBlock

	// Canceling again should panic.
	require.Panics(t, func() {
		c.addShortTermCancel(orderId, goodTilBlock)
	})
}

func TestMemclobCancels_Add_PanicsIfAlreadyExistsInExpiryToOrderIds(t *testing.T) {
	c := getOrderbook()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.cancelExpiryToOrderIds[goodTilBlock] = make(map[types.OrderId]bool)
	c.cancelExpiryToOrderIds[goodTilBlock][orderId] = true

	// Canceling again should panic.
	require.Panics(t, func() {
		c.addShortTermCancel(orderId, goodTilBlock)
	})
}

func TestOrderbook_TakerCrossesBook(t *testing.T) {
	memclob := NewMemClobPriceTimePriority(false)
	memclob.CreateOrderbook(constants.ClobPair_Btc)
	ob := memclob.mustGetOrderbook(constants.ClobPair_Btc.GetClobPairId())
	buy := func(subticks uint64) types.MatchableOrder {
		o := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
		o.Subticks = subticks
		return &o
	}
	sell := func(subticks uint64) types.MatchableOrder {
		o := constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22
		o.Subticks = subticks
		return &o
	}

	require.False(t, ob.takerCrossesBook(buy(100)), "empty book")
	require.False(t, ob.takerCrossesBook(sell(100)), "empty book")

	// Resting ask at 100, resting bid at 90.
	ask := constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22
	ask.Subticks = 100
	ob.mustAddOrderToOrderbook(ask, false)
	bid := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	bid.Subticks = 90
	ob.mustAddOrderToOrderbook(bid, false)

	require.True(t, ob.takerCrossesBook(buy(100)), "buy at the best ask crosses")
	require.True(t, ob.takerCrossesBook(buy(150)), "buy through the ask crosses")
	require.False(t, ob.takerCrossesBook(buy(99)), "buy below the ask rests")
	require.True(t, ob.takerCrossesBook(sell(90)), "sell at the best bid crosses")
	require.True(t, ob.takerCrossesBook(sell(1)), "sell through the bid crosses")
	require.False(t, ob.takerCrossesBook(sell(91)), "sell above the bid rests")
}
