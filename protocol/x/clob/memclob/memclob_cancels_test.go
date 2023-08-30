package memclob

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestMemclobCancels_Remove_SingleCancelInTilBlock(t *testing.T) {
	c := newMemclobCancels()

	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	c.orderIdToExpiry[order.OrderId] = order.GetGoodTilBlock()
	c.expiryToOrderIds[order.GetGoodTilBlock()] = map[types.OrderId]bool{
		order.OrderId: true,
	}

	c.remove(order.OrderId)

	require.Empty(t, c.orderIdToExpiry)
	require.Empty(t, c.expiryToOrderIds)
}

func TestMemclobCancels_Remove_TwoCancelsInTilBlock(t *testing.T) {
	c := newMemclobCancels()

	// TODO(DEC-124): replace with `AddCancel`
	order1 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	c.orderIdToExpiry[order1.OrderId] = order1.GetGoodTilBlock()

	order2 := constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15
	c.orderIdToExpiry[order2.OrderId] = order2.GetGoodTilBlock()

	c.expiryToOrderIds[order2.GetGoodTilBlock()] = map[types.OrderId]bool{
		order1.OrderId: true,
		order2.OrderId: true,
	}

	c.remove(order1.OrderId)

	require.Len(t, c.orderIdToExpiry, 1)
	require.NotContains(t, c.orderIdToExpiry, order1.OrderId)
	require.Len(t, c.expiryToOrderIds, 1)
	require.Len(t, c.expiryToOrderIds[order1.GetGoodTilBlock()], 1)
	require.NotContains(t, c.orderIdToExpiry, order1.OrderId)
}

func TestMemclobCancels_Remove_PanicsIfGoodTilBlockDoesNotExistInOrderIdToExpiry(t *testing.T) {
	c := newMemclobCancels()

	order1 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15

	require.Empty(t, c.orderIdToExpiry)
	require.Empty(t, c.expiryToOrderIds)
	require.Panics(t, func() {
		c.remove(order1.OrderId)
	})
}

func TestMemclobCancels_Remove_PanicsIfGoodTilBlockDoesNotExistInExpiryToOrderIds(t *testing.T) {
	c := newMemclobCancels()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.orderIdToExpiry[orderId] = goodTilBlock

	// Removing cancel should panic.
	require.Panics(t, func() {
		c.remove(orderId)
	})
}

func TestMemclobCancels_Remove_PanicsIfOrderIdDoesNotExistInExpiryToOrderIdsSubmap(t *testing.T) {
	c := newMemclobCancels()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.orderIdToExpiry[orderId] = goodTilBlock
	c.expiryToOrderIds[goodTilBlock] = make(map[types.OrderId]bool)

	// Removing cancel should panic.
	require.Panics(t, func() {
		c.remove(orderId)
	})
}

func TestMemclobCancels_Add_PanicsIfAlreadyExistsInOrderIdToExpiry(t *testing.T) {
	c := newMemclobCancels()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.orderIdToExpiry[orderId] = goodTilBlock

	// Canceling again should panic.
	require.Panics(t, func() {
		c.addShortTermCancel(orderId, goodTilBlock)
	})
}

func TestMemclobCancels_Add_PanicsIfAlreadyExistsInExpiryToOrderIds(t *testing.T) {
	c := newMemclobCancels()

	// Setup.
	orderId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId
	goodTilBlock := uint32(10)
	c.expiryToOrderIds[goodTilBlock] = make(map[types.OrderId]bool)
	c.expiryToOrderIds[goodTilBlock][orderId] = true

	// Canceling again should panic.
	require.Panics(t, func() {
		c.addShortTermCancel(orderId, goodTilBlock)
	})
}
