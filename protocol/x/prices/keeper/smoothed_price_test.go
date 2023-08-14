package keeper_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	emptySmoothedPrices = map[uint32]uint64{}
)

func TestUpdateSmoothedPrices(t *testing.T) {
	tests := map[string]struct {
		smoothedPrices map[uint32]uint64
		indexPrices    []*api.MarketPriceUpdate
		expectedResult map[uint32]uint64
	}{
		"Empty result - No index prices, no smoothed prices": {
			expectedResult: emptySmoothedPrices,
		},
		"Unchanged - No index prices": {
			smoothedPrices: constants.AtTimeTSingleExchangeSmoothedPrices,
			expectedResult: constants.AtTimeTSingleExchangeSmoothedPrices,
		},
		"Mixed updates and additions - mix of present and missing index prices, smoothed prices": {
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			smoothedPrices: map[uint32]uint64{
				constants.MarketId1: constants.Exchange1_Price1_TimeT.Price + 10,
				constants.MarketId2: constants.Exchange2_Price2_TimeT.Price + 50,
				constants.MarketId7: constants.Price1,
			},
			expectedResult: map[uint32]uint64{
				constants.MarketId0: constants.Exchange0_Price4_TimeT.Price,
				constants.MarketId1: constants.Exchange1_Price1_TimeT.Price + 7,
				constants.MarketId2: constants.Exchange2_Price2_TimeT.Price + 35,
				constants.MarketId7: constants.Price1,
			},
		},
		"Initializes smoothed prices with index prices": {
			indexPrices:    constants.AtTimeTSingleExchangePriceUpdate,
			expectedResult: constants.AtTimeTSingleExchangeSmoothedPrices,
		},
		"All updated - multiple existing overlapped index and smoothed prices": {
			indexPrices:    constants.AtTimeTSingleExchangePriceUpdate,
			smoothedPrices: constants.AtTimeTSingleExchangeSmoothedPricesPlus10,
			expectedResult: constants.AtTimeTSingleExchangeSmoothedPricesPlus7,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, marketToSmoothedPrices, mockTimeProvider := keepertest.PricesKeepers(t)
			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)
			for market, smoothedPrice := range tc.smoothedPrices {
				marketToSmoothedPrices.PushSmoothedPrice(market, smoothedPrice)
			}

			mockTimeProvider.On("Now").Return(constants.TimeT)

			// Run.
			err := k.UpdateSmoothedPrices(ctx)
			require.NoError(t, err)

			// Validate.
			require.Equal(t, tc.expectedResult, marketToSmoothedPrices.GetSmoothedPricesForTest())
		})
	}
}
