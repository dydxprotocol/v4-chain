package types

// StandardParams returns the standard feetiers params for long-term operation of the network.
func StandardParams() PerpetualFeeParams {
	return PerpetualFeeParams{
		Tiers: []*PerpetualFeeTier{
			{
				Name:        "1",
				MakerFeePpm: 100,
				TakerFeePpm: 500,
			},
			{
				Name:                      "2",
				AbsoluteVolumeRequirement: 1_000_000_000_000,
				MakerFeePpm:               100,
				TakerFeePpm:               450,
			},
			{
				Name:                      "3",
				AbsoluteVolumeRequirement: 5_000_000_000_000,
				MakerFeePpm:               50,
				TakerFeePpm:               400,
			},
			{
				Name:                      "4",
				AbsoluteVolumeRequirement: 25_000_000_000_000,
				MakerFeePpm:               0,
				TakerFeePpm:               350,
			},
			{
				Name:                      "5",
				AbsoluteVolumeRequirement: 125_000_000_000_000,
				MakerFeePpm:               0,
				TakerFeePpm:               300,
			},
			{
				Name:                           "6",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerFeePpm:                    -50,
				TakerFeePpm:                    250,
			},
			{
				Name:                           "7",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerVolumeShareRequirementPpm: 10_000,
				MakerFeePpm:                    -90,
				TakerFeePpm:                    250,
			},
			{
				Name:                           "8",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerVolumeShareRequirementPpm: 20_000,
				MakerFeePpm:                    -110,
				TakerFeePpm:                    250,
			},
			{
				Name:                           "9",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerVolumeShareRequirementPpm: 40_000,
				MakerFeePpm:                    -110,
				TakerFeePpm:                    250,
			},
		},
	}
}

// PromotionalParams returns the promotional feetiers params used by the network for the first ~120 days.
// The standard params are applied via a delayed message included in the genesis state.
func PromotionalParams() PerpetualFeeParams {
	return PerpetualFeeParams{
		Tiers: []*PerpetualFeeTier{
			{
				Name:        "1",
				MakerFeePpm: -110,
				TakerFeePpm: 500,
			},
			{
				Name:                      "2",
				AbsoluteVolumeRequirement: 1_000_000_000_000,
				MakerFeePpm:               -110,
				TakerFeePpm:               450,
			},
			{
				Name:                      "3",
				AbsoluteVolumeRequirement: 5_000_000_000_000,
				MakerFeePpm:               -110,
				TakerFeePpm:               400,
			},
			{
				Name:                      "4",
				AbsoluteVolumeRequirement: 25_000_000_000_000,
				MakerFeePpm:               -110,
				TakerFeePpm:               350,
			},
			{
				Name:                      "5",
				AbsoluteVolumeRequirement: 125_000_000_000_000,
				MakerFeePpm:               -110,
				TakerFeePpm:               300,
			},
			{
				Name:                           "6",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerFeePpm:                    -110,
				TakerFeePpm:                    250,
			},
			{
				Name:                           "7",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerVolumeShareRequirementPpm: 10_000,
				MakerFeePpm:                    -110,
				TakerFeePpm:                    250,
			},
			{
				Name:                           "8",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerVolumeShareRequirementPpm: 20_000,
				MakerFeePpm:                    -110,
				TakerFeePpm:                    250,
			},
			{
				Name:                           "9",
				AbsoluteVolumeRequirement:      125_000_000_000_000,
				TotalVolumeShareRequirementPpm: 5_000,
				MakerVolumeShareRequirementPpm: 40_000,
				MakerFeePpm:                    -110,
				TakerFeePpm:                    250,
			},
		},
	}
}

// DefaultGenesis returns the default feetiers genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: PromotionalParams(),
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
