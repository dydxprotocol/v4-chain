package keeper_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestCrossingPriceUpdateCutoffPpm(t *testing.T) {
	require.Equal(t, uint32(500_000), keeper.CrossingPriceUpdateCutoffPpm)
}

// Note: the current market prices (i.e. 5 billion) are set in `CreateTestMarketsAndExchanges`.
func TestPerformStatefulPriceUpdateValidation_Valid(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		msgUpdateMarketPrices []*types.MsgUpdateMarketPrices_MarketPrice
		indexPrices           []*api.MarketPriceUpdate
	}{
		"Empty updates": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{},
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
		},
		"Multiple updates": {
			msgUpdateMarketPrices: constants.ValidMarketPriceUpdates,
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
		},
		"Towards index price = true (current < update < index price)": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_005_000_000),
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
		"Index price crossing = true (price increase), old_ticks > 1, new_ticks <= sqrt(old_ticks) = true": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_005_000_000),
			},
			indexPrices: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeId:     constants.ExchangeId0,
							Price:          price_5_004_999_999,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			},
		},
		"Index price crossing = true (price decrease), old_ticks > 1, new_ticks <= sqrt(old_ticks) = true": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_4_995_000_000),
			},
			indexPrices: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeId:     constants.ExchangeId0,
							Price:          price_4_995_000_001,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			},
		},
		"Index price crossing = true (price increase), old_ticks <= 1, new_ticks <= old_ticks = true": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_000_500_000),
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
		"Index price crossing = true (price decrease), old_ticks <= 1, new_ticks <= old_ticks = true": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_4_999_500_000),
			},
			indexPrices: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId0,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeId:     constants.ExchangeId0,
							Price:          price_4_999_750_000,
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
			ctx, k, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run.
			msg := &types.MsgUpdateMarketPrices{
				MarketPriceUpdates: tc.msgUpdateMarketPrices,
			}
			err := k.PerformStatefulPriceUpdateValidation(ctx, msg, true)

			// Validate.
			require.NoError(t, err)
		})
	}
}

func TestPerformStatefulPriceUpdateValidation_SkipNonDeterministicCheck_Valid(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		msgUpdateMarketPrices []*types.MsgUpdateMarketPrices_MarketPrice
		indexPrices           []*api.MarketPriceUpdate
	}{
		"Index price does not exist": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, 11),
			},
			// Skipping price cache update, so the index price does not exist.
		},
		"Index price crossing = true, old_ticks > 1, new_ticks <= sqrt(old_ticks) = false": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_015_000_000),
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
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_015_000_000),
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
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_005_000_000),
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
			ctx, k, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run.
			msg := &types.MsgUpdateMarketPrices{
				MarketPriceUpdates: tc.msgUpdateMarketPrices,
			}
			err := k.PerformStatefulPriceUpdateValidation(ctx, msg, false) // skips non-deterministic checks.

			// Validate.
			require.NoError(t, err)
		})
	}
}

