package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// RemoveSubaccountFromSafetyHeap removes a subaccount from the safety heap
// given a peretual and side.
func (k Keeper) RemoveSubaccountFromSafetyHeap(
	ctx sdk.Context,
	subaccountId types.SubaccountId,
	perpetualId uint32,
	side types.SafetyHeapPositionSide,
) {
	store := k.GetSafetyHeapStore(ctx, perpetualId, side)
	index := k.MustGetSubaccountHeapIndex(store, subaccountId)
	k.MustRemoveElementAtIndex(ctx, store, index)
}

// AddSubaccountToSafetyHeap adds a subaccount to the safety heap
// given a perpetual and side.
func (k Keeper) AddSubaccountToSafetyHeap(
	ctx sdk.Context,
	subaccountId types.SubaccountId,
	perpetualId uint32,
	side types.SafetyHeapPositionSide,
) {
	store := k.GetSafetyHeapStore(ctx, perpetualId, side)
	k.Insert(ctx, store, subaccountId)
}

// Heap methods

// Insert inserts a subaccount into the safety heap.
func (k Keeper) Insert(
	ctx sdk.Context,
	store prefix.Store,
	subaccountId types.SubaccountId,
) {
	// Add the subaccount to the end of the heap.
	length := k.GetSafetyHeapLength(store)
	k.SetSubaccountAtIndex(store, subaccountId, length)

	// Increment the size of the heap.
	k.SetSafetyHeapLength(store, length+1)

	// Heapify up the element at the end of the heap
	// to restore the heap property.
	k.HeapifyUp(ctx, store, length)
}

// MustRemoveElementAtIndex removes the element at the given index
// from the safety heap.
func (k Keeper) MustRemoveElementAtIndex(
	ctx sdk.Context,
	store prefix.Store,
	index uint32,
) {
	length := k.GetSafetyHeapLength(store)
	if index >= length {
		panic(types.ErrSafetyHeapSubaccountNotFoundAtIndex)
	}

	// Swap the element with the last element.
	k.Swap(store, index, length-1)

	// Remove the last element.
	k.DeleteSubaccountAtIndex(store, length-1)
	k.SetSafetyHeapLength(store, length-1)

	// Heapify down and up the element at the given index
	// to restore the heap property.
	if index < length-1 {
		k.HeapifyDown(ctx, store, index)
		k.HeapifyUp(ctx, store, index)
	}
}

// HeapifyUp moves the element at the given index up the heap
// until the heap property is restored.
func (k Keeper) HeapifyUp(
	ctx sdk.Context,
	store prefix.Store,
	index uint32,
) {
	if index == 0 {
		return
	}

	parentIndex := (index - 1) / 2
	if k.Less(ctx, store, index, parentIndex) {
		k.Swap(store, index, parentIndex)
		k.HeapifyUp(ctx, store, parentIndex)
	}
}

// HeapifyDown moves the element at the given index down the heap
// until the heap property is restored.
func (k Keeper) HeapifyDown(
	ctx sdk.Context,
	store prefix.Store,
	index uint32,
) {
	leftIndex, rightIndex := 2*index+1, 2*index+2

	length := k.GetSafetyHeapLength(store)
	if rightIndex < length && k.Less(ctx, store, rightIndex, leftIndex) {
		// Compare the current node with the right child
		// if right child exists and is less than the left child.
		if k.Less(ctx, store, rightIndex, index) {
			k.Swap(store, index, rightIndex)
			k.HeapifyDown(ctx, store, rightIndex)
		}
	} else if leftIndex < length {
		// Compare the current node with the left child
		// if left child exists.
		if k.Less(ctx, store, leftIndex, index) {
			k.Swap(store, index, leftIndex)
			k.HeapifyDown(ctx, store, leftIndex)
		}
	}
}

// Swap swaps the elements at the given indices.
func (k Keeper) Swap(
	store prefix.Store,
	index1 uint32,
	index2 uint32,
) {
	// No-op case
	if index1 == index2 {
		return
	}

	first := k.MustGetSubaccountAtIndex(store, index1)
	second := k.MustGetSubaccountAtIndex(store, index2)
	k.SetSubaccountAtIndex(store, first, index2)
	k.SetSubaccountAtIndex(store, second, index1)
}

// Less returns true if the element at the first index is less than
// the element at the second index.
func (k Keeper) Less(
	ctx sdk.Context,
	store prefix.Store,
	first uint32,
	second uint32,
) bool {
	firstSubaccountId := k.MustGetSubaccountAtIndex(store, first)
	secondSubaccountId := k.MustGetSubaccountAtIndex(store, second)

	firstRisk, err := k.GetNetCollateralAndMarginRequirements(
		ctx,
		types.Update{
			SubaccountId: firstSubaccountId,
		},
	)
	if err != nil {
		panic(err)
	}

	secondRisk, err := k.GetNetCollateralAndMarginRequirements(
		ctx,
		types.Update{
			SubaccountId: secondSubaccountId,
		},
	)
	if err != nil {
		panic(err)
	}

	// Compare the risks of the two subaccounts and sort
	// them in descending order.
	return firstRisk.Cmp(secondRisk) > 0
}
