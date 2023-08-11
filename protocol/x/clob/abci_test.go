package clob_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4/daemons/liquidation/api"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	testapp "github.com/dydxprotocol/v4/testutil/app"
	clobtest "github.com/dydxprotocol/v4/testutil/clob"
	prices "github.com/dydxprotocol/v4/x/prices/types"

	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquidationtypes "github.com/dydxprotocol/v4/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4/mocks"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/clob"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	unixTimeFive    = time.Unix(5, 0)
	unixTimeTen     = time.Unix(10, 0)
	unixTimeFifteen = time.Unix(15, 0)
)

// assertFillAmountAndPruneState verifies that keeper state for
// `OrderAmountFilled`, and `BlockHeightToPotentiallyPrunableOrders` is correct
// based on the provided expectations.
func assertFillAmountAndPruneState(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	expectedFillAmounts map[types.OrderId]satypes.BaseQuantums,
	expectedPruneableBlockHeights map[uint32][]types.OrderId,
	expectedPrunedOrders map[types.OrderId]bool,
) {
	for orderId, expectedFillAmount := range expectedFillAmounts {
		exists, fillAmount, _ := k.GetOrderFillAmount(ctx, orderId)
		require.True(t, exists)
		require.Equal(t, expectedFillAmount, fillAmount)
	}

	for orderId := range expectedPrunedOrders {
		exists, _, _ := k.GetOrderFillAmount(ctx, orderId)
		require.False(t, exists)
	}

	for blockHeight, orderIds := range expectedPruneableBlockHeights {
		// Verify that expected `blockHeightToPotentiallyPrunableOrders` were deleted.
		blockHeightToPotentiallyPrunableOrdersStore := prefix.NewStore(
			ctx.KVStore(storeKey),
			types.KeyPrefix(types.BlockHeightToPotentiallyPrunableOrdersPrefix),
		)

		potentiallyPrunableOrdersBytes := blockHeightToPotentiallyPrunableOrdersStore.Get(
			types.BlockHeightToPotentiallyPrunableOrdersKey(blockHeight),
		)

		var potentiallyPrunableOrders = &types.PotentiallyPrunableOrders{}
		err := potentiallyPrunableOrders.Unmarshal(potentiallyPrunableOrdersBytes)
		require.NoError(t, err)

		require.ElementsMatch(
			t,
			potentiallyPrunableOrders.OrderIds,
			orderIds,
		)
	}
}

