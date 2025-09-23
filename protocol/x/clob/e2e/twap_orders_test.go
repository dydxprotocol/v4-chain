package clob_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/stretchr/testify/require"
)

func TestTwapOrderPlacementAndCatchup(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Create a TWAP order with:
	// - 300 second duration (5 minutes, minimum allowed)
	// - 60 second interval
	// This will create 5 suborders total
	twapOrder := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_Twap,
				ClobPairId:   0,
			},
			Side:     clobtypes.Order_SIDE_BUY,
			Quantums: 100_000_000_000, // 10 BTC
			Subticks: 0,               // market TWAP order
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
				GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300), // 5 minutes from now
			},
			TwapParameters: &clobtypes.TwapParameters{
				Duration:       300, // 5 minutes
				Interval:       60,  // 1 minute
				PriceTolerance: 0,   // 0% slippage off oracle price
			},
		},
	)

	// Place the TWAP order
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, twapOrder) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// Advance block
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Verify the parent TWAP order state
	twapOrderPlacement, found := tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.True(t, found, "TWAP order placement should exist in state")
	require.Equal(t, twapOrder.Order, twapOrderPlacement.Order)
	require.Equal(t, uint32(4), twapOrderPlacement.RemainingLegs) // 5 original legs - 1 triggered
	require.Equal(t, uint64(100_000_000_000), twapOrderPlacement.RemainingQuantums)

	// Verify the first suborder was placed in the memclob
	suborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}

	suborder, found := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.True(t, found, "First suborder should exist in memclob")
	require.Equal(t, uint64(20_000_000_000), suborder.Quantums) // 100B/5 = 20B per leg
	require.Equal(t, uint64(200_000_000), suborder.Subticks)    // $20,000 oracle price
	require.Equal(
		t,
		uint32(ctx.BlockTime().Unix()+3), // 3 seconds from now
		suborder.GoodTilOneof.(*clobtypes.Order_GoodTilBlockTime).GoodTilBlockTime,
	)
	require.Equal(t, clobtypes.Order_SIDE_BUY, suborder.Side) // Same side as parent order

	// Verify initial suborder in trigger store
	nextSuborder, triggerTime, found := tApp.App.ClobKeeper.GetTwapTriggerPlacement(
		ctx,
		suborderId,
	)
	require.True(t, found, "TWAP trigger placement should exist")
	require.Equal(t, suborderId, nextSuborder)
	require.Equal(
		t,
		ctx.BlockTime().Unix()+int64(twapOrder.Order.TwapParameters.Interval),
		triggerTime,
	)

	// Advance block time by 30 seconds
	// Next suborder should not be triggered
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 30),
	})

	// Verify TWAP order state
	twapOrderPlacement, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.True(t, found, "TWAP order placement should still exist")
	require.Equal(t, uint32(4), twapOrderPlacement.RemainingLegs)                   // One leg executed
	require.Equal(t, uint64(100_000_000_000), twapOrderPlacement.RemainingQuantums) // 100B remaining - no quantums filled

	suborder, found = tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.False(t, found, "Suborder should have been removed from memclob due to expiry")

	// Advance block time by 30 seconds to trigger next suborder
	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 30),
	})

	suborder1, found1 := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.True(t, found1, "Second suborder should exist in memclob")
	require.Equal(t, uint64(25_000_000_000), suborder1.Quantums) // 100B/4 = 25B per leg (catching up)
	require.Equal(t, uint64(200_000_000), suborder1.Subticks)    // $20,000 per oracle price (not moved)
	require.Equal(t, clobtypes.Order_SIDE_BUY, suborder1.Side)   // Same side as parent order

	// Verify new suborder in trigger store
	newSuborder, triggerTime, found := tApp.App.ClobKeeper.GetTwapTriggerPlacement(
		ctx,
		suborderId,
	)
	require.True(t, found, "TWAP trigger placement should exist")
	require.Equal(t, suborderId, newSuborder)
	require.Equal(
		t,
		ctx.BlockTime().Unix()+int64(twapOrder.Order.TwapParameters.Interval),
		triggerTime,
	)
}