func TestPerformStatefulPriceUpdateValidation_Error(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		msgUpdateMarketPrices []*types.MsgUpdateMarketPrices_MarketPrice
		indexPrices           []*api.MarketPriceUpdate

		// Expected.
		expectedErr string
	}{
		"Market does not exist": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(99, 11), // Market with id 99 does not exist.
			},
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			expectedErr: errorsmod.Wrapf(
				types.ErrInvalidMarketPriceUpdateDeterministic,
				"market param price (99) does not exist",
			).Error(),
		},
		"Price does not meet min price change": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(
					constants.MarketId0,
					// The update price (=5,000,249,999) doesn't quite meet the min 50 ppm requirement.
					constants.FiveBillion+(constants.FiveBillion*50/uint64(lib.OneMillion))-1,
				),
			},
			indexPrices: constants.AtTimeTSingleExchangePriceUpdate,
			expectedErr: errorsmod.Wrapf(
				types.ErrInvalidMarketPriceUpdateDeterministic,
				"update price (5000249999) for market (0) does not meet min price change requirement"+
					" (50 ppm) based on the current market price (5000000000)",
			).Error(),
		},
		"Index price does not exist": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, 11),
			},
			// Skipping price cache update, so the index price does not exist.
			expectedErr: errorsmod.Wrapf(
				types.ErrIndexPriceNotAvailable,
				"index price for market (0) is not available",
			).Error(),
		},
		"Index price crossing = true, old_ticks > 1, new_ticks <= sqrt(old_ticks) = false": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_015_000_000),
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
			expectedErr: errorsmod.Wrap(
				types.ErrInvalidMarketPriceUpdateNonDeterministic,
				"update price (5015000000) for market (0) crosses the index price (5010000000) with "+
					"current price (5000000000) and deviates from index price (5000000) more than minimum allowed "+
					"(1581138)",
			).Error(),
		},
		"Index price crossing = true, old_ticks <= 1, new_ticks <= old_ticks = false": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_015_000_000),
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
			expectedErr: errorsmod.Wrap(
				types.ErrInvalidMarketPriceUpdateNonDeterministic,
				"update price (5015000000) for market (0) crosses the index price (5000250000) with current "+
					"price (5000000000) and deviates from index price (14750000) more than minimum allowed (250000)",
			).Error(),
		},
		"Price trends in the opposite direction as the index price": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_5_005_000_000),
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
			expectedErr: errorsmod.Wrap(
				types.ErrInvalidMarketPriceUpdateNonDeterministic,
				"update price (5005000000) for market (0) trends in the opposite direction of the index "+
					"price (4999999999) compared to the current price (5000000000)",
			).Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run and Validate.
			msg := &types.MsgUpdateMarketPrices{
				MarketPriceUpdates: tc.msgUpdateMarketPrices,
			}
			err := k.PerformStatefulPriceUpdateValidation(ctx, msg, true)
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestGetMarketsMissingFromPriceUpdates(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		msgUpdateMarketPrices []*types.MsgUpdateMarketPrices_MarketPrice
		indexPrices           []*api.MarketPriceUpdate
		smoothedIndexPrices   map[uint32]uint64

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
			expectedMarketIds: []uint32{
				constants.MarketId0, constants.MarketId1, constants.MarketId2, constants.MarketId3, constants.MarketId4,
			},
		},
		"Non-empty proposed updates, Empty local updates": {
			msgUpdateMarketPrices: constants.ValidMarketPriceUpdates,
			expectedMarketIds:     nil,
		},
		"Non-empty proposed updates, Non-empty local updates, no missing markets": {
			msgUpdateMarketPrices: constants.ValidMarketPriceUpdates,
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
			smoothedIndexPrices:   constants.AtTimeTSingleExchangeSmoothedPrices,
			expectedMarketIds:     nil,
		},
		"Non-empty proposed updates, Non-empty local updates, single missing market": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, constants.Price5),
				types.NewMarketPriceUpdate(constants.MarketId1, constants.Price6),
				types.NewMarketPriceUpdate(constants.MarketId3, constants.Price7),
				types.NewMarketPriceUpdate(constants.MarketId4, constants.Price4),
			},
			indexPrices:         constants.AtTimeTSingleExchangePriceUpdate,
			smoothedIndexPrices: constants.AtTimeTSingleExchangeSmoothedPrices,
			expectedMarketIds:   []uint32{constants.MarketId2},
		},
		"Non-empty proposed updates, Non-empty local updates, multiple missing markets, sorted": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId1, constants.Price6),
			},
			indexPrices:         constants.AtTimeTSingleExchangePriceUpdate,
			smoothedIndexPrices: constants.AtTimeTSingleExchangeSmoothedPrices,
			// The returned market ids must be sorted.
			expectedMarketIds: []uint32{constants.MarketId0, constants.MarketId2, constants.MarketId3, constants.MarketId4},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			keepertest.CreateTestMarkets(t, ctx, k)
			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run.
			missingMarketIds := k.GetMarketsMissingFromPriceUpdates(ctx, tc.msgUpdateMarketPrices)

			// Validate.
			// Using `Equal` here to test for slice ordering.
			require.Equal(t, tc.expectedMarketIds, missingMarketIds)
		})
	}
}
