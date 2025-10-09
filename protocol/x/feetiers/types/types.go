package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeTiersKeeper interface {
	GetLowestMakerFee(ctx sdk.Context) int32
	GetAffiliateRefereeLowestTakerFee(ctx sdk.Context) int32
	GetPerpetualFeePpm(ctx sdk.Context, address string, isTaker bool, clobPairId uint32) int32
	GetPerpetualFeeParams(
		ctx sdk.Context,
	) PerpetualFeeParams
	SetPerpetualFeeParams(
		ctx sdk.Context,
		params PerpetualFeeParams,
	) error
	GetFeeDiscountCampaignParams(
		ctx sdk.Context,
		clobPairId uint32,
	) (params FeeDiscountCampaignParams, err error)
	SetFeeDiscountCampaignParams(
		ctx sdk.Context,
		feeHoliday FeeDiscountCampaignParams,
	) error
	HasAuthority(authority string) bool
}
