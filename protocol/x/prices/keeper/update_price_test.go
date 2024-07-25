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
)

// Note: markets and exchanges are created by `CreateTestMarketsAndExchanges`.
func TestGetValidMarketPriceUpdates(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		indexPrices                   []*api.MarketPriceUpdate
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
		"Empty result: propose price does not meet min price change": {
			indexPrices: []*api.MarketPriceUpdate{invalidMarket2PriceDoesNotMeetMinChangeUpdate},
			expectedMsg: emptyResult,
		},
		"Single market price update": {
			indexPrices: []*api.MarketPriceUpdate{validMarket0Update},
			expectedMsg: validMarket0UpdateResult,
		},
		"Mix of valid and invalid index prices": {
			indexPrices: []*api.MarketPriceUpdate{
				validMarket0Update,
				invalidMarket1PriceIsZeroUpdate,               // Price cannot be 0.
				invalidMarket2PriceDoesNotMeetMinChangeUpdate, // Price does not meet min price change req.
				invalidMarket9DoesNotExistUpdate,              // Market with id 9 does not exist.
			},
			expectedMsg: validMarket0UpdateResult,
		},
		"Mix of valid, invalid, and missing smoothed prices": {
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			expectedMsg: &types.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
					types.NewMarketPriceUpdate(constants.MarketId0, constants.Price4),
					types.NewMarketPriceUpdate(constants.MarketId1, constants.Price1),
					types.NewMarketPriceUpdate(constants.MarketId2, constants.Price2),
					types.NewMarketPriceUpdate(constants.MarketId3, constants.Price3),
					types.NewMarketPriceUpdate(constants.MarketId4, constants.Price3),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			if !tc.skipCreateMarketsAndExchanges {
				keepertest.CreateTestMarkets(t, ctx, k)
			}
			indexPriceCache.UpdatePrices(tc.indexPrices)

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
