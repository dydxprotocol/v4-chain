package types

// Module name and store keys
const (
	// ModuleName defines the module name
	// Use `ratelimit` instead of `ratelimit` to prevent potential key space conflicts with the IBC module.
	ModuleName = "ratelimit"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// DenomCapacityKeyPrefix is the prefix for the key-value store for DenomCapacity
	DenomCapacityKeyPrefix = "DenomCapacity:"

	// LimitParamsKeyPrefix is the prefix for the key-value store for LimitParams
	LimitParamsKeyPrefix = "LimitParams:"
)

// State
const ()
