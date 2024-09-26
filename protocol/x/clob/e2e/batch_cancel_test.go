package clob_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/msgsender"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/off_chain_updates"
	ocutypes "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/off_chain_updates/types"
	indexersharedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/shared/types"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
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

			rate := sdaiservertypes.TestSDAIEventRequests[0].ConversionRate

			msgUpdateSDAIConversionRate := ratelimittypes.MsgUpdateSDAIConversionRate{
				Sender:         constants.Alice_Num0.Owner,
				ConversionRate: rate,
			}

			for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: msgUpdateSDAIConversionRate.Sender,
					Gas:                  1200000,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				&msgUpdateSDAIConversionRate,
			) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

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
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

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

var (
	PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20,
			testapp.DefaultGenesis(),
		),
	)
	PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB30 = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB30,
			testapp.DefaultGenesis(),
		),
	)
	PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20 = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			testapp.DefaultGenesis(),
		),
	)
	PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB20,
			testapp.DefaultGenesis(),
		),
	)
	PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26 = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26,
			testapp.DefaultGenesis(),
		),
	)
	PlaceOrder_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20 = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20,
			testapp.DefaultGenesis(),
		),
	)
	PlaceOrder_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO,
			testapp.DefaultGenesis(),
		),
	)
	PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23 = *clobtypes.NewMsgPlaceOrder(
		testapp.MustScaleOrder(
			constants.Order_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23,
			testapp.DefaultGenesis(),
		),
	)
)

// Tests cancelling multiple orders.
func TestBatchCancelBatchFunctionality(t *testing.T) {
	tests := map[string]struct {
		firstBlockOrders       []clobtypes.MsgPlaceOrder
		firstBlockBatchCancel  []clobtypes.MsgBatchCancel
		secondBlockOrders      map[clobtypes.MsgPlaceOrder]bool
		secondBlockBatchCancel []clobtypes.MsgBatchCancel

		expectedOrderIdsInMemclob          map[clobtypes.OrderId]bool
		expectedCancelExpirationsInMemclob map[clobtypes.OrderId]uint32
	}{
		"Cancel a batch of orders, one not cancelled": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20,
				PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB20, // not cancelled
				PlaceOrder_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20,
				PlaceOrder_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23,
			},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1, 4, 5},
						},
					},
					GoodTilBlock: 25,
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId:     false,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      false,
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB20.Order.OrderId:      true,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: false,
				constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO.OrderId:   false,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23.Order.OrderId:      false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId:      25,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      25,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: 25,
				constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO.OrderId:   25,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23.Order.OrderId:      25,
			},
		},
		"Cancel a batch of orders, one cancel gtb is for an order with greater gtb": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20,
				PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26, // gtb 26 is > cancel gtb
				PlaceOrder_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20,
				PlaceOrder_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23,
			},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1, 4, 5},
						},
					},
					GoodTilBlock: 25,
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId:     false,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      false,
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26.Order.OrderId:      true,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: false,
				constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO.OrderId:   false,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23.Order.OrderId:      false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId: 25,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId: 25,
				// no cancel for client_id 5 since cancel gtb < existing order gtb
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26.Order.OrderId:      0,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: 25,
				constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO.OrderId:   25,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23.Order.OrderId:      25,
			},
		},
		"Cancel two batch of orders, overwriting cancel gtb with higher values": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20,
				PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26,
				PlaceOrder_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20,
				PlaceOrder_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23,
			},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1, 4, 5},
						},
					},
					GoodTilBlock: 25,
				},
			},
			secondBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1, 2, 4, 5},
						},
					},
					GoodTilBlock: 30,
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId:     false,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      false,
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26.Order.OrderId:      false,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: false,
				constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO.OrderId:   false,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23.Order.OrderId:      false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId:      30,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      30,
				PlaceOrder_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26.Order.OrderId:      30,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: 30,
				constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO.OrderId:   30,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23.Order.OrderId:      30,
			},
		},
		"Batch cancels prevent new orders with lower gtb from being placed": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20,
				PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				PlaceOrder_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20,
			},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1},
						},
					},
					GoodTilBlock: 25,
				},
			},
			secondBlockOrders: map[clobtypes.MsgPlaceOrder]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20:      false,
				PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20:      false,
				PlaceOrder_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20: false,
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId:     false,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      false,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId:      25,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      25,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: 25,
			},
		},
		"Batch cancel nonexistent orders": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1, 4, 5},
						},
					},
					GoodTilBlock: 25,
				},
			},
			secondBlockOrders: map[clobtypes.MsgPlaceOrder]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20:      false,
				PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20:      false,
				PlaceOrder_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20: false,
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId:      25,
				constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.OrderId:      25,
				constants.Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20.OrderId: 25,
				constants.Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO.OrderId:   25,
				PlaceOrder_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23.Order.OrderId:      25,
			},
		},
		"Batch cancel does not prevent orders with higher gtb from being placed": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1, 4, 5},
						},
					},
					GoodTilBlock: 25,
				},
			},
			secondBlockOrders: map[clobtypes.MsgPlaceOrder]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB30: true,
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB30.Order.OrderId: true,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB30.Order.OrderId: 25,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							constants.Alice_Num1_100_000USD,
						}
					},
				)
				return genesis
			}).WithCrashingAppCheckTxNonDeterminismChecksEnabled(false).Build()
			ctx := tApp.InitChain()

			rate := sdaiservertypes.TestSDAIEventRequests[0].ConversionRate

			msgUpdateSDAIConversionRate := ratelimittypes.MsgUpdateSDAIConversionRate{
				Sender:         constants.Alice_Num0.Owner,
				ConversionRate: rate,
			}

			for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: msgUpdateSDAIConversionRate.Sender,
					Gas:                  1200000,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				&msgUpdateSDAIConversionRate,
			) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Advance block to 10
			ctx = tApp.AdvanceToBlock(10, testapp.AdvanceToBlockOptions{})

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

			// Advance block to 15
			ctx = tApp.AdvanceToBlock(15, testapp.AdvanceToBlockOptions{})

			// Place second block orders and cancels
			for order, shouldSucceed := range tc.secondBlockOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					require.Equal(t, shouldSucceed, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
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
				if expectedCancelExpirationBlock > 0 {
					require.True(t, exists)
					require.Equal(t, expectedCancelExpirationBlock, cancelExpirationBlock)
				} else {
					require.False(t, exists)
				}
			}
		})
	}
}

