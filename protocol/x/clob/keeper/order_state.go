package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// OrderIdFillState is a struct that represents an order fill amount in state.
type OrderIdFillState struct {
	types.OrderFillState
	OrderId types.OrderId
}

// GetAllOrderFillStates iterates over the keeper store, and returns a slice of all fill amounts known to the keeper.
// This method is called during application startup as a means of hydrating the memclob with the known fill amounts
// in state.
func (k Keeper) GetAllOrderFillStates(ctx sdk.Context) (fillStates []OrderIdFillState) {
	// Retrieve an instance of the store.
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.OrderAmountFilledKeyPrefix),
	)

	// Iterate over all keys with the `OrderAmountFilledKeyPrefx`.
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// Unmarshal the value into an `OrderFillState` struct.
		var orderFillState types.OrderFillState
		k.cdc.MustUnmarshal(iterator.Value(), &orderFillState)

		// Unmarshal the key into an `OrderId` struct.
		var orderId types.OrderId
		k.cdc.MustUnmarshal(iterator.Key(), &orderId)

		// Combine both the key and value into a new struct called `OrderIdFillState` which contains all of the
		// relevant fill information.
		fillStates = append(fillStates, OrderIdFillState{
			OrderFillState: orderFillState,
			OrderId:        orderId,
		})
	}

	return fillStates
}

// SetOrderFillAmount writes the total `fillAmount` and `prunableBlockHeight` of an order to on-chain state.
// TODO(DEC-1219): Determine whether we should continue using `OrderFillState` proto for stateful orders.
func (k Keeper) SetOrderFillAmount(
	ctx sdk.Context,
	orderId types.OrderId,
	fillAmount satypes.BaseQuantums,
	prunableBlockHeight uint32,
) {
	// Define `OrderFillState` based on the provided arguments.
	var orderFillState = types.OrderFillState{
		FillAmount:          uint64(fillAmount),
		PrunableBlockHeight: prunableBlockHeight,
	}

	// Marshal `orderFillState` to bytes.
	orderFillStateBytes := k.cdc.MustMarshal(&orderFillState)

	// Retrieve an instance of the store.
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.OrderAmountFilledKeyPrefix),
	)

	// Write `orderFillStateBytes` to state.
	store.Set(
		orderId.ToStateKey(),
		orderFillStateBytes,
	)

	// Retrieve an instance of the memStore.
	memStore := prefix.NewStore(
		ctx.KVStore(k.memKey),
		[]byte(types.OrderAmountFilledKeyPrefix),
	)

	// Write `orderFillStateBytes` to memStore.
	memStore.Set(
		orderId.ToStateKey(),
		orderFillStateBytes,
	)
}

// GetOrderFillAmount returns the total `fillAmount` and `prunableBlockHeight` from the memStore.
func (k Keeper) GetOrderFillAmount(
	ctx sdk.Context,
	orderId types.OrderId,
) (
	exists bool,
	fillAmount satypes.BaseQuantums,
	prunableBlockHeight uint32,
) {
	memStore := ctx.KVStore(k.memKey)

	// Retrieve an instance of the memStore.
	memPrefixStore := prefix.NewStore(
		memStore,
		[]byte(types.OrderAmountFilledKeyPrefix),
	)

	// Retrieve the `OrderFillState` bytes from the store.
	orderFillStateBytes := memPrefixStore.Get(
		orderId.ToStateKey(),
	)

	// If the `OrderFillState` does not exist, early return.
	if orderFillStateBytes == nil {
		return false, 0, 0
	}

	// Unmarshal the `orderFillStateBytes` into a struct, and return relevant values.
	var orderFillState types.OrderFillState
	k.cdc.MustUnmarshal(orderFillStateBytes, &orderFillState)

	return true, satypes.BaseQuantums(orderFillState.FillAmount), orderFillState.PrunableBlockHeight
}

