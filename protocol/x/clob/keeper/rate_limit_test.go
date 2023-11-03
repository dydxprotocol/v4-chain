package keeper_test

import (
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRateLimitPlaceOrderIsNoopOutsideOfCheckTxAndReCheckTx(t *testing.T) {
	tApp := testApp.NewTestAppBuilder(t).Build()
	checkTxCtx := tApp.AdvanceToBlock(21, testApp.AdvanceToBlockOptions{})
	deliverTxCtx := checkTxCtx.WithIsCheckTx(false).WithIsReCheckTx(false)
	msg := clobtypes.NewMsgPlaceOrder(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price5_GTB20)

	// We expect an error and that the GTB is out of bounds.
	require.Error(
		t,
		tApp.App.ClobKeeper.RateLimitPlaceOrder(checkTxCtx, msg),
		"GoodTilBlock 20 is less than the current blockHeight 22",
	)

	// We don't expect any checks from occurring.
	require.Nil(t, tApp.App.ClobKeeper.RateLimitPlaceOrder(deliverTxCtx, msg))
}

func TestRateLimitCancelOrderIsNoopOutsideOfCheckTxAndReCheckTx(t *testing.T) {
	tApp := testApp.NewTestAppBuilder(t).Build()
	checkTxCtx := tApp.AdvanceToBlock(21, testApp.AdvanceToBlockOptions{})
	deliverTxCtx := checkTxCtx.WithIsCheckTx(false).WithIsReCheckTx(false)
	msg := clobtypes.NewMsgCancelOrderShortTerm(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price5_GTB20.OrderId, 20)

	// We expect an error and that the GTB is out of bounds.
	require.Error(
		t,
		tApp.App.ClobKeeper.RateLimitCancelOrder(checkTxCtx, msg),
		"GoodTilBlock 20 is less than the current blockHeight 22",
	)

	// We don't expect any checks from occurring.
	require.Nil(t, tApp.App.ClobKeeper.RateLimitCancelOrder(deliverTxCtx, msg))
}
