package affiliates

import (
	"runtime/debug"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
)

func EndBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	err := keeper.AggregateAffiliateVolumeForFills(ctx)
	if err != nil {
		log.ErrorLog(ctx, "error aggregating affiliate volume for fills", "error",
			err, "stack", string(debug.Stack()))
	}
}
