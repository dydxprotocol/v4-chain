package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		LimitParamsList: []LimitParams{
			DefaultUsdcRateLimitParams(),
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, limitParams := range gs.LimitParamsList {
		if err := limitParams.Validate(); err != nil {
			return err
		}
	}
	return nil
}
