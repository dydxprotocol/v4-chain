package keeper

import (
	"encoding/binary"
	"math"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetProcessProposerMatchesEvents gets the process proposer matches events from the latest block.
func (k Keeper) GetProcessProposerMatchesEvents(ctx sdk.Context) types.ProcessProposerMatchesEvents {
	// Retrieve an instance of the memory store.
	memStore := ctx.KVStore(k.memKey)

	// Retrieve the `processProposerMatchesEvents` bytes from the store.
	processProposerMatchesEventsBytes := memStore.Get(
		[]byte(types.ProcessProposerMatchesEventsKey),
	)

	// Unmarshal the `processProposerMatchesEvents` into a struct and return it.
	var processProposerMatchesEvents types.ProcessProposerMatchesEvents
	k.cdc.MustUnmarshal(processProposerMatchesEventsBytes, &processProposerMatchesEvents)
	return processProposerMatchesEvents
}

// MustSetProcessProposerMatchesEvents sets the process proposer matches events from the latest block.
// This function panics if:
//   - the current block height does not match the block height of the ProcessProposerMatchesEvents
//   - called outside of deliver TX mode
//   - Any of the ProcessProposerMatchesEvents fields have duplicates.
//
// TODO(DEC-1281): add parameter validation.
func (k Keeper) MustSetProcessProposerMatchesEvents(
	ctx sdk.Context,
	processProposerMatchesEvents types.ProcessProposerMatchesEvents,
) {
	lib.AssertDeliverTxMode(ctx)

	if err := processProposerMatchesEvents.ValidateProcessProposerMatchesEvents(ctx); err != nil {
		panic(err)
	}

	// Retrieve an instance of the memory store.
	memStore := ctx.KVStore(k.memKey)

	// Write `processProposerMatchesEvents` to the `memStore`.
	memStore.Set(
		[]byte(types.ProcessProposerMatchesEventsKey),
		k.cdc.MustMarshal(&processProposerMatchesEvents),
	)
}

// InitializeProcessProposerMatchesEvents initializes the process proposer matches events.
// This function should only be called from the CLOB genesis.
func (k Keeper) InitializeProcessProposerMatchesEvents(
	ctx sdk.Context,
) {
	processProposerMatchesEvents := types.ProcessProposerMatchesEvents{
		BlockHeight: 1,
	}

	memStore := ctx.KVStore(k.memKey)
	memStore.Set(
		[]byte(types.ProcessProposerMatchesEventsKey),
		k.cdc.MustMarshal(&processProposerMatchesEvents),
	)
}

// ResetAllDeliveredOrderIds resets the lists of delivered order ids. This should be reset every block.
func (k Keeper) ResetAllDeliveredOrderIds(ctx sdk.Context) {
	k.ResetUnorderedOrderIds(ctx, types.DeliveredCancelKeyPrefix)
	k.ResetOrderedOrderIds(
		ctx,
		types.OrderedDeliveredLongTermOrderKeyPrefix,
		types.OrderedDeliveredLongTermOrderIndexKey,
	)
	k.ResetOrderedOrderIds(
		ctx,
		types.OrderedDeliveredConditionalOrderKeyPrefix,
		types.OrderedDeliveredConditionalOrderIndexKey,
	)
}

// AddDeliveredLongTermOrderId saves a long term order id to the memstore for processing
// in the next PrepareCheckState. The order of additions is maintained.
func (k Keeper) AddDeliveredLongTermOrderId(ctx sdk.Context, orderId types.OrderId) {
	k.AppendOrderedOrderId(
		ctx,
		types.OrderedDeliveredLongTermOrderKeyPrefix,
		types.OrderedDeliveredLongTermOrderIndexKey,
		orderId,
	)
}

// AddDeliveredConditionalOrderId saves a conditional order id to the memstore for processing
// in the next PrepareCheckState. The order of additions is maintained.
func (k Keeper) AddDeliveredConditionalOrderId(ctx sdk.Context, orderId types.OrderId) {
	k.AppendOrderedOrderId(
		ctx,
		types.OrderedDeliveredConditionalOrderKeyPrefix,
		types.OrderedDeliveredConditionalOrderIndexKey,
		orderId,
	)
}

// AddDeliveredCancelledOrderId saves a cancelled order id to the memstore for processing
// in the next PrepareCheckState. The order of additions is not maintained.
func (k Keeper) AddDeliveredCancelledOrderId(ctx sdk.Context, orderId types.OrderId) {
	k.SetUnorderedOrderId(ctx, types.DeliveredCancelKeyPrefix, orderId)
}

// GetDeliveredLongTermOrderIds gets the ordered list of delivered long term order ids from the memstore.
func (k Keeper) GetDeliveredLongTermOrderIds(ctx sdk.Context) []types.OrderId {
	return k.GetOrderIds(ctx, types.OrderedDeliveredLongTermOrderKeyPrefix)
}

// GetDeliveredConditionalOrderIds gets the ordered list of delivered conditional order ids from the memstore.
func (k Keeper) GetDeliveredConditionalOrderIds(ctx sdk.Context) []types.OrderId {
	return k.GetOrderIds(ctx, types.OrderedDeliveredConditionalOrderKeyPrefix)
}

// GetDeliveredCancelledOrderIds gets the unordered list of delivered cancelled order ids from the memstore.
func (k Keeper) GetDeliveredCancelledOrderIds(ctx sdk.Context) []types.OrderId {
	return k.GetOrderIds(ctx, types.DeliveredCancelKeyPrefix)
}

// HasDeliveredCancelledOrderId returns the existence of an order id in the memstore's list of delivered
// cancelled order ids.
func (k Keeper) HasDeliveredCancelledOrderId(ctx sdk.Context, orderId types.OrderId) bool {
	return k.HasUnorderedOrderId(ctx, types.DeliveredCancelKeyPrefix, orderId)
}

func (k Keeper) ResetOrderedOrderIds(ctx sdk.Context, keyPrefix string, indexKey string) {
	memStore := ctx.KVStore(k.memKey)
	store := prefix.NewStore(memStore, []byte(keyPrefix))
	it := store.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		store.Delete(it.Key())
	}
	memStore.Delete([]byte(indexKey))
}

