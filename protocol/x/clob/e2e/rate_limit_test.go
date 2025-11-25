package clob_test

import (
	"testing"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/types"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestRateLimitingOrders_RateLimitsAreEnforced(t *testing.T) {
	tests := map[string]struct {
		blockRateLimitConfig clobtypes.BlockRateLimitConfiguration
		firstMsg             sdktypes.Msg
		secondMsg            sdktypes.Msg
	}{
		"Short term orders with same subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			secondMsg: &PlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB20,
		},
		"Short term orders with different subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			secondMsg: &PlaceOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTB20,
		},
		"Stateful orders with same subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			secondMsg: &LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
		},
		"Stateful orders with different subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			secondMsg: &LongTermPlaceOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTBT5,
		},
		"Short term order cancellations with same subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &CancelOrder_Alice_Num0_Id0_Clob1_GTB5,
			secondMsg: &CancelOrder_Alice_Num0_Id0_Clob0_GTB20,
		},
		"Short term order cancellations with different subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &CancelOrder_Alice_Num0_Id0_Clob1_GTB5,
			secondMsg: &CancelOrder_Alice_Num1_Id0_Clob0_GTB20,
		},
		"Batch cancellations with same subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     2,
					},
				},
			},
			firstMsg:  &BatchCancel_Alice_Num0_Clob0_1_2_3_GTB5,
			secondMsg: &BatchCancel_Alice_Num0_Clob0_1_2_3_GTB20,
		},
		"Batch cancellations with different subaccounts": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     2,
					},
				},
			},
			firstMsg:  &BatchCancel_Alice_Num0_Clob0_1_2_3_GTB5,
			secondMsg: &BatchCancel_Alice_Num1_Clob0_1_2_3_GTB20,
		},
		"Leverage updates with same address": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxLeverageUpdatesPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &UpdateLeverage_Alice_Num0_PerpId0_Lev5,
			secondMsg: &UpdateLeverage_Alice_Num0_PerpId1_Lev10,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				// Disable non-determinism checks since we mutate keeper state directly.
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *clobtypes.GenesisState) {
							genesisState.BlockRateLimitConfig = tc.blockRateLimitConfig
						},
					)
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = []satypes.Subaccount{
								constants.Alice_Num0_10_000USD,
								constants.Alice_Num1_10_000USD,
							}
						})
					return genesis
				}).Build()
			ctx := tApp.InitChain()

			firstCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.firstMsg),
				},
				tc.firstMsg,
			)
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// First transaction should be allowed.
			resp := tApp.CheckTx(firstCheckTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

			secondCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.secondMsg),
				},
				tc.secondMsg,
			)
			// Rate limit is 1 over two block, second attempt should be blocked.
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "exceeds configured block rate limit")

			// Rate limit of 1 over two blocks should still apply, total should be 3 now (2 in block 2, 1 in block 3).
			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "exceeds configured block rate limit")

			// Rate limit of 1 over two blocks should still apply, total should be 2 now (1 in block 3, 1 in block 4).
			tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "exceeds configured block rate limit")

			// Advancing two blocks should make the total count 0 now and the msg should be accepted.
			tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
		})
	}
}

