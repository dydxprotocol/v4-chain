package clob_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobflags "github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// Buckets hold two orders and refill only from volume, so Alice's $10,000 earns nothing until she trades.
// Bucket state is node-local and does not survive a restart, so the crashing-app determinism check is off,
// as it is for the equity tier memo.
func TestPlaceShortTermOrder_TokenBucket(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).
		WithCrashingAppCheckTxNonDeterminismChecksEnabled(false).
		WithAppOptions(map[string]any{
			clobflags.ShortTermTokenBucketEnabled:            true,
			clobflags.ShortTermTokenBucketMinCapacity:        uint32(2),
			clobflags.ShortTermTokenBucketBurstBlocks:        uint32(0),
			clobflags.ShortTermTokenBucketVolumePerOrder:     uint64(1),
			clobflags.ShortTermTokenBucketCollateralPerOrder: uint64(0),
		}).
		WithGenesisDocFn(func() types.GenesisDoc {
			genesis := testapp.DefaultGenesis()
			testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
				state.Subaccounts = []satypes.Subaccount{
					constants.Alice_Num0_10_000USD,
					constants.Bob_Num0_10_000USD,
				}
			})
			testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
				state.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
				state.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{}
			})
			return genesis
		}).Build()
	ctx := tApp.InitChain()

	place := func(order clobtypes.Order) (resp abcitypes.ResponseCheckTx) {
		for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(order)) {
			resp = tApp.CheckTx(tx)
		}
		return resp
	}
	aliceBuy := func(clientId uint32) clobtypes.Order {
		order := testapp.MustScaleOrder(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20, testapp.DefaultGenesis())
		order.OrderId.ClientId = clientId
		order.GoodTilOneof = &clobtypes.Order_GoodTilBlock{GoodTilBlock: uint32(ctx.BlockHeight()) + 10}
		return order
	}
	requireOK := func(resp abcitypes.ResponseCheckTx) {
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	requireLimited := func(resp abcitypes.ResponseCheckTx) {
		require.Conditionf(t, resp.IsErr, "Expected CheckTx to be rate limited. Response: %+v", resp)
		require.Contains(t, resp.Log, clobtypes.ErrShortTermOrderRateExceeded.Error())
	}

	// Alice's first placement is a resting 0.1 BTC bid. The small test orders fill for under one quote quantum
	// of notional, which would leave her x/stats volume at zero.
	aliceBtcBid := constants.Order_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10_IOC
	aliceBtcBid.OrderId.ClientId = 100
	aliceBtcBid.Quantums = 10_000_000
	aliceBtcBid.TimeInForce = clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED
	requireOK(place(aliceBtcBid))
	requireOK(place(aliceBuy(1)))
	requireLimited(place(aliceBuy(2)))

	// Without volume nothing refills across blocks.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	requireLimited(place(aliceBuy(2)))

	// Bob sells into Alice's bid. The optimistic match in CheckTx earns her nothing yet.
	requireOK(place(constants.Order_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10))
	requireLimited(place(aliceBuy(2)))

	// The block delivers the fill to x/stats, but her bucket is only rated again 100 blocks after its first
	// rating. That rating sees the volume, and the bucket refills to capacity from the following block.
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
	requireLimited(place(aliceBuy(2)))
	ctx = tApp.AdvanceToBlock(102, testapp.AdvanceToBlockOptions{})
	requireLimited(place(aliceBuy(2)))
	ctx = tApp.AdvanceToBlock(103, testapp.AdvanceToBlockOptions{})
	requireOK(place(aliceBuy(3)))
	requireOK(place(aliceBuy(4)))
	requireLimited(place(aliceBuy(5)))

	// With the bucket empty, one reduce-only order a block still goes through. Alice is long from the fill, so
	// a reduce-only IOC sell priced above her own resting bids is valid and simply goes unfilled.
	reduceOnlySell := func(clientId uint32) clobtypes.Order {
		order := aliceBuy(clientId)
		order.Side = clobtypes.Order_SIDE_SELL
		order.Subticks *= 2
		order.TimeInForce = clobtypes.Order_TIME_IN_FORCE_IOC
		order.ReduceOnly = true
		return order
	}
	requireOK(place(reduceOnlySell(6)))
	requireLimited(place(reduceOnlySell(7)))
}
