package types

import "time"

// MarketPriceTimestamp maintains a `MarketId`, `Price` and `LastUpdatedAt`.
type MarketPriceTimestamp struct {
	MarketId      uint32
	Price         uint64
	LastUpdatedAt time.Time
}
