package types

// Module name and store keys.
const (
	// ModuleName defines the module name.
	ModuleName = "listing"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
)

// State.
const (
	// HardCapForMarketsKey is the key to retrieve the hard cap for listed markets.
	HardCapForMarketsKey = "HardCapForMarkets"

	// ListingVaultDepositParamsKey is the key to retrieve the listing vault deposit params.
	ListingVaultDepositParamsKey = "ListingVaultDepositParams"
)
