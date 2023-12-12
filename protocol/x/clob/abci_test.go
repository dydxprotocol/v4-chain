package clob_test

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquidationtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
			[]byte(types.BlockHeightToPotentiallyPrunableOrdersPrefix),
		)

		potentiallyPrunableOrdersBytes := blockHeightToPotentiallyPrunableOrdersStore.Get(
			lib.Uint32ToKey(blockHeight),
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

func TestEndBlocker_Failure(t *testing.T) {
	blockHeight := uint32(5)
	tests := map[string]struct {
		blockTime               time.Time
		expiredStatefulOrderIds []types.OrderId
		setupState              func(ctx sdk.Context, k keepertest.ClobKeepersTestContext, m *mocks.MemClob)

		expectedPanicMessage string
	}{
		"Panics if cancelled order ids and expired order ids overlap": {
			blockTime: unixTimeFifteen,
			expiredStatefulOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				expiredOrder := constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20

				ks.ClobKeeper.SetLongTermOrderPlacement(ctx, expiredOrder, blockHeight)
				ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					expiredOrder.MustGetUnixGoodTilBlockTime(),
					expiredOrder.OrderId,
				)

				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						BlockHeight: blockHeight,
						PlacedStatefulCancellationOrderIds: []types.OrderId{
							expiredOrder.OrderId,
						},
					},
				)
			},
			expectedPanicMessage: fmt.Sprintf(
				"PruneUntriggeredConditionalOrders: duplicate order id %+v in expired and "+
					"cancelled order lists", constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()

			mockIndexerEventManager := &mocks.IndexerEventManager{}

			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
			ctx := ks.Ctx.WithBlockHeight(int64(blockHeight)).WithBlockTime(tc.blockTime)

			for _, orderId := range tc.expiredStatefulOrderIds {
				mockIndexerEventManager.On("AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewStatefulOrderRemovalEvent(
							orderId,
							indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_EXPIRED,
						),
					),
				).Once().Return()
			}

			tc.setupState(ctx, ks, memClob)

			require.PanicsWithValue(
				t,
				tc.expectedPanicMessage,
				func() {
					clob.EndBlocker(
						ctx,
						*ks.ClobKeeper,
					)
				},
			)
		})
	}
}

