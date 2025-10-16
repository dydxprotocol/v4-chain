package keeper

import (
	"encoding/binary"
	"errors"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// GetPerMarketFeeDiscountParams retrieves the PerMarketFeeDiscountParams for a CLOB pair
func (k Keeper) GetPerMarketFeeDiscountParams(
	ctx sdk.Context,
	clobPairId uint32,
) (params types.PerMarketFeeDiscountParams, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketFeeDiscountPrefix))
	b := store.Get(lib.Uint32ToKey(clobPairId))
	if b == nil {
		return params, types.ErrMarketFeeDiscountNotFound
	}
	if err := k.cdc.Unmarshal(b, &params); err != nil {
		return params, err
	}
	return params, nil
}

// SetPerMarketFeeDiscountParams stores per market fee discount parameters
func (k Keeper) SetPerMarketFeeDiscountParams(
	ctx sdk.Context,
	params types.PerMarketFeeDiscountParams,
) error {
	// Validate the params
	err := params.Validate(ctx.BlockTime())
	if err != nil {
		return err
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketFeeDiscountPrefix))
	key := lib.Uint32ToKey(params.ClobPairId)
	value, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(key, value)
	return nil
}

// GetAllMarketFeeDiscountParams returns all configured per-market fee discounts
func (k Keeper) GetAllMarketFeeDiscountParams(
	ctx sdk.Context,
) []types.PerMarketFeeDiscountParams {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketFeeDiscountPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	discountParams := []types.PerMarketFeeDiscountParams{}
	for ; iterator.Valid(); iterator.Next() {
		var marketDiscount types.PerMarketFeeDiscountParams
		if err := k.cdc.Unmarshal(iterator.Value(), &marketDiscount); err != nil {
			// Log error and skip corrupted entry
			clobPairId := binary.BigEndian.Uint32(iterator.Key())
			k.Logger(ctx).Error(
				"failed to unmarshal per-market fee discount",
				"clob_pair_id", clobPairId,
				"error", err,
			)
			continue
		}
		discountParams = append(discountParams, marketDiscount)
	}
	return discountParams
}

// GetDiscountedPpm returns the charge PPM (parts per million) for a CLOB pair.
// If a fee discount is active, it returns the charge PPM.
// If no active fee discount exists, it returns 1,000,000 (100% charge -> no discount).
func (k Keeper) GetDiscountedPpm(
	ctx sdk.Context,
	clobPairId uint32,
) uint32 {
	marketDiscount, err := k.GetPerMarketFeeDiscountParams(ctx, clobPairId)
	if err != nil {
		// If the error is ErrMarketFeeDiscountNotFound, this is normal
		if !errors.Is(err, types.ErrMarketFeeDiscountNotFound) {
			// If it's any other type of error, log it as it's unexpected
			k.Logger(ctx).Error(
				"failed to get per market fee discount params",
				"clob_pair_id", clobPairId,
				"error", err,
			)
		}
		return types.MaxChargePpm
	}

	currentTime := ctx.BlockTime()
	if (currentTime.Equal(marketDiscount.StartTime) || currentTime.After(marketDiscount.StartTime)) &&
		currentTime.Before(marketDiscount.EndTime) {
		return marketDiscount.ChargePpm
	} else {
		return types.MaxChargePpm
	}
}
