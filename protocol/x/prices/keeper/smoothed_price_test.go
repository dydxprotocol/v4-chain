package keeper_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	emptySmoothedPrices = map[uint32]uint64{}
)

func errInterpolator(v0 uint64, v1 uint64, ppm uint32) (uint64, error) {
	return 0, fmt.Errorf("error while interpolating")
}

func alternatingErrInterpolator() func(v0 uint64, v1 uint64, ppm uint32) (uint64, error) {
	var i int
	return func(v0 uint64, v1 uint64, ppm uint32) (uint64, error) {
		i++
		if i%2 == 0 {
			return 0, fmt.Errorf("error while interpolating")
		}
		return lib.Uint64LinearInterpolate(v0, v1, ppm)
	}
}

func TestUpdateSmoothedPrices(t *testing.T) {
	tests := map[string]struct {
		smoothedPrices        map[uint32]uint64
		indexPrices           []*api.MarketPriceUpdate
		expectedResult        map[uint32]uint64
		linearInterpolateFunc func(v0 uint64, v1 uint64, ppm uint32) (uint64, error)
		expectedErr           string
	}{
		"Empty result - No index prices, no smoothed prices": {
			expectedResult:        emptySmoothedPrices,
			linearInterpolateFunc: lib.Uint64LinearInterpolate,
		},
		"Unchanged - No index prices": {
			smoothedPrices:        constants.AtTimeTSingleExchangeSmoothedPrices,
			expectedResult:        constants.AtTimeTSingleExchangeSmoothedPrices,
			linearInterpolateFunc: lib.Uint64LinearInterpolate,
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
			linearInterpolateFunc: lib.Uint64LinearInterpolate,
		},
		"Initializes smoothed prices with index prices": {
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
			expectedResult:        constants.AtTimeTSingleExchangeSmoothedPrices,
			linearInterpolateFunc: lib.Uint64LinearInterpolate,
		},
		"All updated - multiple existing overlapped index and smoothed prices": {
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
			smoothedPrices:        constants.AtTimeTSingleExchangeSmoothedPricesPlus10,
			expectedResult:        constants.AtTimeTSingleExchangeSmoothedPricesPlus7,
			linearInterpolateFunc: lib.Uint64LinearInterpolate,
		},
		"Interpolation errors - returns error": {
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
			smoothedPrices:        constants.AtTimeTSingleExchangeSmoothedPricesPlus10,
			linearInterpolateFunc: errInterpolator,
			expectedErr: "Error updating smoothed price for market 0: error while interpolating\n" +
				"Error updating smoothed price for market 1: error while interpolating\n" +
				"Error updating smoothed price for market 2: error while interpolating",
			expectedResult: constants.AtTimeTSingleExchangeSmoothedPricesPlus10, // no change
		},
		"Single interpolation error - returns error, continues updating other markets": {
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
			smoothedPrices:        constants.AtTimeTSingleExchangeSmoothedPricesPlus10,
			linearInterpolateFunc: alternatingErrInterpolator(),
			expectedErr:           "Error updating smoothed price for market 1: error while interpolating",
			expectedResult: map[uint32]uint64{
				constants.MarketId0: constants.AtTimeTSingleExchangeSmoothedPricesPlus7[constants.MarketId0],  // update
				constants.MarketId1: constants.AtTimeTSingleExchangeSmoothedPricesPlus10[constants.MarketId1], // no change
				constants.MarketId2: constants.AtTimeTSingleExchangeSmoothedPricesPlus7[constants.MarketId2],  // update
			}, // no change
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, marketToSmoothedPrices, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)
			for market, smoothedPrice := range tc.smoothedPrices {
				marketToSmoothedPrices.PushSmoothedPrice(market, smoothedPrice)
			}

			// Run.
			err := k.UpdateSmoothedPrices(ctx, tc.linearInterpolateFunc)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			// Validate.
			require.Equal(t, tc.expectedResult, marketToSmoothedPrices.GetSmoothedPricesForTest())
		})
	}
}
