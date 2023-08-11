package types

import (
	"github.com/dydxprotocol/v4/lib"
)

const (
	// MarketKeyPrefix is the prefix to retrieve all Markets
	MarketKeyPrefix = "Market/value/"
	// NumMarketsKey is the prefix to retrieve the cardinality of Markets
	NumMarketsKey = "Market/num/"
)

// MarketKey returns the store key to retrieve a Market from the id field
func MarketKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytesForState(id)
}