func TestDuplicateTWAPOrderPlacement(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Place a TWAP order to buy 100B quantums over 4 legs
	twapOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 100_000_000_000, // 100B quantums
		Subticks: 200_000_000,     // $20,000 per oracle price
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       320,
			Interval:       80,
			PriceTolerance: 10_000,
		},
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 100),
		},
	}

	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(twapOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(twapOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.False(t, resp.IsOK(), "Expected CheckTx to fail. Response: %+v", resp)
		require.Contains(t, resp.GetLog(), "A stateful order with this OrderId already exists")
	}
}

func TestTWAPOrderWithMatchingOrders(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Place a TWAP order to buy 100B quantums over 4 legs
	twapOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 100_000_000_000, // 100B quantums
		Subticks: 200_000_000,     // $20,000 per oracle price
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       320,
			Interval:       80,
			PriceTolerance: 10_000,
		},
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 100),
		},
	}

	// Place a matching sell order at the same price
	matchingOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_SELL,
		Quantums: 50_000_000_000, // 50B quantums
		Subticks: 200_000_000,    // $20,000 per oracle price
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 3600),
		},
	}

	// Place market order first and then TWAP order
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(matchingOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(twapOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	suborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}

	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 40),
		RequestPrepareProposalTxsOverride: [][]byte{
			testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&clobtypes.Order{
							OrderId:  suborderId,
							Side:     clobtypes.Order_SIDE_BUY,
							Quantums: 25_000_000_000,
							Subticks: 200_000_000,
						},
						[]clobtypes.MakerFill{
							{
								FillAmount:   25_000_000_000,
								MakerOrderId: matchingOrder.OrderId,
							},
						},
					),
				},
			}),
		},
	})

	// Verify TWAP order state is updated
	twapOrderPlacement, found := tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.OrderId,
	)
	require.True(t, found, "TWAP order placement should exist")
	require.Equal(t, uint32(3), twapOrderPlacement.RemainingLegs)
	require.Equal(t, uint64(75_000_000_000), twapOrderPlacement.RemainingQuantums) // 100B - 25B filled

	_, triggerTime, found1 := tApp.App.ClobKeeper.GetTwapTriggerPlacement(
		ctx,
		suborderId,
	)
	require.True(t, found1, "Second suborder should exist in trigger store")

	require.Equal(t, int64(80), triggerTime)

	_, placed_suborder1_found := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.False(t, placed_suborder1_found, "Second suborder should not have been placed yet")

	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 40),
	})

	filled_amount := tApp.App.ClobKeeper.MemClob.GetOrderFilledAmount(ctx, suborderId)
	require.Equal(t, uint64(25_000_000_000), filled_amount.ToUint64())

	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 0),
		RequestPrepareProposalTxsOverride: [][]byte{
			testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&clobtypes.Order{
							OrderId:  suborderId,
							Side:     clobtypes.Order_SIDE_BUY,
							Quantums: 25_000_000_000,
							Subticks: 200_000_000,
						},
						[]clobtypes.MakerFill{
							{
								FillAmount:   25_000_000_000,
								MakerOrderId: matchingOrder.OrderId,
							},
						},
					),
				},
			}),
		},
	})

	// Verify TWAP order state is updated again
	twapOrderPlacement, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.OrderId,
	)
	require.True(t, found, "TWAP order placement should still exist")
	require.Equal(t, uint32(2), twapOrderPlacement.RemainingLegs)
	require.Equal(t, uint64(50_000_000_000), twapOrderPlacement.RemainingQuantums) // 75B - 25B filled

	_, triggerTime, found2 := tApp.App.ClobKeeper.GetTwapTriggerPlacement(
		ctx,
		suborderId,
	)
	require.True(t, found2, "Third suborder should exist in trigger store")
	require.Equal(t, ctx.BlockTime().Unix()+int64(twapOrder.TwapParameters.Interval), triggerTime)
}

