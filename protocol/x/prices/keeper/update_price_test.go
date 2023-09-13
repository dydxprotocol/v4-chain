package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

const (
	fiveBillionAndFiveMillion         = constants.FiveBillion + constants.FiveMillion
	fiveBillionMinusFiveMillionAndOne = constants.FiveBillion - constants.FiveMillion - 1
)

var (
	// MsgUpdateMarketPrices test constants.
	emptyResult = &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{},
	}

	validMarket0UpdateResult = &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
			types.NewMarketPriceUpdate(constants.MarketId0, fiveBillionAndFiveMillion),
		},
	}

	// MarketPriceUpdate test constants.
	validMarket0Update = &api.MarketPriceUpdate{
		MarketId: constants.MarketId0,
		ExchangePrices: []*api.ExchangePrice{
			{
				ExchangeId:     constants.ExchangeId0,
				Price:          fiveBillionAndFiveMillion,
				LastUpdateTime: &constants.TimeT,
			},
		},
	}

	invalidMarket1PriceIsZeroUpdate = &api.MarketPriceUpdate{
		MarketId: constants.MarketId1,
		ExchangePrices: []*api.ExchangePrice{
			{
				ExchangeId:     constants.ExchangeId1,
				Price:          0,
				LastUpdateTime: &constants.TimeT,
			},
		},
	}

	invalidMarket2PriceDoesNotMeetMinChangeUpdate = &api.MarketPriceUpdate{
		MarketId: constants.MarketId2,
		ExchangePrices: []*api.ExchangePrice{
			{
				ExchangeId:     constants.ExchangeId2,
				Price:          constants.FiveBillion + 2, // 5,000,000,002
				LastUpdateTime: &constants.TimeT,
			},
		},
	}

	invalidMarket9DoesNotExistUpdate = constants.Market9_SingleExchange_AtTimeUpdate[0]

	// SmoothedPrice test constants.
	validMarket0SmoothedPrices = map[uint32][]uint64{
		constants.MarketId0: {fiveBillionAndFiveMillion + 1},
	}

	invalidMarket0HistoricalSmoothedPricesCrossesOraclePrice = map[uint32][]uint64{
		constants.MarketId0: {
			fiveBillionAndFiveMillion + 1,     // Valid
			fiveBillionMinusFiveMillionAndOne, // Invalid: crosses oracle price.
		},
	}

	invalidMarket0HistoricalSmoothedPricesDoesNotMeetMinPriceChange = map[uint32][]uint64{
		constants.MarketId0: {
			fiveBillionAndFiveMillion + 1, // Valid
			constants.FiveBillion + 1,     // Invalid: does not meet min price change.
		},
	}

	invalidMarket2HistoricalSmoothedPricesCrossesOraclePrice = map[uint32][]uint64{
		constants.MarketId2: {
			fiveBillionAndFiveMillion + 1,     // Valid
			fiveBillionMinusFiveMillionAndOne, // Invalid: crosses oracle price.
		},
	}

	invalidMarket0SmoothedPriceTrendsAwayFromIndexPrice = map[uint32][]uint64{
		constants.MarketId0: {fiveBillionMinusFiveMillionAndOne},
	}

	invalidMarket2SmoothedPriceTrendsAwayFromIndexPrice = map[uint32][]uint64{
		constants.MarketId2: {constants.FiveBillion - 2},
	}

	market2SmoothedPriceNotProposed = map[uint32][]uint64{
		constants.MarketId2: {constants.FiveBillion + 3},
	}

	market2SmoothedPriceDoesNotMeetMinChangeUpdate = map[uint32][]uint64{
		constants.MarketId2: {constants.FiveBillion + 1},
	}

	invalidMarket9DoesNotExistSmoothedPrice = map[uint32][]uint64{
		constants.MarketId9: {1_000_000},
	}
)

