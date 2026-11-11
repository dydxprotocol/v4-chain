package keeper_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
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

func TestDeletePendingTwapTrigger(t *testing.T) {
	ks := setupTestTWAPOrderState(t)

	suborderId := types.OrderId{
		SubaccountId: constants.Alice_Num0,
		ClientId:     0,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	otherSuborderId := types.OrderId{
		SubaccountId: constants.Bob_Num0,
		ClientId:     0,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}

	ks.ClobKeeper.AddSuborderToTriggerStore(ks.Ctx, suborderId, 5)
	ks.ClobKeeper.AddSuborderToTriggerStore(ks.Ctx, otherSuborderId, 0)

	triggerTime, found := ks.ClobKeeper.GetPendingTwapTrigger(ks.Ctx, suborderId)
	require.True(t, found)
	require.Equal(t, ks.Ctx.BlockTime().Unix()+5, triggerTime)

	require.True(t, ks.ClobKeeper.DeletePendingTwapTrigger(ks.Ctx, suborderId))

	_, _, found = ks.ClobKeeper.GetTwapTriggerPlacement(ks.Ctx, suborderId)
	require.False(t, found)
	_, found = ks.ClobKeeper.GetPendingTwapTrigger(ks.Ctx, suborderId)
	require.False(t, found)
	require.False(t, ks.ClobKeeper.DeletePendingTwapTrigger(ks.Ctx, suborderId))

	remainingId, _, found := ks.ClobKeeper.GetTwapTriggerPlacement(ks.Ctx, otherSuborderId)
	require.True(t, found)
	require.Equal(t, otherSuborderId, remainingId)
	_, found = ks.ClobKeeper.GetPendingTwapTrigger(ks.Ctx, otherSuborderId)
	require.True(t, found)
}

func TestDeleteTWAPOrderPlacement_ClearsTriggerStore(t *testing.T) {
	ks := setupTestTWAPOrderState(t)
	ctx := ks.Ctx.WithBlockTime(time.Unix(1000, 0))

	twapOrder := constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO
	ks.ClobKeeper.SetTWAPOrderPlacement(ctx, twapOrder, 5)

	suborderId := types.OrderId{
		SubaccountId: twapOrder.OrderId.SubaccountId,
		ClientId:     twapOrder.OrderId.ClientId,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrder.OrderId.ClobPairId,
	}

	// Sanity check: the pending trigger exists before deletion.
	_, _, found := ks.ClobKeeper.GetTwapTriggerPlacement(ctx, suborderId)
	require.True(t, found)

	ks.ClobKeeper.DeleteTWAPOrderPlacement(ctx, twapOrder.OrderId)

	_, found = ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.False(t, found)

	_, _, found = ks.ClobKeeper.GetTwapTriggerPlacement(ctx, suborderId)
	require.False(t, found)
	_, found = ks.ClobKeeper.GetPendingTwapTrigger(ctx, suborderId)
	require.False(t, found)
}

func TestIsOrphanTwapSuborder(t *testing.T) {
	ks := setupTestTWAPOrderState(t)
	ctx := ks.Ctx.WithBlockTime(time.Unix(1000, 0))

	// Not a suborder at all.
	require.False(
		t,
		ks.ClobKeeper.IsOrphanTwapSuborder(
			ctx,
			constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
		),
	)

	twapOrder := constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO
	ks.ClobKeeper.SetTWAPOrderPlacement(ctx, twapOrder, 5)

	suborderId := types.OrderId{
		SubaccountId: twapOrder.OrderId.SubaccountId,
		ClientId:     twapOrder.OrderId.ClientId,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrder.OrderId.ClobPairId,
	}

	// Parent placement exists: not orphaned.
	require.False(t, ks.ClobKeeper.IsOrphanTwapSuborder(ctx, suborderId))

	ks.ClobKeeper.DeleteTWAPOrderPlacement(ctx, twapOrder.OrderId)

	// Parent placement is gone: orphaned.
	require.True(t, ks.ClobKeeper.IsOrphanTwapSuborder(ctx, suborderId))
}

// setRawTwapTrigger writes a trigger-store entry without an index entry, as pre-index state did.
func setRawTwapTrigger(
	ks keepertest.ClobKeepersTestContext,
	ctx sdk.Context,
	triggerTime int64,
	suborderId types.OrderId,
) {
	ks.ClobKeeper.GetTWAPTriggerOrderPlacementStore(ctx).Set(
		types.GetTWAPTriggerKey(triggerTime, suborderId),
		[]byte{},
	)
}

func TestIndexTwapTriggers(t *testing.T) {
	ks := setupTestTWAPOrderState(t)
	ctx := ks.Ctx.WithBlockTime(time.Unix(1000, 0))

	// A live parent whose first-due trigger (created by SetTWAPOrderPlacement at the block time) is
	// left unindexed, plus a later duplicate trigger for the same suborder id.
	twapOrder := constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO
	ks.ClobKeeper.SetTWAPOrderPlacement(ctx, twapOrder, 5)
	liveSuborderId := types.OrderId{
		SubaccountId: twapOrder.OrderId.SubaccountId,
		ClientId:     twapOrder.OrderId.ClientId,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrder.OrderId.ClobPairId,
	}
	ks.ClobKeeper.GetTWAPTriggerIndexStore(ctx).Delete(liveSuborderId.ToStateKey())
	setRawTwapTrigger(ks, ctx, 1030, liveSuborderId)
	require.Equal(t, 2, countTriggersForSuborder(ks, ctx, liveSuborderId))

	// An orphan trigger whose parent placement is missing.
	orphanSuborderId := types.OrderId{
		SubaccountId: constants.Carl_Num0,
		ClientId:     7,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   0,
	}
	setRawTwapTrigger(ks, ctx, 1000, orphanSuborderId)

	indexed, deleted := ks.ClobKeeper.IndexTwapTriggers(ctx)
	require.Equal(t, 1, indexed)
	require.Equal(t, 2, deleted)

	// The orphan trigger is gone.
	_, _, found := ks.ClobKeeper.GetTwapTriggerPlacement(ctx, orphanSuborderId)
	require.False(t, found)
	_, found = ks.ClobKeeper.GetPendingTwapTrigger(ctx, orphanSuborderId)
	require.False(t, found)

	// The live suborder keeps only its first-due trigger, and the index points at it.
	require.Equal(t, 1, countTriggersForSuborder(ks, ctx, liveSuborderId))
	_, scannedTime, found := ks.ClobKeeper.GetTwapTriggerPlacement(ctx, liveSuborderId)
	require.True(t, found)
	require.Equal(t, int64(1000), scannedTime)
	indexedTime, found := ks.ClobKeeper.GetPendingTwapTrigger(ctx, liveSuborderId)
	require.True(t, found)
	require.Equal(t, int64(1000), indexedTime)
}

// setupTwapSweepTest builds a keeper with a funded subaccount, perpetual, and CLOB pair so that
// GenerateAndPlaceTriggeredTwapSuborders can actually place generated suborders through the normal
// stateful order placement path (collateralization check included).
func setupTwapSweepTest(t *testing.T) (ks keepertest.ClobKeepersTestContext, ctx sdk.Context) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	memClob.On("CreateOrderbook", mock.Anything).Return()
	indexerEventManager := &mocks.IndexerEventManager{}
	indexerEventManager.On("AddTxnEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	ks = keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)

	require.NoError(t, keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper))
	keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)
	require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ks.Ctx, constants.PerpetualFeeParams))

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

	ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, satypes.Subaccount{
		Id: &constants.Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&constants.Usdc_Asset_100_000,
		},
	})

	clobPair := constants.ClobPair_Btc
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

	ctx = ks.Ctx.WithBlockHeight(6).WithBlockTime(time.Unix(6, 0))
	ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
		Height:    6,
		Timestamp: time.Unix(6, 0),
	})
	ks.ClobKeeper.MustSetProcessProposerMatchesEvents(ctx, types.ProcessProposerMatchesEvents{
		BlockHeight: 6,
	})

	return ks, ctx
}

