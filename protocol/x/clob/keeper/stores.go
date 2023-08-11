package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// getStatefulOrderPlacementStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order placement from state.
func (k Keeper) getStatefulOrderPlacementStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.KeyPrefix(types.StatefulOrderPlacementKeyPrefix),
	)
}

// getStatefulOrderPlacementMemStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order placement from state.
func (k Keeper) getStatefulOrderPlacementMemStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.memKey),
		types.KeyPrefix(types.StatefulOrderPlacementKeyPrefix),
	)
}

// getStatefulOrdersTimeSliceStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order time slice from state.
func (k Keeper) getStatefulOrdersTimeSliceStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.KeyPrefix(types.StatefulOrdersTimeSlicePrefix),
	)
}

// getTransientStore fetches a transient store used for reading and
// updating the transient store.
func (k Keeper) getTransientStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.transientStoreKey)
}
