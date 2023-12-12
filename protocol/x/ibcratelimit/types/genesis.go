package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	// TODO(CORE-824): Implement keepers
	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO(CORE-824): Implement keepers
	return nil
}
