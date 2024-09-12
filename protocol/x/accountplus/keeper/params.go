package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.ParamsKeyPrefix))
	if bz == nil {
		return types.Params{
			IsSmartAccountActive: false,
		}
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.ParamsKeyPrefix), bz)
}

// GetIsSmartAccountActive returns the value of the isSmartAccountActive parameter.
// If the value has not been set, it will return false.
func (k *Keeper) GetIsSmartAccountActive(ctx sdk.Context) bool {
	return k.GetParams(ctx).IsSmartAccountActive
}

// SetActiveState sets the active state of the smart account module.
func (k Keeper) SetActiveState(ctx sdk.Context, active bool) {
	params := k.GetParams(ctx)
	params.IsSmartAccountActive = active
	k.SetParams(ctx, params)
}
