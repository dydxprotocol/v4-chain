package process

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// ProcessPricesKeeper defines the expected Prices keeper used for `ProcessProposal`.
type ProcessPricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *types.MsgUpdateMarketPrices,
		performNonDeterministicValidation bool,
	) error
}
