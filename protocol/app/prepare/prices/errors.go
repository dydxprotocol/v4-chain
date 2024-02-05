package prices

import (
	"fmt"
)

type InvalidPriceError struct {
	MarketID uint64
	Reason   string
}

func (e *InvalidPriceError) Error() string {
	return fmt.Sprintf("invalid price for market %d: %s", e.MarketID, e.Reason)
}
