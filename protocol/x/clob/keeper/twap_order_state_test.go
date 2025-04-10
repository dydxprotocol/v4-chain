package keeper_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTWAPOrderTriggerStoreOrderingSameOrderId(t *testing.T) {
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
