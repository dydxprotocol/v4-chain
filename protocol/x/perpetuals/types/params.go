package types

// Validate validates perpetual module's parameters.
func (params Params) Validate() error {
	if params.FundingRateClampFactorPpm == 0 {
		return ErrFundingRateClampFactorPpmIsZero
	}
	if params.PremiumVoteClampFactorPpm == 0 {
		return ErrPremiumVoteClampFactorPpmIsZero
	}

	return nil
}
