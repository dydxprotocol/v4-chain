package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

const (
	paramsKey = "Params"
)

// GetParams returns the Params in state.
func (k Keeper) GetParams(
	ctx sdk.Context,
) (
	params types.Params,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(paramsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

func (k Keeper) GetWindowDuration(ctx sdk.Context) time.Duration {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(paramsKey))
	var params types.Params
	k.cdc.MustUnmarshal(b, &params)
	return params.WindowDuration
}

// SetParams updates the Params in state.
// Returns an error iff validation fails.
func (k Keeper) SetParams(
	ctx sdk.Context,
	params types.Params,
) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(paramsKey), b)

	return nil
}