// Note: markets and exchanges are created by `CreateTestMarketsAndExchanges`.
func TestGetValidMarketPriceUpdates(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		indexPrices []*api.MarketPriceUpdate
		// historicalSmoothedIndexPrice prices for each market are expected to be ordered from most recent to least
		// recent.
		historicalSmoothedIndexPrices map[uint32][]uint64
		skipCreateMarketsAndExchanges bool

		// Expected.
		expectedMsg *types.MsgUpdateMarketPrices
	}{
		"Empty result: no markets": {
			skipCreateMarketsAndExchanges: true,
			expectedMsg:                   emptyResult,
		},
		"Empty result: no index prices": {
			indexPrices: []*api.MarketPriceUpdate{},
			expectedMsg: emptyResult,
		},
		"Empty result: price is zero": {
			indexPrices: []*api.MarketPriceUpdate{invalidMarket1PriceIsZeroUpdate},
			expectedMsg: emptyResult,
		},
		"Empty result: no overlap between markets and index prices": {
			indexPrices: []*api.MarketPriceUpdate{invalidMarket9DoesNotExistUpdate},
			expectedMsg: emptyResult,
		},
		"Single result: index price used when no smoothed prices": {
			indexPrices: []*api.MarketPriceUpdate{validMarket0Update},
			expectedMsg: validMarket0UpdateResult,
		},
		"Single result: no overlap between markets for index prices and smoothed prices": {
			indexPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedIndexPrices: invalidMarket9DoesNotExistSmoothedPrice,
			expectedMsg:                   validMarket0UpdateResult,
		},
		"Empty result: propose price is index price, does not meet min price change": {
			indexPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedIndexPrices: market2SmoothedPriceNotProposed,
			expectedMsg:                   emptyResult,
		},
		"Empty result: propose price is smoothed price, does not meet min price change": {
			indexPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedIndexPrices: market2SmoothedPriceDoesNotMeetMinChangeUpdate,
			expectedMsg:                   emptyResult,
		},
		"Empty result: propose price is good, but historical smoothed price does not meet min price change": {
			indexPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedIndexPrices: invalidMarket0HistoricalSmoothedPricesDoesNotMeetMinPriceChange,
			expectedMsg:                   emptyResult,
		},
		"Empty result: propose price is good, but historical smoothed price crosses oracle price": {
			indexPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedIndexPrices: invalidMarket0HistoricalSmoothedPricesCrossesOraclePrice,
			expectedMsg:                   emptyResult,
		},
		"Empty result: proposed price is smoothed price, meets min change but trends away from index price": {
			indexPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedIndexPrices: invalidMarket0SmoothedPriceTrendsAwayFromIndexPrice,
			expectedMsg:                   emptyResult,
		},
		"Empty result: proposed price does not meet min change and historical smoothed price crosses oracle price": {
			indexPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedIndexPrices: invalidMarket2HistoricalSmoothedPricesCrossesOraclePrice,
			expectedMsg:                   emptyResult,
		},
		"Empty result: proposed price does not meet min change and smoothed price is trending away from index price": {
			indexPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedIndexPrices: invalidMarket2SmoothedPriceTrendsAwayFromIndexPrice,
			expectedMsg:                   emptyResult,
		},
		"Single market price update": {
			indexPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedIndexPrices: validMarket0SmoothedPrices,
			expectedMsg:                   validMarket0UpdateResult,
		},
		"Multiple market price updates, some from smoothed price and some from index price": {
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			historicalSmoothedIndexPrices: map[uint32][]uint64{
				constants.MarketId0: {constants.Price4 - 1},
				constants.MarketId1: {constants.Price1 + 1},
				constants.MarketId2: {constants.Price2},
			},
			expectedMsg: &types.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
					types.NewMarketPriceUpdate(constants.MarketId0, constants.Price4),
					types.NewMarketPriceUpdate(constants.MarketId1, constants.Price1+1),
					types.NewMarketPriceUpdate(constants.MarketId2, constants.Price2),
				},
			},
		},
		"Mix of valid and invalid index prices": {
			indexPrices: []*api.MarketPriceUpdate{
				validMarket0Update,
				invalidMarket1PriceIsZeroUpdate,               // Price cannot be 0.
				invalidMarket2PriceDoesNotMeetMinChangeUpdate, // Price does not meet min price change req.
				invalidMarket9DoesNotExistUpdate,              // Market with id 9 does not exist.
			},
			historicalSmoothedIndexPrices: map[uint32][]uint64{
				constants.MarketId0: {validMarket0Update.ExchangePrices[0].Price},
				constants.MarketId1: {constants.Price4},
				constants.MarketId2: {constants.Price2},
				constants.MarketId9: {constants.Price4},
			},
			expectedMsg: validMarket0UpdateResult,
		},
		"Mix of valid, invalid, and missing smoothed prices": {
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			historicalSmoothedIndexPrices: map[uint32][]uint64{
				constants.MarketId0: {constants.Price4}, // Same as index price.
				constants.MarketId1: {0},                // Invalid price, so index price is used.
				constants.MarketId9: {constants.Price1}, // Invalid market.
			},
			expectedMsg: &types.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
					types.NewMarketPriceUpdate(constants.MarketId0, constants.Price4),
					types.NewMarketPriceUpdate(constants.MarketId1, constants.Price1),
					types.NewMarketPriceUpdate(constants.MarketId2, constants.Price2),
				},
			},
		},
		"Mix of valid, invalid, and invalid historical smoothed prices": {
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			historicalSmoothedIndexPrices: map[uint32][]uint64{
				constants.MarketId0: {
					constants.Price4,          // Same as index price.
					fiveBillionAndFiveMillion, // Invalid: crosses oracle price.
				},
				constants.MarketId1: {constants.Price1}, // Valid: same as index price.
				constants.MarketId9: {constants.Price1}, // Invalid market.
			},
			expectedMsg: &types.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
					types.NewMarketPriceUpdate(constants.MarketId1, constants.Price1),
					types.NewMarketPriceUpdate(constants.MarketId2, constants.Price2),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, marketSmoothedPrices, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			if !tc.skipCreateMarketsAndExchanges {
				keepertest.CreateTestMarkets(t, ctx, k)
			}
			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Smoothed prices are listed in reverse chronological order for test case constant legibility.
			// Therefore, add them in reverse order to the `marketSmoothedPrices` cache.
			for market, historicalSmoothedPrices := range tc.historicalSmoothedIndexPrices {
				for i := len(historicalSmoothedPrices) - 1; i >= 0; i-- {
					marketSmoothedPrices.PushSmoothedPrice(market, historicalSmoothedPrices[i])
				}
			}

			// Run.
			result := k.GetValidMarketPriceUpdates(ctx)

			// Validate.
			require.Equal(t, tc.expectedMsg, result)
			// TODO(DEC-532): validate on either metrics or logging.
			// Validating metrics might be difficult because it's hard to mock `telemetry`.
			// Alternatively, we can add mock logging in `ctx`.
		})
	}
}