func TestEndBlocker_Success(t *testing.T) {
	prunedOrderIdOne := types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0}
	prunedOrderIdTwo := types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1}
	orderIdThree := types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 2}
	blockHeight := uint32(5)

	seenPlaceOrderOne := types.MsgPlaceOrder{
		Order: types.Order{
			OrderId:      prunedOrderIdOne,
			GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight},
		},
	}

	tests := map[string]struct {
		// Setup.
		setupState func(ctx sdk.Context, k *keeper.Keeper)

		// Expectations.
		expectedFillAmounts                  map[types.OrderId]satypes.BaseQuantums
		expectedOffchainUpdateOrders         map[types.OrderId]bool
		expectedPruneableBlockHeights        map[uint32][]types.OrderId
		expectedPrunedOrders                 map[types.OrderId]bool
		expectedPrunedSeenPlaceOrders        map[types.MsgPlaceOrder]bool
		expectedStatefulPlacementInState     map[types.OrderId]bool
		expectedStatefulOrderTimeSlice       map[time.Time][]types.OrderId
		expectedProcessProposerMatchesEvents types.ProcessProposerMatchesEvents
	}{
		"Prunes existing Short-Term orders and seen place orders correctly": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Set `prunedOrderIdOne` and `prunedOrderIdTwo` as existing orders which already have fill amounts.
				k.SetOrderFillAmount(
					ctx,
					prunedOrderIdOne,
					100,
					blockHeight,
				)

				// Set `prunedOrderIdTwo` to be prunable at the next block height (this takes precedent of the blockHeight
				// set in `AddOrdersForPruning`).
				k.SetOrderFillAmount(
					ctx,
					prunedOrderIdTwo,
					100,
					blockHeight+1,
				)

				// This order should not be pruned.
				k.SetOrderFillAmount(
					ctx,
					orderIdThree,
					150,
					blockHeight+10,
				)

				// Set both of these orders as prunable at the current `blockHeight` so we can assert that they were pruned
				// correctly.
				k.AddOrdersForPruning(
					ctx,
					[]types.OrderId{prunedOrderIdOne, prunedOrderIdTwo},
					blockHeight,
				)

				// Set order as seen so we can assert that it was pruned correctly.
				k.AddSeenPlaceOrder(
					ctx,
					seenPlaceOrderOne,
				)

				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						BlockHeight:                blockHeight,
						OrdersIdsFilledInLastBlock: []types.OrderId{prunedOrderIdTwo, orderIdThree},
					},
				)
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				prunedOrderIdTwo: 100,
				orderIdThree:     150,
			},
			expectedOffchainUpdateOrders: map[types.OrderId]bool{
				prunedOrderIdTwo: true,
				orderIdThree:     true,
			},
			expectedPrunedOrders: map[types.OrderId]bool{
				prunedOrderIdOne: true,
			},
			expectedPrunedSeenPlaceOrders: map[types.MsgPlaceOrder]bool{
				seenPlaceOrderOne: false,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight:                blockHeight,
				OrdersIdsFilledInLastBlock: []types.OrderId{prunedOrderIdTwo, orderIdThree},
			},
		},
		"Removes expired stateful orders and updates process proposer matches events": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// These orders should get removed.
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					5,
					blockHeight,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					blockHeight,
				)
				k.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
				)
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
					5,
					blockHeight,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					blockHeight,
				)
				k.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
				)

				// This order should not be pruned.
				k.SetOrderFillAmount(
					ctx,
					orderIdThree,
					150,
					blockHeight+10,
				)

				// This order should not get removed.
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					5,
					blockHeight,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					blockHeight,
				)
				k.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				)

				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						BlockHeight:                blockHeight,
						OrdersIdsFilledInLastBlock: []types.OrderId{orderIdThree},
					},
				)
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				orderIdThree: 150,
			},
			expectedOffchainUpdateOrders: map[types.OrderId]bool{
				orderIdThree: true,
			},
			expectedStatefulPlacementInState: map[types.OrderId]bool{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId:  false,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId:  false,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId: true,
			},
			expectedStatefulOrderTimeSlice: map[time.Time][]types.OrderId{
				unixTimeTen: {},
				unixTimeFifteen: {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
				},
				BlockHeight:                blockHeight,
				OrdersIdsFilledInLastBlock: []types.OrderId{orderIdThree},
			},
		},
		"Stateful order placements are not overwritten": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// This order should not be pruned.
				k.SetOrderFillAmount(
					ctx,
					orderIdThree,
					150,
					blockHeight+10,
				)

				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders: []types.Order{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						},
						OrdersIdsFilledInLastBlock: []types.OrderId{orderIdThree},
						ExpiredStatefulOrderIds:    []types.OrderId{},
						BlockHeight:                blockHeight,
					},
				)
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				orderIdThree: 150,
			},
			expectedOffchainUpdateOrders: map[types.OrderId]bool{
				orderIdThree: true,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				},
				OrdersIdsFilledInLastBlock: []types.OrderId{orderIdThree},
				BlockHeight:                blockHeight,
			},
		},
		"Does not send order update message offchain message for a stateful order fill that got cancelled": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
					20,
					blockHeight+10,
				)
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
					20,
					blockHeight+10,
				)
				k.MustSetProcessProposerMatchesEvents(ctx, types.ProcessProposerMatchesEvents{
					OperationsProposedInLastBlock: []types.Operation{
						types.NewOrderPlacementOperation(
							constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,
						),
						types.NewOrderPlacementOperation(
							constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
						),
						types.NewMatchOperation(
							&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
							[]types.MakerFill{
								{
									FillAmount:   20,
									MakerOrderId: constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
								},
							},
						),
						// An order cancellation for Bob's order.
						types.NewOrderCancellationOperation(
							types.NewMsgCancelOrderStateful(
								constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
								15,
							),
						),
					},
					OrdersIdsFilledInLastBlock: []types.OrderId{
						constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
						constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
					},
					BlockHeight: blockHeight,
				})
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OperationsProposedInLastBlock: []types.Operation{
					types.NewOrderPlacementOperation(
						constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,
					),
					types.NewOrderPlacementOperation(
						constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					),
					types.NewMatchOperation(
						&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
						[]types.MakerFill{
							{
								FillAmount:   20,
								MakerOrderId: constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
							},
						},
					),
					// An order cancellation for Bob's order.
					types.NewOrderCancellationOperation(
						types.NewMsgCancelOrderStateful(
							constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
							15,
						),
					),
				},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
					constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(): 20,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId():    20,
			},
			expectedOffchainUpdateOrders: map[types.OrderId]bool{
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(): true,
				// For Bob's order, there should be 20 fill amount but no indexer OrderUpdate SendOffchainData call.
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(): false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()

			mockIndexerEventManager := &mocks.IndexerEventManager{}

			ctx,
				keeper,
				_,
				_,
				_,
				_,
				storeKey,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
			ctx = ctx.WithBlockHeight(int64(blockHeight)).WithBlockTime(unixTimeTen)

			if tc.setupState != nil {
				tc.setupState(ctx, keeper)
			}

			// Add expectations for Indexer off-chain order updates.
			offchainUpdates := types.NewOffchainUpdates()
			for orderId, fillAmount := range tc.expectedFillAmounts {
				message, success := off_chain_updates.CreateOrderUpdateMessage(
					ctx.Logger(),
					orderId,
					fillAmount,
				)
				require.Equal(t, true, success)
				offchainUpdates.AddUpdateMessage(orderId, message)
				// Test must indicate if an offchain update must be sent.
				require.Contains(t, tc.expectedOffchainUpdateOrders, orderId)
				if tc.expectedOffchainUpdateOrders[orderId] {
					mockIndexerEventManager.On("SendOffchainData", message).Return()
				}
			}

			// Assert that the indexer events for Expired Stateful Orders were emitted.
			for _, orderId := range tc.expectedProcessProposerMatchesEvents.ExpiredStatefulOrderIds {
				mockIndexerEventManager.On("AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewStatefulOrderExpirationEvent(
							orderId,
						),
					),
				).Once().Return()
			}

			clob.EndBlocker(
				ctx,
				*keeper,
			)

			assertFillAmountAndPruneState(
				t,
				keeper,
				ctx,
				storeKey,
				tc.expectedFillAmounts,
				tc.expectedPruneableBlockHeights,
				tc.expectedPrunedOrders,
			)

			for placeOrder, expectedSeen := range tc.expectedPrunedSeenPlaceOrders {
				require.Equal(
					t,
					expectedSeen,
					keeper.HasSeenPlaceOrder(ctx, placeOrder),
				)
			}

			require.True(t, memClob.AssertExpectations(t))

			require.True(
				t,
				unixTimeTen.Equal(
					keeper.MustGetBlockTimeForLastCommittedBlock(ctx),
				),
			)

			for orderId, exists := range tc.expectedStatefulPlacementInState {
				_, found := keeper.GetStatefulOrderPlacement(ctx, orderId)
				require.Equal(t, exists, found)
			}

			for time, expected := range tc.expectedStatefulOrderTimeSlice {
				actual := keeper.GetStatefulOrdersTimeSlice(ctx, time)
				require.Equal(t, expected, actual)
			}

			require.Equal(
				t,
				tc.expectedProcessProposerMatchesEvents,
				keeper.GetProcessProposerMatchesEvents(ctx),
			)

			// Assert that the necessary off-chain indexer events have been added.
			mockIndexerEventManager.AssertExpectations(t)
		})
	}
}

