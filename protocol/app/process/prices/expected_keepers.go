package prices

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PricesKeeper defines the expected Prices keeper used for `DefaultMarketPriceDecoder`.
type PricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *pricestypes.MsgUpdateMarketPrices,
		performNonDeterministicValidation bool,
	) error
}