// AddOrdersForPruning creates or updates a slice of `orderIds` to state for potential future pruning from state.
// These orders will be checked for pruning from state at `prunableBlockHeight`. If the `orderIds` slice provided
// contains duplicates, the duplicates will be ignored.
func (k Keeper) AddOrdersForPruning(ctx sdk.Context, orderIds []types.OrderId, prunableBlockHeight uint32) {
	// Retrieve an instance of the store.
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.BlockHeightToPotentiallyPrunableOrdersPrefix),
	)

	// Retrieve the `PotentiallyPrunableOrders` bytes from the store.
	potentiallyPrunableOrdersBytes := store.Get(
		lib.Uint32ToKey(prunableBlockHeight),
	)

	var potentiallyPrunableOrdersSet = make(map[types.OrderId]bool)
	var potentiallyPrunableOrders = types.PotentiallyPrunableOrders{}
	var potentiallyPrunableOrderIds = make([]types.OrderId, len(orderIds))

	// Initialize `potentiallyPrunableOrderIds` with the provided `orderIds`.
	// Copy to avoid mutating the provided `orderIds`.
	copy(potentiallyPrunableOrderIds, orderIds)

	// If the state already contains `potentiallyPrunableOrders` for this `prunableBlockHeight`, add them to the list of
	// `potentiallyPrunableOrderIds`.
	if potentiallyPrunableOrdersBytes != nil {
		k.cdc.MustUnmarshal(potentiallyPrunableOrdersBytes, &potentiallyPrunableOrders)
		potentiallyPrunableOrderIds = append(potentiallyPrunableOrders.OrderIds, potentiallyPrunableOrderIds...)
	}

	// Iterate over all `potentiallyPrunableOrderIds` and place them in the set in order to dedupe them.
	for _, orderId := range potentiallyPrunableOrderIds {
		potentiallyPrunableOrdersSet[orderId] = true
	}

	// Iterate over the set and build a list of `dedupedOrderIds`.
	var dedupedOrderIds = make([]types.OrderId, 0, len(potentiallyPrunableOrdersSet))
	for orderId := range potentiallyPrunableOrdersSet {
		dedupedOrderIds = append(dedupedOrderIds, orderId)
	}

	// Sort the orderIds so that the state write is deterministic.
	types.MustSortAndHaveNoDuplicates(dedupedOrderIds)

	// Set the new `dedupedOrderIds` on the `potentiallyPrunableOrders`.
	potentiallyPrunableOrders.OrderIds = dedupedOrderIds

	// Marshal `prunableOrders` back to bytes.
	potentiallyPrunableOrdersBytes = k.cdc.MustMarshal(&potentiallyPrunableOrders)

	// Write `prunableOrders` to state for the appropriate block height.
	store.Set(
		lib.Uint32ToKey(prunableBlockHeight),
		potentiallyPrunableOrdersBytes,
	)
}

// PruneOrdersForBlockHeight checks all orders for prunability given the provided `blockHeight`.
// If an order is deemed prunable at this `blockHeight`, then it is pruned.
// Note: An order is only deemed prunable if the `prunableBlockHeight` on the `OrderFillState` is less than or equal
// to the provided `blockHeight` passed this method. Returns a slice of unique `OrderIds` which were pruned from state.
func (k Keeper) PruneOrdersForBlockHeight(ctx sdk.Context, blockHeight uint32) (prunedOrderIds []types.OrderId) {
	// Retrieve an instance of the stores.
	blockHeightToPotentiallyPrunableOrdersStore := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.BlockHeightToPotentiallyPrunableOrdersPrefix),
	)

	// Retrieve the raw bytes of the `prunableOrders`.
	potentiallyPrunableOrderBytes := blockHeightToPotentiallyPrunableOrdersStore.Get(
		lib.Uint32ToKey(blockHeight),
	)

	// If there are no prunable orders for this block, then there is nothing to do. Early return.
	if potentiallyPrunableOrderBytes == nil {
		return
	}

	var potentiallyPrunableOrders types.PotentiallyPrunableOrders
	k.cdc.MustUnmarshal(potentiallyPrunableOrderBytes, &potentiallyPrunableOrders)

	for _, orderId := range potentiallyPrunableOrders.OrderIds {
		// Check if the order can be pruned, and prune if so.
		exists, _, prunableBlockHeight := k.GetOrderFillAmount(ctx, orderId)
		if exists && prunableBlockHeight <= blockHeight {
			k.RemoveOrderFillAmount(ctx, orderId)
			prunedOrderIds = append(prunedOrderIds, orderId)

			if prunableBlockHeight < blockHeight {
				k.Logger(ctx).Error(fmt.Sprintf(
					"prunableBlockHeight %v is less than blockHeight %v in PruneOrdersForBlockHeight, this should never happen.",
					prunableBlockHeight,
					blockHeight,
				))
			}
		}
	}

	// Delete the key for prunable orders at this block height.
	blockHeightToPotentiallyPrunableOrdersStore.Delete(
		lib.Uint32ToKey(blockHeight),
	)

	return prunedOrderIds
}

// RemoveOrderFillAmount removes the fill amount of an Order from state and the memstore.
// This function is a no-op if no order fill amount exists in state and the mem store with `orderId`.
func (k Keeper) RemoveOrderFillAmount(ctx sdk.Context, orderId types.OrderId) {
	// Delete the fill amount from the state store.
	orderAmountFilledStore := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.OrderAmountFilledKeyPrefix),
	)

	orderAmountFilledStore.Delete(orderId.ToStateKey())

	// Delete the fill amount from the mem store.
	memStore := prefix.NewStore(
		ctx.KVStore(k.memKey),
		[]byte(types.OrderAmountFilledKeyPrefix),
	)
	memStore.Delete(orderId.ToStateKey())
}

// PruneStateFillAmountsForShortTermOrders prunes Short-Term order fill amounts from state that are pruneable
// at the block height of the most recently committed block.
func (k Keeper) PruneStateFillAmountsForShortTermOrders(
	ctx sdk.Context,
) {
	blockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())

	// Prune all fill amounts from state which have a pruneable block height of the current `blockHeight`.
	k.PruneOrdersForBlockHeight(ctx, blockHeight)
}