func TestLiquidateSubaccounts(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs          []types.ClobPair
		existingOrders []types.Order

		// Parameters.
		liquidatableSubaccounts []satypes.SubaccountId

		// Expectations.
		expectedPlacedOrders  []*types.MsgPlaceOrder
		expectedMatchedOrders []*types.ClobMatch
	}{
		"Liquidates liquidatable subaccount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				constants.Order_Carl_Num0_Id4_Clob0_Buy05BTC_Price40000,
			},

			liquidatableSubaccounts: []satypes.SubaccountId{
				constants.Dave_Num0,
			},

			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				},
				{
					Order: constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
								FillAmount:   50_000_000,
							},
							{
								MakerOrderId: constants.OrderId_Alice_Num0_ClientId1_Clob0,
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
		},
		"Does not liquidate a non-liquidatable subaccount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				constants.Order_Carl_Num0_Id4_Clob0_Buy05BTC_Price40000,
			},

			liquidatableSubaccounts: []satypes.SubaccountId{
				constants.Carl_Num0,
			},

			expectedPlacedOrders:  []*types.MsgPlaceOrder{},
			expectedMatchedOrders: []*types.ClobMatch{},
		},
		"Can liquidate multiple liquidatable subaccounts and skips liquidatable subaccounts": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Carl_Num1_01BTC_Long_4600USD_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				constants.Order_Carl_Num0_Id5_Clob0_Buy2BTC_Price50000,
			},
			liquidatableSubaccounts: []satypes.SubaccountId{
				constants.Carl_Num0,
				constants.Carl_Num1,
				constants.Dave_Num0,
			},

			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				},
				{
					Order: constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				},
				{
					Order: constants.Order_Carl_Num0_Id5_Clob0_Buy2BTC_Price50000,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   10_000_000,
						Liquidated:  constants.Carl_Num1,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
								FillAmount:   10_000_000,
							},
						},
					},
				),
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
								FillAmount:   40_000_000,
							},
							{
								MakerOrderId: constants.OrderId_Alice_Num0_ClientId1_Clob0,
								FillAmount:   25_000_000,
							},
							{
								MakerOrderId: constants.OrderId_Alice_Num0_ClientId2_Clob0,
								FillAmount:   35_000_000,
							},
						},
					},
				),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis tmtypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						*genesisState = constants.TestPricesGenesisState
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = tc.perpetuals
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *types.GenesisState) {
						genesisState.ClobPairs = tc.clobs
						genesisState.LiquidationsConfig = types.LiquidationsConfig_Default
					},
				)
				return genesis
			}).WithTesting(t).Build()

			ctx, app := tApp.AdvanceToBlock(2)
			// Create all existing orders.
			existingOrderMsgs := make([]types.MsgPlaceOrder, len(tc.existingOrders))
			for i, order := range tc.existingOrders {
				existingOrderMsgs[i] = types.MsgPlaceOrder{Order: order}
			}
			for _, checkTx := range testapp.MustMakeCheckTxs(ctx, app, existingOrderMsgs...) {
				require.True(t, app.CheckTx(checkTx).IsOK())
			}

			// Update the liquidatable subaccount IDs.
			_, err := app.Server.LiquidateSubaccounts(ctx, &api.LiquidateSubaccountsRequest{
				SubaccountIds: tc.liquidatableSubaccounts,
			})
			require.NoError(t, err)

			// TODO(DEC-1971): Replace these test assertions with new verifications on operations queue.
			// Verify test expectations.
			// ctx, app = tApp.AdvanceToBlock(3)
			// placedOrders, matchedOrders := app.ClobKeeper.MemClob.GetPendingFills(ctx)
			// require.Equal(t, tc.expectedPlacedOrders, placedOrders, "Placed orders lists are not equal")
			// require.Equal(t, tc.expectedMatchedOrders, matchedOrders, "Matched orders lists are not equal")
		})
	}
}

