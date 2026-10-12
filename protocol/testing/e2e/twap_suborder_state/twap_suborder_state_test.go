package twap_suborder_state_test

import (
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// TestExternalTwapSuborder_RejectedAtCheckTx verifies that a user-signed MsgPlaceOrder carrying the
// internal-only TWAP suborder order flag is rejected at CheckTx. Only the in-protocol scheduler
// produces suborders; an externally submitted one would otherwise rest as a parentless order and
// fail matching on a later block.
func TestExternalTwapSuborder_RejectedAtCheckTx(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	suborder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     777,
			OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 20_000_000_000,
		Subticks: 200_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300),
		},
	}

	checkTx := testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(suborder))[0]
	resp := tApp.CheckTx(checkTx)

	require.False(
		t,
		resp.IsOK(),
		"Expected CheckTx to reject an externally-submitted TWAP suborder. Response: %+v",
		resp,
	)
	require.Equal(t, clobtypes.ErrInvalidOrderFlag.ABCICode(), resp.Code)
}

// TestExternalTwapSuborder_RejectedAtProcessProposal verifies that even if a user-signed
// MsgPlaceOrder for a TWAP suborder is somehow included by a (malicious or buggy) proposer's block
// proposal, ProcessProposal rejects the whole block via `Msg.ValidateBasic`. This is a second line
// of defense behind the CheckTx rejection.
func TestExternalTwapSuborder_RejectedAtProcessProposal(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	tApp.InitChain()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	suborder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     778,
			OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 20_000_000_000,
		Subticks: 200_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300),
		},
	}

	// A TWAP suborder id is a stateful order id, so `MsgPlaceOrder` for it is not filtered out of
	// "OtherTxs" at the decode stage; it must be caught by `ValidateBasic` during `txs.Validate()`.
	suborderCheckTx := testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(suborder))[0]

	proposalTxs := [][]byte{
		constants.ValidEmptyMsgProposedOperationsTxBytes,
		suborderCheckTx.Tx,
		constants.MsgAcknowledgeBridges_NoEvents_TxBytes,
		constants.EmptyMsgAddPremiumVotesTxBytes,
		constants.EmptyMsgUpdateMarketPricesTxBytes,
	}

	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		RequestProcessProposalTxsOverride: proposalTxs,
		ValidateRespProcess: func(_ sdk.Context, resp abcitypes.ResponseProcessProposal) (stopProcessing bool) {
			require.Equal(t, abcitypes.ResponseProcessProposal_REJECT, resp.Status)
			return true // stop further processing of this (rejected) block in the test harness.
		},
	})
}

// TestTwapCancelReplaceSameClientId_NoPanic is the full regression scenario: a TWAP parent is
// placed with a resting maker so its first leg actually matches, the parent is then cancelled while
// its next leg's trigger is still pending, and a brand new TWAP order re-using the same ClientId
// (and therefore the same suborder id) is placed immediately after. Pre-fix, the cancelled parent's
// leftover trigger is inherited by the new TWAP and collides with its own trigger once both come
// due, corrupting TWAP state and eventually leaving a parentless suborder resting in state that
// panics `ProcessSingleMatch` on a later match. Post-fix, this should never panic and the chain
// should keep advancing normally.
func TestTwapCancelReplaceSameClientId_NoPanic(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Bob rests a large maker SELL so that TWAP suborders (BUY, market/oracle price) can actually
	// cross and match.
	matchingOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_SELL,
		Quantums: 100_000_000_000,
		Subticks: 200_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 3600),
		},
	}

	twapOrder := *clobtypes.NewMsgPlaceOrder(clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     7,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 100_000_000_000,
		Subticks: 0, // market TWAP order
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300),
		},
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       300,
			Interval:       60,
			PriceTolerance: 0,
		},
	})

	suborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     7,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}

	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(matchingOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected maker CheckTx to succeed. Response: %+v", resp)
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, twapOrder) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected TWAP CheckTx to succeed. Response: %+v", resp)
	}

	// Block 2: both orders are delivered. The first suborder leg (20B @ 200_000_000) triggers
	// immediately in this block's EndBlocker and rests on the book, crossing Bob's maker.
	require.NotPanics(t, func() {
		ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	})

	// Block 3: settle the match for the first leg.
	require.NotPanics(t, func() {
		ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
			BlockTime: ctx.BlockTime().Add(time.Second),
			RequestPrepareProposalTxsOverride: [][]byte{
				testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
					OperationsQueue: []clobtypes.OperationRaw{
						clobtestutils.NewMatchOperationRaw(
							&clobtypes.Order{
								OrderId:  suborderId,
								Side:     clobtypes.Order_SIDE_BUY,
								Quantums: 20_000_000_000,
								Subticks: 200_000_000,
							},
							[]clobtypes.MakerFill{
								{
									FillAmount:   20_000_000_000,
									MakerOrderId: matchingOrder.OrderId,
								},
							},
						),
					},
				}),
			},
		})
	})

	twapOrderPlacement, found := tApp.App.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.Order.OrderId)
	require.True(t, found, "TWAP order placement should exist")
	require.Equal(t, uint64(80_000_000_000), twapOrderPlacement.RemainingQuantums)

	// Cancel the parent while the second-leg trigger is still pending.
	cancelMsg := clobtypes.MsgCancelOrder{
		OrderId: twapOrder.Order.OrderId,
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 1000),
		},
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, cancelMsg) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected cancel CheckTx to succeed. Response: %+v", resp)
	}
	require.NotPanics(t, func() {
		ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
	})

	_, found = tApp.App.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.Order.OrderId)
	require.False(t, found, "Cancelled TWAP parent should be removed from state")

	// Re-place a brand new TWAP order re-using the same ClientId (and therefore the same suborder
	// id as the cancelled parent).
	newTwapOrder := *clobtypes.NewMsgPlaceOrder(clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     7,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 40_000_000_000,
		Subticks: 0,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300),
		},
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       300,
			Interval:       60,
			PriceTolerance: 0,
		},
	})
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, newTwapOrder) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected re-placed TWAP CheckTx to succeed. Response: %+v", resp)
	}

	// Advance through several more sweeps. Pre-fix, the cancelled parent's leftover trigger for
	// `suborderId` collides with the new TWAP's own trigger for the same id once both come due,
	// eventually leaving a parentless suborder resting in state that panics a later match.
	require.NotPanics(t, func() {
		ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})
	})
	require.NotPanics(t, func() {
		ctx = tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{
			BlockTime: ctx.BlockTime().Add(time.Second * 60),
		})
	})
	require.NotPanics(t, func() {
		ctx = tApp.AdvanceToBlock(7, testapp.AdvanceToBlockOptions{
			BlockTime: ctx.BlockTime().Add(time.Second * 60),
		})
	})

	// The new TWAP order should be tracking its own placement correctly (three legs triggered: at
	// blocks 5, 6, and 7), not corrupted by any inherited orphan trigger collision.
	newTwapOrderPlacement, found := tApp.App.ClobKeeper.GetTwapOrderPlacement(ctx, newTwapOrder.Order.OrderId)
	require.True(t, found, "New TWAP order placement should still exist")
	require.Equal(t, uint32(2), newTwapOrderPlacement.RemainingLegs)

	// The chain kept producing blocks throughout.
	require.Equal(t, int64(7), tApp.App.LastBlockHeight())
}