func TestTWAPOrderWithMatchingOrdersWhereTWAPOrderIsMaker(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Place a TWAP order to buy 100B quantums over 4 legs
	twapOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 100_000_000_000, // 100B quantums
		Subticks: 200_000_000,     // $20,000 per oracle price
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       320,
			Interval:       80,
			PriceTolerance: 10_000,
		},
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 100),
		},
	}

	// Place a matching sell order at the same price
	matchingOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_SELL,
		Quantums: 50_000_000_000, // 50B quantums
		Subticks: 200_000_000,    // $20,000 per oracle price
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 3600),
		},
	}

	// Place the TWAP order so the first suborder is resting on the book
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(twapOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	suborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}

	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Place market order sell order as a taker order
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(matchingOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 40),
		RequestPrepareProposalTxsOverride: [][]byte{
			testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&clobtypes.Order{
							OrderId:  matchingOrder.OrderId,
							Side:     clobtypes.Order_SIDE_SELL,
							Quantums: 50_000_000_000,
							Subticks: 200_000_000,
						},
						[]clobtypes.MakerFill{
							{
								FillAmount:   25_000_000_000,
								MakerOrderId: suborderId,
							},
						},
					),
				},
			}),
		},
	})

	// Verify TWAP order state is updated
	twapOrderPlacement, found := tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.OrderId,
	)
	require.True(t, found, "TWAP order placement should exist")
	require.Equal(t, uint32(3), twapOrderPlacement.RemainingLegs)
	require.Equal(t, uint64(75_000_000_000), twapOrderPlacement.RemainingQuantums) // 100B - 25B filled

	_, triggerTime, found1 := tApp.App.ClobKeeper.GetTwapTriggerPlacement(
		ctx,
		suborderId,
	)
	require.True(t, found1, "Second suborder should exist in trigger store")

	require.Equal(t, int64(80), triggerTime)

	_, placed_suborder1_found := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.False(t, placed_suborder1_found, "Second suborder should not have been placed yet")

	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 40),
	})

	filled_amount := tApp.App.ClobKeeper.MemClob.GetOrderFilledAmount(ctx, suborderId)
	require.Equal(t, uint64(25_000_000_000), filled_amount.ToUint64())

	ctx = tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 0),
		RequestPrepareProposalTxsOverride: [][]byte{
			testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&clobtypes.Order{
							OrderId:  suborderId,
							Side:     clobtypes.Order_SIDE_BUY,
							Quantums: 25_000_000_000,
							Subticks: 200_000_000,
						},
						[]clobtypes.MakerFill{
							{
								FillAmount:   25_000_000_000,
								MakerOrderId: matchingOrder.OrderId,
							},
						},
					),
				},
			}),
		},
	})

	// Verify TWAP order state is updated again
	twapOrderPlacement, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.OrderId,
	)
	require.True(t, found, "TWAP order placement should still exist")
	require.Equal(t, uint32(2), twapOrderPlacement.RemainingLegs)
	require.Equal(t, uint64(50_000_000_000), twapOrderPlacement.RemainingQuantums) // 75B - 25B filled

	_, triggerTime, found2 := tApp.App.ClobKeeper.GetTwapTriggerPlacement(
		ctx,
		suborderId,
	)
	require.True(t, found2, "Third suborder should exist in trigger store")
	require.Equal(t, ctx.BlockTime().Unix()+int64(twapOrder.TwapParameters.Interval), triggerTime)
}

