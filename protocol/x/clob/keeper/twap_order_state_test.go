package keeper_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTWAPOrderTriggerStoreOrderingBasedOffTimestamp(t *testing.T) {
	ks := setupTestTWAPOrderState(t)
	// Create test orders with same order ID but different timestamps
	suborderId := types.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}

	// Set different trigger offsets to test timestamp ordering
	triggerOffsets := []int64{10, 5, 15, 0}

	// Add orders to trigger store with different timestamps
	for _, offset := range triggerOffsets {
		ks.ClobKeeper.AddSuborderToTriggerStore(ks.Ctx, suborderId, offset)
	}

	// Get all orders from trigger store and verify ordering
	store := ks.ClobKeeper.GetTWAPTriggerOrderPlacementStore(ks.Ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	// Expected order based on timestamp (ascending)
	expectedOrder := []int64{0, 5, 10, 15}
	index := 0

	// validate timestamps are in order
	for ; iterator.Valid(); iterator.Next() {
		timestamp := types.TimeFromTriggerKey(iterator.Key())
		require.Equal(t, ks.Ctx.BlockTime().Unix()+expectedOrder[index], timestamp)
		index++
	}
	require.Equal(t, len(expectedOrder), index)
}

func TestTWAPOrderTriggerStoreOrdering(t *testing.T) {
	ks := setupTestTWAPOrderState(t)

	// Create test orders with different timestamps and order IDs
	// In practice, we do not expect multiple instances of the
	// same suborderId in the trigger store, but this case is
	// constructed as such to test the ordering of the keystore
	// is working as expected (timestamp + orderId bytes)
	suborderIds := []types.OrderId{
		{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   types.OrderIdFlags_TwapSuborder,
			ClobPairId:   0,
		},
		{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   types.OrderIdFlags_TwapSuborder,
			ClobPairId:   0,
		},
		{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   types.OrderIdFlags_TwapSuborder,
			ClobPairId:   1,
		},
		{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   types.OrderIdFlags_TwapSuborder,
			ClobPairId:   0,
		},
	}

	// Set different trigger offsets to test timestamp ordering
	triggerOffsets := []int64{0, 5, 5, 5}

	// Add orders to trigger store with different timestamps
	for i, suborderId := range suborderIds {
		ks.ClobKeeper.AddSuborderToTriggerStore(ks.Ctx, suborderId, triggerOffsets[i])
	}

	//ks.ClobKeeper.GetTwapTriggerPlacement(ks.Ctx, suborderIds[0])

	// Get all orders from trigger store and verify ordering
	store := ks.ClobKeeper.GetTWAPTriggerOrderPlacementStore(ks.Ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	// Expected order based on timestamp (ascending) and then orderId
	expectedOrder := []types.OrderId{
		suborderIds[0], // offset 0
		suborderIds[3], // offset 5 - bob subaccountId < alice subaccountId
		suborderIds[1], // offset 5 - alice clobPairId 0 < 1
		suborderIds[2], // offset 5
	}

	// Verify the order of retrieved orders matches expected order
	index := 0
	for ; iterator.Valid(); iterator.Next() {
		var orderId types.OrderId
		// The key is [timestamp (8 bytes)][orderId state key]
		// We unmarshal just the orderId portion
		key := iterator.Key()
		orderIdBytes := key[8:]
		ks.Cdc.MustUnmarshal(orderIdBytes, &orderId)
		require.Equal(t, expectedOrder[index], orderId)
		index++
	}
	require.Equal(t, len(expectedOrder), index)
}

func TestTWAPOrderKeyBytes(t *testing.T) {
	orderId1 := types.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	orderId2 := types.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   1,
	}

	key1 := types.GetTWAPTriggerKey(5, orderId1)
	key2 := types.GetTWAPTriggerKey(5, orderId2)

	// Print the actual bytes to see the ordering
	fmt.Printf("Key1: %v\n", key1)
	fmt.Printf("Key2: %v\n", key2)

	// Compare the keys
	result := bytes.Compare(key1, key2)
	require.True(t, result < 0) // key1 should come before key2

	ks := setupTestTWAPOrderState(t)
	ks.ClobKeeper.AddSuborderToTriggerStore(ks.Ctx, orderId1, 5)
	ks.ClobKeeper.AddSuborderToTriggerStore(ks.Ctx, orderId2, 5)

	store := ks.ClobKeeper.GetTWAPTriggerOrderPlacementStore(ks.Ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	expectedOrderId := []types.OrderId{
		orderId1,
		orderId2,
	}

	index := 0
	for ; iterator.Valid(); iterator.Next() {
		var orderId types.OrderId
		ks.Cdc.MustUnmarshal(iterator.Key()[8:], &orderId)
		require.Equal(t, expectedOrderId[index], orderId)
		index++
	}
	require.Equal(t, len(expectedOrderId), index)
}

func setupTestTWAPOrderState(t *testing.T) (ks keepertest.ClobKeepersTestContext) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks = keepertest.NewClobKeepersTestContextWithUninitializedMemStore(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)
	return ks
}

