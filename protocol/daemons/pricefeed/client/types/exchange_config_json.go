package types

import (
	"fmt"
)

// ExchangeConfigJson demarshals the exchange configuration json for a particular market.
// The result is a list of parameters that define how the market is resolved on
// each supported exchange.
//
// This struct stores data in an intermediate form as it's being assigned to various
// `ExchangeMarketConfig` objects, which are keyed by exchange id. These objects are not kept
// past the time the `GetAllMarketParams` API response is parsed, and do not contain an id
// because the id is expected to be known at the time the object is in use.
type ExchangeConfigJson struct {
	Exchanges []ExchangeMarketConfigJson `json:"exchanges"`
}

// Validate validates the exchange configuration json, checking that required fields are defined
// and that market and exchange names correspond to valid markets and exchanges.
func (ecj *ExchangeConfigJson) Validate(
	exchangeNames []ExchangeId,
	marketNames map[string]MarketId,
) error {
	if len(ecj.Exchanges) == 0 {
		return fmt.Errorf("exchanges cannot be empty")
	}

	for _, exchange := range ecj.Exchanges {
		err := exchange.Validate(exchangeNames, marketNames)
		if err != nil {
			return fmt.Errorf("invalid exchange: %w", err)
		}
	}
	return nil
}
