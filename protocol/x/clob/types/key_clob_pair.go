package types

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	// ClobPairKeyPrefix is the prefix to retrieve all ClobPair
	ClobPairKeyPrefix = "ClobPair/value/"
	// NumClobPairsKey is the prefix to retrieve the cardinality of ClobPairs
	NumClobPairsKey = "ClobPair/num/"
)

// ClobPairKey returns the store key to retrieve a ClobPair from the index fields
func ClobPairKey(
	id ClobPairId,
) []byte {
	return lib.Uint32ToBytesForState(id.ToUint32())
}
