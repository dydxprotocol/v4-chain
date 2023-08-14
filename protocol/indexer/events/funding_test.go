package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/stretchr/testify/require"
)

func TestNewFundingEvent(t *testing.T) {
	tests := map[string]struct {
		updateType   events.FundingEventV1_Type
		updates      []events.FundingUpdateV1
		txnHash      lib.TxHash
		newEventFunc func(updates []events.FundingUpdateV1) *events.FundingEventV1
	}{
		"premium samples": {
			updateType: events.FundingEventV1_TYPE_PREMIUM_SAMPLE,
			updates: []events.FundingUpdateV1{
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
			updateType: events.FundingEventV1_TYPE_FUNDING_RATE_AND_INDEX,
			updates: []events.FundingUpdateV1{
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
			expectedFundingEventProto := &events.FundingEventV1{
				Type:    tc.updateType,
				Updates: tc.updates,
			}
			require.Equal(t, expectedFundingEventProto, FundingEvent)
		})
	}
}
