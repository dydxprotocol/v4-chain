package keeper_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
		timestamp := binary.BigEndian.Uint64(iterator.Key()[0:8])
		require.Equal(t, ks.Ctx.BlockTime().Unix()+expectedOrder[index], int64(timestamp))
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
