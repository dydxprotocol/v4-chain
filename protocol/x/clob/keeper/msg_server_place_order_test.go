package keeper_test

import (
	"context"
	"testing"
	"time"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
		RemovedOrderIds             []types.OrderId
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
		"Returns an error when order has already been removed": {
			StatefulOrderPlacement: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			RemovedOrderIds: []types.OrderId{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
			ExpectedError: types.ErrStatefulOrderPreviouslyRemoved,
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

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ks.Ctx, constants.PerpetualFeeParams))

			// Create Perpetual.
			perpetual := constants.BtcUsd_100PercentMarginRequirement
			_, err := ks.PerpetualsKeeper.CreatePerpetual(
				ks.Ctx,
				perpetual.Params.Ticker,
				perpetual.Params.MarketId,
				perpetual.Params.AtomicResolution,
				perpetual.Params.DefaultFundingPpm,
				perpetual.Params.LiquidityTier,
			)
			require.NoError(t, err)

			// Create ClobPair.
			clobPair := constants.ClobPair_Btc
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ks.Ctx,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				clobPair.Status,
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
				BlockHeight:                        6,
				PlacedStatefulCancellationOrderIds: tc.PlacedStatefulCancellations,
				RemovedStatefulOrderIds:            tc.RemovedOrderIds,
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
		"Succeeds with long term order": {
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
		"Succeeds with conditional order": {
			StatefulOrderPlacement: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
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

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParams))

			// Create Perpetual.
			perpetual := constants.BtcUsd_100PercentMarginRequirement
			_, err := ks.PerpetualsKeeper.CreatePerpetual(
				ctx,
				perpetual.Params.Ticker,
				perpetual.Params.MarketId,
				perpetual.Params.AtomicResolution,
				perpetual.Params.DefaultFundingPpm,
				perpetual.Params.LiquidityTier,
			)
			require.NoError(t, err)

			// Create ClobPair.
			clobPair := constants.ClobPair_Btc
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ctx,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				clobPair.Status,
			)
			require.NoError(t, err)

			// Setup IndexerEventManager mock.
			if tc.StatefulOrderPlacement.IsConditionalOrder() {
				indexerEventManager.On(
					"AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewConditionalOrderPlacementEvent(
							tc.StatefulOrderPlacement,
						),
					),
				).Return().Once()
			} else {
				indexerEventManager.On(
					"AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewLongTermOrderPlacementEvent(
							tc.StatefulOrderPlacement,
						),
					),
				).Return().Once()
			}

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
			var placements []types.OrderId
			if tc.StatefulOrderPlacement.IsConditionalOrder() {
				placements = events.GetPlacedConditionalOrderIds()
			} else {
				placements = events.GetPlacedLongTermOrderIds()
			}
			require.Len(t, placements, 1)
			require.Equal(t, placements[0], tc.StatefulOrderPlacement.OrderId)

			// Run mock assertions.
			indexerEventManager.AssertExpectations(t)
		})
	}
}
