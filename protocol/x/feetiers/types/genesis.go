package types

// DefaultGenesis returns the default feetiers genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: PerpetualFeeParams{
			Tiers: []*PerpetualFeeTier{
				{
					Name:        "1",
					MakerFeePpm: -110,
					TakerFeePpm: 500,
				},
				{
					Name:                      "2",
					AbsoluteVolumeRequirement: 1_000_000,
					MakerFeePpm:               -110,
					TakerFeePpm:               450,
				},
				{
					Name:                      "3",
					AbsoluteVolumeRequirement: 5_000_000,
					MakerFeePpm:               -110,
					TakerFeePpm:               400,
				},
				{
					Name:                      "4",
					AbsoluteVolumeRequirement: 25_000_000,
					MakerFeePpm:               -110,
					TakerFeePpm:               350,
				},
				{
					Name:                      "5",
					AbsoluteVolumeRequirement: 125_000_000,
					MakerFeePpm:               -110,
					TakerFeePpm:               300,
				},
				{
					Name:                           "6",
					AbsoluteVolumeRequirement:      125_000_000,
					TotalVolumeShareRequirementPpm: 10_000,
					MakerFeePpm:                    -110,
					TakerFeePpm:                    250,
				},
				{
					Name:                           "7",
					AbsoluteVolumeRequirement:      125_000_000,
					TotalVolumeShareRequirementPpm: 10_000,
					MakerVolumeShareRequirementPpm: 20_000,
					MakerFeePpm:                    -110,
					TakerFeePpm:                    250,
				},
				{
					Name:                           "8",
					AbsoluteVolumeRequirement:      125_000_000,
					TotalVolumeShareRequirementPpm: 10_000,
					MakerVolumeShareRequirementPpm: 50_000,
					MakerFeePpm:                    -110,
					TakerFeePpm:                    250,
				},
				{
					Name:                           "9",
					AbsoluteVolumeRequirement:      125_000_000,
					TotalVolumeShareRequirementPpm: 10_000,
					MakerVolumeShareRequirementPpm: 100_000,
					MakerFeePpm:                    -110,
					TakerFeePpm:                    250,
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