// countTriggersForSuborder counts the number of trigger-store entries for the given suborder id
// without mutating the store.
func countTriggersForSuborder(ks keepertest.ClobKeepersTestContext, ctx sdk.Context, suborderId types.OrderId) int {
	store := ks.ClobKeeper.GetTWAPTriggerOrderPlacementStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	count := 0
	for ; iterator.Valid(); iterator.Next() {
		var id types.OrderId
		ks.Cdc.MustUnmarshal(iterator.Key()[8:], &id)
		if id == suborderId {
			count++
		}
	}
	return count
}

func TestGenerateAndPlaceTriggeredTwapSuborders_DedupsDuplicateTriggers(t *testing.T) {
	ks, ctx := setupTwapSweepTest(t)

	twapOrder := constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO
	ks.ClobKeeper.SetTWAPOrderPlacement(ctx, twapOrder, 6)

	suborderId := types.OrderId{
		SubaccountId: twapOrder.OrderId.SubaccountId,
		ClientId:     twapOrder.OrderId.ClientId,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrder.OrderId.ClobPairId,
	}

	// SetTWAPOrderPlacement already added one due trigger (offset 0, i.e. == block time).
	// Add a second, earlier-due trigger for the exact same suborder id.
	ks.ClobKeeper.AddSuborderToTriggerStore(ctx, suborderId, -5)

	// Sanity check: two due triggers exist for the same suborder id before the sweep.
	require.Equal(t, 2, countTriggersForSuborder(ks, ctx, suborderId))

	ks.ClobKeeper.GenerateAndPlaceTriggeredTwapSuborders(ctx)

	// Exactly one suborder was placed - no collision panic from the duplicate trigger.
	placement, found := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, suborderId)
	require.True(t, found)
	require.Equal(t, suborderId, placement.Order.OrderId)

	// The duplicate trigger key was dropped; only the next-interval trigger remains, and it is indexed.
	require.Equal(t, 1, countTriggersForSuborder(ks, ctx, suborderId))
	_, nextTriggerTime, found := ks.ClobKeeper.GetTwapTriggerPlacement(ctx, suborderId)
	require.True(t, found)
	require.Equal(t, ctx.BlockTime().Unix()+int64(twapOrder.TwapParameters.Interval), nextTriggerTime)
	indexedTime, found := ks.ClobKeeper.GetPendingTwapTrigger(ctx, suborderId)
	require.True(t, found)
	require.Equal(t, nextTriggerTime, indexedTime)

	// The parent placement survives with legs decremented by exactly one.
	twapPlacement, found := ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.True(t, found)
	require.Equal(t, twapOrder.GetTotalLegsTWAPOrder()-1, twapPlacement.RemainingLegs)
}

