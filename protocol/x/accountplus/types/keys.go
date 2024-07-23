package types

// Module name and store keys.
const (
	// The Account module uses "acc" as its module name.
	// KVStore keys cannot have other keys as prefixes so we prepend "dydx" to "accountplus"
	ModuleName = "dydxaccountplus"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
)