func TestEndBlocker_Success(t *testing.T) {
	prunedOrderIdOne := types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0}
	prunedOrderIdTwo := types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1}
	orderIdThree := types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 2}
	blockHeight := uint32(5)

	tests := map[string]struct {
		// Setup.
		setupState func(ctx sdk.Context, k keepertest.ClobKeepersTestContext, m *mocks.MemClob)
		blockTime  time.Time

		// Expectations.
		expectedFillAmounts                  map[types.OrderId]satypes.BaseQuantums
		expectedPruneableBlockHeights        map[uint32][]types.OrderId
		expectedPrunedOrders                 map[types.OrderId]bool
		expectedStatefulPlacementInState     map[types.OrderId]bool
		expectedStatefulOrderTimeSlice       map[time.Time][]types.OrderId
		expectedProcessProposerMatchesEvents types.ProcessProposerMatchesEvents
		expectedUntriggeredConditionalOrders map[types.ClobPairId]*keeper.UntriggeredConditionalOrders
		expectedTriggeredConditionalOrderIds []types.OrderId
	}{
		"Prunes existing Short-Term orders and seen place orders correctly": {
			blockTime: unixTimeTen,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				// Set `prunedOrderIdOne` and `prunedOrderIdTwo` as existing orders which already have fill amounts.
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					prunedOrderIdOne,
					100,
					blockHeight,
				)

				// Set `prunedOrderIdTwo` to be prunable at the next block height (this takes precedent of the blockHeight
				// set in `AddOrdersForPruning`).
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					prunedOrderIdTwo,
					100,
					blockHeight+1,
				)

				// This order should not be pruned.
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					orderIdThree,
					150,
					blockHeight+10,
				)

				// Set both of these orders as prunable at the current `blockHeight` so we can assert that they were pruned
				// correctly.
				ks.ClobKeeper.AddOrdersForPruning(
					ctx,
					[]types.OrderId{prunedOrderIdOne, prunedOrderIdTwo},
					blockHeight,
				)

				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						BlockHeight:               blockHeight,
						OrderIdsFilledInLastBlock: []types.OrderId{prunedOrderIdTwo, orderIdThree},
					},
				)
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				prunedOrderIdTwo: 100,
				orderIdThree:     150,
			},
			expectedPrunedOrders: map[types.OrderId]bool{
				prunedOrderIdOne: true,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight:               blockHeight,
				OrderIdsFilledInLastBlock: []types.OrderId{prunedOrderIdTwo, orderIdThree},
			},
		},
		"Prunes expired and cancelled untriggered conditional orders from UntriggeredConditionalorders": {
			blockTime: unixTimeFifteen,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				// expired orders
				orderToPrune1 := constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20
				orderToPrune2 := constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25
				orderToPrune3 := constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit20
				// cancelled order
				orderToPrune4 := constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30_TakeProfit20

				// add expired orders to state, cancelled orders already removed in DeliverTx
				orders := []types.Order{
					orderToPrune1,
					orderToPrune2,
					orderToPrune3,
				}
				for _, order := range orders {
					ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, 0)
					ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
						ctx,
						order.MustGetUnixGoodTilBlockTime(),
						order.OrderId,
					)
				}

				ks.ClobKeeper.UntriggeredConditionalOrders = map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{
					constants.ClobPair_Btc.GetClobPairId(): {
						OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{orderToPrune1, orderToPrune2, orderToPrune4},
						OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
					},
					constants.ClobPair_Eth.GetClobPairId(): {
						OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
						OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{orderToPrune3},
					},
				}

				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulCancellationOrderIds: []types.OrderId{
							orderToPrune4.OrderId,
						},
						BlockHeight: blockHeight,
					},
				)
			},
			expectedUntriggeredConditionalOrders: map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{},
			expectedStatefulPlacementInState: map[types.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId:   false,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25.OrderId:  false,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit20.OrderId: false,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30_TakeProfit20.OrderId: false,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulCancellationOrderIds: []types.OrderId{
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30_TakeProfit20.OrderId,
				},
				ExpiredStatefulOrderIds: []types.OrderId{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit20.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		`Adds newly-placed conditional order to UntriggeredConditionalOrders, but does not add
		cancelled order`: {
			blockTime: unixTimeTen,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				orders := []types.Order{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
					constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
					constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				}
				for _, order := range orders {
					ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, blockHeight)
				}

				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedConditionalOrderIds: []types.OrderId{
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
						},
						PlacedStatefulCancellationOrderIds: []types.OrderId{
							constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
						},
						BlockHeight: blockHeight,
					},
				)
			},
			expectedUntriggeredConditionalOrders: map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{
				constants.ClobPair_Btc.GetClobPairId(): {
					OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
					},
					OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedConditionalOrderIds: []types.OrderId{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
				},
				PlacedStatefulCancellationOrderIds: []types.OrderId{
					constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		`Polls triggered conditional orders from UntriggeredConditionalOrders, update state and
		ProcessProposerMatchesEvents`: {
			blockTime: unixTimeTen,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				// Update perpetual prices
				err := ks.PricesKeeper.UpdateMarketPrices(ctx, []*prices.MsgUpdateMarketPrices_MarketPrice{
					{
						MarketId: constants.ClobPair_Btc.Id,
						Price: types.SubticksToPrice(
							types.Subticks(10),
							constants.BtcUsdExponent,
							constants.ClobPair_Btc,
							constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.AtomicResolution,
							lib.QuoteCurrencyAtomicResolution,
						),
					},
				})
				require.NoError(t, err)

				err = ks.PricesKeeper.UpdateMarketPrices(ctx, []*prices.MsgUpdateMarketPrices_MarketPrice{
					{
						MarketId: constants.ClobPair_Eth.Id,
						Price: types.SubticksToPrice(
							types.Subticks(35),
							constants.EthUsdExponent,
							constants.ClobPair_Eth,
							constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.AtomicResolution,
							lib.QuoteCurrencyAtomicResolution,
						),
					},
				})
				require.NoError(t, err)

				ks.ClobKeeper.UntriggeredConditionalOrders = map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{
					constants.ClobPair_Btc.GetClobPairId(): {
						// 10 oracle price subticks triggers 3 orders here.
						OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
							constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
							constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
							constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
						},
						// 10 oracle price subticks triggers no orders here.
						OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
							constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
							constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
							constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
						},
					},
					constants.ClobPair_Eth.GetClobPairId(): {
						// 35 oracle price subticks triggers no orders here.
						OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
							constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit30,
						},
						// 35 oracle price subticks triggers one order here.
						OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
							constants.ConditionalOrder_Alice_Num0_Id3_Clob1_Buy25_Price10_GTBT15_StopLoss20,
						},
					},
				}

				for _, untrigCondOrders := range ks.ClobKeeper.UntriggeredConditionalOrders {
					for _, conditionalOrder := range untrigCondOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice {
						ks.ClobKeeper.SetLongTermOrderPlacement(ctx, conditionalOrder, blockHeight)
					}
					for _, conditionalOrder := range untrigCondOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice {
						ks.ClobKeeper.SetLongTermOrderPlacement(ctx, conditionalOrder, blockHeight)
					}
				}

				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						BlockHeight: blockHeight,
					},
				)
			},

			expectedTriggeredConditionalOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob1_Buy25_Price10_GTBT15_StopLoss20.OrderId,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
				ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id3_Clob1_Buy25_Price10_GTBT15_StopLoss20.OrderId,
				},
			},
			expectedUntriggeredConditionalOrders: map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{
				constants.ClobPair_Btc.GetClobPairId(): {
					OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
					},
					OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
						constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
						constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
						constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
					},
				},
				constants.ClobPair_Eth.GetClobPairId(): {
					OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit30,
					},
					OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
				},
			},
		},
		"Removes expired stateful orders and updates process proposer matches events": {
			blockTime: unixTimeTen,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				// These orders should get removed.
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					5,
					blockHeight,
				)
				ks.ClobKeeper.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					blockHeight,
				)
				ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
				)
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
					5,
					blockHeight,
				)
				ks.ClobKeeper.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					blockHeight,
				)
				ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
				)

				// This order should not be pruned.
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					orderIdThree,
					150,
					blockHeight+10,
				)

				// This order should not get removed.
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					5,
					blockHeight,
				)
				ks.ClobKeeper.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					blockHeight,
				)
				ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				)

				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						BlockHeight:               blockHeight,
						OrderIdsFilledInLastBlock: []types.OrderId{orderIdThree},
					},
				)
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				orderIdThree: 150,
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
				BlockHeight:               blockHeight,
				OrderIdsFilledInLastBlock: []types.OrderId{orderIdThree},
			},
		},
		"Stateful order placements are not overwritten": {
			blockTime: unixTimeTen,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				// This order should not be pruned.
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					orderIdThree,
					150,
					blockHeight+10,
				)

				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedLongTermOrderIds: []types.OrderId{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
						},
						OrderIdsFilledInLastBlock: []types.OrderId{orderIdThree},
						ExpiredStatefulOrderIds:   []types.OrderId{},
						BlockHeight:               blockHeight,
					},
				)
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				orderIdThree: 150,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedLongTermOrderIds: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				},
				OrderIdsFilledInLastBlock: []types.OrderId{orderIdThree},
				BlockHeight:               blockHeight,
			},
		},
		"Does not send order update message offchain message for a stateful order fill that got cancelled": {
			blockTime: unixTimeTen,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext, m *mocks.MemClob) {
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
					20,
					blockHeight+10,
				)
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
					20,
					blockHeight+10,
				)
				ks.ClobKeeper.MustSetProcessProposerMatchesEvents(ctx, types.ProcessProposerMatchesEvents{
					PlacedStatefulCancellationOrderIds: []types.OrderId{
						constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
					},
					OrderIdsFilledInLastBlock: []types.OrderId{
						constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
						constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
					},
					BlockHeight: blockHeight,
				})
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulCancellationOrderIds: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
				},
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId(),
					constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.GetOrderId(): 20,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.GetOrderId():    20,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()

			mockIndexerEventManager := &mocks.IndexerEventManager{}

			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				types.InsuranceFundModuleAddress,
				constants.Usdc.Denom,
			).Return(
				sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(new(big.Int))),
			)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, mockIndexerEventManager)
			ctx := ks.Ctx.WithBlockHeight(int64(blockHeight)).WithBlockTime(tc.blockTime)

			// Set up prices keeper markets with default prices.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers on perpetuals keeper.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			// Set up clob keeper perpetuals and clob pairs.
			for _, p := range []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			} {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
				)
				require.NoError(t, err)
			}
			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			memClob.On("CreateOrderbook", ctx, constants.ClobPair_Btc).Return()

			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			mockIndexerEventManager.On("AddTxnEvent",
				ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						0,
						0,
						constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Ticker,
						constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.MarketId,
						constants.ClobPair_Btc.Status,
						constants.ClobPair_Btc.QuantumConversionExponent,
						constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.AtomicResolution,
						constants.ClobPair_Btc.SubticksPerTick,
						constants.ClobPair_Btc.StepBaseQuantums,
						constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.LiquidityTier,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)
			memClob.On("CreateOrderbook", ctx, constants.ClobPair_Eth).Return()
			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			mockIndexerEventManager.On("AddTxnEvent",
				ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						1,
						1,
						constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.Ticker,
						constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.MarketId,
						constants.ClobPair_Eth.Status,
						constants.ClobPair_Eth.QuantumConversionExponent,
						constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.AtomicResolution,
						constants.ClobPair_Eth.SubticksPerTick,
						constants.ClobPair_Eth.StepBaseQuantums,
						constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.LiquidityTier,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ctx,
				constants.ClobPair_Eth.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Eth),
				satypes.BaseQuantums(constants.ClobPair_Eth.StepBaseQuantums),
				constants.ClobPair_Eth.QuantumConversionExponent,
				constants.ClobPair_Eth.SubticksPerTick,
				constants.ClobPair_Eth.Status,
			)
			require.NoError(t, err)

			if tc.setupState != nil {
				tc.setupState(ctx, ks, memClob)
			}

			// Assert that the indexer events for Expired Stateful Orders were emitted.
			for _, orderId := range tc.expectedProcessProposerMatchesEvents.ExpiredStatefulOrderIds {
				mockIndexerEventManager.On("AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewStatefulOrderRemovalEvent(
							orderId,
							indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_EXPIRED,
						),
					),
				).Once().Return()
			}

			// Assert that the indexer events for triggered conditional orders were emitted.
			for _, orderId := range tc.expectedTriggeredConditionalOrderIds {
				mockIndexerEventManager.On("AddTxnEvent",
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewConditionalOrderTriggeredEvent(
							orderId,
						),
					),
				).Once().Return()
			}

			clob.EndBlocker(
				ctx,
				*ks.ClobKeeper,
			)

			assertFillAmountAndPruneState(
				t,
				ks.ClobKeeper,
				ctx,
				ks.StoreKey,
				tc.expectedFillAmounts,
				tc.expectedPruneableBlockHeights,
				tc.expectedPrunedOrders,
			)

			require.True(t, memClob.AssertExpectations(t))

			for orderId, exists := range tc.expectedStatefulPlacementInState {
				_, found := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, orderId)
				require.Equal(t, exists, found)
			}

			for time, expected := range tc.expectedStatefulOrderTimeSlice {
				actual := ks.ClobKeeper.GetStatefulOrdersTimeSlice(ctx, time)
				require.Equal(t, expected, actual)
			}

			actualProcessProposerMatchesEvents := ks.ClobKeeper.GetProcessProposerMatchesEvents(ctx)
			// Sort the conditional order ids triggered in the last block for
			// comparison to expected triggered conditional orders.
			sort.Sort(types.SortedOrders(actualProcessProposerMatchesEvents.ConditionalOrderIdsTriggeredInLastBlock))
			require.Equal(
				t,
				tc.expectedProcessProposerMatchesEvents,
				actualProcessProposerMatchesEvents,
			)

			// Triggered conditional order placements should be shifted from the untriggered store to the triggered store.
			for _, triggeredConditionalOrderId := range actualProcessProposerMatchesEvents.
				ConditionalOrderIdsTriggeredInLastBlock {
				// TODO(CLOB-746) Once R/W methods are created, substitute those methods here.
				triggeredConditionalOrderMemstore := ks.ClobKeeper.GetTriggeredConditionalOrderPlacementMemStore(ctx)
				untriggeredConditionalOrderMemstore := ks.ClobKeeper.GetUntriggeredConditionalOrderPlacementMemStore(ctx)
				exists := triggeredConditionalOrderMemstore.Has(triggeredConditionalOrderId.ToStateKey())
				require.True(t, exists)
				exists = untriggeredConditionalOrderMemstore.Has(triggeredConditionalOrderId.ToStateKey())
				require.False(t, exists)
			}

			if tc.expectedUntriggeredConditionalOrders != nil {
				require.Equal(
					t,
					tc.expectedUntriggeredConditionalOrders,
					ks.ClobKeeper.UntriggeredConditionalOrders,
				)
			}

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
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis tmtypes.GenesisDoc) {
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
			}).Build()

			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// Create all existing orders.
			existingOrderMsgs := make([]types.MsgPlaceOrder, len(tc.existingOrders))
			for i, order := range tc.existingOrders {
				existingOrderMsgs[i] = types.MsgPlaceOrder{Order: order}
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, existingOrderMsgs...) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			// Update the liquidatable subaccount IDs.
			_, err := tApp.App.Server.LiquidateSubaccounts(ctx, &api.LiquidateSubaccountsRequest{
				LiquidatableSubaccountIds: tc.liquidatableSubaccounts,
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

	ks := keepertest.NewClobKeepersTestContext(
		t, memClob, &mocks.BankKeeper{}, indexer_manager.NewIndexerEventManagerNoop())

	// Set the process proposer matches events from the last block.
	ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
		ks.Ctx.WithBlockHeight(int64(blockHeight)),
		processProposerMatchesEvents,
	)

	// Ensure that we panic if our current block height doesn't match the block height stored with the events
	require.Panics(t, func() {
		clob.PrepareCheckState(
			ks.Ctx.WithBlockHeight(int64(blockHeight+1)),
			ks.ClobKeeper,
			liquidationtypes.NewDaemonLiquidationInfo(),
		)
	})
}

