package epochs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
)

// BeginBlocker executes all ABCI BeginBlock logic respective to the epochs module.
func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	epochs := keeper.GetAllEpochInfo(ctx)
	// Iterate through all epoch infos, calls MaybeStartNextEpoch() which
	// initializes and/or increments the epoch if applicable.
	for _, epoch := range epochs {
		if _, err := keeper.MaybeStartNextEpoch(ctx, epoch.GetEpochInfoName()); err != nil {
			panic(err)
		}
	}
}