func TestPrepareCheckState_WithProcessProposerMatchesEventsWithBadBlockHeight(t *testing.T) {
	blockHeight := uint32(6)
	processProposerMatchesEvents := types.ProcessProposerMatchesEvents{BlockHeight: blockHeight}

	// Setup keeper state.
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()

	ctx, clobKeeper, _, _, _, _, _, _ := keepertest.ClobKeepers(
		t, memClob, &mocks.BankKeeper{}, indexer_manager.NewIndexerEventManagerNoop())

	// Set the process proposer matches events from the last block.
	clobKeeper.MustSetProcessProposerMatchesEvents(
		ctx.WithBlockHeight(int64(blockHeight)),
		processProposerMatchesEvents,
	)

	// Ensure that we panic if our current block height doesn't match the block height stored with the events
	require.Panics(t, func() {
		clob.PrepareCheckState(
			ctx.WithBlockHeight(int64(blockHeight+1)),
			*clobKeeper,
			memClob,
			liquidationtypes.NewLiquidatableSubaccountIds(),
		)
	})
}

func TestCommitBlocker_WithProcessProposerMatchesEventsWithBadBlockHeight(t *testing.T) {
	blockHeight := uint32(6)
	processProposerMatchesEvents := types.ProcessProposerMatchesEvents{BlockHeight: blockHeight}

	// Setup keeper state.
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()

	ctx, clobKeeper, _, _, _, _, _, _ := keepertest.ClobKeepers(
		t, memClob, &mocks.BankKeeper{}, indexer_manager.NewIndexerEventManagerNoop())

	// Set the process proposer matches events from the last block.
	clobKeeper.MustSetProcessProposerMatchesEvents(
		ctx.WithBlockHeight(int64(blockHeight)),
		processProposerMatchesEvents,
	)

	// Ensure that we panic if our current block height doesn't match the block height stored with the events
	require.Panics(t, func() {
		clob.PrepareCheckState(
			ctx.WithBlockHeight(int64(blockHeight+1)),
			*clobKeeper,
			memClob,
			liquidationtypes.NewLiquidatableSubaccountIds(),
		)
	})
}

