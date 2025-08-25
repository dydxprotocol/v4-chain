package keeper_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	comettypes "github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
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
			StatefulOrderPlacement: constants.LongTermOrder_Alice_Num0_Id1_Clob0_Buy1BTC_Price50000_GTBT15,
			ExpectedError:          types.ErrStatefulOrderCollateralizationCheckFailed,
		},
		"Returns an error when equity tier check fails": {
			// Bob has TNC of $0.
			StatefulOrderPlacement: constants.LongTermOrder_Bob_Num0_Id2_Clob0_Buy15_Price5_GTBT10,
			ExpectedError:          types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
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
		"Returns an error when TWAP order has already been removed": {
			StatefulOrderPlacement: constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO,
			RemovedOrderIds: []types.OrderId{
				constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO.OrderId,
			},
			ExpectedError: types.ErrStatefulOrderPreviouslyRemoved,
		},
		"Returns an error when TWAP suborder size is too small": {
			StatefulOrderPlacement: constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy10_Price35_GTB20_RO,
			ExpectedError:          types.ErrInvalidPlaceOrder,
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

			mockLogger := &mocks.Logger{}
			mockLogger.On("With",
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			).Return(mockLogger)
			if errors.Is(tc.ExpectedError, types.ErrStatefulOrderCollateralizationCheckFailed) {
				mockLogger.On("Info",
					mock.Anything,
					mock.Anything,
					mock.AnythingOfType("*errors.wrappedError"),
				).Return()
			} else {
				mockLogger.On("Error",
					"Error placing order",
					mock.Anything,
					mock.Anything,
				).Return()
			}
			ks.Ctx = ks.Ctx.WithLogger(mockLogger)

			require.NoError(t, keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper))
			// Create test markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ks.Ctx, constants.PerpetualFeeParams))

			// Create Perpetual.
			perpetual := constants.BtcUsd_100PercentMarginRequirement
			_, err := ks.PerpetualsKeeper.CreatePerpetual(
				ks.Ctx,
				perpetual.Params.Id,
				perpetual.Params.Ticker,
				perpetual.Params.MarketId,
				perpetual.Params.AtomicResolution,
				perpetual.Params.DefaultFundingPpm,
				perpetual.Params.LiquidityTier,
				perpetual.Params.MarketType,
			)
			require.NoError(t, err)

			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, constants.Alice_Num0_10_000USD)

			err = ks.ClobKeeper.InitializeEquityTierLimit(
				ks.Ctx,
				types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(20_000_000),
							Limit:          5,
						},
					},
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(20_000_000),
							Limit:          5,
						},
					},
				},
			)
			require.NoError(t, err)

			// Create ClobPair.
			clobPair := constants.ClobPair_Btc
			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			indexerEventManager.On("AddTxnEvent",
				ks.Ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						clobtest.MustPerpetualId(clobPair),
						clobPair.Id,
						perpetual.Params.Ticker,
						perpetual.Params.MarketId,
						clobPair.Status,
						clobPair.QuantumConversionExponent,
						perpetual.Params.AtomicResolution,
						clobPair.SubticksPerTick,
						clobPair.StepBaseQuantums,
						perpetual.Params.LiquidityTier,
						perpetual.Params.MarketType,
						perpetual.Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ks.Ctx,
				clobPair.Id,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				clobPair.Status,
			)
			require.NoError(t, err)

			ctx := ks.Ctx.WithBlockHeight(6)
			ctx = ctx.WithBlockTime(time.Unix(int64(6), 0))
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    6,
				Timestamp: time.Unix(int64(6), 0),
			})

			for _, order := range tc.StatefulOrders {
				ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, 5)
				ks.ClobKeeper.AddStatefulOrderIdExpiration(
					ctx,
					order.MustGetUnixGoodTilBlockTime(),
					order.GetOrderId(),
				)
			}

			for _, orderId := range tc.PlacedStatefulCancellations {
				ks.ClobKeeper.AddDeliveredCancelledOrderId(
					ctx,
					orderId,
				)
			}

			processProposerMatchesEvents := types.ProcessProposerMatchesEvents{
				BlockHeight:             6,
				RemovedStatefulOrderIds: tc.RemovedOrderIds,
			}
			ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
				ctx,
				processProposerMatchesEvents,
			)

			// Run MsgHandler for placement.
			_, err = msgServer.PlaceOrder(ctx, &types.MsgPlaceOrder{Order: tc.StatefulOrderPlacement})
			require.ErrorIs(t, err, tc.ExpectedError)

			mockLogger.AssertExpectations(t)
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
		"Succeeds with twap order": {
			StatefulOrderPlacement: constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO,
			Subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Bob_Num0,
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

			require.NoError(t, keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper))

			ctx := ks.Ctx.WithBlockHeight(2)
			ctx = ctx.WithBlockTime(time.Unix(int64(2), 0))
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    2,
				Timestamp: time.Unix(int64(2), 0),
			})

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
				perpetual.Params.Id,
				perpetual.Params.Ticker,
				perpetual.Params.MarketId,
				perpetual.Params.AtomicResolution,
				perpetual.Params.DefaultFundingPpm,
				perpetual.Params.LiquidityTier,
				perpetual.Params.MarketType,
			)
			require.NoError(t, err)

			// Create ClobPair.
			clobPair := constants.ClobPair_Btc
			indexerEventManager.On("AddTxnEvent",
				ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						0,
						0,
						perpetual.Params.Ticker,
						perpetual.Params.MarketId,
						clobPair.Status,
						clobPair.QuantumConversionExponent,
						perpetual.Params.AtomicResolution,
						clobPair.SubticksPerTick,
						clobPair.StepBaseQuantums,
						perpetual.Params.LiquidityTier,
						perpetual.Params.MarketType,
						perpetual.Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ctx,
				clobPair.Id,
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
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewConditionalOrderPlacementEvent(
							tc.StatefulOrderPlacement,
						),
					),
				).Return().Once()
			} else if tc.StatefulOrderPlacement.IsTwapOrder() {
				indexerEventManager.On(
					"AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewTwapOrderPlacementEvent(
							tc.StatefulOrderPlacement,
						),
					),
				).Return().Once()
			} else {
				indexerEventManager.On(
					"AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
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
			if tc.StatefulOrderPlacement.IsTwapOrder() {
				twap_order, found := ks.ClobKeeper.GetTwapOrderPlacement(ctx, tc.StatefulOrderPlacement.GetOrderId())
				require.Equal(t, tc.StatefulOrderPlacement, twap_order.Order)
				require.True(t, found)

				suborder := tc.StatefulOrderPlacement
				suborder.OrderId.OrderFlags = types.OrderIdFlags_TwapSuborder

				twap_suborder, timestamp, found_suborder := ks.ClobKeeper.GetTwapTriggerPlacement(ctx, suborder.OrderId)
				require.Equal(t, suborder.OrderId, twap_suborder)
				require.True(t, found_suborder)
				require.Equal(t, timestamp, ctx.BlockTime().Unix())
			} else {
				_, found := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.StatefulOrderPlacement.GetOrderId())
				require.True(t, found)
			}

			// Ensure placement exists in memstore.
			var placements []types.OrderId
			if !tc.StatefulOrderPlacement.IsTwapOrder() {
				if tc.StatefulOrderPlacement.IsConditionalOrder() {
					placements = ks.ClobKeeper.GetDeliveredConditionalOrderIds(ctx)
				} else {
					placements = ks.ClobKeeper.GetDeliveredLongTermOrderIds(ctx)
				}
				require.Len(t, placements, 1)
				require.Equal(t, placements[0], tc.StatefulOrderPlacement.OrderId)
			}

			// Run mock assertions.
			indexerEventManager.AssertExpectations(t)
		})
	}
}