func TestCommitBlocker_WithProcessProposerMatchesEventsWithBadBlockHeight(t *testing.T) {
	blockHeight := uint32(6)
	processProposerMatchesEvents := types.ProcessProposerMatchesEvents{BlockHeight: blockHeight}

	// Setup keeper state.
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()

	ks := keepertest.NewClobKeepersTestContext(
		t, memClob, &mocks.BankKeeper{}, indexer_manager.NewIndexerEventManagerNoop())

	// Set the process proposer matches events from the last block.
	ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
		ks.Ctx.WithBlockHeight(int64(blockHeight)),
		processProposerMatchesEvents,
	)

	// Ensure that we panic if our current block height doesn't match the block height stored with the events
	require.Panics(t, func() {
		clob.PrepareCheckState(
			ks.Ctx.WithBlockHeight(int64(blockHeight+1)),
			ks.ClobKeeper,
			liquidationtypes.NewDaemonLiquidationInfo(),
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
						PlacedLongTermOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
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
						PlacedLongTermOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						BlockHeight: lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
					},
				)
				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedLongTermOrderIds: []types.OrderId{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
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

			ks := keepertest.NewClobKeepersTestContext(
				t, memClob, &mocks.BankKeeper{}, indexer_manager.NewIndexerEventManagerNoop())
			ctx := ks.Ctx.WithBlockHeight(int64(20)).WithBlockTime(unixTimeTen)

			if tc.setupState != nil {
				tc.setupState(ks.Ctx, ks.ClobKeeper)
			}

			clob.BeginBlocker(
				ctx,
				ks.ClobKeeper,
			)

			// Assert expecations.
			require.Equal(
				t,
				types.ProcessProposerMatchesEvents{BlockHeight: 20},
				ks.ClobKeeper.GetProcessProposerMatchesEvents(ctx),
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
		expectedOperationsQueue []types.InternalOperation
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

			expectedOperationsQueue: []types.InternalOperation{},
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
				BlockHeight:            4,
				PlacedLongTermOrderIds: []types.OrderId{},
			},
			placedOperations: []types.Operation{
				// This order lands on the book.
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.MustGetOrder(),
				),
				// This order crosses the first order and matches for 5.
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.MustGetOrder(),
				),
				// This replacement order crosses the first order and matches for 5.
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.MustGetOrder(),
				),
			},
			liquidatableSubaccounts: []satypes.SubaccountId{},

			expectedOperationsQueue: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.MustGetOrder(),
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.MustGetOrder(),
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20.OrderId,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.MustGetOrder(),
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
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
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				types.InsuranceFundModuleAddress,
				constants.Usdc.Denom,
			).Return(sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(new(big.Int))))
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockHeight(int64(tc.processProposerMatchesEvents.BlockHeight))
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParams))

			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
				)
				require.NoError(t, err)
			}

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobs {
				_, err = ks.ClobKeeper.CreatePerpetualClobPair(
					ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			// Initialize the liquidations config.
			err = ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default)
			require.NoError(t, err)

			// Create all pre-existing stateful orders in state.
			for _, order := range tc.preExistingStatefulOrders {
				ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, 100)
			}

			// Set the ProcessProposerMatchesEvents in state.
			ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
				ctx.WithIsCheckTx(false),
				tc.processProposerMatchesEvents,
			)

			// Set the blocktime of the last committed block.
			ctx = ctx.WithBlockTime(unixTimeFive)
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Timestamp: unixTimeFive,
			})

			// Initialize the memclob with each placed operation using a forked version of state,
			// and ensure the forked state is not committed to the base state.
			setupCtx, _ := ctx.CacheContext()
			for _, operation := range tc.placedOperations {
				switch operation.Operation.(type) {
				case *types.Operation_ShortTermOrderPlacement:
					order := operation.GetShortTermOrderPlacement()
					tempCtx, writeCache := setupCtx.CacheContext()
					tempCtx = tempCtx.WithTxBytes(order.Order.GetOrderHash().ToBytes())
					_, _, err := ks.ClobKeeper.PlaceShortTermOrder(
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
				case *types.Operation_ShortTermOrderCancellation:
					orderCancellation := operation.GetShortTermOrderCancellation()
					err := ks.ClobKeeper.CancelShortTermOrder(ctx, orderCancellation)
					if err != nil {
						panic(err)
					}
				case *types.Operation_PreexistingStatefulOrder:
					orderId := operation.GetPreexistingStatefulOrder()
					orderPlacement, found := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, *orderId)
					if !found {
						panic(
							fmt.Sprintf(
								"Order ID %+v not found in state.",
								orderId,
							),
						)
					}
					tempCtx, writeCache := setupCtx.CacheContext()
					_, _, err := ks.ClobKeeper.PlaceShortTermOrder(
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
			liquidatableSubaccountIds := liquidationtypes.NewDaemonLiquidationInfo()
			liquidatableSubaccountIds.UpdateLiquidatableSubaccountIds(tc.liquidatableSubaccounts)

			// Run the test.
			clob.PrepareCheckState(
				ctx,
				ks.ClobKeeper,
				liquidatableSubaccountIds,
			)

			// Verify test expectations.
			require.NoError(t, err)
			operationsQueue, _ := memClob.GetOperationsToReplay(ctx)

			require.Equal(t, tc.expectedOperationsQueue, operationsQueue)

			memclob.AssertMemclobHasOrders(
				t,
				ctx,
				memClob,
				tc.expectedBids,
				tc.expectedAsks,
			)

			memclob.AssertMemclobHasShortTermTxBytes(
				t,
				ctx,
				memClob,
				tc.expectedOperationsQueue,
				tc.expectedBids,
				tc.expectedAsks,
			)
		})
	}
}
