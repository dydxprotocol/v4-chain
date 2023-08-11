package types

import (
	"time"
)

// PriceTimestamp maintains a price and its last update timestamp.
type PriceTimestamp struct {
	LastUpdateTime time.Time
	Price          uint64
}

// NewPriceTimestamp creates a new PriceTimestamp.
func NewPriceTimestamp() *PriceTimestamp {
	return &PriceTimestamp{}
}

// UpdatePrice updates the price if the given update has a greater timestamp. Returns true if
// updating succeeds. Otherwise, returns false.
func (pt *PriceTimestamp) UpdatePrice(price uint64, newUpdateTime *time.Time) bool {
	if newUpdateTime.After(pt.LastUpdateTime) {
		pt.LastUpdateTime = *newUpdateTime
		pt.Price = price

		return true
	}

	return false
}

// GetValidPrice returns (price, true) if the last update time is greater than or
// equal to the given cutoff time. Otherwise returns (0, false).
func (pt *PriceTimestamp) GetValidPrice(cutoffTime time.Time) (uint64, bool) {
	if pt.LastUpdateTime.Before(cutoffTime) {
		return 0, false
	}
	return pt.Price, true
}
