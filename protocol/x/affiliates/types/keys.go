package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "affiliates"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	ReferredByKeyPrefix = "ReferredBy:"

	ReferredVolumeKeyPrefix = "ReferredVolume:"

	AffiliateTiersKey = "AffiliateTiers"

	AffiliateWhitelistKey = "AffiliateWhitelist"
)
