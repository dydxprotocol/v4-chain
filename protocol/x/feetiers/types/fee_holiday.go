package types

import (
	"time"
)

const (
	// Maximum duration for a fee holiday (30 days in seconds)
	MaxFeeHolidayDuration = 30 * 24 * 60 * 60
)

func (m *FeeHolidayParams) Validate(currentTime time.Time) error {
	// Validate time range (start < end)
	if m.StartTimeUnix >= m.EndTimeUnix {
		return ErrInvalidTimeRange
	}

	// Validate reasonable time range (max 30 days)
	duration := m.EndTimeUnix - m.StartTimeUnix
	if duration > MaxFeeHolidayDuration {
		return ErrInvalidTimeRange
	}

	// Check that end time is in the future
	if m.EndTimeUnix <= currentTime.Unix() {
		return ErrInvalidTimeRange
	}

	return nil
}
