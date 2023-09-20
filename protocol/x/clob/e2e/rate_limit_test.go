package clob_test

import (
	"bytes"
	"testing"

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
		blockRateLimitConifg clobtypes.BlockRateLimitConfiguration
		firstMsg             sdktypes.Msg
		secondMsg            sdktypes.Msg
	}{
		"Short term orders": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			secondMsg: &PlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB20,
		},
		"Stateful orders": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
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
		"Short term order cancellations": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrderCancellationsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &CancelOrder_Alice_Num0_Id0_Clob1_GTB5,
			secondMsg: &CancelOrder_Alice_Num0_Id0_Clob0_GTB20,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.BlockRateLimitConfig = tc.blockRateLimitConifg
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			firstCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(tc.firstMsg),
				},
				tc.firstMsg,
			)
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// First transaction should be allowed.
			require.True(t, tApp.CheckTx(firstCheckTx).IsOK())

			secondCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(tc.secondMsg),
				},
				tc.secondMsg,
			)
			// Rate limit is 1 over two block, second attempt should be blocked.
			resp := tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "Rate of 2 exceeds configured block rate limit")

			// Rate limit of 1 over two blocks should still apply, total should be 3 now (2 in block 2, 1 in block 3).
			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "Rate of 3 exceeds configured block rate limit")

			// Rate limit of 1 over two blocks should still apply, total should be 2 now (1 in block 3, 1 in block 4).
			tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "Rate of 2 exceeds configured block rate limit")

			// Advancing two blocks should make the total count 0 now and the msg should be accepted.
			tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
		})
	}
}

func TestCancellationAndMatchInTheSameBlock_Regression(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().Build()

	LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
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

	PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price10_GTB20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     1,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell7_Price10_GTB20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
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
		ValidateDeliverTxs: func(_ sdktypes.Context, _ abcitypes.RequestDeliverTx, _ abcitypes.ResponseDeliverTx) bool {
			// Don't halt the chain since it's expected that the order will be removed after getting fully filled,
			// so the subsequent cancellation will be invalid.
			return false
		},
	})
}

func TestStatefulCancellation_Deduplication(t *testing.T) {
	LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
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
			tApp := testapp.NewTestAppBuilder().Build()
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
	LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
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
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
			}).WithTesting(t).Build()
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

func TestRateLimitingOrders_StatefulOrderRateLimitsAreAcrossMarkets(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
	}).WithTesting(t).Build()
	ctx := tApp.InitChain()

	firstMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	)

	// Second order should not be allowed in 2nd block and allowed in 4th block.
	secondMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5),
			AccSequenceNumberForSigning: 2,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
	)

	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	// First order should be allowed and second should be rejected.
	require.True(t, tApp.CheckTx(firstMarketCheckTx).IsOK())
	resp := tApp.CheckTx(secondMarketCheckTx)
	require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
	require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
	require.Contains(t, resp.Log, "Rate of 2 exceeds configured block rate limit")

	// Retrying in the 4th block should succeed since the rate limits should have been pruned.
	tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
	resp = tApp.CheckTx(secondMarketCheckTx)
	require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
}

func TestRateLimitingOrders_StatefulOrdersDuringDeliverTxAreRateLimited(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
	}).WithTesting(t).Build()
	ctx := tApp.InitChain()

	firstMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	)
	secondMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5),
			AccSequenceNumberForSigning: 2,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
	)

	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{firstMarketCheckTx.Tx, secondMarketCheckTx.Tx},
		ValidateDeliverTxs: func(
			context sdktypes.Context,
			request abcitypes.RequestDeliverTx,
			response abcitypes.ResponseDeliverTx,
		) (haltChain bool) {
			if bytes.Equal(request.Tx, firstMarketCheckTx.Tx) {
				require.Conditionf(t, response.IsOK, "Expected DeliverTx to succeed. Response: %+v", response)
			} else {
				require.Conditionf(t, response.IsErr, "Expected DeliverTx to error. Response: %+v", response)
				require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), response.Code)
				require.Contains(t, response.Log, "Rate of 2 exceeds configured block rate limit")
			}
			return false
		},
	})

	// Advance to block 3 which should cause the delivered stateful order to still be rejected from block 2.
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{secondMarketCheckTx.Tx},
		ValidateDeliverTxs: func(
			ctx sdktypes.Context,
			request abcitypes.RequestDeliverTx,
			response abcitypes.ResponseDeliverTx,
		) (haltchain bool) {
			require.Conditionf(t, response.IsErr, "Expected DeliverTx to error. Response: %+v", response)
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), response.Code)
			require.Contains(t, response.Log, "Rate of 3 exceeds configured block rate limit")
			return false
		},
	})

	// Advance to block 4 should clear out the delivered transactions in 2 and 3 allowing them to be
	// delivered in block 5.
	tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
	tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{secondMarketCheckTx.Tx},
	})
}
