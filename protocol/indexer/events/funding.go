package events

// NewPremiumSamplesEvent creates a FundingEvent representing a list of new funding premium
// samples generated at the end of each `funding-sample` epoch.
func NewPremiumSamplesEvent(
	newSamplesForEvent []FundingUpdate,
) *FundingEvent {
	return newFundingEvent(
		newSamplesForEvent,
		FundingEvent_TYPE_PREMIUM_SAMPLE,
	)
}

// NewFundingRatesAndIndicesEvent creates a FundingEvent representing a list of new
// funding rates generated at the end of each `funding-tick` epoch and funding indices
// accordingly updated with `funding rate * price`.
func NewFundingRatesAndIndicesEvent(
	newFundingRatesAndIndicesForEvent []FundingUpdate,
) *FundingEvent {
	return newFundingEvent(
		newFundingRatesAndIndicesForEvent,
		FundingEvent_TYPE_FUNDING_RATE_AND_INDEX,
	)
}

func newFundingEvent(
	newUpdatesForEvent []FundingUpdate,
	updateType FundingEvent_Type,
) *FundingEvent {
	return &FundingEvent{
		Updates: newUpdatesForEvent,
		Type:    updateType,
	}
}
