package types

import (
	"fmt"
)

// DefaultGenesis returns the default Prices genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		MarketParams: []MarketParam{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated key for Markets.
	marketParamKeyMap := make(map[uint32]struct{})
	expectedMarketId := uint32(0)
	for _, marketParam := range gs.MarketParams {
		if _, exists := marketParamKeyMap[marketParam.Id]; exists {
			return fmt.Errorf("duplicated market param id")
		}
		marketParamKeyMap[marketParam.Id] = struct{}{}

		if marketParam.Id != expectedMarketId {
			return fmt.Errorf("found gap in market param id")
		}
		expectedMarketId = expectedMarketId + 1

		if err := marketParam.Validate(); err != nil {
			return err
		}
	}

	if len(gs.MarketParams) != len(gs.MarketPrices) {
		return fmt.Errorf("expected the same number of market prices and market params")
	}

	for i, marketPrice := range gs.MarketPrices {
		if err := marketPrice.ValidateFromParam(gs.MarketParams[i]); err != nil {
			return err
		}
	}

	return nil
}