// TestOrphanTwapSuborderInState_SkippedAtHydration simulates the pre-fix bug state: a TWAP suborder
// committed to state with no corresponding parent TWAP placement (e.g. left behind by a mishandled
// trigger collision, or a parent cancelled/completed without clearing its suborder). Hydration
// (`InitStatefulOrders`, run on every node startup/restart) must skip such an orphan rather than
// re-placing it onto the memclob, where a later fill would panic trying to recover its (missing)
// parent.
func TestOrphanTwapSuborderInState_SkippedAtHydration(t *testing.T) {
	// Non-determinism checks are disabled because this test injects state directly into the
	// primary app's committed store (bypassing normal tx delivery) to simulate a pre-existing
	// orphan; the mirrored check apps used for those checks never observe that direct write.
	tApp := testapp.NewTestAppBuilder(t).WithNonDeterminismChecksEnabled(false).Build()
	ctx := tApp.InitChain()

	orphanSuborderId := clobtypes.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     555,
		OrderFlags:   clobtypes.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	orphanSuborder := clobtypes.Order{
		OrderId:  orphanSuborderId,
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 10_000_000_000,
		Subticks: 200_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 300),
		},
	}

	// Write the orphan directly into committed state, bypassing the normal (guarded) placement
	// path, to simulate the pre-fix bug state.
	header := tApp.GetHeader()
	uncachedCtx := tApp.App.NewUncachedContext(false, header)
	tApp.App.ClobKeeper.SetLongTermOrderPlacement(uncachedCtx, orphanSuborder, uint32(header.Height))

	parentTwapId := clobtypes.OrderId{
		SubaccountId: orphanSuborderId.SubaccountId,
		ClientId:     orphanSuborderId.ClientId,
		OrderFlags:   clobtypes.OrderIdFlags_Twap,
		ClobPairId:   orphanSuborderId.ClobPairId,
	}
	_, found := tApp.App.ClobKeeper.GetTwapOrderPlacement(uncachedCtx, parentTwapId)
	require.False(t, found, "Sanity check: parent TWAP placement should not exist")

	// Hydration (as performed on every node startup/restart, before any ABCI calls are received)
	// must skip the orphan rather than placing it onto the memclob.
	require.NotPanics(t, func() {
		tApp.App.ClobKeeper.InitStatefulOrders(uncachedCtx)
	})

	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(orphanSuborderId)
	require.False(t, found, "Orphan TWAP suborder should be skipped during hydration, not placed on the memclob")

	// The chain must also keep advancing normally afterwards.
	require.NotPanics(t, func() {
		tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	})

	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(orphanSuborderId)
	require.False(t, found, "Orphan TWAP suborder should remain absent from the memclob")
}
