package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

const (
	fiveBillionAndFiveMillion = constants.FiveBillion + constants.FiveMillion
)

var (
	testAddress    = constants.AliceAccAddress
	testAddressStr = testAddress.String()
	emptyResult    = &types.MsgUpdateMarketPrices{
		Proposer:           testAddressStr,
		MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{},
	}

	validMarket0Update = &api.MarketPriceUpdate{
		MarketId: constants.MarketId0,
		ExchangePrices: []*api.ExchangePrice{
			{
				ExchangeFeedId: constants.ExchangeFeedId0,
				Price:          fiveBillionAndFiveMillion,
				LastUpdateTime: &constants.TimeT,
			},
		},
	}

	invalidMarket1PriceIsZeroUpdate = &api.MarketPriceUpdate{
		MarketId: constants.MarketId1,
		ExchangePrices: []*api.ExchangePrice{
			{
				ExchangeFeedId: constants.ExchangeFeedId1,
				Price:          0,
				LastUpdateTime: &constants.TimeT,
			},
		},
	}

	invalidMarket2PriceDoesNotMeetMinChangeUpdate = &api.MarketPriceUpdate{
		MarketId: constants.MarketId2,
		ExchangePrices: []*api.ExchangePrice{
			{
				ExchangeFeedId: constants.ExchangeFeedId2,
				Price:          constants.FiveBillion + 1, // 5,000,000,001
				LastUpdateTime: &constants.TimeT,
			},
		},
	}

	invalidMarket9DoesNotExistUpdate = constants.Market9_SingleExchange_AtTimeUpdate[0]
)

// Note: markets and exchanges are created by `CreateTestMarketsAndExchangeFeeds`.
func TestGetValidMarketPriceUpdates(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		indexPrices                   []*api.MarketPriceUpdate
		skipCreateMarketsAndExchanges bool

		// Expected.
		expectedMsg *types.MsgUpdateMarketPrices
	}{
		"Single market price update": {
			indexPrices: []*api.MarketPriceUpdate{validMarket0Update},
			expectedMsg: &types.MsgUpdateMarketPrices{
				Proposer: testAddressStr,
				MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
					types.NewMarketPriceUpdate(0, fiveBillionAndFiveMillion),
				},
			},
		},
		"Multiple market price updates": {
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			expectedMsg: &types.MsgUpdateMarketPrices{
				Proposer: testAddressStr,
				MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
					types.NewMarketPriceUpdate(constants.MarketId0, constants.Price4),
					types.NewMarketPriceUpdate(constants.MarketId1, constants.Price1),
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
			expectedMsg: &types.MsgUpdateMarketPrices{
				Proposer: testAddressStr,
				MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
					types.NewMarketPriceUpdate(0, fiveBillionAndFiveMillion),
				},
			},
		},
		"Empty result: no markets": {
			skipCreateMarketsAndExchanges: true,
			expectedMsg:                   emptyResult,
		},
		"Empty result: no index prices": {
			indexPrices: []*api.MarketPriceUpdate{},
			expectedMsg: emptyResult,
		},
		"Empty result: no overlap between markets and index prices": {
			indexPrices: []*api.MarketPriceUpdate{invalidMarket9DoesNotExistUpdate},
			expectedMsg: emptyResult,
		},
		"Empty result: index price does not meet min price change": {
			indexPrices: []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			expectedMsg: emptyResult,
		},
		"Empty result: price is zero": {
			indexPrices: []*api.MarketPriceUpdate{invalidMarket1PriceIsZeroUpdate},
			expectedMsg: emptyResult,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, mockTimeProvider := keepertest.PricesKeepers(t)
			if !tc.skipCreateMarketsAndExchanges {
				keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, k)
			}
			indexPriceCache.UpdatePrices(tc.indexPrices)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			// Run.
			result := k.GetValidMarketPriceUpdates(ctx, testAddress)

			// Validate.
			require.Equal(t, tc.expectedMsg, result)
			// TODO(DEC-532): validate on either metrics or logging.
			// Validating metrics might be difficult because it's hard to mock `telemetry`.
			// Alternatively, we can add mock logging in `ctx`.
		})
	}
}
