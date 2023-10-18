package types

import (
	"fmt"
)

// ClientExchangeQueryConfigOverrides is a struct representation of the client exchange query config overrides passed
// into the application as daemon flag arguments. The struct contains a list of deltas that can be applied to the
// default exchange query configs.
type ClientExchangeQueryConfigOverrides struct {
	ExchangeQueryConfigs []*ExchangeQueryConfig `json:"exchange_query_configs"`
}

// Validate validates the exchange query config override passed into the application as daemon flag arguments.
// These configs are not expected to have the full complement of fields defined, since zero values are considered
// to be unset. Thus, a valid delta may not be a valid config. However, each config does need to have a populated,
// valid exchange id.
func (ceqc *ClientExchangeQueryConfigOverrides) Validate(validExchanges map[string]struct{}) error {
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

// ApplyClientExchangeQueryConfigOverride applies the client exchange query config overrides specified by the price
// daemon startup flags to the exchange query configs.
func ApplyClientExchangeQueryConfigOverride(
	exchangeQueryConfigs map[ExchangeId]*ExchangeQueryConfig,
	clientExchangeQueryConfigOverrides *ClientExchangeQueryConfigOverrides,
) (
	updatedConfigs map[ExchangeId]*ExchangeQueryConfig,
	err error,
) {
	// Compute list of valid exchanges here from the input `exchangeQueryConfigs` map in order to avoid import
	// loops by referencing the default config directly.
	validExchanges := make(map[ExchangeId]struct{}, len(exchangeQueryConfigs))
	updatedConfigs = make(map[ExchangeId]*ExchangeQueryConfig, len(exchangeQueryConfigs))
	for exchangeId := range exchangeQueryConfigs {
		validExchanges[exchangeId] = struct{}{}
		updatedConfigs[exchangeId] = exchangeQueryConfigs[exchangeId].Copy()
	}

	for _, eqc := range clientExchangeQueryConfigOverrides.ExchangeQueryConfigs {
		config, ok := exchangeQueryConfigs[eqc.ExchangeId]
		if !ok {
			return nil, fmt.Errorf("invalid exchange id %v", eqc.ExchangeId)
		}
		updatedConfig, err := config.ApplyDeltaAndValidate(eqc, validExchanges)
		if err != nil {
			return nil, err
		}
		updatedConfigs[eqc.ExchangeId] = updatedConfig
	}
	return updatedConfigs, nil
}
