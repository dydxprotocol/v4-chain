package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "vest"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// VestEntryKeyPrefix is the prefix used when storing a VestEntry in the state.
	VestEntryKeyPrefix = "Entry:"
)

// Module accounts
const (
	// CommunityTreasuryAccountName defines the root string for community treasury module account.
	CommunityTreasuryAccountName = "community_treasury"

	// CommunityVesterAccountName defines the root string for community vester module account.
	CommunityVesterAccountName = "community_vester"
)
