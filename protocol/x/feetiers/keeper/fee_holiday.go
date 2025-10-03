package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// GetFeeHolidayParams returns fee holiday configuration for a given clob pair id
func (k Keeper) GetFeeHolidayParams(
	ctx sdk.Context,
	clobPairId uint32,
) (params types.FeeHolidayParams, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.FeeHolidayPrefix))
	b := store.Get(lib.Uint32ToKey(clobPairId))

	if b == nil {
		return params, types.ErrFeeHolidayNotFound
	}

	k.cdc.MustUnmarshal(b, &params)
	return params, nil
}

// SetFeeHolidayParams stores fee holiday configuration
func (k Keeper) SetFeeHolidayParams(
	ctx sdk.Context,
	feeHoliday types.FeeHolidayParams,
) error {
	// Validate the params
	err := feeHoliday.Validate(ctx.BlockTime())
	if err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.FeeHolidayPrefix))
	key := lib.Uint32ToKey(feeHoliday.ClobPairId)
	value := k.cdc.MustMarshal(&feeHoliday)
	store.Set(key, value)
	return nil
}

// GetAllFeeHolidayParams returns all configured fee holidays
func (k Keeper) GetAllFeeHolidayParams(
	ctx sdk.Context,
) []types.FeeHolidayParams {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.FeeHolidayPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	feeHolidays := []types.FeeHolidayParams{}
	for ; iterator.Valid(); iterator.Next() {
		var feeHoliday types.FeeHolidayParams
		k.cdc.MustUnmarshal(iterator.Value(), &feeHoliday)
		feeHolidays = append(feeHolidays, feeHoliday)
	}

	return feeHolidays
}

// IsFeeHolidayActive checks if fee holiday is currently active for a CLOB pair
func (k Keeper) IsFeeHolidayActive(
	ctx sdk.Context,
	clobPairId uint32,
) bool {
	feeHoliday, err := k.GetFeeHolidayParams(ctx, clobPairId)
	if err != nil {
		return false
	}

	currentTime := ctx.BlockTime().Unix()
	return currentTime >= feeHoliday.StartTimeUnix && currentTime < feeHoliday.EndTimeUnix
}
