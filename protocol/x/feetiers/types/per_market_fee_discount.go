package types

import (
	"time"
)

const (
	// Maximum duration for a fee discount period (90 days in seconds)
	MaxFeeDiscountDuration = 90 * 24 * 60 * 60

	// Maximum ppm value for fee discount
	MaxChargePpm = 1000000
)

// Validate checks if the PerMarketFeeDiscountParams are valid
func (m *PerMarketFeeDiscountParams) Validate(currentTime time.Time) error {
	// Validate time range (start < end)
	if m.StartTimeUnix >= m.EndTimeUnix {
		return ErrInvalidTimeRange
	}

	// Validate reasonable time range (max 90 days)
	duration := m.EndTimeUnix - m.StartTimeUnix
	if duration > MaxFeeDiscountDuration {
		return ErrInvalidTimeRange
	}

	// Validate charge_ppm is within valid range (0-1000000)
	if m.ChargePpm > MaxChargePpm {
		return ErrInvalidChargePpm
	}

	// Check that end time is in the future
	if m.EndTimeUnix <= currentTime.Unix() {
		return ErrInvalidTimeRange
	}

	return nil
}
