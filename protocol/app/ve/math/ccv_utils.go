package voteweighted

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTotalPower(
	ctx sdk.Context,
	ccvStore CCValidatorStore,
) math.Int {
	total := math.NewInt(0)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, ccVal := range ccvStore.GetAllCCValidator(sdkCtx) {
		total = total.Add(math.NewInt(ccVal.Power))
	}
	return total
}