func TestSetTWAPOrderPlacement(t *testing.T) {
	tests := map[string]struct {
		order             types.Order
		blockHeight       uint32
		expectedTotalLegs uint32
		expectedQuantums  uint64
	}{
		"successfully sets TWAP order with 5 minute duration and 1 minute intervals": {
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_Twap,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 1000,
				TwapParameters: &types.TwapParameters{
					Duration: 300, // 5 minutes
					Interval: 60,  // 1 minute
				},
			},
			blockHeight:       100,
			expectedTotalLegs: 5,
			expectedQuantums:  1000,
		},
		"successfully sets TWAP order with 1 hour duration and 5 minute intervals": {
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     2,
					OrderFlags:   types.OrderIdFlags_Twap,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_SELL,
				Quantums: 2000,
				TwapParameters: &types.TwapParameters{
					Duration: 3600, // 1 hour
					Interval: 300,  // 5 minutes
				},
			},
			blockHeight:       200,
			expectedTotalLegs: 12,
			expectedQuantums:  2000,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContextWithUninitializedMemStore(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)

			// Set block time for consistent testing
			ctx := ks.Ctx.WithBlockTime(time.Unix(1000, 0))

			// Set the TWAP order placement
			ks.ClobKeeper.SetTWAPOrderPlacement(ctx, tc.order, tc.blockHeight)

			// Verify the order was stored correctly
			storedOrder, found := ks.ClobKeeper.GetTwapOrderPlacement(ctx, tc.order.OrderId)
			require.True(t, found, "TWAP order should be found in store")
			require.Equal(t, tc.order, storedOrder.Order, "stored order should match input order")
			require.Equal(t,
				tc.expectedTotalLegs,
				storedOrder.RemainingLegs,
				"remaining legs should equal total legs initially",
			)
			require.Equal(t,
				tc.expectedQuantums,
				storedOrder.RemainingQuantums,
				"remaining quantums should equal initial quantums",
			)

			// Verify the first suborder was created in trigger store
			suborderId := types.OrderId{
				SubaccountId: tc.order.OrderId.SubaccountId,
				ClientId:     tc.order.OrderId.ClientId,
				OrderFlags:   types.OrderIdFlags_TwapSuborder,
				ClobPairId:   tc.order.OrderId.ClobPairId,
			}
			triggerPlacement, triggerTime, found := ks.ClobKeeper.GetTwapTriggerPlacement(ctx, suborderId)

			require.True(t, found, "trigger placement should be found")
			require.Equal(t, suborderId, triggerPlacement, "trigger placement should match suborderId")
			require.Equal(t, int64(1000), triggerTime, "trigger time should match block time")
		})
	}
}

