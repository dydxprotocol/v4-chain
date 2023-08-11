package types

import "time"

// DefaultGenesis returns the default stats genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{
			WindowDuration: time.Duration(30 * 24 * time.Hour),
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	return nil
}
