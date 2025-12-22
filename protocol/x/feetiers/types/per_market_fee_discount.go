package types

import (
	"time"
)

const (
	// Maximum ppm value for fee discount
	MaxChargePpm = 1_000_000
)

// Validate checks if the PerMarketFeeDiscountParams are valid
func (m *PerMarketFeeDiscountParams) Validate(currentTime time.Time) error {
	// Validate time range (start < end)
	if !m.StartTime.Before(m.EndTime) {
		return ErrInvalidTimeRange
	}

	// Validate charge_ppm is within valid range (0-1000000)
	if m.ChargePpm > MaxChargePpm {
		return ErrInvalidChargePpm
	}

	// Check that end time is in the future
	if !m.EndTime.After(currentTime) {
		return ErrInvalidTimeRange
	}

	return nil
}