func (k Keeper) ResetUnorderedOrderIds(ctx sdk.Context, keyPrefix string) {
	memStore := ctx.KVStore(k.memKey)
	store := prefix.NewStore(memStore, []byte(keyPrefix))
	it := store.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		store.Delete(it.Key())
	}
}

func (k Keeper) AppendOrderedOrderId(
	ctx sdk.Context,
	keyPrefix string,
	indexKey string,
	orderId types.OrderId,
) {
	memStore := ctx.KVStore(k.memKey)
	index := uint32(0)
	if bytes := memStore.Get([]byte(indexKey)); bytes != nil {
		index = binary.BigEndian.Uint32(bytes)
	}
	store := prefix.NewStore(memStore, []byte(keyPrefix))
	store.Set(lib.Uint32ToKey(index), k.cdc.MustMarshal(&orderId))

	if index == math.MaxUint32 {
		panic("store key index overflow")
	}
	memStore.Set([]byte(indexKey), lib.Uint32ToKey(index+1))
}

func (k Keeper) GetOrderIds(
	ctx sdk.Context,
	keyPrefix string,
) []types.OrderId {
	ret := []types.OrderId{}
	memStore := ctx.KVStore(k.memKey)
	store := prefix.NewStore(memStore, []byte(keyPrefix))
	it := store.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		var orderId types.OrderId
		k.cdc.MustUnmarshal(it.Value(), &orderId)
		ret = append(ret, orderId)
	}
	return ret
}

func (k Keeper) HasUnorderedOrderId(
	ctx sdk.Context,
	keyPrefix string,
	orderId types.OrderId,
) bool {
	memStore := ctx.KVStore(k.memKey)
	store := prefix.NewStore(memStore, []byte(keyPrefix))
	return store.Has(orderId.ToStateKey())
}

func (k Keeper) SetUnorderedOrderId(
	ctx sdk.Context,
	keyPrefix string,
	orderId types.OrderId,
) {
	memStore := ctx.KVStore(k.memKey)
	store := prefix.NewStore(memStore, []byte(keyPrefix))
	store.Set(orderId.ToStateKey(), k.cdc.MustMarshal(&orderId))
}
