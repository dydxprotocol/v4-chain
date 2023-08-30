package types

import "time"

const (
	// MaxPriceAge defines the duration in which a price update is valid for.
	MaxPriceAge = time.Duration(30_000_000_000) // 30 sec, duration uses nanoseconds.
)
