package types

// Validate validates perpetual module's parameters.
func (params Params) Validate() error {
	if err := ValidateFundingRateClampFactorPpm(params.FundingRateClampFactorPpm); err != nil {
		return err
	}
	if err := ValidatePremiumVoteClampFactorPpm(params.PremiumVoteClampFactorPpm); err != nil {
		return err
	}
	return nil
}

// ValidateFundingRateClampFactorPpm validates that `fundingRateClampFactorPpm` is not zero.
func ValidateFundingRateClampFactorPpm(num uint32) error {
	if num == 0 {
		return ErrFundingRateClampFactorPpmIsZero
	}
	return nil
}

// ValidatePremiumVoteClampFactorPpm validates that `premiumVoteClampFactorPpm` is not zero.
func ValidatePremiumVoteClampFactorPpm(num uint32) error {
	if num == 0 {
		return ErrPremiumVoteClampFactorPpmIsZero
	}
	return nil
}