func TestCombinedPlaceCancelBatchCancel_RateLimitsAreEnforced(t *testing.T) {
	tests := map[string]struct {
		blockRateLimitConfig clobtypes.BlockRateLimitConfiguration
		firstBatch           []sdktypes.Msg
		secondBatch          []sdktypes.Msg
		thirdBatch           []sdktypes.Msg
		firstBatchSuccess    []bool
		secondBatchSuccess   []bool
		thirdBatchSuccess    []bool
		lastOrder            sdktypes.Msg
	}{
		"Combination Place, Cancel, BatchCancel orders": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     6, // TODO FIX THIS AFTER SETTLE ON A NUM
					},
				},
			},
			firstBatch: []sdktypes.Msg{
				&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20, // 1-weight success @ 1
				&PlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB20, // 1-weight success @ 2
				&CancelOrder_Alice_Num0_Id0_Clob0_GTB20,             // 1-weight success @ 3
			},
			firstBatchSuccess: []bool{
				true,
				true,
				true,
			},
			secondBatch: []sdktypes.Msg{
				&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB23, // 1-weight success @ 4
				&CancelOrder_Alice_Num1_Id0_Clob0_GTB20,             // 1-weight success @ 5
				&BatchCancel_Alice_Num0_Clob0_1_2_3_GTB20,           // 2-weight failure @ 7
				&CancelOrder_Alice_Num0_Id0_Clob0_GTB23,             // 1-weight failure @ 8
			},
			secondBatchSuccess: []bool{
				true,
				true,
				false,
				false,
			},
			// advance one block, subtract 3 for a count of 5
			thirdBatch: []sdktypes.Msg{
				&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB24, // 1-weight success @ 6
				&BatchCancel_Alice_Num0_Clob0_1_2_3_GTB20,           // 2-weight failure @ 8
				&CancelOrder_Alice_Num0_Id0_Clob0_GTB20,             // 1-weight failure @ 9
			},
			thirdBatchSuccess: []bool{
				true,
				false,
				false,
			},
			// advance one block, subtract 5 for a count of 4
			lastOrder: &BatchCancel_Alice_Num1_Clob0_1_2_3_GTB20, // 2-weight pass @ 6
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				// Disable non-determinism checks since we mutate keeper state directly.
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *clobtypes.GenesisState) {
							genesisState.BlockRateLimitConfig = tc.blockRateLimitConfig
						},
					)
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = []satypes.Subaccount{
								constants.Alice_Num0_10_000USD,
								constants.Alice_Num1_10_000USD,
							}
						})
					return genesis
				}).Build()
			ctx := tApp.InitChain()

			firstCheckTxArray := []abcitypes.RequestCheckTx{}
			for _, msg := range tc.firstBatch {
				checkTx := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), msg),
					},
					msg,
				)
				firstCheckTxArray = append(firstCheckTxArray, checkTx)
			}
			secondCheckTxArray := []abcitypes.RequestCheckTx{}
			for _, msg := range tc.secondBatch {
				checkTx := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), msg),
					},
					msg,
				)
				secondCheckTxArray = append(secondCheckTxArray, checkTx)
			}
			thirdCheckTxArray := []abcitypes.RequestCheckTx{}
			for _, msg := range tc.thirdBatch {
				checkTx := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), msg),
					},
					msg,
				)
				thirdCheckTxArray = append(thirdCheckTxArray, checkTx)
			}

			tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// First batch of transactions.
			for idx, checkTx := range firstCheckTxArray {
				resp := tApp.CheckTx(checkTx)
				shouldSucceed := tc.firstBatchSuccess[idx]
				if shouldSucceed {
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				} else {
					require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
					require.Contains(t, resp.Log, "exceeds configured block rate limit")
				}
			}
			// Advance one block
			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			// Second batch of transactions.
			for idx, checkTx := range secondCheckTxArray {
				resp := tApp.CheckTx(checkTx)
				shouldSucceed := tc.secondBatchSuccess[idx]
				if shouldSucceed {
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				} else {
					require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
					require.Contains(t, resp.Log, "exceeds configured block rate limit")
				}
			}
			// Advance one block
			tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			// Third batch of transactions.
			for idx, checkTx := range thirdCheckTxArray {
				resp := tApp.CheckTx(checkTx)
				shouldSucceed := tc.thirdBatchSuccess[idx]
				if shouldSucceed {
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				} else {
					require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
					require.Contains(t, resp.Log, "exceeds configured block rate limit")
				}
			}
			// Advance one block
			tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})
			lastCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.lastOrder),
				},
				tc.lastOrder,
			)
			resp := tApp.CheckTx(lastCheckTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
		})
	}
}

