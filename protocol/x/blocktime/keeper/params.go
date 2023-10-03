package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
)

// GetParams returns the DowntimeParams in state.
func (k Keeper) GetDowntimeParams(
	ctx sdk.Context,
) (
	params types.DowntimeParams,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.DowntimeParamsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// SetParams updates the Params in state.
// Returns an error iff validation fails.
func (k Keeper) SetDowntimeParams(
	ctx sdk.Context,
	params types.DowntimeParams,
) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.DowntimeParamsKey), b)

	// For each new duration, we assume the worst case. For new durations that are smaller than all existing
	// durations, we'll use the current block's info. Note that at genesis, this is true for all durations.
	newAllDowntimeInfo := types.AllDowntimeInfo{}
	for _, duration := range params.Durations {
		newAllDowntimeInfo.Infos = append(newAllDowntimeInfo.Infos, &types.AllDowntimeInfo_DowntimeInfo{
			Duration: duration,
			BlockInfo: types.BlockInfo{
				Height:    uint32(ctx.BlockHeight()),
				Timestamp: ctx.BlockTime(),
			},
		})
	}

	// Assuming the worst case means assuming that each previously recorded downtime lasted as long as possible.
	// So for each new duration, we take the downtime of the largest existing duration that is smaller.
	allDowntimeInfo := k.GetAllDowntimeInfo(ctx)
	for _, info := range newAllDowntimeInfo.Infos {
		for _, oldInfo := range allDowntimeInfo.Infos {
			if info.Duration >= oldInfo.Duration {
				info.BlockInfo = oldInfo.BlockInfo
			} else {
				break
			}
		}
	}
	k.SetAllDowntimeInfo(ctx, &newAllDowntimeInfo)
	return nil
}
