package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	fiveBillionAndFiveMillion = constants.FiveBillion + constants.FiveMillion

	testPriceValidUpdate               = fiveBillionAndFiveMillion
	testPriceDoesNotMeetMinPriceChange = constants.FiveBillion + 2
)

var (
	testMarketParamPrice = types.MarketParamPrice{
		Param: constants.TestMarketParams[0], // minPriceChangePpm of 50 - need 5 million to meet min change.
		Price: constants.TestMarketPrices[0], // Price initialized to 5 billion.
	}
)

func TestShouldProposePrice(t *testing.T) {
	tests := map[string]struct {
		proposalPrice       uint64
		expectShouldPropose bool
		expectReasons       []proposeCancellationReason
	}{
		"Should not propose: proposal price does not meet min price change": {
			proposalPrice:       testPriceDoesNotMeetMinPriceChange,
			expectShouldPropose: false,
			expectReasons: []proposeCancellationReason{
				{
					Reason: metrics.ProposedPriceDoesNotMeetMinPriceChange,
					Value:  true,
				},
			},
		},
		"Should propose": {
			proposalPrice:       testPriceValidUpdate,
			expectShouldPropose: true,
			expectReasons: []proposeCancellationReason{
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
			)
			require.Equal(t, tc.expectShouldPropose, actualShouldPropose)
			require.Equal(t, tc.expectReasons, actualReasons)
		})
	}
}
