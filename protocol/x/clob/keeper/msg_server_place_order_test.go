package keeper_test

import (
	"context"
	"testing"
	"time"

	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	clobtest "github.com/dydxprotocol/v4/testutil/clob"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder_PanicIfShortTermOrder(t *testing.T) {
	msgPlaceOrder := types.MsgPlaceOrder{Order: constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16}
	require.Panicsf(
		t,
		func() {
			msgServer := keeper.NewMsgServerImpl(nil)
			//nolint: errcheck
			msgServer.PlaceOrder(context.Background(), &msgPlaceOrder)
		},
		"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
		msgPlaceOrder.Order.OrderId,
	)
}

func TestPlaceOrder_Error(t *testing.T) {
	tests := map[string]struct {
		StatefulOrders              []types.Order
		StatefulOrderPlacement      types.Order
		PlacedStatefulCancellations []types.OrderId
		ExpectedError               error
	}{
		"Returns an error when validation fails": {
			StatefulOrderPlacement: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			ExpectedError:          types.ErrTimeExceedsGoodTilBlockTime,
		},
		"Returns an error when collateralization check fails": {
			StatefulOrderPlacement: constants.LongTermOrder_Bob_Num0_Id2_Clob0_Buy15_Price5_GTBT10,
			ExpectedError:          types.ErrStatefulOrderCollateralizationCheckFailed,
		},
		"Returns an error when order replacement is attempted": {
			StatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
			StatefulOrderPlacement: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			ExpectedError:          types.ErrStatefulOrderAlreadyExists,
		},
		"Returns an error when order has already been cancelled": {
			StatefulOrderPlacement: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			PlacedStatefulCancellations: []types.OrderId{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
			ExpectedError: types.ErrStatefulOrderPreviouslyCancelled,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize mocks, context, msgServer.
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything, mock.Anything)
			indexerEventManager := &mocks.IndexerEventManager{}

			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)
			msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)

			// Create test markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Create Perpetual.
			perpetual := constants.BtcUsd_100PercentMarginRequirement
			_, err := ks.PerpetualsKeeper.CreatePerpetual(
				ks.Ctx,
				perpetual.Ticker,
				perpetual.MarketId,
				perpetual.AtomicResolution,
				perpetual.DefaultFundingPpm,
				perpetual.LiquidityTier,
			)
			require.NoError(t, err)

			// Create ClobPair.
			clobPair := constants.ClobPair_Btc
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ks.Ctx,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				clobPair.Status,
				clobPair.MakerFeePpm,
				clobPair.TakerFeePpm,
			)
			require.NoError(t, err)

			ctx := ks.Ctx.WithBlockHeight(6)
			ctx = ctx.WithBlockTime(time.Unix(int64(6), 0))
			ks.ClobKeeper.SetBlockTimeForLastCommittedBlock(ctx)

			for _, order := range tc.StatefulOrders {
				ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, 5)
				ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					order.MustGetUnixGoodTilBlockTime(),
					order.GetOrderId(),
				)
			}

			processProposerMatchesEvents := types.ProcessProposerMatchesEvents{
				BlockHeight:                 6,
				PlacedStatefulCancellations: tc.PlacedStatefulCancellations,
			}
			ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
				ctx,
				processProposerMatchesEvents,
			)

			// Run MsgHandler for placement.
			_, err = msgServer.PlaceOrder(ctx, &types.MsgPlaceOrder{Order: tc.StatefulOrderPlacement})
			require.ErrorIs(t, err, tc.ExpectedError)
		})
	}
}

func TestPlaceOrder_Success(t *testing.T) {
	tests := map[string]struct {
		StatefulOrderPlacement types.Order
		Subaccounts            []satypes.Subaccount
	}{
		"Succeeds": {
			StatefulOrderPlacement: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.MustGetOrder(),
			Subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
				},
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize mocks, context, msgServer.
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything, mock.Anything)
			indexerEventManager := &mocks.IndexerEventManager{}

			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)
			msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)

			ctx := ks.Ctx.WithBlockHeight(2)
			ctx = ctx.WithBlockTime(time.Unix(int64(2), 0))
			ks.ClobKeeper.SetBlockTimeForLastCommittedBlock(ctx)

			// Create test markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create all subaccounts.
			for _, subaccount := range tc.Subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			// Create Perpetual.
			perpetual := constants.BtcUsd_100PercentMarginRequirement
			_, err := ks.PerpetualsKeeper.CreatePerpetual(
				ctx,
				perpetual.Ticker,
				perpetual.MarketId,
				perpetual.AtomicResolution,
				perpetual.DefaultFundingPpm,
				perpetual.LiquidityTier,
			)
			require.NoError(t, err)

			// Create ClobPair.
			clobPair := constants.ClobPair_Btc
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ctx,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				clobPair.Status,
				clobPair.MakerFeePpm,
				clobPair.TakerFeePpm,
			)
			require.NoError(t, err)

			// Setup IndexerEventManager mock.
			indexerEventManager.On(
				"AddTxnEvent",
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexer_manager.GetB64EncodedEventMessage(
					indexerevents.NewStatefulOrderPlacementEvent(
						tc.StatefulOrderPlacement,
					),
				),
			).Return().Once()

			// Add BlockHeight to `ProcessProposerMatchesEvents`. This is normally done in `BeginBlock`.
			ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
				ctx,
				types.ProcessProposerMatchesEvents{
					BlockHeight: lib.MustConvertIntegerToUint32(2),
				},
			)

			// Run MsgHandler for placement.
			_, err = msgServer.PlaceOrder(ctx, &types.MsgPlaceOrder{Order: tc.StatefulOrderPlacement})
			require.NoError(t, err)

			// Ensure stateful order placement exists in state.
			_, found := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.StatefulOrderPlacement.GetOrderId())
			require.True(t, found)

			// Ensure placement exists in `ProcessProposerMatchesEvents`.
			events := ks.ClobKeeper.GetProcessProposerMatchesEvents(ctx)
			placements := events.GetPlacedStatefulOrders()
			require.Len(t, placements, 1)
			require.Equal(t, placements[0], tc.StatefulOrderPlacement)

			// Run mock assertions.
			indexerEventManager.AssertExpectations(t)
		})
	}
}
