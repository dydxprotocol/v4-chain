package keeper

import (
	"encoding/binary"
	"errors"

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

	if err := k.cdc.Unmarshal(b, &params); err != nil {
		return params, err
	}
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
	value, err := k.cdc.Marshal(&feeHoliday)
	if err != nil {
		return err
	}
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
		if err := k.cdc.Unmarshal(iterator.Value(), &feeHoliday); err != nil {
			// Log error and skip corrupted entry
			clobPairId := binary.BigEndian.Uint32(iterator.Key())
			k.Logger(ctx).Error(
				"failed to unmarshal fee holiday",
				"clob_pair_id", clobPairId,
				"error", err,
			)
			continue
		}
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
		if !errors.Is(err, types.ErrFeeHolidayNotFound) {
			k.Logger(ctx).Error(
				"failed to get fee holiday params",
				"clob_pair_id", clobPairId,
				"error", err,
			)
		}
		return false
	}

	currentTime := ctx.BlockTime().Unix()
	return currentTime >= feeHoliday.StartTimeUnix && currentTime < feeHoliday.EndTimeUnix
}
