package types

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	// MarketParamKeyPrefix is the prefix to retrieve all MarketParams
	MarketParamKeyPrefix = "Market/param/"
	// MarketPriceKeyPrefix is the prefix to retrieve all MarketPrices
	MarketPriceKeyPrefix = "Market/price/"
)

// MarketKey returns the store key to retrieve a MarketParam or MarketPrice from the id field
func MarketKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytesForState(id)
}
