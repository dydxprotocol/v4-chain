package types

import "math"

func (m *PerpetualFeeParams) Validate() error {
	if len(m.Tiers) == 0 {
		return ErrNoTiersExist
	}

	if m.Tiers[0].AbsoluteVolumeRequirement != 0 ||
		m.Tiers[0].TotalVolumeShareRequirementPpm != 0 ||
		m.Tiers[0].MakerVolumeShareRequirementPpm != 0 {
		return ErrInvalidFirstTierRequirements
	}

	for i := 1; i < len(m.Tiers); i++ {
		prevTier := m.Tiers[i-1]
		currTier := m.Tiers[i]
		if prevTier.AbsoluteVolumeRequirement > currTier.AbsoluteVolumeRequirement ||
			prevTier.TotalVolumeShareRequirementPpm > currTier.TotalVolumeShareRequirementPpm ||
			prevTier.MakerVolumeShareRequirementPpm > currTier.MakerVolumeShareRequirementPpm {
			return ErrTiersOutOfOrder
		}
	}

	lowestMakerFee := int32(math.MaxInt32)
	lowestTakerFee := int32(math.MaxInt32)
	for _, tier := range m.Tiers {
		if tier.MakerFeePpm < lowestMakerFee {
			lowestMakerFee = tier.MakerFeePpm
		}
		if tier.TakerFeePpm < lowestTakerFee {
			lowestTakerFee = tier.TakerFeePpm
		}
	}

	// Prevent overflow
	if int64(lowestMakerFee)+int64(lowestTakerFee) < 0 {
		return ErrInvalidFee
	}

	return nil
}
