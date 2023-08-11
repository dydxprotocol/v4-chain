package events

import (
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
)

// NewPremiumSamplesEvent creates a FundingEvent representing a list of new funding premium
// samples generated at the end of each `funding-sample` epoch.
func NewPremiumSamplesEvent(
	newSamplesForEvent []perptypes.FundingPremium,
) *FundingEvent {
	return newFundingEvent(
		newSamplesForEvent,
		FundingEvent_TYPE_PREMIUM_SAMPLE,
	)
}

// NewFundingRatesEvent creates a FundingEvent representing a list of new funding rates
// generated at the end of each `funding-tick` epoch.
func NewFundingRatesEvent(
	newFundingRatesForEvent []perptypes.FundingPremium,
) *FundingEvent {
	return newFundingEvent(
		newFundingRatesForEvent,
		FundingEvent_TYPE_FUNDING_RATE,
	)
}

func newFundingEvent(
	newValuesForEvent []perptypes.FundingPremium,
	valueType FundingEvent_Type,
) *FundingEvent {
	return &FundingEvent{
		Values: newValuesForEvent,
		Type:   valueType,
	}
}
