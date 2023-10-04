package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestNewMarketToPrice_IsEmpty(t *testing.T) {
	mtp := types.NewMarketToPrice()

	require.Empty(t, mtp.MarketToPriceTimestamp)
}

func TestUpdatePrice_Valid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice(constants.Market9_TimeT_Price1)

	require.Len(t, mtp.MarketToPriceTimestamp, 1)

	marketPriceTimestamp := mtp.GetAllPrices()[0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateValid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice(constants.Market9_TimeTMinusThreshold_Price2)
	mtp.UpdatePrice(constants.Market9_TimeT_Price1)

	require.Len(t, mtp.MarketToPriceTimestamp, 1)

	marketPriceTimestamp := mtp.GetAllPrices()[0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateInvalid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice(constants.Market9_TimeT_Price1)
	mtp.UpdatePrice(constants.Market9_TimeTMinusThreshold_Price2)

	require.Len(t, mtp.MarketToPriceTimestamp, 1)

	marketPriceTimestamp := mtp.GetAllPrices()[0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateForTwoMarketsValid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice(constants.Market9_TimeT_Price1)
	mtp.UpdatePrice(constants.Market8_TimeTMinusThreshold_Price2)

	require.Len(t, mtp.MarketToPriceTimestamp, 2)

	marketPriceTimestamp := mtp.MarketToPriceTimestamp[constants.MarketId9]
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdateTime)

	marketPriceTimestamp2 := mtp.MarketToPriceTimestamp[constants.MarketId8]
	require.Equal(t, constants.Price2, marketPriceTimestamp2.Price)
	require.Equal(t, constants.TimeTMinusThreshold, marketPriceTimestamp2.LastUpdateTime)
}

func TestGetValidPriceForMarket_Mixed(t *testing.T) {
	tests := map[string]struct {
		initialPrices []*types.MarketPriceTimestamp
		market        types.MarketId
		cutoffTime    time.Time

		expectedPrice  uint64
		expectedExists bool
	}{
		"valid price": {
			initialPrices: []*types.MarketPriceTimestamp{
				constants.Market7_TimeT_Price1,
			},
			market:     constants.MarketId7,
			cutoffTime: constants.TimeTMinus1,

			expectedPrice:  constants.Price1,
			expectedExists: true,
		},
		"invalid - stale price": {
			initialPrices: []*types.MarketPriceTimestamp{
				constants.Market7_TimeT_Price1,
			},
			market:     constants.MarketId7,
			cutoffTime: constants.TimeTPlus1,

			expectedExists: false,
		},
		"invalid - no market price": {
			initialPrices: []*types.MarketPriceTimestamp{
				constants.Market7_TimeT_Price1,
			},
			market:     constants.MarketId8, // different market from initial price
			cutoffTime: constants.TimeT,

			expectedExists: false,
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			// Setup.
			mtp := types.NewMarketToPrice()
			for _, marketPriceTimestamp := range tc.initialPrices {
				mtp.UpdatePrice(marketPriceTimestamp)
			}

			// Execute.
			price, ok := mtp.GetValidPriceForMarket(tc.market, tc.cutoffTime)

			// Assert.
			if tc.expectedExists {
				require.True(t, ok)
				require.Equal(t, tc.expectedPrice, price)
			} else {
				require.False(t, ok)
				require.Zero(t, price)
			}
		})
	}
}
