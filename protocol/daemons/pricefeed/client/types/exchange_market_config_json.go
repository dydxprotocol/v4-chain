package types

import (
	"fmt"
)

// ExchangeMarketConfigJson captures per-exchange information for resolving a market, including
// the ticker and conversion details. It demarshals JSON parameters from the chain for a
// particular market on a specific exchange.
type ExchangeMarketConfigJson struct {
	ExchangeName   string `json:"exchangeName"`
	Ticker         string `json:"ticker"`
	AdjustByMarket string `json:"adjustByMarket,omitempty"`
	Invert         bool   `json:"invert,omitempty"`
}

// Validate validates the exchange market configuration json. It returns an error if the
// configuration is invalid.
func (emcj *ExchangeMarketConfigJson) Validate(
	exchangeIds []ExchangeId,
	marketNames map[string]MarketId,
) error {
	// Build a map with exchange names as keys for ease of membership testing. The exchange names
	// in the config should match the exchange ids exactly.
	exchangeNameMap := make(map[ExchangeId]struct{})
	for _, exchangeName := range exchangeIds {
		exchangeNameMap[exchangeName] = struct{}{}
	}

	if emcj.ExchangeName == "" {
		return fmt.Errorf("exchange name cannot be empty")
	}
	if _, exists := exchangeNameMap[emcj.ExchangeName]; !exists {
		return fmt.Errorf("exchange name '%v' is not valid", emcj.ExchangeName)
	}
	if emcj.Ticker == "" {
		return fmt.Errorf("ticker cannot be empty")
	}
	if emcj.AdjustByMarket != "" {
		if _, exists := marketNames[emcj.AdjustByMarket]; !exists {
			return fmt.Errorf("adjustment market '%v' is not valid", emcj.AdjustByMarket)
		}
	}
	return nil
}
