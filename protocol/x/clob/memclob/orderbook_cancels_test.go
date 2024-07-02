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
