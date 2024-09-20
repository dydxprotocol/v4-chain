package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// GetPerpetualFeeParams returns the PerpetualFeeParams in state.
func (k Keeper) GetPerpetualFeeParams(
	ctx sdk.Context,
) (
	params types.PerpetualFeeParams,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.PerpetualFeeParamsKey))
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

	lowestMakerFee := GetLowestMakerFeeFromTiers(params.Tiers)
	lowestTakerFee := GetAffiliateRefereeLowestTakerFeeFromTiers(params.Tiers)

	unconditionalRevShareConfig, err := k.revShareKeeper.GetUnconditionalRevShareConfigParams(ctx)
	if err != nil {
		return err
	}

	marketMapperRevShareParams := k.revShareKeeper.GetMarketMapperRevenueShareParams(ctx)
	if err != nil {
		return err
	}

	valid := k.revShareKeeper.ValidateRevShareSafety(ctx, unconditionalRevShareConfig,
		marketMapperRevShareParams, lowestTakerFee, lowestMakerFee)
	if !valid {
		return errorsmod.Wrapf(
			types.ErrRevShareSafetyViolation,
			"rev share safety violation",
		)
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.PerpetualFeeParamsKey), b)

	return nil
}
