package keeper

import (
	"encoding/binary"
	"errors"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// GetFeeDiscountCampaignParams retrieves the FeeDiscountCampaignParams for a CLOB pair
func (k Keeper) GetFeeDiscountCampaignParams(
	ctx sdk.Context,
	clobPairId uint32,
) (params types.FeeDiscountCampaignParams, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.FeeDiscountCampaignPrefix))
	b := store.Get(lib.Uint32ToKey(clobPairId))

	if b == nil {
		return params, types.ErrFeeDiscountCampaignNotFound
	}

	if err := k.cdc.Unmarshal(b, &params); err != nil {
		return params, err
	}

	return params, nil
}

// SetFeeDiscountCampaignParams stores fee discount campaign configuration
func (k Keeper) SetFeeDiscountCampaignParams(
	ctx sdk.Context,
	params types.FeeDiscountCampaignParams,
) error {
	// Validate the params
	err := params.Validate(ctx.BlockTime())
	if err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.FeeDiscountCampaignPrefix))
	key := lib.Uint32ToKey(params.ClobPairId)
	value, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(key, value)
	return nil
}

// GetAllFeeDiscountCampaignParams returns all configured fee discount campaigns
func (k Keeper) GetAllFeeDiscountCampaignParams(
	ctx sdk.Context,
) []types.FeeDiscountCampaignParams {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.FeeDiscountCampaignPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	campaigns := []types.FeeDiscountCampaignParams{}
	for ; iterator.Valid(); iterator.Next() {
		var campaign types.FeeDiscountCampaignParams
		if err := k.cdc.Unmarshal(iterator.Value(), &campaign); err != nil {
			// Log error and skip corrupted entry
			clobPairId := binary.BigEndian.Uint32(iterator.Key())
			k.Logger(ctx).Error(
				"failed to unmarshal fee discount campaign",
				"clob_pair_id", clobPairId,
				"error", err,
			)
			continue
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns
}

// GetDiscountPpm returns the charge PPM (parts per million) for a CLOB pair.
// If a discount campaign is active, it returns the campaign's charge PPM.
// If no active campaign exists, it returns 1,000,000 (100% charge -> no discount).
func (k Keeper) GetDiscountPpm(
	ctx sdk.Context,
	clobPairId uint32,
) uint32 {
	campaign, err := k.GetFeeDiscountCampaignParams(ctx, clobPairId)
	if err != nil {
		// If the error is ErrFeeDiscountCampaignNotFound, this is normal
		if !errors.Is(err, types.ErrFeeDiscountCampaignNotFound) {
			// If it's any other type of error, log it as it's unexpected
			k.Logger(ctx).Error(
				"failed to get fee discount campaign params",
				"clob_pair_id", clobPairId,
				"error", err,
			)
		}
		return types.MaxChargePpm
	}

	currentTime := ctx.BlockTime().Unix()
	if currentTime >= campaign.StartTimeUnix && currentTime < campaign.EndTimeUnix {
		return campaign.ChargePpm
	}

	return types.MaxChargePpm
}
