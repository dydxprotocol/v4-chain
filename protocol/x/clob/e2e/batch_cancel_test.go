package clob_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func TestBatchCancel(t *testing.T) {
	tests := map[string]struct {
		firstBlockOrders       []clobtypes.MsgPlaceOrder
		firstBlockBatchCancel  *clobtypes.MsgBatchCancel
		secondBlockOrders      []clobtypes.MsgPlaceOrder
		secondBlockBatchCancel *clobtypes.MsgBatchCancel

		expectedOrderIdsInMemclob          map[clobtypes.OrderId]bool
		expectedCancelExpirationsInMemclob map[clobtypes.OrderId]uint32
		expectedOrderFillAmounts           map[clobtypes.OrderId]uint64
	}{
		"Cancel unfilled short term order": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
			},
			firstBlockBatchCancel: &clobtypes.MsgBatchCancel{
				SubaccountId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId.SubaccountId,
				ShortTermCancels: []clobtypes.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds:  []uint32{0},
					},
				},
				GoodTilBlock: 5,
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
			fmt.Println("submitting the batch cancel")
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *tc.firstBlockBatchCancel) {
				fmt.Println("request checkTx", checkTx)
				tApp.CheckTx(checkTx)
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
			if tc.secondBlockBatchCancel != nil {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *tc.secondBlockBatchCancel) {
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
