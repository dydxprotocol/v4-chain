package clob_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

// TestClobTx_MixedMsgsRejectedInDeliverTx proves that a transaction mixing a long-term
// `MsgPlaceOrder` with a `MsgSend` is rejected in `DeliverTx` as well as in `CheckTx`:
// `ValidateMsgsInClobTx` runs ahead of the DeliverTx/simulate short-circuit in
// `ClobDecorator.AnteHandle`, so neither the order placement nor the bank send executes.
func TestClobTx_MixedMsgsRejectedInDeliverTx(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	placeOrderMsg := &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5
	sendMsg := &banktypes.MsgSend{
		FromAddress: constants.AliceAccAddress.String(),
		ToAddress:   constants.BobAccAddress.String(),
		Amount: sdk.Coins{
			sdk.NewCoin(constants.TestNativeTokenDenom, sdkmath.OneInt()),
		},
	}

	checkTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.AliceAccAddress.String(),
		},
		placeOrderMsg,
		sendMsg,
	)

	recipientBalanceBefore := tApp.App.BankKeeper.GetBalance(
		ctx,
		constants.BobAccAddress,
		constants.TestNativeTokenDenom,
	)

	// Control: CheckTx already rejects this mixed clob tx.
	checkTxResp := tApp.CheckTx(checkTx)
	require.Falsef(t, checkTxResp.IsOK(), "Expected CheckTx to reject mixed clob tx. Response: %+v", checkTxResp)
	require.Equal(t, sdkerrors.ErrInvalidRequest.ABCICode(), checkTxResp.Code)

	// Deliver the tx directly in a block, bypassing CheckTx/mempool filtering, to prove that
	// `DeliverTx` independently enforces the same msg-composition invariant.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{checkTx.Tx},
		ValidateFinalizeBlock: func(
			ctx sdk.Context,
			request abcitypes.RequestFinalizeBlock,
			response abcitypes.ResponseFinalizeBlock,
		) (stopProcessing bool) {
			require.Len(t, response.TxResults, 1)
			txResult := response.TxResults[0]
			require.Truef(t, txResult.IsErr(), "Expected DeliverTx to reject mixed clob tx. Result: %+v", txResult)
			require.Equal(t, sdkerrors.ErrInvalidRequest.ABCICode(), txResult.Code)
			return false
		},
	})

	// The block still commits despite the rejected tx.
	require.Equal(t, int64(2), ctx.BlockHeight())

	// The long-term order must not have been placed.
	_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, placeOrderMsg.Order.OrderId)
	require.False(t, found)

	// The bank send must not have moved any funds.
	recipientBalanceAfter := tApp.App.BankKeeper.GetBalance(
		ctx,
		constants.BobAccAddress,
		constants.TestNativeTokenDenom,
	)
	require.Equal(t, recipientBalanceBefore, recipientBalanceAfter)
}
