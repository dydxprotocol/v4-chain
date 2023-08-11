package types

import (
	"errors"
	"fmt"
)

// DefaultGenesis returns the default Prices genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Markets:       []Market{},
		ExchangeFeeds: []ExchangeFeed{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated key for Markets.
	marketKeyMap := make(map[uint32]struct{})
	expectedMarketId := uint32(0)
	for _, market := range gs.Markets {
		if _, exists := marketKeyMap[market.Id]; exists {
			return fmt.Errorf("duplicated market id")
		}
		marketKeyMap[market.Id] = struct{}{}

		if market.Id != expectedMarketId {
			return fmt.Errorf("found gap in market id")
		}
		expectedMarketId = expectedMarketId + 1

		if len(market.Pair) == 0 {
			return errors.New("Pair must be non-empty string")
		}
	}

	// Check for duplicated key for ExchangeFeeds.
	exchangeKeyMap := make(map[uint32]struct{})
	expectedExchangeFeedId := uint32(0)
	for _, exchange := range gs.ExchangeFeeds {
		if _, exists := exchangeKeyMap[exchange.Id]; exists {
			return fmt.Errorf("duplicated exchange feed id")
		}
		exchangeKeyMap[exchange.Id] = struct{}{}

		if exchange.Id != expectedExchangeFeedId {
			return fmt.Errorf("found gap in exchange feed id")
		}
		expectedExchangeFeedId = expectedExchangeFeedId + 1
	}

	return nil
}
