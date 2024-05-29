package memclob

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"github.com/zyedidia/generic/list"
)

func TestGetSubaccountOrders(t *testing.T) {
	tests := map[string]struct {
		// State.
		memclobOrders []types.Order

		// GetSubaccountOrders parameters.
		clobPairId   uint32
		subaccountId satypes.SubaccountId
		side         types.Order_Side

		// Expectations.
		expectedOpenOrders []types.Order
		expectedErr        error
	}{
		"Returns nothing when there are no open orders": {
			memclobOrders: []types.Order{},

			clobPairId:   0,
			subaccountId: constants.Alice_Num0,
			side:         types.Order_SIDE_BUY,

			expectedOpenOrders: []types.Order{},
			expectedErr:        nil,
		},
		"Returns nothing when a subaccount has no open orders, but orders exist on the CLOB": {
			memclobOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			clobPairId:   0,
			subaccountId: constants.Alice_Num1,
			side:         types.Order_SIDE_BUY,

			expectedOpenOrders: []types.Order{},
			expectedErr:        nil,
		},
		`Returns nothing when a subaccount has no open orders on that side,
		but orders exist on the other side of that CLOB`: {
			memclobOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			clobPairId:   0,
			subaccountId: constants.Alice_Num0,
			side:         types.Order_SIDE_SELL,

			expectedOpenOrders: []types.Order{},
			expectedErr:        nil,
		},
		"Returns a users open order on a side": {
			memclobOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			clobPairId:   0,
			subaccountId: constants.Alice_Num0,
			side:         types.Order_SIDE_BUY,

			expectedOpenOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			expectedErr: nil,
		},
		"Does not return orders on the other side of the book": {
			memclobOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},

			clobPairId:   0,
			subaccountId: constants.Alice_Num0,
			side:         types.Order_SIDE_BUY,

			expectedOpenOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			expectedErr: nil,
		},
		"Does not return orders from other subaccounts": {
			memclobOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
			},

			clobPairId:   0,
			subaccountId: constants.Alice_Num0,
			side:         types.Order_SIDE_BUY,

			expectedOpenOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			expectedErr: nil,
		},
		"Does not return orders from other CLOBs": {
			memclobOrders: []types.Order{
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},

			clobPairId:   1,
			subaccountId: constants.Alice_Num1,
			side:         types.Order_SIDE_SELL,

			expectedOpenOrders: []types.Order{
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			expectedErr: nil,
		},
		"Returns multiple of a users open orders on a side, ignoring irrelevant orders": {
			memclobOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,
				constants.Order_Bob_Num0_Id2_Clob1_Sell12_Price13_GTB20,
			},

			clobPairId:   1,
			subaccountId: constants.Bob_Num0,
			side:         types.Order_SIDE_SELL,

			expectedOpenOrders: []types.Order{
				constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,
				constants.Order_Bob_Num0_Id2_Clob1_Sell12_Price13_GTB20,
			},
			expectedErr: nil,
		},
		"Returns an error when an invalid side is passed as a parameter": {
			memclobOrders: []types.Order{},

			clobPairId:   0,
			subaccountId: constants.Alice_Num0,
			side:         types.Order_SIDE_UNSPECIFIED,

			expectedOpenOrders: []types.Order{},
			expectedErr:        types.ErrInvalidOrderSide,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memclob := NewMemClobPriceTimePriority(false)

			// Populate memclob state.
			// 1. Create the map containing all open CLOB orders for the CLOB we're fetching.
			memclob.orderbooks[types.ClobPairId(tc.clobPairId)] = &Orderbook{
				SubaccountOpenClobOrders: make(
					map[satypes.SubaccountId]map[types.Order_Side]map[types.OrderId]bool,
				),
				orderIdToLevelOrder:       make(map[types.OrderId]*list.Node[types.ClobOrder]),
				blockExpirationsForOrders: make(map[uint32]map[types.OrderId]bool),
			}

			// 2. Add all orders to the `SubaccountOpenClobOrders` and `orderIdToLevelOrder` map.
			for _, order := range tc.memclobOrders {
				// 2.1. Create the orderbook, if it doesn't exist.
				if _, exists := memclob.orderbooks[order.GetClobPairId()]; !exists {
					memclob.orderbooks[order.GetClobPairId()] = &Orderbook{
						SubaccountOpenClobOrders: make(
							map[satypes.SubaccountId]map[types.Order_Side]map[types.OrderId]bool,
						),
						orderIdToLevelOrder:       make(map[types.OrderId]*list.Node[types.ClobOrder]),
						blockExpirationsForOrders: make(map[uint32]map[types.OrderId]bool),
					}
				}
				openClobOrders := memclob.orderbooks[order.GetClobPairId()].SubaccountOpenClobOrders

				// 2.2. Create the map containing all of a subaccount's open orders on this CLOB,
				// if it doesn't exist.
				subaccountId := order.OrderId.SubaccountId
				if _, exists := openClobOrders[subaccountId]; !exists {
					openClobOrders[subaccountId] = make(map[types.Order_Side]map[types.OrderId]bool)
				}
				userOpenClobOrders := openClobOrders[subaccountId]

				// 2.3. Create the map containing all of a subaccount's open orders on this CLOB
				// on this side, if it doesn't exist.
				if _, exists := userOpenClobOrders[order.Side]; !exists {
					userOpenClobOrders[order.Side] = make(map[types.OrderId]bool)
				}
				userOpenClobOrdersSide := userOpenClobOrders[order.Side]

				userOpenClobOrdersSide[order.OrderId] = true

				// 2.4. Add the order to the `orderIdToLevelOrder` map.
				memclob.orderbooks[order.GetClobPairId()].orderIdToLevelOrder[order.OrderId] = &types.LevelOrder{
					Value: types.ClobOrder{
						Order: order,
					},
				}
			}

			orders, err := memclob.GetSubaccountOrders(
				types.ClobPairId(tc.clobPairId),
				tc.subaccountId,
				tc.side,
			)

			if tc.expectedErr != nil {
				require.ErrorIs(t, tc.expectedErr, err)
			} else {
				require.Nil(t, err)
				expectedOrdersMap := make(map[types.OrderId]types.Order, len(tc.expectedOpenOrders))
				for _, order := range tc.expectedOpenOrders {
					expectedOrdersMap[order.OrderId] = order
				}

				for _, order := range orders {
					orderId := order.OrderId
					expectedOrder, exists := expectedOrdersMap[orderId]
					require.True(t, exists, "Order exists that was not in expected orders")
					require.Equal(t, expectedOrder, order)
					delete(expectedOrdersMap, orderId)
				}
				require.Empty(t, expectedOrdersMap, fmt.Sprintf("%d expected order(s) don't exist", len(expectedOrdersMap)))
			}
		})
	}
}

func TestGetSubaccountOrders_OrderNotFoundPanics(t *testing.T) {
	memclob := NewMemClobPriceTimePriority(false)
	memclob.orderbooks[0] = &Orderbook{
		SubaccountOpenClobOrders: map[satypes.SubaccountId]map[types.Order_Side]map[types.OrderId]bool{

			constants.Alice_Num0: {
				types.Order_SIDE_BUY: {
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId: true,
				},
			},
		},
	}

	require.Panics(t, func() {
		//nolint: errcheck
		memclob.GetSubaccountOrders(
			0,
			constants.Alice_Num0,
			types.Order_SIDE_BUY,
		)
	})
}
