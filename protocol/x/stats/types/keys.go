package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "stats"

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "transient_" + ModuleName

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// ParamsKey defines the key for the params
	ParamsKey = "params"
)
