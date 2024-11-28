package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
)

func (k Keeper) GetSynchronyParams(ctx sdk.Context) types.SynchronyParams {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte(types.SynchronyParamsKey))

	if bytes == nil {
		return types.DefaultSynchronyParams()
	}

	var params types.SynchronyParams
	k.cdc.MustUnmarshal(bytes, &params)
	return params
}

func (k Keeper) SetSynchronyParams(ctx sdk.Context, params types.SynchronyParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.SynchronyParamsKey), k.cdc.MustMarshal(&params))
}

func (k Keeper) GetBlockDelay(ctx sdk.Context) time.Duration {
	return k.GetSynchronyParams(ctx).NextBlockDelay
}
