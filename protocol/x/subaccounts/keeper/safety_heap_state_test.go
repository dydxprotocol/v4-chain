package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetSubaccountAtIndex(t *testing.T) {
	// Setup keeper state and test parameters.
	ctx, subaccountsKeeper, _, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, false)

	// Write a couple of subaccounts to the store.
	store := subaccountsKeeper.GetSafetyHeapStore(ctx, 0, types.Long)
	subaccountsKeeper.SetSubaccountAtIndex(store, constants.Alice_Num0, 0)
	subaccountsKeeper.SetSubaccountAtIndex(store, constants.Bob_Num0, 1)
	subaccountsKeeper.SetSubaccountAtIndex(store, constants.Carl_Num0, 2)

	// Get.
	require.Equal(t, constants.Alice_Num0, subaccountsKeeper.MustGetSubaccountAtIndex(store, 0))
	require.Equal(t, constants.Bob_Num0, subaccountsKeeper.MustGetSubaccountAtIndex(store, 1))
	require.Equal(t, constants.Carl_Num0, subaccountsKeeper.MustGetSubaccountAtIndex(store, 2))

	// Overwrite the subaccount at index 1.
	subaccountsKeeper.SetSubaccountAtIndex(store, constants.Alice_Num1, 0)
	require.Equal(t, constants.Alice_Num1, subaccountsKeeper.MustGetSubaccountAtIndex(store, 0))
	subaccountsKeeper.SetSubaccountAtIndex(store, constants.Dave_Num0, 1)
	require.Equal(t, constants.Dave_Num0, subaccountsKeeper.MustGetSubaccountAtIndex(store, 1))

	// Delete
	subaccountsKeeper.DeleteSubaccountAtIndex(store, 2)
	subaccountsKeeper.DeleteSubaccountAtIndex(store, 1)

	// Getting non existent subaccount.
	_, found := subaccountsKeeper.GetSubaccountAtIndex(store, 1)
	require.False(t, found)

	require.PanicsWithError(
		t,
		types.ErrSafetyHeapSubaccountNotFoundAtIndex.Error(),
		func() {
			subaccountsKeeper.MustGetSubaccountAtIndex(store, 1)
		},
	)

	_, found = subaccountsKeeper.GetSubaccountAtIndex(store, 2)
	require.False(t, found)

	require.PanicsWithError(
		t,
		types.ErrSafetyHeapSubaccountNotFoundAtIndex.Error(),
		func() {
			subaccountsKeeper.MustGetSubaccountAtIndex(store, 2)
		},
	)
}

func TestGetSetSubaccountHeapIndex(t *testing.T) {
	// Setup keeper state and test parameters.
	ctx, subaccountsKeeper, _, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, false)

	// Write a couple of subaccounts to the store.
	store := subaccountsKeeper.GetSafetyHeapStore(ctx, 0, types.Long)
	subaccountsKeeper.SetSubaccountHeapIndex(store, constants.Alice_Num0, 0)
	subaccountsKeeper.SetSubaccountHeapIndex(store, constants.Bob_Num0, 1)
	subaccountsKeeper.SetSubaccountHeapIndex(store, constants.Carl_Num0, 2)

	// Get.
	require.Equal(t, uint32(0), subaccountsKeeper.MustGetSubaccountHeapIndex(store, constants.Alice_Num0))
	require.Equal(t, uint32(1), subaccountsKeeper.MustGetSubaccountHeapIndex(store, constants.Bob_Num0))
	require.Equal(t, uint32(2), subaccountsKeeper.MustGetSubaccountHeapIndex(store, constants.Carl_Num0))

	// Overwrite the subaccount at index 1.
	subaccountsKeeper.SetSubaccountHeapIndex(store, constants.Alice_Num0, 3)
	require.Equal(t, uint32(3), subaccountsKeeper.MustGetSubaccountHeapIndex(store, constants.Alice_Num0))
	subaccountsKeeper.SetSubaccountHeapIndex(store, constants.Carl_Num0, 4)
	require.Equal(t, uint32(4), subaccountsKeeper.MustGetSubaccountHeapIndex(store, constants.Carl_Num0))

	// Delete
	subaccountsKeeper.DeleteSubaccountHeapIndex(store, constants.Alice_Num0)
	subaccountsKeeper.DeleteSubaccountHeapIndex(store, constants.Carl_Num0)

	// Getting non existent subaccount.
	_, found := subaccountsKeeper.GetSubaccountHeapIndex(store, constants.Alice_Num0)
	require.False(t, found)

	require.PanicsWithError(
		t,
		types.ErrSafetyHeapSubaccountIndexNotFound.Error(),
		func() {
			subaccountsKeeper.MustGetSubaccountHeapIndex(store, constants.Alice_Num0)
		},
	)

	_, found = subaccountsKeeper.GetSubaccountHeapIndex(store, constants.Carl_Num0)
	require.False(t, found)

	require.PanicsWithError(
		t,
		types.ErrSafetyHeapSubaccountIndexNotFound.Error(),
		func() {
			subaccountsKeeper.MustGetSubaccountHeapIndex(store, constants.Carl_Num0)
		},
	)
}

func TestGetSetSafetyHeapLength(t *testing.T) {
	// Setup keeper state and test parameters.
	ctx, subaccountsKeeper, _, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, false)

	// Write a couple of subaccounts to the store.
	store := subaccountsKeeper.GetSafetyHeapStore(ctx, 0, types.Long)

	require.Equal(t, uint32(0), subaccountsKeeper.GetSafetyHeapLength(store))

	subaccountsKeeper.SetSafetyHeapLength(store, 7)

	require.Equal(t, uint32(7), subaccountsKeeper.GetSafetyHeapLength(store))
}
