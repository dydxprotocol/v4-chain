package clob_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestTwapOrderPlacementAndTrigger(t *testing.T) {
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
			Subticks: 500_000_000,     // $50,000 per BTC limit
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
				GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300), // 5 minutes from now
			},
			TwapConfig: &clobtypes.TwapOrderConfig{
				Duration:                300, // 5 minutes
				Interval:                60,  // 1 minute
				SlippagePercent:         0,   // 0% slippage off oracle price
				GoodTillBlockTimeOffset: 3,   // 3 seconds offset
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
	require.Equal(t, uint32(4), twapOrderPlacement.RemainingLegs) // 300/60 = 5 legs
	require.Equal(t, uint64(100_000_000_000), twapOrderPlacement.RemainingQuantums)
	require.Equal(t, uint32(2), twapOrderPlacement.BlockHeight)

	// Verify the first suborder was placed in the memclob
	suborderId0 := clobtypes.OrderId{
		SubaccountId:   constants.Alice_Num0,
		ClientId:       0,
		OrderFlags:     clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:     0,
		SequenceNumber: 0,
	}
	suborder, found := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId0)
	require.True(t, found, "First suborder should exist in memclob")
	require.Equal(t, uint64(20_000_000_000), suborder.Quantums) // 100B/5 = 20B per leg
	require.Equal(t, uint64(200_000_000), suborder.Subticks)    // $20,000 per oracle price
	require.Equal(
		t,
		uint32(ctx.BlockTime().Unix()+3), // 3 seconds from now
		suborder.GoodTilOneof.(*clobtypes.Order_GoodTilBlockTime).GoodTilBlockTime,
	)
	require.Equal(t, clobtypes.Order_SIDE_BUY, suborder.Side) // Same side as parent order

	// Verify initial suborder in trigger store
	triggerPlacements, found := tApp.App.ClobKeeper.GetTwapTriggerPlacements(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.True(t, found, "TWAP trigger placement should exist")
	require.Equal(t, 1, len(triggerPlacements), "Should have one suborder in trigger store")

	initialSuborder := triggerPlacements[0]
	require.Equal(t, clobtypes.OrderIdFlags_TwapSuborder, initialSuborder.Order.OrderId.OrderFlags)
	require.Equal(t, uint32(1), initialSuborder.Order.OrderId.SequenceNumber)
	require.Equal(t, uint64(0), initialSuborder.Order.Quantums)
	require.Equal(t, uint64(ctx.BlockTime().Unix()+60), initialSuborder.TriggerBlockTime)

	// Advance block time by 60 seconds to trigger next suborder
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: ctx.BlockTime().Add(time.Second * 60),
	})

	// Verify TWAP order state is updated
	twapOrderPlacement, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.True(t, found, "TWAP order placement should still exist")
	require.Equal(t, uint32(3), twapOrderPlacement.RemainingLegs)                   // One leg executed
	require.Equal(t, uint64(100_000_000_000), twapOrderPlacement.RemainingQuantums) // 100B remaining - no quantums filled

	suborderId1 := clobtypes.OrderId{
		SubaccountId:   constants.Alice_Num0,
		ClientId:       0,
		OrderFlags:     clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:     0,
		SequenceNumber: 1,
	}

	_, found0 := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId0)
	require.False(t, found0, "First suborder should have been removed from memclob")

	suborder1, found1 := tApp.App.ClobKeeper.MemClob.GetOrder(suborderId1)
	require.True(t, found1, "Second suborder should exist in memclob")
	require.Equal(t, uint64(25_000_000_000), suborder1.Quantums) // 100B/4 = 25B per leg (catching up)
	require.Equal(t, uint64(200_000_000), suborder1.Subticks)    // $20,000 per oracle price (not moved)
	require.Equal(t, clobtypes.Order_SIDE_BUY, suborder1.Side)   // Same side as parent order

	// Verify new suborder in trigger store
	triggerPlacements, found = tApp.App.ClobKeeper.GetTwapTriggerPlacements(
		ctx,
		twapOrder.Order.OrderId,
	)
	require.True(t, found, "TWAP trigger placement should exist")
	require.Equal(t, 1, len(triggerPlacements), "Should only contain updated suborder in trigger store")

	newSuborder := triggerPlacements[0]
	require.Equal(t, clobtypes.OrderIdFlags_TwapSuborder, newSuborder.Order.OrderId.OrderFlags)
	require.Equal(t, uint32(2), newSuborder.Order.OrderId.SequenceNumber)
	require.Equal(t, uint64(0), newSuborder.Order.Quantums)
	require.Equal(t, uint64(ctx.BlockTime().Unix()+60), newSuborder.TriggerBlockTime)
}