func TestBeginBlocker_Success(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setupState func(ctx sdk.Context, k *keeper.Keeper)
	}{
		"Initializes next block's process proposer matches events overwriting state that was set": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders: []types.Order{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						BlockHeight: lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
					},
				)
			},
		},
		"Initializes next block's process proposer matches events overwriting state that was set multiple times": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders: []types.Order{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						BlockHeight: lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
					},
				)
				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders: []types.Order{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
						},
						BlockHeight: lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
					},
				)
			},
		},
		"Initializes next block's process proposer matches events from genesis state": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()

			ctx,
				keeper,
				_,
				_,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexer_manager.NewIndexerEventManagerNoop())
			ctx = ctx.WithBlockHeight(int64(20)).WithBlockTime(unixTimeTen)

			if tc.setupState != nil {
				tc.setupState(ctx, keeper)
			}

			clob.BeginBlocker(
				ctx,
				*keeper,
			)

			// Assert expecations.
			require.Equal(
				t,
				types.ProcessProposerMatchesEvents{BlockHeight: 20},
				keeper.GetProcessProposerMatchesEvents(ctx),
			)
			require.True(t, memClob.AssertExpectations(t))
		})
	}
}

// TODO(CLOB-231): Add more test coverage to `PrepareCheckState` method.
func TestPrepareCheckState(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []*perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs                        []types.ClobPair
		preExistingStatefulOrders    []types.Order
		processProposerMatchesEvents types.ProcessProposerMatchesEvents
		// Memclob state.
		placedOperations []types.Operation

		// Parameters.
		liquidatableSubaccounts []satypes.SubaccountId

		// Expectations.
		expectedOperationsQueue []types.Operation
		expectedBids            []memclob.OrderWithRemainingSize
		expectedAsks            []memclob.OrderWithRemainingSize
	}{
		"Nothing on memclob or in state": {
			perpetuals:                []*perptypes.Perpetual{},
			subaccounts:               []satypes.Subaccount{},
			clobs:                     []types.ClobPair{},
			preExistingStatefulOrders: []types.Order{},
			processProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: 4,
			},
			placedOperations: []types.Operation{},

			liquidatableSubaccounts: []satypes.SubaccountId{},

			expectedOperationsQueue: []types.Operation{},
			expectedBids:            []memclob.OrderWithRemainingSize{},
			expectedAsks:            []memclob.OrderWithRemainingSize{},
		},
		`Regression: Local validator replays two matches of exactly the same size for the same OrderId to the memclob.
     ReplayOperations should not panic as the MatchOperations should have unique taker OrderHashes therefore the
		nonce for the second match should not already exist.`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_NoMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Alice_Num1_10_000USD,
			},
			clobs:                     []types.ClobPair{constants.ClobPair_Btc},
			preExistingStatefulOrders: []types.Order{},
			processProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight:          4,
				PlacedStatefulOrders: []types.Order{},
			},
			placedOperations: []types.Operation{
				// This order lands on the book.
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.MustGetOrder(),
				),
				// This order crosses the first order and matches for 5.
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.MustGetOrder(),
				),
				// This replacement order crosses the first order and matches for 5.
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.MustGetOrder(),
				),
			},
			liquidatableSubaccounts: []satypes.SubaccountId{},

			expectedOperationsQueue: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.MustGetOrder(),
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.MustGetOrder(),
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.MustGetOrder(),
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId,
						},
					},
				),
			},
			expectedBids: []memclob.OrderWithRemainingSize{},
			expectedAsks: []memclob.OrderWithRemainingSize{},
		},
		// TODO(CLOB-249): Fix test once sanity check is added back.
		`Regression: Proposer includes an order from subaccount 0 which is removed due to collateralization by a
		subsequent operation included by the proposer. Subsequent operations included by the proposer then recollateralize
		subaccount 0. The local validator attempts to replay the previously removed order from subaccount 0 as a taker. 
		The order already has a nonce so is not replaced during ReplayOperations. Placement of the order during Replay
		should therefore not cause a panic.`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_NoMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.User1_Num1_100_000USD,
			},
			clobs:                     []types.ClobPair{constants.ClobPair_Btc},
			preExistingStatefulOrders: []types.Order{},
			processProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: 4,
				PlacedStatefulOrders: []types.Order{
					// Place the order that will be removed due to collateralization.
					constants.LongTermOrder_User1_Num0_Id2_Clob0_Sell02BTC_Price20_GTB15.MustGetOrder(),

					// Place a very bad order for Alice_Num0. This order will match and bring Subaccount 0's collateralization
					// such that it won't be sufficiently collateralized for the first order.
					constants.LongTermOrder_User1_Num0_Id1_Clob0_Sell02BTC_Price10_GTB15.MustGetOrder(),
					constants.LongTermOrder_User1_Num1_Id1_Clob0_Buy02BTC_Price10_GTB15.MustGetOrder(),

					// This order from Alice_Num1 removes the first order from AliceNum0 due to collateralization.
					// (This order also gets removed for self-trade due to the final order in this slice,
					// though that is unrelated to what is being primarily tested here, but worth noting).
					constants.LongTermOrder_Alice_Num1_Id2_Clob0_Buy10_Price40_GTBT10.MustGetOrder(),

					// Place the inverse of the previous match so that Alice_Num0's balance is restored and sufficiently
					// collateralized to support the first order. This will allow the order in `placedOperations` to be valid
					// when we go to replay.
					constants.LongTermOrder_User1_Num1_Id2_Clob0_Sell02BTC_Price10_GTB15.MustGetOrder(),
					constants.LongTermOrder_User1_Num0_Id3_Clob0_Buy02BTC_Price10_GTB15.MustGetOrder(),
				},
			},
			placedOperations: []types.Operation{
				// This order lands on the book.
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num1_Id3_Clob0_Buy10_Price40_GTBT10.MustGetOrder(),
				),
				// This order is not placed, because it has already been placed as part of
				// the `PlacedStatefulOrders` above.
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num0_Id2_Clob0_Sell02BTC_Price20_GTB15.MustGetOrder(),
				),
			},
			liquidatableSubaccounts: []satypes.SubaccountId{},

			expectedOperationsQueue: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num0_Id2_Clob0_Sell02BTC_Price20_GTB15.MustGetOrder(),
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num0_Id1_Clob0_Sell02BTC_Price10_GTB15.MustGetOrder(),
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num1_Id1_Clob0_Buy02BTC_Price10_GTB15.MustGetOrder(),
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_User1_Num1_Id1_Clob0_Buy02BTC_Price10_GTB15,
					[]types.MakerFill{
						{
							FillAmount:   20_000_000,
							MakerOrderId: constants.LongTermOrder_User1_Num0_Id1_Clob0_Sell02BTC_Price10_GTB15.OrderId,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id2_Clob0_Buy10_Price40_GTBT10.MustGetOrder(),
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num1_Id2_Clob0_Sell02BTC_Price10_GTB15.MustGetOrder(),
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num0_Id3_Clob0_Buy02BTC_Price10_GTB15.MustGetOrder(),
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_User1_Num0_Id3_Clob0_Buy02BTC_Price10_GTB15,
					[]types.MakerFill{
						{
							FillAmount:   20_000_000,
							MakerOrderId: constants.LongTermOrder_User1_Num1_Id2_Clob0_Sell02BTC_Price10_GTB15.OrderId,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num1_Id3_Clob0_Buy10_Price40_GTBT10.MustGetOrder(),
				),
			},
			expectedBids: []memclob.OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_User1_Num1_Id3_Clob0_Buy10_Price40_GTBT10,
					RemainingSize: 10,
				},
			},
			expectedAsks: []memclob.OrderWithRemainingSize{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			ctx,
				clobKeeper,
				pricesKeeper,
				assetsKeeper,
				perpetualsKeeper,
				subaccountsKeeper,
				_,
				_ := keepertest.ClobKeepers(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx = ctx.WithIsCheckTx(true).WithBlockHeight(int64(tc.processProposerMatchesEvents.BlockHeight))
			// Create the default markets.
			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Ticker,
					p.MarketId,
					p.AtomicResolution,
					p.DefaultFundingPpm,
					p.LiquidityTier,
				)
				require.NoError(t, err)
			}

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				subaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobs {
				_, err = clobKeeper.CreatePerpetualClobPair(
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
			}

			// Initialize the liquidations config.
			err = clobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default)
			require.NoError(t, err)

			// Create all pre-existing stateful orders in state.
			for _, order := range tc.preExistingStatefulOrders {
				clobKeeper.SetStatefulOrderPlacement(ctx, order, 100)
			}

			// Set the ProcessProposerMatchesEvents in state.
			clobKeeper.MustSetProcessProposerMatchesEvents(
				ctx.WithIsCheckTx(false),
				tc.processProposerMatchesEvents,
			)

			// Set the block time on the context and of the last committed block.
			ctx = ctx.WithBlockTime(unixTimeFive)
			clobKeeper.SetBlockTimeForLastCommittedBlock(ctx)

			// Initialize the memclob with each placed operation using a forked version of state,
			// and ensure the forked state is not committed to the base state.
			setupCtx, _ := ctx.CacheContext()
			for _, operation := range tc.placedOperations {
				switch operation.Operation.(type) {
				case *types.Operation_OrderPlacement:
					order := operation.GetOrderPlacement()
					tempCtx, writeCache := setupCtx.CacheContext()
					_, _, err := clobKeeper.CheckTxPlaceOrder(
						tempCtx,
						order,
					)
					if err == nil {
						writeCache()
					} else {
						isExpectedErr := errors.Is(err, types.ErrFokOrderCouldNotBeFullyFilled) ||
							errors.Is(err, types.ErrPostOnlyWouldCrossMakerOrder)
						if !isExpectedErr {
							panic(
								fmt.Errorf(
									"Expected error ErrFokOrderCouldNotBeFullyFilled or ErrPostOnlyWouldCrossMakerOrder, got %w",
									err,
								),
							)
						}
					}
				case *types.Operation_OrderCancellation:
					orderCancellation := operation.GetOrderCancellation()
					err := clobKeeper.CheckTxCancelOrder(ctx, orderCancellation)
					if err != nil {
						panic(err)
					}
				case *types.Operation_PreexistingStatefulOrder:
					orderId := operation.GetPreexistingStatefulOrder()
					orderPlacement, found := clobKeeper.GetStatefulOrderPlacement(ctx, *orderId)
					if !found {
						panic(
							fmt.Sprintf(
								"Order ID %+v not found in state.",
								orderId,
							),
						)
					}
					tempCtx, writeCache := setupCtx.CacheContext()
					_, _, err := clobKeeper.CheckTxPlaceOrder(
						tempCtx,
						types.NewMsgPlaceOrder(orderPlacement.Order),
					)
					if err == nil {
						writeCache()
					} else {
						isExpectedErr := errors.Is(err, types.ErrFokOrderCouldNotBeFullyFilled) ||
							errors.Is(err, types.ErrPostOnlyWouldCrossMakerOrder)
						if !isExpectedErr {
							panic(isExpectedErr)
						}
					}
				}
			}

			// Set the liquidatable subaccount IDs.
			liquidatableSubaccountIds := liquidationtypes.NewLiquidatableSubaccountIds()
			liquidatableSubaccountIds.UpdateSubaccountIds(tc.liquidatableSubaccounts)

			// Run the test.
			clob.PrepareCheckState(
				ctx,
				*clobKeeper,
				memClob,
				liquidatableSubaccountIds,
			)

			// Verify test expectations.
			require.NoError(t, err)
			operationsQueue := memClob.GetOperations(ctx)

			require.Equal(t, tc.expectedOperationsQueue, operationsQueue)

			memclob.AssertMemclobHasOrders(
				t,
				ctx,
				memClob,
				tc.expectedBids,
				tc.expectedAsks,
			)
		})
	}
}
