package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "assets"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// DenomToIdKeyPrefix is the prefix to retrieve denom-to-asset-id mappings
	DenomToIdKeyPrefix = "denom_to_id/"

	// AssetKeyPrefix is the prefix to retrieve all Assets
	AssetKeyPrefix = "asset/"
)
