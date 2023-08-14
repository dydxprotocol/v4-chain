package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	marketId1                      = 1
	price1                         = uint64(9_000_000_000)
	testSmoothedPriceHistoryLength = 5
)

func TestNewMarketToSmoothedPrices_IsEmpty(t *testing.T) {
	mtsp := types.NewMarketToSmoothedPrices(types.SmoothedPriceTrackingBlockHistoryLength)
	require.Empty(t, mtsp.GetSmoothedPricesForTest())
}

func TestSetSmoothedPrice(t *testing.T) {
	mtsp := types.NewMarketToSmoothedPrices(types.SmoothedPriceTrackingBlockHistoryLength)

	mtsp.PushSmoothedPrice(marketId1, price1)
	actualPrice, ok := mtsp.GetSmoothedPrice(marketId1)

	require.True(t, ok)
	require.Equal(t, actualPrice, price1)

	// Set a new price for the same market enough times to cause the ring buffer to loop. This is to sanity
	// check ring buffer logic.
	for i := 0; i < int(types.SmoothedPriceTrackingBlockHistoryLength); i++ {
		updatePrice := price1 + uint64(i+1)*uint64(1_000_000_000)

		mtsp.PushSmoothedPrice(marketId1, updatePrice)
		actualPrice, ok := mtsp.GetSmoothedPrice(marketId1)

		require.True(t, ok)
		require.Equal(t, actualPrice, updatePrice)
	}
}

func TestGetHistoricalSmoothedPrices(t *testing.T) {
	tests := map[string]struct {
		prices         []uint64
		expectedPrices []uint64
	}{
		"no prices": {
			prices:         []uint64{},
			expectedPrices: []uint64{},
		},
		"one price": {
			prices:         []uint64{price1},
			expectedPrices: []uint64{price1},
		},
		"two prices": {
			prices:         []uint64{price1, price1 + 1},
			expectedPrices: []uint64{price1 + 1, price1},
		},
		"num prices >> SmoothedPriceTrackingBlockHistoryLength": {
			prices:         []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			expectedPrices: []uint64{12, 11, 10, 9, 8},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mtsp := types.NewMarketToSmoothedPrices(testSmoothedPriceHistoryLength)
			for _, price := range tc.prices {
				mtsp.PushSmoothedPrice(marketId1, price)
			}
			require.Equal(t, tc.expectedPrices, mtsp.GetHistoricalSmoothedPrices(marketId1))
		})
	}
}
