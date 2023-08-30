package perpetuals

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// EndBlocker executes all ABCI EndBlock logic respective to the perpetuals module.
func EndBlocker(ctx sdk.Context, k types.PerpetualsKeeper) {
	// We don't expect the following two calls to take effect in the same block,
	// since according to their genesis set-up, `funding-tick` happens every exact
	// hour while `funding-sample` happens every minute on the half-minute.
	// This could change if the `EpochInfo` parameters are changed by governance.
	// If they do take effect in the same block, `funding-sample` should be processed
	// first so that new samples are processed in `MaybeProcessNewFundingTickEpoch`.
	k.MaybeProcessNewFundingSampleEpoch(ctx)
	k.MaybeProcessNewFundingTickEpoch(ctx)
}
