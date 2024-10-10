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
		updateMarketPrice *types.MarketPriceUpdate
		daemonPrices      []*api.MarketPriceUpdate
	}{
		"daemon price does not exist": {
			updateMarketPrice: &types.MarketPriceUpdate{
				MarketId:  constants.MarketId0,
				SpotPrice: 11,
				PnlPrice:  11,
			},
			// Skipping price cache update, so the daemon price does not exist.
		},
		"daemon price crossing = true, old_ticks > 1, new_ticks <= sqrt(old_ticks) = false": {
			updateMarketPrice: &types.MarketPriceUpdate{
				MarketId:  constants.MarketId0,
				SpotPrice: price_5_015_000_000,
				PnlPrice:  price_5_015_000_000,
			},

			daemonPrices: []*api.MarketPriceUpdate{
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
		"daemon price crossing = true, old_ticks <= 1, new_ticks <= old_ticks = false": {
			updateMarketPrice: &types.MarketPriceUpdate{
				MarketId:  constants.MarketId0,
				SpotPrice: price_5_015_000_000,
				PnlPrice:  price_5_015_000_000,
			},
			daemonPrices: []*api.MarketPriceUpdate{
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
		"daemon price trends in the opposite direction of update price from current price": {
			updateMarketPrice: &types.MarketPriceUpdate{
				MarketId:  constants.MarketId0,
				SpotPrice: price_5_005_000_000,
				PnlPrice:  price_5_005_000_000,
			},
			daemonPrices: []*api.MarketPriceUpdate{
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
			ctx, k, _, daemonPriceCache, _, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			daemonPriceCache.UpdatePrices(tc.daemonPrices)

			// Run.
			isSpotValid, isPnlValid := k.PerformStatefulPriceUpdateValidation(ctx, tc.updateMarketPrice) // skips non-deterministic checks.

			// Validate.
			require.True(t, isSpotValid)
			require.True(t, isPnlValid)
		})
	}
}
