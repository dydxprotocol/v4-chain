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
	ReferredByKeyPrefix = "RB:"

	ReferredVolumeInWindowKeyPrefix = "RVW:"

	AffiliateTiersKey = "AT"

	AffiliateWhitelistKey = "AW"

	AffiliateParametersKey = "AP"

	AffiliateOverridesKey = "AO"
)
