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

	// VestEntryKeyPrefix is the prefix used when storing a VestEntry in the state.
	VestEntryKeyPrefix = "vest_entry"

	// CommunityTreasury defines the root string for community treasury module account.
	CommunityTreasury = "community_treasury"

	// CommunityVester defines the root string for community vester module account.
	CommunityVester = "community_vester"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// VestEntryKey returns the store key (using the vester account) to retrieve a vest entry from state.
func VestEntryKey(
	vesterAccount string,
) []byte {
	return []byte(vesterAccount)
}
