package keeper_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

const (
	fiveBillionAndFiveMillion         = constants.FiveBillion + constants.FiveMillion
	fiveBillionMinusFiveMillionAndOne = constants.FiveBillion - constants.FiveMillion - 1
)

var (
	emptyResult = []*types.MarketSpotPriceUpdate{}

	validMarket0UpdateResult = []*types.MarketSpotPriceUpdate{
		{
			MarketId:  constants.MarketId0,
			SpotPrice: fiveBillionAndFiveMillion,
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

	invalidMarket0SmoothedPriceTrendsAwayFromDaemonPrice = map[uint32][]uint64{
		constants.MarketId0: {fiveBillionMinusFiveMillionAndOne},
	}

	invalidMarket2SmoothedPriceTrendsAwayFromDaemonPrice = map[uint32][]uint64{
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
		daemonPrices []*api.MarketPriceUpdate
		// historicalSmoothedDaemonPrice prices for each market are expected to be ordered from most recent to least
		// recent.
		historicalSmoothedDaemonPrices map[uint32][]uint64
		skipCreateMarketsAndExchanges  bool

		// Expected.
		expectedMsg []*types.MarketSpotPriceUpdate
	}{
		"Empty result: no markets": {
			skipCreateMarketsAndExchanges: true,
			expectedMsg:                   emptyResult,
		},
		"Empty result: no daemon prices": {
			daemonPrices: []*api.MarketPriceUpdate{},
			expectedMsg:  emptyResult,
		},
		"Empty result: price is zero": {
			daemonPrices: []*api.MarketPriceUpdate{invalidMarket1PriceIsZeroUpdate},
			expectedMsg:  emptyResult,
		},
		"Empty result: no overlap between markets and daemon prices": {
			daemonPrices: []*api.MarketPriceUpdate{invalidMarket9DoesNotExistUpdate},
			expectedMsg:  emptyResult,
		},
		"Single result: daemon price used when no smoothed prices": {
			daemonPrices: []*api.MarketPriceUpdate{validMarket0Update},
			expectedMsg:  validMarket0UpdateResult,
		},
		"Single result: no overlap between markets for daemon prices and smoothed prices": {
			daemonPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedDaemonPrices: invalidMarket9DoesNotExistSmoothedPrice,
			expectedMsg:                    validMarket0UpdateResult,
		},
		"Empty result: propose price is daemon price, does not meet min price change": {
			daemonPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedDaemonPrices: market2SmoothedPriceNotProposed,
			expectedMsg:                    emptyResult,
		},
		"Empty result: propose price is smoothed price, does not meet min price change": {
			daemonPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedDaemonPrices: market2SmoothedPriceDoesNotMeetMinChangeUpdate,
			expectedMsg:                    emptyResult,
		},
		"Empty result: propose price is good, but historical smoothed price does not meet min price change": {
			daemonPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedDaemonPrices: invalidMarket0HistoricalSmoothedPricesDoesNotMeetMinPriceChange,
			expectedMsg:                    emptyResult,
		},
		"Empty result: propose price is good, but historical smoothed price crosses oracle price": {
			daemonPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedDaemonPrices: invalidMarket0HistoricalSmoothedPricesCrossesOraclePrice,
			expectedMsg:                    emptyResult,
		},
		"Empty result: proposed price is smoothed price, meets min change but trends away from daemon price": {
			daemonPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedDaemonPrices: invalidMarket0SmoothedPriceTrendsAwayFromDaemonPrice,
			expectedMsg:                    emptyResult,
		},
		"Empty result: proposed price does not meet min change and historical smoothed price crosses oracle price": {
			daemonPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedDaemonPrices: invalidMarket2HistoricalSmoothedPricesCrossesOraclePrice,
			expectedMsg:                    emptyResult,
		},
		"Empty result: proposed price does not meet min change and smoothed price is trending away from daemon price": {
			daemonPrices:                   []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			historicalSmoothedDaemonPrices: invalidMarket2SmoothedPriceTrendsAwayFromDaemonPrice,
			expectedMsg:                    emptyResult,
		},
		"Single market price update": {
			daemonPrices:                   []*api.MarketPriceUpdate{validMarket0Update},
			historicalSmoothedDaemonPrices: validMarket0SmoothedPrices,
			expectedMsg:                    validMarket0UpdateResult,
		},
		"Multiple market price updates, some from smoothed price and some from daemon price": {
			daemonPrices: constants.AtTimeTSingleExchangePriceUpdate,
			historicalSmoothedDaemonPrices: map[uint32][]uint64{
				constants.MarketId0: {constants.Price4 - 1},
				constants.MarketId1: {constants.Price1 + 1},
				constants.MarketId2: {constants.Price2},
			},
			expectedMsg: []*types.MarketSpotPriceUpdate{
				types.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price4),
				types.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price1+1),
				types.NewMarketSpotPriceUpdate(constants.MarketId2, constants.Price2),
				types.NewMarketSpotPriceUpdate(constants.MarketId3, constants.Price3),
				types.NewMarketSpotPriceUpdate(constants.MarketId4, constants.Price3),
			},
		},
		"Mix of valid and invalid daemon prices": {
			daemonPrices: []*api.MarketPriceUpdate{
				validMarket0Update,
				invalidMarket1PriceIsZeroUpdate,               // Price cannot be 0.
				invalidMarket2PriceDoesNotMeetMinChangeUpdate, // Price does not meet min price change req.
				invalidMarket9DoesNotExistUpdate,              // Market with id 9 does not exist.
			},
			historicalSmoothedDaemonPrices: map[uint32][]uint64{
				constants.MarketId0: {validMarket0Update.ExchangePrices[0].Price},
				constants.MarketId1: {constants.Price4},
				constants.MarketId2: {constants.Price2},
				constants.MarketId9: {constants.Price4},
			},
			expectedMsg: validMarket0UpdateResult,
		},
		"Mix of valid, invalid, and missing smoothed prices": {
			daemonPrices: constants.AtTimeTSingleExchangePriceUpdate,
			historicalSmoothedDaemonPrices: map[uint32][]uint64{
				constants.MarketId0: {constants.Price4}, // Same as daemon price.
				constants.MarketId1: {0},                // Invalid price, so daemon price is used.
				constants.MarketId9: {constants.Price1}, // Invalid market.
			},
			expectedMsg: []*types.MarketSpotPriceUpdate{
				types.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price4),
				types.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price1),
				types.NewMarketSpotPriceUpdate(constants.MarketId2, constants.Price2),
				types.NewMarketSpotPriceUpdate(constants.MarketId3, constants.Price3),
				types.NewMarketSpotPriceUpdate(constants.MarketId4, constants.Price3),
			},
		},
		"Mix of valid, invalid, and invalid historical smoothed prices": {
			daemonPrices: constants.AtTimeTSingleExchangePriceUpdate,
			historicalSmoothedDaemonPrices: map[uint32][]uint64{
				constants.MarketId0: {
					constants.Price4,          // Same as daemon price.
					fiveBillionAndFiveMillion, // Invalid: crosses oracle price.
				},
				constants.MarketId1: {constants.Price1}, // Valid: same as daemon price.
				constants.MarketId9: {constants.Price1}, // Invalid market.
			},
			expectedMsg: []*types.MarketSpotPriceUpdate{
				types.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price1),
				types.NewMarketSpotPriceUpdate(constants.MarketId2, constants.Price2),
				types.NewMarketSpotPriceUpdate(constants.MarketId3, constants.Price3),
				types.NewMarketSpotPriceUpdate(constants.MarketId4, constants.Price3),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, daemonPriceCache, marketSmoothedPrices, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			if !tc.skipCreateMarketsAndExchanges {
				keepertest.CreateTestMarkets(t, ctx, k)
			}
			daemonPriceCache.UpdatePrices(tc.daemonPrices)

			// Smoothed prices are listed in reverse chronological order for test case constant legibility.
			// Therefore, add them in reverse order to the `marketSmoothedPrices` cache.
			for market, historicalSmoothedPrices := range tc.historicalSmoothedDaemonPrices {
				for i := len(historicalSmoothedPrices) - 1; i >= 0; i-- {
					marketSmoothedPrices.PushSmoothedSpotPrice(market, historicalSmoothedPrices[i])
				}
			}

			// Run.
			result := k.GetValidMarketSpotPriceUpdates(ctx)

			// Validate.
			require.Equal(t, tc.expectedMsg, result)
			// TODO(DEC-532): validate on either metrics or logging.
			// Validating metrics might be difficult because it's hard to mock `telemetry`.
			// Alternatively, we can add mock logging in `ctx`.
		})
	}
}