func TestGenerateSuborder(t *testing.T) {
	tests := map[string]struct {
		twapOrderPlacement types.TwapOrderPlacement
		blockTime          int64
		clobPair           types.ClobPair
		expectedOrder      types.Order
	}{
		"buy order with positive price tolerance": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_BUY,
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 500_000, // 50% tolerance
					},
				},
				RemainingLegs:     5,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 200,
				Subticks: 1500,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003, // blockTime + TWAP_SUBORDER_GOOD_TIL_BLOCK_TIME_OFFSET
				},
			},
		},
		"sell order with negative price tolerance": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_SELL,
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 500_000, // 50% tolerance
					},
				},
				RemainingLegs:     5,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_SELL,
				Quantums: 200,
				Subticks: 500,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003, // blockTime + TWAP_SUBORDER_GOOD_TIL_BLOCK_TIME_OFFSET
				},
			},
		},
		"buy order with 2.5% price tolerance and rounded up subticks": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_BUY,
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 25_000, // 2.5% tolerance
					},
				},
				RemainingLegs:     5,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 200,
				Subticks: 1030, // 1000 * (1 + 0.025) rounded up to nearest 10
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003,
				},
			},
		},
		"sell order with 2.5% price tolerance and rounded down subticks": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_SELL,
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 25_000, // 2.5% tolerance
					},
				},
				RemainingLegs:     5,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_SELL,
				Quantums: 200,
				Subticks: 970, // 1000 * (1 - 0.025) rounded down to nearest 10
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003,
				},
			},
		},
		"buy order with 2.5% price tolerance with rounded up subticks and quantums": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_BUY,
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 25_000, // 2.5% tolerance
					},
				},
				RemainingLegs:     6,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 165,  // 1000 / 6 = 166.67 rounded down to nearest 5 (step base quantums)
				Subticks: 1030, // 1000 * (1 + 0.025) rounded up to nearest 10
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003,
				},
			},
		},
		"buy order with 2.5% price tolerance with rounded up subticks and max catchup quantums": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_BUY,
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 25_000, // 2.5% tolerance
					},
				},
				RemainingLegs:     1,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 500,  // (1000 / 6) * 3
				Subticks: 1030, // 1000 * (1 + 0.025) rounded up to nearest 10
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003,
				},
			},
		},
		"limit twap buy order with 50% price tolerance": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_BUY,
					Subticks: 800, // limit price for twap order
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 500_000, // 50% tolerance
					},
				},
				RemainingLegs:     5,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 200,
				Subticks: 800,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003,
				},
			},
		},
		"limit twap sell order with 50% price tolerance": {
			twapOrderPlacement: types.TwapOrderPlacement{
				Order: types.Order{
					OrderId: types.OrderId{
						SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_Twap,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_SELL,
					Subticks: 1200, // limit price for twap order
					Quantums: 1000,
					TwapParameters: &types.TwapParameters{
						Duration:       360,
						Interval:       60,
						PriceTolerance: 500_000, // 50% tolerance
					},
				},
				RemainingLegs:     5,
				RemainingQuantums: 1000,
			},
			blockTime: 1000,
			clobPair: types.ClobPair{
				Id:                        0,
				StepBaseQuantums:          100,
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedOrder: types.Order{
				OrderId: types.OrderId{
					SubaccountId: satypes.SubaccountId{Owner: "owner", Number: 1},
					ClientId:     1,
					OrderFlags:   types.OrderIdFlags_TwapSuborder,
					ClobPairId:   0,
				},
				Side:     types.Order_SIDE_SELL,
				Quantums: 200,
				Subticks: 1200,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 1003,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters
			memClob := &mocks.MemClob{}
			indexerEventManager := &mocks.IndexerEventManager{}
			clobKeeper := &mocks.MemClobKeeper{}

			memClob.On("GetClobKeeper").Return(&clobKeeper)
			memClob.On("SetClobKeeper", mock.Anything).Return()
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContextWithUninitializedMemStore(
				t,
				memClob,
				&mocks.BankKeeper{},
				indexerEventManager,
			)

			indexerEventManager.On(
				"AddTxnEvent",
				ks.Ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				mock.Anything,
			).Return()

			// Set block time for consistent testing
			ctx := ks.Ctx.WithBlockTime(time.Unix(tc.blockTime, 0))

			keepertest.CreateTestPricesAndPerpetualMarkets(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				ks.PricesKeeper,
				[]perptypes.Perpetual{
					*perptest.GeneratePerpetual(
						perptest.WithId(tc.twapOrderPlacement.Order.OrderId.ClobPairId),
						perptest.WithMarketId(tc.twapOrderPlacement.Order.OrderId.ClobPairId),
					),
				},
				[]pricestypes.MarketParamPrice{
					*pricestest.GenerateMarketParamPrice(
						pricestest.WithId(tc.twapOrderPlacement.Order.OrderId.ClobPairId),
						pricestest.WithPriceValue(100_000), // subticks = 1_000
					),
				},
			)

			clobPair := *clobtest.GenerateClobPair(
				clobtest.WithId(tc.twapOrderPlacement.Order.OrderId.ClobPairId),
				clobtest.WithPerpetualId(tc.twapOrderPlacement.Order.OrderId.ClobPairId),
			)

			keepertest.CreateTestClobPairs(t, ks.Ctx, ks.ClobKeeper, []types.ClobPair{clobPair})

			// Generate the suborder
			suborderId := types.OrderId{
				SubaccountId: tc.twapOrderPlacement.Order.OrderId.SubaccountId,
				ClientId:     tc.twapOrderPlacement.Order.OrderId.ClientId,
				OrderFlags:   types.OrderIdFlags_TwapSuborder,
				ClobPairId:   tc.twapOrderPlacement.Order.OrderId.ClobPairId,
			}
			generatedOrder, isGenerated := ks.ClobKeeper.GenerateSuborder(
				ctx,
				suborderId,
				tc.twapOrderPlacement,
				tc.blockTime,
			)

			// Verify the generated order matches expectations
			require.True(t, isGenerated)
			require.Equal(t, tc.expectedOrder.OrderId, generatedOrder.OrderId)
			require.Equal(t, tc.expectedOrder.Side, generatedOrder.Side)
			require.Equal(t, tc.expectedOrder.GoodTilOneof, generatedOrder.GoodTilOneof)
			require.Equal(t, tc.expectedOrder.Subticks, generatedOrder.Subticks)
			require.Equal(t, tc.expectedOrder.Quantums, generatedOrder.Quantums)
		})
	}
}
