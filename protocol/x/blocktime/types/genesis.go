package types

import "time"

// DefaultGenesis returns the default blocktime genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DowntimeParams{
			Durations: []time.Duration{
				5 * time.Minute,
				30 * time.Minute,
			},
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
