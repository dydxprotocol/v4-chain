package prices

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PreBlocker executes all ABCI PreBlock logic respective to the clob module.
func PreBlocker(
	ctx sdk.Context,
	keeper types.PricesKeeper,
) {
	keeper.Hydrate(ctx)
}
