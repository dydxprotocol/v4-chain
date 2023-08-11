package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/feetiers/types"
)

const (
	paramsKey = "PerpetualFeeParams"
)

// GetPerpetualFeeParams returns the PerpetualFeeParams in state.
func (k Keeper) GetPerpetualFeeParams(
	ctx sdk.Context,
) (
	params types.PerpetualFeeParams,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(paramsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// SetPerpetualFeeParams updates the PerpetualFeeParams in state.
// Returns an error iff validation fails.
func (k Keeper) SetPerpetualFeeParams(
	ctx sdk.Context,
	params types.PerpetualFeeParams,
) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(paramsKey), b)

	return nil
}
