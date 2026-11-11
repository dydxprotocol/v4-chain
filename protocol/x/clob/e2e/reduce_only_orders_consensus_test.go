package clob_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// TestReduceOnlyOrder_ProposerOverfillRejected verifies that a proposer-supplied match cannot overfill a
// signed reduce-only order past the victim's position size in DeliverTx. An honest memclob would
// resize this reduce-only fill down to the victim's position size before ever proposing it, so this
// raw operations queue (built directly, bypassing CheckTx/memclob) can only occur if a malicious
// proposer submits it verbatim. The match should be rejected with
// ErrReduceOnlyWouldChangePositionSide, the victim's position should be untouched, and the block
// should still commit.
func TestReduceOnlyOrder_ProposerOverfillRejected(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					constants.Carl_Num0_100000USD,
					// Alice is long 1 BTC and will place a signed reduce-only sell order.
					constants.Alice_Num1_1BTC_Long_500_000USD,
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *perptypes.GenesisState) {
				genesisState.Params = constants.PerpetualsGenesisParams
				genesisState.LiquidityTiers = constants.LiquidityTiers
				genesisState.Perpetuals = []perptypes.Perpetual{
					constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					constants.EthUsd_20PercentInitial_10PercentMaintenance,
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	carlBuyOrder := clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     150_000_000, // 1.5 BTC, large enough to cover the overfill below
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}

	// Alice is long 1 BTC (100_000_000 quantums). This signed reduce-only order is sized at 1.5 BTC,
	// larger than her position, exactly as the memclob would allow it onto the book pre-resize.
	aliceReduceOnlySellOrder := clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     150_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}

	blockAdvancement := testapp.BlockAdvancementWithErrors{
		BlockAdvancement: testapp.BlockAdvancement{
			ShortTermOrdersAndOperations: []interface{}{
				carlBuyOrder,
				aliceReduceOnlySellOrder,
				// The proposer supplies a raw match that fills Alice's reduce-only order for
				// 120_000_000 (1.2 BTC), which exceeds her 100_000_000 (1 BTC) position but does not
				// exceed either order's quantums (150_000_000), so it is only caught by the DeliverTx
				// reduce-only guard.
				clobtestutils.NewMatchOperationRaw(
					&aliceReduceOnlySellOrder,
					[]clobtypes.MakerFill{
						{
							FillAmount:   120_000_000,
							MakerOrderId: carlBuyOrder.OrderId,
						},
					},
				),
			},
		},
		ExpectedDeliverTxErrors: testapp.TxIndexesToErrors{
			0: clobtypes.ErrReduceOnlyWouldChangePositionSide.Error(),
		},
	}

	ctx = blockAdvancement.AdvanceToBlock(ctx, 2, tApp, t)

	// The block still commits despite the rejected ops tx.
	require.Equal(t, int64(2), ctx.BlockHeight())

	// Alice's position is unchanged: the overfilled match must not be applied.
	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num1)
	require.Equal(
		t,
		constants.Alice_Num1_1BTC_Long_500_000USD.PerpetualPositions,
		aliceSubaccount.PerpetualPositions,
	)
	require.Equal(
		t,
		constants.Alice_Num1_1BTC_Long_500_000USD.AssetPositions,
		aliceSubaccount.AssetPositions,
	)

	// Carl (the would-be maker) must not have picked up any position either.
	carlSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Carl_Num0)
	require.Equal(t, constants.Carl_Num0_100000USD.PerpetualPositions, carlSubaccount.PerpetualPositions)
	require.Equal(t, constants.Carl_Num0_100000USD.AssetPositions, carlSubaccount.AssetPositions)

	// Neither order should have any recorded fill amount.
	for _, orderId := range []clobtypes.OrderId{aliceReduceOnlySellOrder.OrderId, carlBuyOrder.OrderId} {
		exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
		require.False(t, exists)
		require.Equal(t, satypes.BaseQuantums(0), fillAmount)
	}
}
