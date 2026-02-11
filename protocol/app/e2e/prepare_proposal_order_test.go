package e2e_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// TestPrepareProposal_CancelTxsBeforeOtherTxs verifies that PrepareProposal reorders
// transactions so that CLOB cancel-only txs appear before all other tx types in the
// "Other" group of the proposal.
//
// Note: only stateful (long-term/conditional) CLOB orders appear in the mempool and
// the "Other" group. Short-term CLOB orders are excluded from the mempool by CometBFT
// and are handled through the CLOB memclob mechanism only.
func TestPrepareProposal_CancelTxsBeforeOtherTxs(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Step 1: Place a long-term order via CheckTx so it exists in state.
	longTermPlaceOrder := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
		testapp.DefaultGenesis(),
	))
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, longTermPlaceOrder) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed for place order. Response: %+v", resp)
	}

	// Advance to block 2 to commit the long-term order to state.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Step 2: Submit the crossing sell and cancel via CheckTx (in "wrong" order —
	// sell first, cancel second) so they enter the mempool naturally.
	// PrepareProposal should reorder so the cancel comes before the sell.

	// Crossing sell order (signed by Bob — crosses Alice's buy at Price 10 on Clob0).
	crossingSellOrder := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell5_Price10_GTBT10,
		testapp.DefaultGenesis(),
	))
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, crossingSellOrder) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed for crossing sell. Response: %+v", resp)
	}

	// Cancel order tx (stateful, signed by Alice — same order we just placed).
	cancelMsg := constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, cancelMsg) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed for cancel. Response: %+v", resp)
	}

	// Step 3: Advance to block 3. The mempool has [sell, cancel].
	// PrepareProposal should reorder so the cancel-only tx comes first in the "Other" group.
	// Because the cancel is reordered before the sell, Alice's buy is canceled before the
	// crossing sell is applied, so no match occurs.
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		ValidateRespPrepare: func(ctx sdk.Context, resp abcitypes.ResponsePrepareProposal) (haltChain bool) {
			// Response layout: [ProposedOperationsTx, OtherTxs..., AcknowledgeBridgesTx, AddPremiumVotesTx, UpdateMarketPricesTx]
			// With our 2 override txs, we expect: 1 (ops) + 2 (other) + 3 (bridge, funding, prices) = 6.
			require.GreaterOrEqual(t, len(resp.Txs), 5,
				"proposal should have at least ops + 1 other + 3 tail txs, got %d", len(resp.Txs))
			otherTxs := resp.Txs[1 : len(resp.Txs)-3]
			require.NotEmpty(t, otherTxs, "Other group should contain our override txs")

			decoder := tApp.App.TxConfig().TxDecoder()
			seenNonCancel := false
			for i, txBytes := range otherTxs {
				tx, err := decoder(txBytes)
				require.NoError(t, err, "decode other tx at index %d", i)
				isCancel := prepare.IsTxCancelOnly(tx)
				if isCancel {
					require.False(t, seenNonCancel,
						"cancel-only tx at index %d must not appear after non-cancel tx; all cancels must come first", i)
				} else {
					seenNonCancel = true
				}
			}
			return false
		},
	})

	// Step 4: Verify post-block state.
	// Alice's buy should have been canceled — no longer in state.
	_, aliceOrderExists := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, longTermPlaceOrder.Order.OrderId)
	require.False(t, aliceOrderExists, "Alice's buy order should have been canceled")

	// Bob's crossing sell should rest on the book — it was placed but had nothing to match
	// (Alice's buy was already canceled by the time the sell was applied).
	_, bobOrderExists := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, crossingSellOrder.Order.OrderId)
	require.True(t, bobOrderExists, "Bob's sell order should rest on the book (no match after cancel)")
}
