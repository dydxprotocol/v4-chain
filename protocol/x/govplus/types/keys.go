package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "govplus"

	// StoreKey defines the primary module store key
	// This is not govplus because then StoreKey gov (from x/gov) would be a prefix of it,
	// and that is not allowed.
	StoreKey = "dydxgovplus"
)
