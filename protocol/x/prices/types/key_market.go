package types

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	// MarketParamKeyPrefix is the prefix to retrieve all MarketParams
	MarketParamKeyPrefix = "Market/param/"
	// MarketPriceKeyPrefix is the prefix to retrieve all MarketPrices
	MarketPriceKeyPrefix = "Market/price/"
	// NumMarketsKey is the prefix to retrieve the number of markets. All markets
	// should have an existing MarketParam and MarketPrice with a shared id between
	// 0 and NumMarkets - 1.
	NumMarketsKey = "NumMarkets"
)

// MarketKey returns the store key to retrieve a MarketParam or MarketPrice from the id field
func MarketKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytesForState(id)
}