func TestCancellationAndMatchInTheSameBlock_Regression(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()

	LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0, ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_LongTerm,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
		},
		testapp.DefaultGenesis(),
	))
	LCancelOrder_Alice_Num0_Id0_Clob0_GTBT20 := *clobtypes.NewMsgCancelOrderStateful(
		LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.Order.OrderId,
		20,
	)

	PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price10_GTB20 := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     1,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell7_Price10_GTB20 := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     7,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))

	tApp.InitChain()
	ctx := tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
	) {
		resp := tApp.CheckTx(msg)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price10_GTB20,
	) {
		resp := tApp.CheckTx(msg)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		LCancelOrder_Alice_Num0_Id0_Clob0_GTBT20,
	) {
		resp := tApp.CheckTx(msg)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		PlaceOrder_Bob_Num0_Id0_Clob0_Sell7_Price10_GTB20,
	) {
		resp := tApp.CheckTx(msg)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	// We shouldn't be overfilling orders and the line below shouldn't panic.
	_ = tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{
		ValidateFinalizeBlock: func(
			_ sdktypes.Context,
			_ abcitypes.RequestFinalizeBlock,
			_ abcitypes.ResponseFinalizeBlock,
		) bool {
			// Don't halt the chain since it's expected that the order will be removed after getting fully filled,
			// so the subsequent cancellation will be invalid.
			return false
		},
	})
}

func TestStatefulCancellation_Deduplication(t *testing.T) {
	LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0, ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_LongTerm,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		},
		testapp.DefaultGenesis(),
	))

	tests := map[string]struct {
		advanceAfterPlaceOrder  bool
		advanceAfterCancelOrder bool
	}{
		"Cancels in same block as order placed": {},
		"Cancels in next block after order was placed": {
			advanceAfterPlaceOrder: true,
		},
		"Cancels in subsequent blocks after order was placed": {
			advanceAfterPlaceOrder:  true,
			advanceAfterCancelOrder: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx, tApp.App, LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}
			if tc.advanceAfterPlaceOrder {
				ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			}
			// First cancellation should pass since the order should be known.
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgCancelOrderStateful(
				LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.Order.OrderId,
				11,
			)) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			if tc.advanceAfterCancelOrder {
				// Don't deliver the transactions ensuring that it is re-added via Recheck
				ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
					DeliverTxsOverride: make([][]byte, 0),
				})
			}

			// Subsequent cancellations should fail
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgCancelOrderStateful(
				LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.Order.OrderId,
				12,
			)) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Contains(t, resp.Log, "An uncommitted stateful order cancellation with this OrderId already exists")
			}

			if tc.advanceAfterCancelOrder {
				// Don't deliver the transactions ensuring that it is re-added via Recheck
				ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
					DeliverTxsOverride: make([][]byte, 0),
				})
			}

			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgCancelOrderStateful(
				LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.Order.OrderId,
				13,
			)) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Contains(t, resp.Log, "An uncommitted stateful order cancellation with this OrderId already exists")
			}
		})
	}
}

func TestStatefulOrderPlacement_Deduplication(t *testing.T) {
	LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0, ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_LongTerm,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		},
		testapp.DefaultGenesis(),
	))

	tests := map[string]struct {
		advanceBlock bool
	}{
		"Duplicates in same block": {},
		"Duplicates in subsequent blocks": {
			advanceBlock: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				// On block advancement we will lose the mempool causing stateful orders in the mempool
				// to be dropped and thus they won't be rechecked.
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(!tc.advanceBlock).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						// Disable the default rate limit of 2 stateful orders per block so we can test with
						// more than 2 orders.
						func(genesisState *clobtypes.GenesisState) {
							genesisState.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{}
						},
					)
					return genesis
				}).Build()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// First placement should pass since the order is unknown.
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx, tApp.App, LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20) {
				require.True(t, tApp.CheckTx(checkTx).IsOK())
			}

			if tc.advanceBlock {
				// Don't deliver the transaction ensuring that it is re-added via Recheck
				ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
					DeliverTxsOverride: make([][]byte, 0),
				})
			}

			// Subsequent placements should fail
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx, tApp.App, LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Contains(t, resp.Log, "An uncommitted stateful order with this OrderId already exists")
			}

			if tc.advanceBlock {
				// Don't deliver the transaction ensuring that it is re-added via Recheck
				ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
					DeliverTxsOverride: make([][]byte, 0),
				})
			}

			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx, tApp.App, LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Contains(t, resp.Log, "An uncommitted stateful order with this OrderId already exists")
			}
		})
	}
}

