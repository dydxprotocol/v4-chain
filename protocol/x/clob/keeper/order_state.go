package keeper

import (
	"bytes"
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
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
}

// GetOrderFillAmount returns the total `fillAmount` and `prunableBlockHeight` from state.
func (k Keeper) GetOrderFillAmount(
	ctx sdk.Context,
	orderId types.OrderId,
) (
	exists bool,
	fillAmount satypes.BaseQuantums,
	prunableBlockHeight uint32,
) {
	store := ctx.KVStore(k.storeKey)

	prefixStore := prefix.NewStore(
		store,
		[]byte(types.OrderAmountFilledKeyPrefix),
	)

	// Retrieve the `OrderFillState` bytes from the store.
	orderFillStateBytes := prefixStore.Get(
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

// GetPruneableOrdersStore gets a prefix store for pruneable orders at a given height.
// The full format for these keys is <PrunableOrdersKeyPrefix><height>:<order_id>.
func (k Keeper) GetPruneableOrdersStore(ctx sdk.Context, height uint32) prefix.Store {
	var buf bytes.Buffer
	buf.Write([]byte(types.PrunableOrdersKeyPrefix))
	buf.Write(lib.Uint32ToKey(height))
	buf.Write([]byte(":"))
	return prefix.NewStore(ctx.KVStore(k.storeKey), buf.Bytes())
}

// AddOrdersForPruning creates or updates `orderIds` to state for potential future pruning from state.
func (k Keeper) AddOrdersForPruning(ctx sdk.Context, orderIds []types.OrderId, prunableBlockHeight uint32) {
	store := k.GetPruneableOrdersStore(ctx, prunableBlockHeight)
	for _, orderId := range orderIds {
		store.Set(
			orderId.ToStateKey(),
			k.cdc.MustMarshal(&orderId),
		)
	}
}

// Deprecated: Do not use. Retained for testing purposes.
// LegacyAddOrdersForPruning is the old key-per-height format of storing orders to prune.
func (k Keeper) LegacyAddOrdersForPruning(ctx sdk.Context, orderIds []types.OrderId, prunableBlockHeight uint32) {
	// Retrieve an instance of the store.
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.LegacyBlockHeightToPotentiallyPrunableOrdersPrefix),
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
	potentiallyPrunableOrdersStore := k.GetPruneableOrdersStore(ctx, blockHeight)
	it := potentiallyPrunableOrdersStore.Iterator(nil, nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var orderId types.OrderId
		k.cdc.MustUnmarshal(it.Value(), &orderId)
		exists, _, prunableBlockHeight := k.GetOrderFillAmount(ctx, orderId)
		if exists && prunableBlockHeight <= blockHeight {
			k.RemoveOrderFillAmount(ctx, orderId)
			prunedOrderIds = append(prunedOrderIds, orderId)

			if prunableBlockHeight < blockHeight {
				log.ErrorLog(ctx,
					"prunableBlockHeight is less than blockHeight in PruneOrdersForBlockHeight, this should never happen.",
					log.PrunableBlockHeight,
					prunableBlockHeight,
				)
			}
		}
		potentiallyPrunableOrdersStore.Delete(it.Key())
	}

	return prunedOrderIds
}

// MigratePruneableOrders is used to migrate prunable orders from key-per-height to key-per-order format.
func (k Keeper) MigratePruneableOrders(ctx sdk.Context) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.LegacyBlockHeightToPotentiallyPrunableOrdersPrefix), // nolint:staticcheck
	)
	it := store.Iterator(nil, nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		if it.Value() == nil {
			continue
		}

		height := binary.BigEndian.Uint32(it.Key())
		var potentiallyPrunableOrders types.PotentiallyPrunableOrders
		k.cdc.MustUnmarshal(it.Value(), &potentiallyPrunableOrders)
		k.AddOrdersForPruning(ctx, potentiallyPrunableOrders.OrderIds, height)
		store.Delete(it.Key())
	}
}

// RemoveOrderFillAmount removes the fill amount of an Order from state.
// This function is a no-op if no order fill amount exists in state with `orderId`.
func (k Keeper) RemoveOrderFillAmount(ctx sdk.Context, orderId types.OrderId) {
	// Delete the fill amount from the state store.
	orderAmountFilledStore := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.OrderAmountFilledKeyPrefix),
	)

	orderAmountFilledStore.Delete(orderId.ToStateKey())

	// If grpc stream is on, zero out the fill amount.
	if k.GetFullNodeStreamingManager().Enabled() {
		allUpdates := types.NewOffchainUpdates()
		if message, success := off_chain_updates.CreateOrderUpdateMessage(
			ctx,
			orderId,
			0, // Total filled quantums is zero because it's been pruned from state.
		); success {
			allUpdates.AddUpdateMessage(orderId, message)
		}
		k.SendOrderbookUpdates(ctx, allUpdates)
	}
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
