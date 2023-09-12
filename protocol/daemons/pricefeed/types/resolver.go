package types

// Resolver is a function type that "resolves" a slice of values to a single value.
// The function also returns an error if there was an error in resolving the value.
type Resolver func([]uint64) (uint64, error)
