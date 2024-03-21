package prices

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PricesKeeper defines the expected Prices keeper used for `PrepareProposal`.
type PricesKeeper interface {
	GetValidMarketPriceUpdates(ctx sdk.Context) *pricestypes.MsgUpdateMarketPrices
}

// DefaultPriceUpdateGenerator is the default implementation of the PriceUpdateGenerator interface.
// This implementation retrieves the `MsgUpdateMarketPrices` from the `PricesKeeper`, i.e the default
// behavior for dydx v4's PrepareProposalHandler.
type DefaultPriceUpdateGenerator struct {
	// pk is an adapter for the
	pk PricesKeeper
}

// NewDefaultPriceUpdateGenerator returns a new DefaultPriceUpdateGenerator.
func NewDefaultPriceUpdateGenerator(keeper PricesKeeper) *DefaultPriceUpdateGenerator {
	return &DefaultPriceUpdateGenerator{
		pk: keeper,
	}
}

func (pug *DefaultPriceUpdateGenerator) GetValidMarketPriceUpdates(
	ctx sdk.Context, _ []byte) (*pricestypes.MsgUpdateMarketPrices, error) {
	msgUpdateMarketPrices := pug.pk.GetValidMarketPriceUpdates(ctx)
	if msgUpdateMarketPrices == nil {
		return nil, fmt.Errorf("MsgUpdateMarketPrices cannot be nil")
	}
	return msgUpdateMarketPrices, nil
}
