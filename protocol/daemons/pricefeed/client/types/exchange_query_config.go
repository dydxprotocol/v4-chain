package types

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// ExchangeQueryConfig contains configuration values for querying an exchange, passed in on startup.
// The configuration values include
//  1. `ExchangeId`
//  2. `IntervalMs` delay between task-loops where each task-loop sends API requests to an exchange
//  3. `TimeoutMs` max time to wait on an API call to an exchange
//  4. `MaxQueries` max number of API calls made to an exchange per task-loop. This parameter is used
//     for rate limiting requests to the exchange.
//
// For single-market API exchanges, the price fetcher will send approximately
// MaxQueries API responses into the exchange's buffered channel once every IntervalMs milliseconds.
// Note: the `ExchangeQueryConfig` will be used in the map of `{ exchangeId, `ExchangeQueryConfig` }`
// that dictates how the pricefeed client queries for market prices.
type ExchangeQueryConfig struct {
	ExchangeId ExchangeId `json:"exchange_id" validate:"required"`
	// Disabled is used to disable querying an exchange. Use `Disabled` instead of `Enabled` so the default
	// value is an operating exchange.
	Disabled   bool   `json:"disabled"`
	IntervalMs uint32 `json:"interval_ms" validate:"gt=0"`
	TimeoutMs  uint32 `json:"timeout_ms" validate:"gt=0"`
	MaxQueries uint32 `json:"max_queries" validate:"gt=0"`
}

// cache the set of all valid exchange ids.
var (
	jsonValidator *validator.Validate
)

// getValidator returns a cached validator for validating json fields.
func getValidator() *validator.Validate {
	if jsonValidator == nil {
		jsonValidator = validator.New()
	}
	return jsonValidator
}

// Copy returns a copy of the exchange query config.
func (eqc *ExchangeQueryConfig) Copy() *ExchangeQueryConfig {
	return &ExchangeQueryConfig{
		ExchangeId: eqc.ExchangeId,
		Disabled:   eqc.Disabled,
		IntervalMs: eqc.IntervalMs,
		TimeoutMs:  eqc.TimeoutMs,
		MaxQueries: eqc.MaxQueries,
	}
}

// Validate validates the exchange query config.
// Note: validateExchanges must be passed in as an argument to avoid an import cycle.
func (eqc *ExchangeQueryConfig) Validate(validExchanges map[ExchangeId]struct{}) error {
	if _, ok := validExchanges[eqc.ExchangeId]; !ok {
		return fmt.Errorf("invalid exchange id %v", eqc.ExchangeId)
	}
	err := getValidator().Struct(eqc)
	if err != nil {
		return fmt.Errorf("invalid exchange query config: %w", err)
	}
	return nil
}

// ValidateDelta validates the exchange query config override passed into the application as daemon flag arguments.
// In this case, we expect that zero fields are unset and only validate that the id is valid.
// Note: validateExchanges must be passed in as an argument to avoid an import cycle.
func (eqc *ExchangeQueryConfig) ValidateDelta(validExchanges map[ExchangeId]struct{}) error {
	if _, ok := validExchanges[eqc.ExchangeId]; !ok {
		return fmt.Errorf("invalid exchange id %v", eqc.ExchangeId)
	}
	return nil
}

// ApplyDelta applies the delta to the exchange query config and validates the result. It does not mutate
// the original exchange query config.
func (eqc *ExchangeQueryConfig) ApplyDelta(
	delta *ExchangeQueryConfig,
) (
	updatedConfig *ExchangeQueryConfig,
	err error,
) {
	if delta.ExchangeId != eqc.ExchangeId {
		return nil, fmt.Errorf("exchange id mismatch: %v, %v", delta.ExchangeId, eqc.ExchangeId)
	}

	updatedConfig = eqc.Copy()

	// Always update disabled status.
	updatedConfig.Disabled = delta.Disabled

	// Only update other fields if the value is specified. We consider 0 to be invalid / not set.
	if delta.IntervalMs != 0 {
		updatedConfig.IntervalMs = delta.IntervalMs
	}

	if delta.TimeoutMs != 0 {
		updatedConfig.TimeoutMs = delta.TimeoutMs
	}

	if delta.MaxQueries != 0 {
		updatedConfig.MaxQueries = delta.MaxQueries
	}

	return updatedConfig, nil
}
