package types

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	// DenomToIdKeyPrefix is the prefix to retrieve denom-to-asset-id mappings.
	DenomToIdKeyPrefix = "Asset/denom_to_id/"
	// AssetKeyPrefix is the prefix to retrieve all Assets
	AssetKeyPrefix = "Asset/value/"
	// NumAssetsKey is the prefix to retrieve the cardinality of Assets
	NumAssetsKey = "Asset/num/"
)

// AssetKey returns the store key to retrieve an Asset from the id field
func AssetKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytesForState(id)
}