func TestRateLimitingOrders_StatefulOrdersDuringDeliverTxAreNotRateLimited(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 2,
							Limit:     1,
						},
					},
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	firstMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	)
	secondMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning:        constants.Alice_Num0.Owner,
			AccSequenceNumberForSigning: 2,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
	)

	// We expect both to be accepted even though the rate limit is 1.
	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{firstMarketCheckTx.Tx, secondMarketCheckTx.Tx},
	})
}

func TestRateLimitingShortTermOrders_GuardedAgainstReplayAttacks(t *testing.T) {
	tests := map[string]struct {
		blockRateLimitConfig clobtypes.BlockRateLimitConfiguration
		replayLessGTB        sdktypes.Msg
		replayGreaterGTB     sdktypes.Msg
		firstValidGTB        sdktypes.Msg
		secondValidGTB       sdktypes.Msg
	}{
		"Short term order placements": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 1,
						Limit:     1,
					},
				},
			},
			replayLessGTB:    &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
			replayGreaterGTB: &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB47,
			firstValidGTB:    &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			secondValidGTB:   &PlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB20,
		},
		"Short term order cancellations": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 1,
						Limit:     1,
					},
				},
			},
			replayLessGTB:    &CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			replayGreaterGTB: &CancelOrder_Alice_Num0_Id0_Clob0_GTB47,
			firstValidGTB:    &CancelOrder_Alice_Num0_Id0_Clob0_GTB20,
			secondValidGTB:   &CancelOrder_Alice_Num0_Id1_Clob0_GTB20,
		},
		"Batch cancellations": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 1,
						Limit:     2,
					},
				},
			},
			replayLessGTB:    &BatchCancel_Alice_Num0_Clob0_1_2_3_GTB5,
			replayGreaterGTB: &BatchCancel_Alice_Num0_Clob0_1_2_3_GTB47,
			firstValidGTB:    &BatchCancel_Alice_Num0_Clob0_1_2_3_GTB20,
			secondValidGTB:   &BatchCancel_Alice_Num0_Clob1_1_2_3_GTB20,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.BlockRateLimitConfig = tc.blockRateLimitConfig
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							constants.Alice_Num0_10_000USD,
							constants.Alice_Num1_10_000USD,
						}
					})
				return genesis
			}).Build()
			ctx := tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})

			// First tx fails due to GTB being too low.
			replayLessGTBTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.replayLessGTB),
				},
				tc.replayLessGTB,
			)
			resp := tApp.CheckTx(replayLessGTBTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrHeightExceedsGoodTilBlock.ABCICode(), resp.Code)

			// Second tx fails due to GTB being too high.
			replayGreaterGTBTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.replayGreaterGTB),
				},
				tc.replayGreaterGTB,
			)
			resp = tApp.CheckTx(replayGreaterGTBTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrGoodTilBlockExceedsShortBlockWindow.ABCICode(), resp.Code)

			firstCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.firstValidGTB),
				},
				tc.firstValidGTB,
			)
			// First transaction should be allowed due to GTB being valid. The first two tx do not count towards
			// the rate limit.
			resp = tApp.CheckTx(firstCheckTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

			secondCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.secondValidGTB),
				},
				tc.secondValidGTB,
			)
			// Rate limit is 1, second attempt should be blocked.
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "exceeds configured block rate limit")
		})
	}
}

