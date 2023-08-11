package types

import (
	"github.com/dydxprotocol/v4/lib"
)

const (
	// ExchangeFeedKeyPrefix is the prefix to retrieve all ExchangeFeeds
	ExchangeFeedKeyPrefix = "ExchangeFeed/value/"
	// NumExchangeFeedsKey is the prefix to retrieve the cardinality of ExchangeFeeds
	NumExchangeFeedsKey = "ExchangeFeed/num/"
)

// ExchangeFeedKey returns the store key to retrieve an ExchangeFeed from the id fields
func ExchangeFeedKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytesForState(id)
}
