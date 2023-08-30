package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	fiveBillionAndFiveMillion         = constants.FiveBillion + constants.FiveMillion
	fiveBillionMinusFiveMillionAndOne = constants.FiveBillion - constants.FiveMillion - 1
	fiveBillionAndTenMillion          = constants.FiveBillion + 2*constants.FiveMillion

	testPriceValidUpdate                    = fiveBillionAndFiveMillion
	testPriceLargeValidUpdate               = fiveBillionAndTenMillion
	testPriceDoesNotMeetMinPriceChange      = constants.FiveBillion + 2
	testPriceCrossesOraclePrice             = fiveBillionMinusFiveMillionAndOne
	testPriceCrossesAndDoesNotMeetMinChange = constants.FiveBillion - 1
)

var (
	testMarketParamPrice = types.MarketParamPrice{
		Param: constants.TestMarketParams[0], // minPriceChangePpm of 50 - need 5 million to meet min change.
		Price: constants.TestMarketPrices[0], // Price initialized to 5 billion.
	}
)

func TestShouldProposePrice(t *testing.T) {
	tests := map[string]struct {
		proposalPrice            uint64
		indexPrice               uint64
		historicalSmoothedPrices []uint64
		expectShouldPropose      bool
		expectReasons            []proposeCancellationReason
	}{
		"Should not propose: proposal price is smoothed price, crosses index price": {
			proposalPrice: testPriceCrossesOraclePrice,
			indexPrice:    testPriceLargeValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceCrossesOraclePrice,
				testPriceValidUpdate,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				// These are both true because the proposed price is the most recent smoothed price.
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  true,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
		"Should not propose: proposal price is smoothed price, does not meet min price change": {
			proposalPrice: testPriceDoesNotMeetMinPriceChange,
			indexPrice:    testPriceLargeValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceDoesNotMeetMinPriceChange,
				testPriceValidUpdate,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
			},
		},
		"Should not propose: proposal price is index price, does not meet min price change": {
			proposalPrice: testPriceDoesNotMeetMinPriceChange,
			indexPrice:    testPriceDoesNotMeetMinPriceChange,
			historicalSmoothedPrices: []uint64{
				testPriceLargeValidUpdate,
				testPriceValidUpdate,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
			},
		},
		"Should not propose: a historical smoothed price crosses index price": {
			proposalPrice: testPriceValidUpdate,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceDoesNotMeetMinPriceChange,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
		"Should not propose: multiple historical smoothed prices cross index price": {
			proposalPrice: testPriceValidUpdate,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceCrossesOraclePrice,
				testPriceCrossesOraclePrice,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
		"Should not propose: a historical smoothed price does not meet min price change": {
			proposalPrice: testPriceValidUpdate,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceDoesNotMeetMinPriceChange,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
		"Should not propose: multiple historical smoothed prices do not meet min price change": {
			proposalPrice: testPriceValidUpdate,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceDoesNotMeetMinPriceChange,
				testPriceDoesNotMeetMinPriceChange,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
		"Should not propose: historical smoothed price crosses and does not meet min price change": {
			proposalPrice: testPriceValidUpdate,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceCrossesAndDoesNotMeetMinChange,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
		"Should not propose: proposal price crosses and does not meet min price change": {
			proposalPrice: testPriceCrossesAndDoesNotMeetMinChange,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceLargeValidUpdate,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  true,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
			},
		},
		"Should not propose: multiple historical smoothed prices issues": {
			proposalPrice: testPriceValidUpdate,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceDoesNotMeetMinPriceChange,
				testPriceCrossesOraclePrice,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
		"Should not propose: multiple issues": {
			proposalPrice: testPriceDoesNotMeetMinPriceChange,
			indexPrice:    testPriceValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceDoesNotMeetMinPriceChange,
				testPriceCrossesOraclePrice,
			},
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
			},
		},
		"Should propose": {
			proposalPrice: testPriceValidUpdate,
			indexPrice:    testPriceLargeValidUpdate,
			historicalSmoothedPrices: []uint64{
				testPriceValidUpdate,
				testPriceLargeValidUpdate,
				testPriceValidUpdate,
			},
			expectShouldPropose: true,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.RecentSmoothedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceCrossesOraclePrice,
					Value:  false,
				},
				{
					Reason: metrics.RecentSmoothedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  false,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualShouldPropose, actualReasons := shouldProposePrice(
				tc.proposalPrice,
				testMarketParamPrice,
				tc.indexPrice,
				tc.historicalSmoothedPrices,
			)
			require.Equal(t, tc.expectShouldPropose, actualShouldPropose)
			require.Equal(t, tc.expectReasons, actualReasons)
		})
	}
}
