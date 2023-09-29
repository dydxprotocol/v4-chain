package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
)

// GetParams returns the Params in state.
func (k Keeper) GetParams(
	ctx sdk.Context,
) (
	params types.Params,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.ParamsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
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
	store.Set([]byte(types.ParamsKey), b)

	return nil
}
