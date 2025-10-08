package types

import (
	"time"
)

const (
	// Maximum duration for a fee holiday (90 days in seconds)
	MaxFeeDiscountCampaignDuration = 90 * 24 * 60 * 60

	// Maximum ppm value for fee discount
	MaxChargePpm = 1000000
)

// Validate checks if the FeeDiscountCampaignParams are valid
func (m *FeeDiscountCampaignParams) Validate(currentTime time.Time) error {
	// Validate time range (start < end)
	if m.StartTimeUnix >= m.EndTimeUnix {
		return ErrInvalidTimeRange
	}

	// Validate reasonable time range (max 30 days)
	duration := m.EndTimeUnix - m.StartTimeUnix
	if duration > MaxFeeDiscountCampaignDuration {
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