func TestTwapOrderStopsPlacingSubordersWhenCollateralIsDepleted(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Create a TWAP order: 100_000_000_000 quantums (10 BTC), 10 suborders of 1 BTC each
	twapOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 100_000_000_000, // 10 BTC
		Subticks: 0,               // market TWAP order
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300), // 5 minutes from now
		},
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       300, // 5 minutes
			Interval:       75,  // 75 seconds
			PriceTolerance: 0,   // 0% slippage off oracle price
		},
	}

	// Place the TWAP order
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(twapOrder)) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// --- First suborder trigger ---
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 75),
	})

	suborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}

	// Place a matching sell order for the first suborder
	matchingOrder1 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     1,
			OrderFlags:   0,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: uint32(4),
		},
		Side:     clobtypes.Order_SIDE_SELL,
		Quantums: 25_000_000_000,
		Subticks: 100_000_000,
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(matchingOrder1),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// Simulate fill of first suborder
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 0),
		RequestPrepareProposalTxsOverride: [][]byte{
			testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&clobtypes.Order{
							OrderId:  suborderId,
							Side:     clobtypes.Order_SIDE_BUY,
							Quantums: 25_000_000_000,
							Subticks: 100_000_000,
						},
						[]clobtypes.MakerFill{
							{
								FillAmount:   25_000_000_000,
								MakerOrderId: matchingOrder1.OrderId,
							},
						},
					),
				},
			}),
		},
	})

	// --- Second suborder trigger ---
	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 75),
	})

	// Place a matching sell order for the second suborder
	matchingOrder2 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     2,
			OrderFlags:   0,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: uint32(6),
		},
		Side:     clobtypes.Order_SIDE_SELL,
		Quantums: 25_000_000_000, // 1 BTC
		Subticks: 100_000_000,
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(matchingOrder2),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// Simulate fill of second suborder
	suborderId2 := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 0),
		RequestPrepareProposalTxsOverride: [][]byte{
			testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&clobtypes.Order{
							OrderId:  suborderId2,
							Side:     clobtypes.Order_SIDE_BUY,
							Quantums: 25_000_000_000,
							Subticks: 100_000_000,
						},
						[]clobtypes.MakerFill{
							{
								FillAmount:   25_000_000_000,
								MakerOrderId: matchingOrder2.OrderId,
							},
						},
					),
				},
			}),
		},
	})

	withdrawal := &sendingtypes.MsgWithdrawFromSubaccount{
		Sender:    constants.Alice_Num0,
		Recipient: constants.BobAccAddress.String(), // send to bob
		AssetId:   constants.Usdc.Id,
		Quantums:  99_999_995_000_000_000, // remaining balance + 95B to be below collat requirements for next suborder
	}

	CheckTx_MsgWithdrawFromSubaccount := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: withdrawal.Sender.Owner,
			Gas:                  200_000,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		withdrawal,
	)
	tApp.CheckTx(CheckTx_MsgWithdrawFromSubaccount)
	ctx = tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{})

	// --- Third suborder trigger (should fail to place due to insufficient collateral) ---
	ctx = tApp.AdvanceToBlock(7, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 75),
	})

	// There should be no new suborder placed for Alice
	suborderId3 := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	_, found := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId3)
	require.False(t, found, "No new suborder should be placed after running out of collateral")

	// The TWAP order should be removed from the store (since it failed to place a suborder)
	_, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.False(t, found, "TWAP order placement should be deleted after failed suborder due to insufficient collateral")
}

func TestTwapOrderCancellation(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Create a TWAP order with:
	// - 300 second duration (5 minutes)
	// - 60 second interval
	// This will create 5 suborders total
	twapOrder := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_Twap,
				ClobPairId:   0,
			},
			Side:     clobtypes.Order_SIDE_BUY,
			Quantums: 100_000_000_000, // 10 BTC
			Subticks: 0,               // market TWAP order
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
				GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300), // 5 minutes from now
			},
			TwapParameters: &clobtypes.TwapParameters{
				Duration:       300, // 5 minutes
				Interval:       60,  // 1 minute
				PriceTolerance: 0,   // 0% slippage off oracle price
			},
		},
	)

	// Place the TWAP order
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, twapOrder) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// Advance block time to trigger first suborder
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 60),
	})

	// Verify first suborder was created
	suborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	suborder, found := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.True(t, found, "First suborder should exist in memclob")
	require.Equal(t, uint64(20_000_000_000), suborder.Quantums) // 100B/5 = 20B per leg

	// Cancel the TWAP order
	cancelMsg := clobtypes.MsgCancelOrder{
		OrderId: twapOrder.Order.OrderId,
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{
			// 5 minutes and 30 seconds from now
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 330),
		},
	}

	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, cancelMsg) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// Advance block to deliver cancellation message
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

	// Verify TWAP order was removed from state
	_, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.False(t, found, "TWAP order placement should be removed")

	// Verify suborder was removed from memclob
	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.False(t, found, "Suborder should be removed from memclob")

	// Verify no more suborders will be triggered
	tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 60),
	})

	// Verify no new suborder was created
	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.False(t, found, "No new suborder should be created after cancellation")
}

