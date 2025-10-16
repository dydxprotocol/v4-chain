package types

import (
	"time"
)

const (
	// Maximum duration for a fee discount period (90 days in seconds)
	MaxFeeDiscountDuration = 90 * 24 * 60 * 60

	// Maximum ppm value for fee discount
	MaxChargePpm = 1_000_000
)

// Validate checks if the PerMarketFeeDiscountParams are valid
func (m *PerMarketFeeDiscountParams) Validate(currentTime time.Time) error {
	// Validate time range (start < end)
	if !m.StartTime.Before(m.EndTime) {
		return ErrInvalidTimeRange
	}

	// Validate reasonable time range (max 90 days)
	duration := m.EndTime.Sub(m.StartTime)
	if duration.Seconds() > MaxFeeDiscountDuration {
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
