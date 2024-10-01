package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeTiersKeeper interface {
	GetLowestMakerFee(ctx sdk.Context) int32
	GetAffiliateRefereeLowestTakerFee(ctx sdk.Context) int32
	GetPerpetualFeePpm(ctx sdk.Context, address string, isTaker bool) int32
	GetPerpetualFeeParams(
		ctx sdk.Context,
	) PerpetualFeeParams
	SetPerpetualFeeParams(
		ctx sdk.Context,
		params PerpetualFeeParams,
	) error
	HasAuthority(authority string) bool
}