func TestGenerateAndPlaceTriggeredTwapSuborders_AlreadyExistsKeepsParent(t *testing.T) {
	ks, ctx := setupTwapSweepTest(t)

	twapOrder := constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO
	ks.ClobKeeper.SetTWAPOrderPlacement(ctx, twapOrder, 6)

	suborderId := types.OrderId{
		SubaccountId: twapOrder.OrderId.SubaccountId,
		ClientId:     twapOrder.OrderId.ClientId,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrder.OrderId.ClobPairId,
	}

	// Force the sweep's suborder placement to collide with ErrStatefulOrderAlreadyExists by
	// pre-placing a long-term order under the exact suborder id the sweep will generate.
	collidingOrder := twapOrder
	collidingOrder.OrderId = suborderId
	ks.ClobKeeper.SetLongTermOrderPlacement(ctx, collidingOrder, 6)
	ks.ClobKeeper.AddStatefulOrderIdExpiration(ctx, collidingOrder.MustGetUnixGoodTilBlockTime(), suborderId)

	ks.ClobKeeper.GenerateAndPlaceTriggeredTwapSuborders(ctx)

	// The parent placement survives the collision; the skipped non-final leg still counts and its
	// quantity is caught up on the next interval.
	twapPlacement, found := ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.True(t, found)
	require.Equal(t, twapOrder.GetTotalLegsTWAPOrder()-1, twapPlacement.RemainingLegs)

	// The next trigger is queued rather than the parent being torn down.
	require.Equal(t, 1, countTriggersForSuborder(ks, ctx, suborderId))
	_, nextTriggerTime, found := ks.ClobKeeper.GetTwapTriggerPlacement(ctx, suborderId)
	require.True(t, found)
	require.Equal(t, ctx.BlockTime().Unix()+int64(twapOrder.TwapParameters.Interval), nextTriggerTime)

	// The pre-existing colliding placement (which caused the collision) is untouched.
	existing, found := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, suborderId)
	require.True(t, found)
	require.Equal(t, collidingOrder, existing.Order)
}

