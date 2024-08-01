package keeper_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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

// Note: the current market prices (i.e. 5 billion) are set in `CreateTestMarketsAndExchanges`.
func TestUpdateMarketPrices_Valid(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		msgUpdateMarketPrices []*types.MsgUpdateMarketPrices_MarketPrice
		indexPrices           []*api.MarketPriceUpdate

		// Expected.
		expectedMarketToPriceInState map[uint32]uint64
	}{
		"Empty updates": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{},
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: constants.FiveBillion,  // no change
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
			},
		},
		"Multiple updates": {
			msgUpdateMarketPrices: constants.ValidMarketPriceUpdates,
			indexPrices:           constants.AtTimeTSingleExchangePriceUpdate,
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: constants.Price5,
				constants.MarketId1: constants.Price6,
				constants.MarketId2: constants.Price7,
				constants.MarketId3: constants.Price4,
				constants.MarketId4: constants.Price3,
			},
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
							ExchangeId:     constants.BinanceExchangeName,
							Price:          price_5_010_000_000,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			},
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_5_005_000_000,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
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
							ExchangeId:     constants.BinanceExchangeName,
							Price:          price_5_004_999_999,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			},
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_5_005_000_000,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
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
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_4_995_000_000,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
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
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_5_000_500_000,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
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
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_4_999_500_000,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			msgServer := keeper.NewMsgServerImpl(k)
			keepertest.CreateTestMarkets(t, ctx, k)

			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run.
			_, err := msgServer.UpdateMarketPrices(
				ctx,
				&types.MsgUpdateMarketPrices{
					MarketPriceUpdates: tc.msgUpdateMarketPrices,
				})
			allMarketPricesAfterUpdate := k.GetAllMarketPrices(ctx)

			// Validate.
			require.NoError(t, err)
			require.Len(t, tc.expectedMarketToPriceInState, len(allMarketPricesAfterUpdate))
			for _, marketPrice := range allMarketPricesAfterUpdate {
				expectedPrice, exists := tc.expectedMarketToPriceInState[marketPrice.Id]
				require.True(t, exists)
				require.Equal(t, expectedPrice, marketPrice.Price)
			}
		})
	}
}

func TestUpdateMarketPrices_SkipNonDeterministicCheck_Valid(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		msgUpdateMarketPrices []*types.MsgUpdateMarketPrices_MarketPrice
		indexPrices           []*api.MarketPriceUpdate

		// Expected.
		expectedMarketToPriceInState map[uint32]uint64
	}{
		"Index price does not exist, but still updates state": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, 11),
			},
			// Skipping price cache update, so the index price does not exist.
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: 11,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
			},
		},
		"Index price trends in the opposite direction of update price from current price, but still updates state": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(constants.MarketId0, price_4_995_000_001),
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
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_4_995_000_001,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
			},
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
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_5_015_000_000,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
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
			expectedMarketToPriceInState: map[uint32]uint64{
				constants.MarketId0: price_5_015_000_000,
				constants.MarketId1: constants.ThreeBillion, // no change
				constants.MarketId2: constants.FiveBillion,  // no change
				constants.MarketId3: constants.FiveBillion,  // no change
				constants.MarketId4: constants.ThreeBillion, // no change
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)

			msgServer := keeper.NewMsgServerImpl(k)
			keepertest.CreateTestMarkets(t, ctx, k)

			indexPriceCache.UpdatePrices(tc.indexPrices)

			// Run.
			_, err := msgServer.UpdateMarketPrices(
				ctx,
				&types.MsgUpdateMarketPrices{
					MarketPriceUpdates: tc.msgUpdateMarketPrices,
				})
			allMarketPricesAfterUpdate := k.GetAllMarketPrices(ctx)

			// Validate.
			require.NoError(t, err)
			require.Len(t, tc.expectedMarketToPriceInState, len(allMarketPricesAfterUpdate))
			for _, marketPrice := range allMarketPricesAfterUpdate {
				expectedPrice, exists := tc.expectedMarketToPriceInState[marketPrice.Id]
				require.True(t, exists)
				require.Equal(t, expectedPrice, marketPrice.Price)
			}
		})
	}
}

func TestUpdateMarketPrices_Error(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		msgUpdateMarketPrices []*types.MsgUpdateMarketPrices_MarketPrice

		// Expected.
		expectedErr error
	}{
		"Market does not exist": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(99, 11), // Market with id 99 does not exist.
			},
			expectedErr: errorsmod.Wrapf(
				types.ErrInvalidMarketPriceUpdateDeterministic,
				"market param price (99) does not exist",
			),
		},
		"Price does not meet min price change": {
			msgUpdateMarketPrices: []*types.MsgUpdateMarketPrices_MarketPrice{
				types.NewMarketPriceUpdate(
					constants.MarketId0,
					// The update price (=5,000,249,999) doesn't quite meet the min 50 ppm requirement.
					constants.FiveBillion+(constants.FiveBillion*50/uint64(lib.OneMillion))-1,
				),
			},
			expectedErr: errorsmod.Wrapf(
				types.ErrInvalidMarketPriceUpdateDeterministic,
				"update price (5000249999) for market (0) does not meet min price change requirement"+
					" (50 ppm) based on the current market price (5000000000)",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, k, _, _, mockTimeKeeper, _, _ := keepertest.PricesKeepers(t)
			mockTimeKeeper.On("Now").Return(constants.TimeT)
			msgServer := keeper.NewMsgServerImpl(k)
			keepertest.CreateTestMarkets(t, ctx, k)

			// Run and Validate.
			require.PanicsWithError(
				t,
				tc.expectedErr.Error(),
				func() {
					_, _ = msgServer.UpdateMarketPrices(
						ctx,
						&types.MsgUpdateMarketPrices{
							MarketPriceUpdates: tc.msgUpdateMarketPrices,
						})
				},
			)
		})
	}
}

func TestUpdateMarketPrices_Panic(t *testing.T) {
	// Init.
	ctx, _, _, _, _, _, _ := keepertest.PricesKeepers(t)
	mockKeeper := &mocks.PricesKeeper{}
	msgServer := keeper.NewMsgServerImpl(mockKeeper)

	testError := errors.New("panic like there's no tomorrow")
	testUpdates := []*types.MsgUpdateMarketPrices_MarketPrice{}
	testMsg := &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: testUpdates,
	}

	// Setup.
	mockKeeper.On("Logger", ctx).Return(log.NewNopLogger())
	mockKeeper.On("PerformStatefulPriceUpdateValidation", ctx, testMsg, false).Return(nil)
	mockKeeper.On("UpdateMarketPrices", ctx, testUpdates).Return(testError)

	// Run and Validate.
	require.PanicsWithError(
		t,
		testError.Error(),
		func() {
			_, _ = msgServer.UpdateMarketPrices(ctx, testMsg)
		},
	)
}
