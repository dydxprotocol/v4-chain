package sim_helpers

type GenesisParameters[T any] struct {
	// Reasonable indicates a value that is similar to real-world conditions.
	Reasonable T
	// Valid indicates a value that can be handled by the chain without panicking.
	Valid T
}
