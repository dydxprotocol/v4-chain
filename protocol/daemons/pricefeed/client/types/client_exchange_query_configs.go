package types

import "fmt"

// ClientExchangeQueryConfigs is a struct representation of the client exchange query config overrides passed into the
// application as daemon flag arguments. The struct contains a list of deltas that can be applied to the default
// exchange query configs.
type ClientExchangeQueryConfigs struct {
	ExchangeQueryConfigs []*ExchangeQueryConfig `json:"exchange_query_configs"`
}

// ValidateDelta validates the exchange query config override passed into the application as daemon flag arguments.
// These configs are not expected to have the full complement of fields defined, since zero values are considered
// to be unset. However, each config does need to have a populated, valid exchange id.
func (ceqc *ClientExchangeQueryConfigs) ValidateDelta(validExchanges map[string]struct{}) error {
	seenExchanges := make(map[ExchangeId]struct{})
	for _, eqc := range ceqc.ExchangeQueryConfigs {
		if err := eqc.ValidateDelta(validExchanges); err != nil {
			return err
		}
		if _, ok := seenExchanges[eqc.ExchangeId]; ok {
			return fmt.Errorf("duplicate exchange id %v", eqc.ExchangeId)
		}
		seenExchanges[eqc.ExchangeId] = struct{}{}
	}
	return nil
}

// Validate validates the client exchange query configs.
func (ceqc *ClientExchangeQueryConfigs) Validate(validExchanges map[string]struct{}) error {
	seenExchanges := make(map[ExchangeId]struct{})
	for _, eqc := range ceqc.ExchangeQueryConfigs {
		if err := eqc.Validate(validExchanges); err != nil {
			return err
		}
		if _, ok := seenExchanges[eqc.ExchangeId]; ok {
			return fmt.Errorf("duplicate exchange id %v", eqc.ExchangeId)
		}
		seenExchanges[eqc.ExchangeId] = struct{}{}
	}
	return nil
}

// ApplyClientExchangeQueryConfigOverride applies the client exchange query config overrides specified by the price
// daemon startup flags to the exchange query configs.
func ApplyClientExchangeQueryConfigOverride(
	exchangeQueryConfigs map[ExchangeId]*ExchangeQueryConfig,
	clientExchangeQueryConfigOverrides *ClientExchangeQueryConfigs,
) (
	updatedConfigs map[ExchangeId]*ExchangeQueryConfig,
	err error,
) {
	updatedConfigs = make(map[ExchangeId]*ExchangeQueryConfig)
	for exchangeId, eqc := range exchangeQueryConfigs {
		updatedConfigs[exchangeId] = eqc.Copy()
	}
	for _, eqc := range clientExchangeQueryConfigOverrides.ExchangeQueryConfigs {
		if _, ok := updatedConfigs[eqc.ExchangeId]; !ok {
			return nil, fmt.Errorf("invalid exchange id %v", eqc.ExchangeId)
		}
		if err := updatedConfigs[eqc.ExchangeId].ApplyDelta(eqc); err != nil {
			return nil, err
		}
	}
	return updatedConfigs, nil
}
