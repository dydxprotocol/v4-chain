package clob_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
