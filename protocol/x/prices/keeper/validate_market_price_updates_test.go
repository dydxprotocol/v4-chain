package keeper_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

const (
	price_5_015_000_000 = constants.FiveBillion + (3 * constants.FiveMillion)
	price_5_010_000_000 = constants.FiveBillion + (2 * constants.FiveMillion)
	price_5_005_000_000 = constants.FiveBillion + constants.FiveMillion
	price_5_004_999_999 = price_5_005_000_000 - 1
	price_5_000_500_000 = constants.FiveBillion + constants.OneMillion/2
	price_5_000_250_000 = constants.FiveBillion + constants.OneMillion/4
	price_4_999_750_000 = constants.FiveBillion - constants.OneMillion/4
	price_4_999_500_000 = constants.FiveBillion - constants.OneMillion/2
	price_4_995_000_001 = constants.FiveBillion - constants.FiveMillion + 1
	price_4_995_000_000 = constants.FiveBillion - constants.FiveMillion
)

func TestCrossingPriceUpdateCutoffPpm(t *testing.T) {
	require.Equal(t, uint32(500_000), keeper.CrossingPriceUpdateCutoffPpm)
}

func TestPerformStatefulPriceUpdateValidation_SkipNonDeterministicCheck_Valid(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		updateMarketPrices []*types.MarketPriceUpdates_MarketPriceUpdate
		indexPrices        []*api.MarketPriceUpdate
	}{
		"Index price does not exist": {
			updateMarketPrices: []*types.MarketPriceUpdates_MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					Price:    11,
				},
			},
			// Skipping price cache update, so the index price does not exist.
		},
		"Index price crossing = true, old_ticks > 1, new_ticks <= sqrt(old_ticks) = false": {
			updateMarketPrices: []*types.MarketPriceUpdates_MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					Price:    price_5_015_000_000,
				},
			},
			indexPrices: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeId:     constants.ExchangeId0,
							Price:          price_5_010_000_000,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			},
		},
		"Index price crossing = true, old_ticks <= 1, new_ticks <= old_ticks = false": {
			updateMarketPrices: []*types.MarketPriceUpdates_MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					Price:    price_5_015_000_000,
				},
			},
			indexPrices: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeId:     constants.ExchangeId0,
							Price:          price_5_000_250_000,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			},
		},
		"Index price trends in the opposite direction of update price from current price": {
			updateMarketPrices: []*types.MarketPriceUpdates_MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					Price:    price_5_005_000_000,
				},
			},
			indexPrices: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeId:     constants.ExchangeId0,
							Price:          constants.FiveBillion - 1,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, _, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run.
			msg := &types.MarketPriceUpdates{
				MarketPriceUpdates: tc.updateMarketPrices,
			}
			err := k.PerformStatefulPriceUpdateValidation(ctx, msg) // skips non-deterministic checks.

			// Validate.
			require.NoError(t, err)
		})
	}
}

func TestGetMarketsMissingFromPriceUpdates(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		updateMarketPrices  []*types.MarketPriceUpdates_MarketPriceUpdate
		indexPrices         []*api.MarketPriceUpdate
		smoothedIndexPrices map[uint32]uint64

		// Expected.
		expectedMarketIds []uint32
	}{
		"Empty proposed updates, Empty local updates": {
			expectedMarketIds: nil,
		},
		"Empty proposed updates, Non-empty local updates": {
			indexPrices:         constants.AtTimeTSingleExchangePriceUpdate,
			smoothedIndexPrices: constants.AtTimeTSingleExchangeSmoothedPrices,
			// The returned market ids must be sorted.
			expectedMarketIds: []uint32{constants.MarketId0, constants.MarketId1, constants.MarketId2},
		},
		"Non-empty proposed updates, Empty local updates": {
			updateMarketPrices: constants.ValidMarketPriceUpdates,
			expectedMarketIds:  nil,
		},
		"Non-empty proposed updates, Non-empty local updates, no missing markets": {
			updateMarketPrices:  constants.ValidMarketPriceUpdates,
			indexPrices:         constants.AtTimeTSingleExchangePriceUpdate,
			smoothedIndexPrices: constants.AtTimeTSingleExchangeSmoothedPrices,
			expectedMarketIds:   nil,
		},
		"Non-empty proposed updates, Non-empty local updates, single missing market": {
			updateMarketPrices: []*types.MarketPriceUpdates_MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					Price:    constants.Price5,
				},
				{
					MarketId: constants.MarketId1,
					Price:    constants.Price6,
				},
			},
			indexPrices:         constants.AtTimeTSingleExchangePriceUpdate,
			smoothedIndexPrices: constants.AtTimeTSingleExchangeSmoothedPrices,
			expectedMarketIds:   []uint32{constants.MarketId2},
		},
		"Non-empty proposed updates, Non-empty local updates, multiple missing markets, sorted": {
			updateMarketPrices: []*types.MarketPriceUpdates_MarketPriceUpdate{
				{
					MarketId: constants.MarketId1,
					Price:    constants.Price6,
				},
			},
			indexPrices:         constants.AtTimeTSingleExchangePriceUpdate,
			smoothedIndexPrices: constants.AtTimeTSingleExchangeSmoothedPrices,
			// The returned market ids must be sorted.
			expectedMarketIds: []uint32{constants.MarketId0, constants.MarketId2},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, marketToSmoothedPrices, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			for market, price := range tc.smoothedIndexPrices {
				marketToSmoothedPrices.PushSmoothedPrice(market, price)
			}
			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run.
			missingMarketIds := k.GetMarketsMissingFromPriceUpdates(ctx, tc.updateMarketPrices)

			// Validate.
			// Using `Equal` here to test for slice ordering.
			require.Equal(t, tc.expectedMarketIds, missingMarketIds)
		})
	}
}