func TestHandleMsgPlaceOrder(t *testing.T) {
	testOrder := &types.Order{
		OrderId: types.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   types.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         types.Order_SIDE_BUY,
		Quantums:     100,
		Subticks:     10_000,
		GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	tests := map[string]struct {
		// Whether order is internal.
		isInternalOrder bool
		// Quantums of USDC that subaccount has.
		assetQuantums int64
		// Whether a cancellation exists for the order.
		cancellationExists bool
		// Whether a removal exists for the order.
		removalExists bool
		// Whether equity tier limit exists.
		equityTierLimitExists bool

		// Expected error.
		expectedError error
	}{
		"Success - Place an Internal Order, Validations are Skipped": {
			isInternalOrder:       true,
			assetQuantums:         -1_000_000_000,
			cancellationExists:    false,
			removalExists:         false,
			equityTierLimitExists: true,
		},
		"Error - Place an Internal Order, Order Already Cancelled": {
			isInternalOrder:       true,
			assetQuantums:         -1_000_000_000,
			cancellationExists:    true,
			removalExists:         false,
			equityTierLimitExists: true,
			expectedError:         types.ErrStatefulOrderPreviouslyCancelled,
		},
		"Error - Place an Internal Order, Order Already Removed": {
			isInternalOrder:       true,
			assetQuantums:         -1_000_000_000,
			cancellationExists:    false,
			removalExists:         true,
			equityTierLimitExists: true,
			expectedError:         types.ErrStatefulOrderPreviouslyRemoved,
		},
		"Success - Place an External Order, All Validations Pass": {
			isInternalOrder:       false,
			assetQuantums:         1_000_000_000,
			cancellationExists:    false,
			removalExists:         false,
			equityTierLimitExists: true,
		},
		"Error - Place an External Order, Order Already Cancelled": {
			isInternalOrder:       false,
			assetQuantums:         1_000_000_000,
			cancellationExists:    true,
			removalExists:         false,
			equityTierLimitExists: true,
			expectedError:         types.ErrStatefulOrderPreviouslyCancelled,
		},
		"Error - Place an External Order, Order Already Removed": {
			isInternalOrder:       false,
			assetQuantums:         1_000_000_000,
			cancellationExists:    false,
			removalExists:         true,
			equityTierLimitExists: true,
			expectedError:         types.ErrStatefulOrderPreviouslyRemoved,
		},
		"Error - Place an External Order, Equity Tier Limit Reached": {
			isInternalOrder:       false,
			assetQuantums:         1,
			cancellationExists:    false,
			removalExists:         false,
			equityTierLimitExists: true,
			expectedError:         types.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit,
		},
		"Error - Place an External Order, Collateralization Check Failed": {
			isInternalOrder:       false,
			assetQuantums:         -1_000_000_000,
			cancellationExists:    false,
			removalExists:         false,
			equityTierLimitExists: false,
			expectedError:         types.ErrStatefulOrderCollateralizationCheckFailed,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis comettypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &testOrder.OrderId.SubaccountId,
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										assettypes.AssetUsdc.Id,
										big.NewInt(tc.assetQuantums),
									),
								},
							},
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *types.GenesisState) {
						if !tc.equityTierLimitExists {
							genesisState.EquityTierLimitConfig = types.EquityTierLimitConfiguration{
								ShortTermOrderEquityTiers: genesisState.EquityTierLimitConfig.ShortTermOrderEquityTiers,
								StatefulOrderEquityTiers:  nil,
							}
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain().WithIsCheckTx(false)
			k := tApp.App.ClobKeeper

			// Add order to placed cancellations / removals if specified.
			ppme := k.GetProcessProposerMatchesEvents(ctx)
			if tc.cancellationExists {
				k.AddDeliveredCancelledOrderId(ctx, testOrder.OrderId)
			}
			if tc.removalExists {
				ppme.RemovedStatefulOrderIds = []types.OrderId{testOrder.OrderId}
			}
			k.MustSetProcessProposerMatchesEvents(ctx, ppme)

			// Place order.
			err := k.HandleMsgPlaceOrder(
				ctx,
				&types.MsgPlaceOrder{
					Order: *testOrder,
				},
				tc.isInternalOrder,
			)
			if tc.expectedError == nil {
				require.NoError(t, err)

				// Ensure order placement exists in state.
				placement, found := k.GetLongTermOrderPlacement(ctx, testOrder.OrderId)
				require.True(t, found)
				require.Equal(t, *testOrder, placement.Order)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())

				// Ensure order placement does not exist in state.
				_, found := k.GetLongTermOrderPlacement(ctx, testOrder.OrderId)
				require.False(t, found)
			}
		})
	}
}