// TestGenerateAndPlaceTriggeredTwapSuborders_FinalLegCollisionRetries covers a collision on the
// final leg: the previous suborder was removed in the same block, so the placement fails with
// ErrStatefulOrderPreviouslyRemoved. The parent keeps its last leg, retries it on the next interval,
// and only then completes.
func TestGenerateAndPlaceTriggeredTwapSuborders_FinalLegCollisionRetries(t *testing.T) {
	ks, ctx := setupTwapSweepTest(t)

	twapOrder := constants.TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO
	ks.ClobKeeper.SetTWAPOrderPlacement(ctx, twapOrder, 6)
	suborderId := types.OrderId{
		SubaccountId: twapOrder.OrderId.SubaccountId,
		ClientId:     twapOrder.OrderId.ClientId,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrder.OrderId.ClobPairId,
	}
	interval := int64(twapOrder.TwapParameters.Interval)

	// Advance the parent to its final leg.
	placement, found := ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.True(t, found)
	for placement.RemainingLegs > 1 {
		ks.ClobKeeper.DecrementTwapOrderRemainingLegs(ctx, placement)
		placement, found = ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
		require.True(t, found)
	}

	// The previous suborder was removed in this block, so placing the same suborder id collides.
	ks.ClobKeeper.MustSetProcessProposerMatchesEvents(ctx, types.ProcessProposerMatchesEvents{
		BlockHeight:             6,
		RemovedStatefulOrderIds: []types.OrderId{suborderId},
	})
	ks.ClobKeeper.GenerateAndPlaceTriggeredTwapSuborders(ctx)

	placement, found = ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.True(t, found)
	require.Equal(t, uint32(1), placement.RemainingLegs)
	_, found = ks.ClobKeeper.GetLongTermOrderPlacement(ctx, suborderId)
	require.False(t, found)
	nextTriggerTime, found := ks.ClobKeeper.GetPendingTwapTrigger(ctx, suborderId)
	require.True(t, found)
	require.Equal(t, ctx.BlockTime().Unix()+interval, nextTriggerTime)

	// Next interval, with the removal no longer pending: the final leg is placed and the parent
	// completes.
	ctx = ctx.WithBlockTime(time.Unix(nextTriggerTime, 0))
	ks.ClobKeeper.MustSetProcessProposerMatchesEvents(ctx, types.ProcessProposerMatchesEvents{
		BlockHeight: 6,
	})
	ks.ClobKeeper.GenerateAndPlaceTriggeredTwapSuborders(ctx)

	_, found = ks.ClobKeeper.GetLongTermOrderPlacement(ctx, suborderId)
	require.True(t, found)
	placement, found = ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.True(t, found)
	require.Equal(t, uint32(0), placement.RemainingLegs)

	// The following interval, the completed parent is deleted together with its pending trigger.
	ctx = ctx.WithBlockTime(time.Unix(nextTriggerTime+interval, 0))
	ks.ClobKeeper.GenerateAndPlaceTriggeredTwapSuborders(ctx)

	_, found = ks.ClobKeeper.GetTwapOrderPlacement(ctx, twapOrder.OrderId)
	require.False(t, found)
	require.Equal(t, 0, countTriggersForSuborder(ks, ctx, suborderId))
	_, found = ks.ClobKeeper.GetPendingTwapTrigger(ctx, suborderId)
	require.False(t, found)
}