func TestRateLimitingOrders_StatefulOrdersNotCountedDuringRecheck(t *testing.T) {
	blockRateLimitConfig := clobtypes.BlockRateLimitConfiguration{
		MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 2,
				Limit:     2,
			},
		},
	}
	firstMsg := &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5
	secondMsg := &LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5

	tApp := testapp.NewTestAppBuilder(t).
		// Disable non-determinism checks since we mutate keeper state directly.
		WithNonDeterminismChecksEnabled(false).
		WithGenesisDocFn(func() (genesis types.GenesisDoc) {
			genesis = testapp.DefaultGenesis()
			testapp.UpdateGenesisDocWithAppStateForModule(
				&genesis,
				func(genesisState *clobtypes.GenesisState) {
					genesisState.BlockRateLimitConfig = blockRateLimitConfig
				},
			)
			testapp.UpdateGenesisDocWithAppStateForModule(
				&genesis,
				func(genesisState *satypes.GenesisState) {
					genesisState.Subaccounts = []satypes.Subaccount{
						constants.Alice_Num0_10_000USD,
						constants.Alice_Num1_10_000USD,
					}
				})
			return genesis
		}).Build()
	ctx := tApp.InitChain()

	firstCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), firstMsg),
		},
		firstMsg,
	)
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// First transaction should be allowed.
	resp := tApp.CheckTx(firstCheckTx)
	require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		// First transaction did not get proposed in a block.
		// Txn remains in the mempool and will get rechecked.
		DeliverTxsOverride: [][]byte{},
	})

	// Rate limit is 2 over two block, second attempt should be allowed.
	secondCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), secondMsg),
		},
		secondMsg,
	)
	resp = tApp.CheckTx(secondCheckTx)
	require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
}

func TestRateLimitingUpdateLeverage_RateLimitsAreEnforced(t *testing.T) {
	tests := map[string]struct {
		blockRateLimitConfig clobtypes.BlockRateLimitConfiguration
		firstMsg             sdktypes.Msg
		secondMsg            sdktypes.Msg
		expectSecondOK       bool
	}{
		"UpdateLeverage with same address is rate limited": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxLeverageUpdatesPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:       &UpdateLeverage_Alice_Num0_PerpId0_Lev5,
			secondMsg:      &UpdateLeverage_Alice_Num0_PerpId1_Lev10,
			expectSecondOK: false,
		},
		"UpdateLeverage with same address but looser limits is not rate limited": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxLeverageUpdatesPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 1,
						Limit:     1,
					},
				},
			},
			firstMsg:       &UpdateLeverage_Alice_Num0_PerpId0_Lev5,
			secondMsg:      &UpdateLeverage_Alice_Num0_PerpId1_Lev10,
			expectSecondOK: true,
		},
		"UpdateLeverage with different addresses is not rate limited": {
			blockRateLimitConfig: clobtypes.BlockRateLimitConfiguration{
				MaxLeverageUpdatesPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:       &UpdateLeverage_Alice_Num0_PerpId0_Lev5,
			secondMsg:      &UpdateLeverage_Bob_Num0_PerpId0_Lev5,
			expectSecondOK: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *clobtypes.GenesisState) {
							genesisState.BlockRateLimitConfig = tc.blockRateLimitConfig
						},
					)
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = []satypes.Subaccount{
								constants.Alice_Num0_10_000USD,
								constants.Alice_Num1_10_000USD,
								constants.Bob_Num0_10_000USD,
							}
						})
					return genesis
				}).Build()
			ctx := tApp.InitChain()

			firstCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.firstMsg),
				},
				tc.firstMsg,
			)

			// First transaction should be allowed.
			resp := tApp.CheckTx(firstCheckTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

			// Advance to next block to trigger rate limit pruning
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			secondCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), tc.secondMsg),
				},
				tc.secondMsg,
			)
			resp = tApp.CheckTx(secondCheckTx)
			if tc.expectSecondOK {
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			} else {
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
				require.Contains(t, resp.Log, "exceeds configured block rate limit")
			}
		})
	}
}
