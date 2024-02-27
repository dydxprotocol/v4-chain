package prices

import (
	"fmt"
)

// InvalidPriceError represents an error thrown when a price retrieved from a vote-extension is invalid.
// - MarketID: the market-id of the market with the invalid price
// - Reason: the reason the price is invalid
type InvalidPriceError struct {
	MarketID uint64
	Reason   string
}

func (e *InvalidPriceError) Error() string {
	return fmt.Sprintf("invalid price for market %d: %s", e.MarketID, e.Reason)
}
