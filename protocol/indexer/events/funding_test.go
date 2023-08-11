package events_test

import (
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"
	"testing"

	"github.com/dydxprotocol/v4/indexer/events"
	"github.com/stretchr/testify/require"

	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
)

func TestNewFundingEvent(t *testing.T) {
	tests := map[string]struct {
		valueType    events.FundingEvent_Type
		values       []perptypes.FundingPremium
		txnHash      lib.TxHash
		newEventFunc func(values []perptypes.FundingPremium) *events.FundingEvent
	}{
		"premium samples": {
			valueType: events.FundingEvent_TYPE_PREMIUM_SAMPLE,
			values: []perptypes.FundingPremium{
				{
					PerpetualId: 0,
					PremiumPpm:  1000,
				},
				{
					PerpetualId: 1,
					PremiumPpm:  0,
				},
			},
			txnHash:      constants.TestTxHashString,
			newEventFunc: events.NewPremiumSamplesEvent,
		},
		"funding rates": {
			valueType: events.FundingEvent_TYPE_FUNDING_RATE,
			values: []perptypes.FundingPremium{
				{
					PerpetualId: 0,
					PremiumPpm:  -1000,
				},
				{
					PerpetualId: 1,
					PremiumPpm:  0,
				},
				{
					PerpetualId: 2,
					PremiumPpm:  5000,
				},
			},
			txnHash:      constants.TestTxHashString,
			newEventFunc: events.NewFundingRatesEvent,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			FundingEvent := tc.newEventFunc(tc.values)
			expectedFundingEventProto := &events.FundingEvent{
				Type:   tc.valueType,
				Values: tc.values,
			}
			require.Equal(t, expectedFundingEventProto, FundingEvent)
		})
	}
}
