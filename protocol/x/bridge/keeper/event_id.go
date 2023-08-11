package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
)

const (
	nextAcknowledgedEventIdKey = "NextAcknowledgedEventId"
)

// GetNextAcknowledgedEventId returns the `nextAcknowledgedEventIdKey` from state.
func (k Keeper) GetNextAcknowledgedEventId(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	var rawBytes []byte = store.Get([]byte(nextAcknowledgedEventIdKey))
	return lib.BytesToUint32(rawBytes)
}

// SetNextAcknowledgedEventId sets the `nextAcknowledgedEventIdKey` in state.
func (k Keeper) SetNextAcknowledgedEventId(
	ctx sdk.Context,
	id uint32,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(nextAcknowledgedEventIdKey), lib.Uint32ToBytes(id))
}
