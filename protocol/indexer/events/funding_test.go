package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"

	"github.com/dydxprotocol/v4/indexer/events"
	"github.com/stretchr/testify/require"
)

func TestNewFundingEvent(t *testing.T) {
	tests := map[string]struct {
		updateType   events.FundingEvent_Type
		updates      []events.FundingUpdate
		txnHash      lib.TxHash
		newEventFunc func(updates []events.FundingUpdate) *events.FundingEvent
	}{
		"premium samples": {
			updateType: events.FundingEvent_TYPE_PREMIUM_SAMPLE,
			updates: []events.FundingUpdate{
				{
					PerpetualId:     0,
					FundingValuePpm: 1000,
				},
				{
					PerpetualId:     1,
					FundingValuePpm: 0,
				},
			},
			txnHash:      constants.TestTxHashString,
			newEventFunc: events.NewPremiumSamplesEvent,
		},
		"funding rates and indices": {
			updateType: events.FundingEvent_TYPE_FUNDING_RATE_AND_INDEX,
			updates: []events.FundingUpdate{
				{
					PerpetualId:     0,
					FundingValuePpm: -1000,
					FundingIndex:    dtypes.NewInt(0),
				},
				{
					PerpetualId:     1,
					FundingValuePpm: 0,
					FundingIndex:    dtypes.NewInt(1000),
				},
				{
					PerpetualId:     2,
					FundingValuePpm: 5000,
					FundingIndex:    dtypes.NewInt(-1000),
				},
			},
			txnHash:      constants.TestTxHashString,
			newEventFunc: events.NewFundingRatesAndIndicesEvent,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			FundingEvent := tc.newEventFunc(tc.updates)
			expectedFundingEventProto := &events.FundingEvent{
				Type:    tc.updateType,
				Updates: tc.updates,
			}
			require.Equal(t, expectedFundingEventProto, FundingEvent)
		})
	}
}
