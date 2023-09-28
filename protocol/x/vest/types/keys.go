package types

const (
	// ModuleName defines the module name
	ModuleName = "vest"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vest"

	// VestEntryKey is the prefix used when storing a VestEntry in the state.
	VestEntryKey = "vest_entry"

	// CommunityTreasuryAccountName defines the root string for community treasury module account.
	CommunityTreasuryAccountName = "community_treasury"

	// CommunityVesterAccountName defines the root string for community vester module account.
	CommunityVesterAccountName = "community_vester"
)