func TestTwapOrderWithThreeSuborders(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Create a TWAP order with:
	// - 300 second duration
	// - 100 second interval
	// This will create exactly 3 suborders
	twapOrder := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_Twap,
				ClobPairId:   0,
			},
			Side:     clobtypes.Order_SIDE_BUY,
			Quantums: 90_000_000_000,
			Subticks: 0,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
				GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 90),
			},
			TwapParameters: &clobtypes.TwapParameters{
				Duration:       300,
				Interval:       100,
				PriceTolerance: 0,
			},
		},
	)

	// Place the TWAP order
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, twapOrder) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// Advance block
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Verify the parent TWAP order state
	twapOrderPlacement, found := tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.True(t, found, "TWAP order placement should exist in state")
	require.Equal(t, twapOrder.Order, twapOrderPlacement.Order)
	require.Equal(t, uint32(2), twapOrderPlacement.RemainingLegs) // 3 original legs - 1 triggered
	require.Equal(t, uint64(90_000_000_000), twapOrderPlacement.RemainingQuantums)

	// Verify first suborder was created
	suborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	suborder, found := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.True(t, found, "First suborder should exist in memclob")
	require.Equal(t, uint64(30_000_000_000), suborder.Quantums) // 90B/3 = 30B per leg

	// Advance block time by 100 seconds to trigger second suborder
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 100),
	})

	// Verify second suborder was created
	suborder, found = tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.True(t, found, "Second suborder should exist in memclob")
	require.Equal(t, uint64(45_000_000_000), suborder.Quantums) // 90B/2 = 45B per leg

	// Advance block time by 100 seconds to trigger third suborder
	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 100),
	})

	// Verify third suborder was created
	suborder, found = tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.True(t, found, "Third suborder should exist in memclob")
	require.Equal(t, uint64(90_000_000_000), suborder.Quantums) // 90B/1 = 90B for last leg

	// Advance block time by 100 seconds to complete TWAP order
	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 100),
	})

	// Verify TWAP order was removed from state
	_, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.False(t, found, "TWAP order placement should be removed")

	// Verify no more suborders will be triggered
	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(suborderId)
	require.False(t, found, "No new suborder should be created after completion")

	// Verify no entries in trigger store
	_, _, found = tApp.App.ClobKeeper.GetTwapTriggerPlacement(
		ctx,
		suborderId,
	)
	require.False(t, found, "No trigger placement should exist after completion")
}

func TestTwapOrderStatefulOrderCount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Create a TWAP order with:
	// - 300 second duration
	// - 100 second interval
	// This will create exactly 3 suborders
	twapOrder := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_Twap,
				ClobPairId:   0,
			},
			Side:     clobtypes.Order_SIDE_BUY,
			Quantums: 90_000_000_000,
			Subticks: 0,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
				GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300),
			},
			TwapParameters: &clobtypes.TwapParameters{
				Duration:       300,
				Interval:       100,
				PriceTolerance: 0,
			},
		},
	)

	// Place the TWAP order
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, twapOrder) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
	}

	// Advance block
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Verify stateful order count is 1 (only the parent TWAP order)
	statefulOrderCount := tApp.App.ClobKeeper.GetStatefulOrderCount(ctx, constants.Alice_Num0)
	require.Equal(t, uint32(1), statefulOrderCount, "Stateful order count should be 1")

	// Advance block time by 100 seconds to trigger second suborder
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 100),
	})

	// Verify stateful order count is still 1
	statefulOrderCount = tApp.App.ClobKeeper.GetStatefulOrderCount(ctx, constants.Alice_Num0)
	require.Equal(t, uint32(1), statefulOrderCount, "Stateful order count should remain 1")

	// Advance block time by 100 seconds to trigger third suborder
	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 100),
	})

	// Verify stateful order count is still 1
	statefulOrderCount = tApp.App.ClobKeeper.GetStatefulOrderCount(ctx, constants.Alice_Num0)
	require.Equal(t, uint32(1), statefulOrderCount, "Stateful order count should remain 1")

	// Advance block time by 100 seconds to complete TWAP order
	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 100),
	})

	// Verify stateful order count is 0 after completion
	statefulOrderCount = tApp.App.ClobKeeper.GetStatefulOrderCount(ctx, constants.Alice_Num0)
	require.Equal(t, uint32(0), statefulOrderCount, "Stateful order count should be 0 after completion")
}
