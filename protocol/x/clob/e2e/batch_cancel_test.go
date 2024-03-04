package clob_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// These tests are the same as the e2e tests for single order cancellations.
func TestBatchCancelSingleCancelFunctionality(t *testing.T) {
	tests := map[string]struct {
		firstBlockOrders       []clobtypes.MsgPlaceOrder
		firstBlockBatchCancel  []clobtypes.MsgBatchCancel
		secondBlockOrders      []clobtypes.MsgPlaceOrder
		secondBlockBatchCancel []clobtypes.MsgBatchCancel

		expectedOrderIdsInMemclob          map[clobtypes.OrderId]bool
		expectedCancelExpirationsInMemclob map[clobtypes.OrderId]uint32
		expectedOrderFillAmounts           map[clobtypes.OrderId]uint64
	}{
		"Cancel unfilled short term order": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
			},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 5,
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 0,
			},
		},
		"Batch cancel partially filled short term order in same block": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
				*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
					clobtypes.Order{
						OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
						Side:         clobtypes.Order_SIDE_SELL,
						Quantums:     4,
						Subticks:     10,
						GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
					},
					testapp.DefaultGenesis(),
				)),
			},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 5,
				},
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 40,
			},
		},
		"Cancel partially filled short term order in next block": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
				*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
					clobtypes.Order{
						OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
						Side:         clobtypes.Order_SIDE_SELL,
						Quantums:     4,
						Subticks:     10,
						GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
					},
					testapp.DefaultGenesis(),
				)),
			},
			secondBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 5,
				},
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 40,
			},
		},
		"Cancel succeeds for fully-filled order": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			},
			secondBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 5,
				},
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 50,
			},
		},
		"Cancel with GTB < existing order GTB does not remove order from memclob": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			},
			secondBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 5,
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId: true,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
		},
		"Cancel with GTB < existing cancel GTB is not placed on memclob": {
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 5,
				},
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 3,
				},
			},

			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
		},
		"Cancel with GTB > existing cancel GTB is placed on memclob": {
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 5,
				},
				{
					SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0},
						},
					},
					GoodTilBlock: 6,
				},
			},

			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 6,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			// Place first block orders and cancels
			for _, order := range tc.firstBlockOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			for _, batch := range tc.firstBlockBatchCancel {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, batch) {
					tApp.CheckTx(checkTx)
				}
			}

			// Advance block
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Place second block orders and cancels
			for _, order := range tc.secondBlockOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			for _, batch := range tc.secondBlockBatchCancel {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, batch) {
					tApp.CheckTx(checkTx)
				}
			}

			// Verify expectations
			for orderId, shouldHaveOrder := range tc.expectedOrderIdsInMemclob {
				_, exists := tApp.App.ClobKeeper.MemClob.GetOrder(ctx, orderId)
				require.Equal(t, shouldHaveOrder, exists)
			}
			for orderId, expectedCancelExpirationBlock := range tc.expectedCancelExpirationsInMemclob {
				cancelExpirationBlock, exists := tApp.App.ClobKeeper.MemClob.GetCancelOrder(ctx, orderId)
				require.True(t, exists)
				require.Equal(t, expectedCancelExpirationBlock, cancelExpirationBlock)
			}
			for orderId, expectedFillAmount := range tc.expectedOrderFillAmounts {
				_, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
			}
		})
	}
}
