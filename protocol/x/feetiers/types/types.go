package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeTiersKeeper interface {
	GetLowestMakerFee(ctx sdk.Context) int32
	GetAffiliateRefereeLowestTakerFee(ctx sdk.Context) int32
	GetPerpetualFeePpm(ctx sdk.Context, address string, isTaker bool, feeTierOverrideIdx uint32, clobPairId uint32) int32
	GetPerpetualFeeParams(
		ctx sdk.Context,
	) PerpetualFeeParams
	SetPerpetualFeeParams(
		ctx sdk.Context,
		params PerpetualFeeParams,
	) error
	GetPerMarketFeeDiscountParams(
		ctx sdk.Context,
		clobPairId uint32,
	) (params PerMarketFeeDiscountParams, err error)
	SetPerMarketFeeDiscountParams(
		ctx sdk.Context,
		feeHoliday PerMarketFeeDiscountParams,
	) error
	SetStakingTiers(
		ctx sdk.Context,
		stakingTiers []*StakingTier,
	) error
	HasAuthority(authority string) bool
}
