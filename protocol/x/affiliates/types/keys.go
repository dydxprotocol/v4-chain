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

<<<<<<< HEAD
	ReferredVolumeKeyPrefix = "RV:"
=======
	ReferredVolumeInWindowKeyPrefix = "RVW:"
>>>>>>> 1b536022 (Integrate commission and overrides to fee tier calculation (#3117))

	AffiliateTiersKey = "AT"

	AffiliateWhitelistKey = "AW"
)
