package types

// DefaultGenesis returns the default feetiers genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: PerpetualFeeParams{
			Tiers: []*PerpetualFeeTier{
				{
					Name: "1",
				},
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
