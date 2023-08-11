package events

// NewPremiumSamplesEvent creates a FundingEvent representing a list of new funding premium
// samples generated at the end of each `funding-sample` epoch.
func NewPremiumSamplesEvent(
	newSamplesForEvent []FundingUpdateV1,
) *FundingEventV1 {
	return newFundingEvent(
		newSamplesForEvent,
		FundingEventV1_TYPE_PREMIUM_SAMPLE,
	)
}

// NewFundingRatesAndIndicesEvent creates a FundingEvent representing a list of new
// funding rates generated at the end of each `funding-tick` epoch and funding indices
// accordingly updated with `funding rate * price`.
func NewFundingRatesAndIndicesEvent(
	newFundingRatesAndIndicesForEvent []FundingUpdateV1,
) *FundingEventV1 {
	return newFundingEvent(
		newFundingRatesAndIndicesForEvent,
		FundingEventV1_TYPE_FUNDING_RATE_AND_INDEX,
	)
}

func newFundingEvent(
	newUpdatesForEvent []FundingUpdateV1,
	updateType FundingEventV1_Type,
) *FundingEventV1 {
	return &FundingEventV1{
		Updates: newUpdatesForEvent,
		Type:    updateType,
	}
}