// Tests emitting offchain updates.
func TestBatchCancelOffchainUpdates(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	batchCancel := clobtypes.MsgBatchCancel{
		SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
		ShortTermCancels: []clobtypes.OrderBatch{
			{
				ClobPairId: 0,
				ClientIds:  []uint32{0, 3},
			},
			{
				ClobPairId: 1,
				ClientIds:  []uint32{1, 2},
			},
		},
		GoodTilBlock: 25,
	}

	CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num1.Owner,
		},
		&PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num1.Owner,
		},
		&PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
	)
	CheckTx_BatchCancel := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num1.Owner,
		},
		&batchCancel,
	)

	tests := map[string]struct {
		firstBlockOrders      []clobtypes.MsgPlaceOrder
		firstBlockBatchCancel []clobtypes.MsgBatchCancel

		expectedOffchainUpdates []msgsender.Message
	}{
		"Cancel a batch of orders and check offchain updates": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20,
				PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			firstBlockBatchCancel: []clobtypes.MsgBatchCancel{
				{
					SubaccountId: PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId.SubaccountId,
					ShortTermCancels: []clobtypes.OrderBatch{
						{
							ClobPairId: 0,
							ClientIds:  []uint32{0, 3},
						},
						{
							ClobPairId: 1,
							ClientIds:  []uint32{1, 2},
						},
					},
					GoodTilBlock: 25,
				},
			},
			expectedOffchainUpdates: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.Tx),
				}),
				// 4 removals from the batch remove operation
				off_chain_updates.MustCreateOrderRemoveMessageWithReason(
					ctx,
					// Order id for the first cancel in the batch.
					clobtypes.OrderId{
						SubaccountId: constants.Alice_Num1,
						ClientId:     0,
						OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
						ClobPairId:   0,
					},
					indexersharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
					ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_BatchCancel.Tx),
				}),
				off_chain_updates.MustCreateOrderRemoveMessageWithReason(
					ctx,
					// Order id for the second cancel in the batch.
					clobtypes.OrderId{
						SubaccountId: constants.Alice_Num1,
						ClientId:     3,
						OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
						ClobPairId:   0,
					},
					indexersharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
					ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_BatchCancel.Tx),
				}),
				off_chain_updates.MustCreateOrderRemoveMessageWithReason(
					ctx,
					// Order id for the third cancel in the batch.
					clobtypes.OrderId{
						SubaccountId: constants.Alice_Num1,
						ClientId:     1,
						OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
						ClobPairId:   1,
					},
					indexersharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
					ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_BatchCancel.Tx),
				}),
				off_chain_updates.MustCreateOrderRemoveMessageWithReason(
					ctx,
					// Order id for the fourth cancel in the batch.
					clobtypes.OrderId{
						SubaccountId: constants.Alice_Num1,
						ClientId:     2,
						OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
						ClobPairId:   1,
					},
					indexersharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
					ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_BatchCancel.Tx),
				}),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							constants.Alice_Num1_100_000USD,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				return genesis
			}).WithCrashingAppCheckTxNonDeterminismChecksEnabled(false).WithAppOptions(appOpts).Build()

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequests[0].ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			// Advance block to 10
			ctx = tApp.AdvanceToBlock(10, testapp.AdvanceToBlockOptions{})

			// Clear any messages produced prior to these checkTx calls.
			msgSender.Clear()

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

			// Verify offchain messages
			require.ElementsMatch(
				t,
				tc.expectedOffchainUpdates[:4],
				msgSender.GetOffchainMessages()[:4],
			)
		})
	}
}
