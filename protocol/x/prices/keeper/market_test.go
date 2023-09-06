package keeper_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateMarket(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)

	marketParam, err := keeper.CreateMarket(
		ctx,
		types.MarketParam{
			Id:                 0,
			Pair:               constants.BtcUsdPair,
			Exponent:           int32(-6),
			ExchangeConfigJson: "test_config_placeholder",
			MinExchanges:       2,
			MinPriceChangePpm:  uint32(9_999),
		},
		types.MarketPrice{
			Id:       0,
			Exponent: int32(-6),
			Price:    constants.FiveBillion,
		},
	)

	require.NoError(t, err)

	marketPrice, err := keeper.GetMarketPrice(ctx, marketParam.Id)
	require.NoError(t, err)

	// Verify expected param.
	require.Equal(t, uint32(0), marketParam.Id)
	require.Equal(t, constants.BtcUsdPair, marketParam.Pair)
	require.Equal(t, int32(-6), marketParam.Exponent)
	require.Equal(t, "test_config_placeholder", marketParam.ExchangeConfigJson)
	require.Equal(t, uint32(2), marketParam.MinExchanges)
	require.Equal(t, uint32(9999), marketParam.MinPriceChangePpm)

	// Verify expected price of 0 created.
	require.Equal(t, uint32(0), marketPrice.Id)
	require.Equal(t, int32(-6), marketPrice.Exponent)
	require.Equal(t, constants.FiveBillion, marketPrice.Price)

	// Verify expected market event.
	keepertest.AssertMarketCreateEventInIndexerBlock(t, keeper, ctx, marketParam)
}

func TestMarketIsRecentlyAdded(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT).Once()

	// Nonexistent markets should not be recently added.
	require.False(t, keeper.IsRecentlyAdded(0))

	keepertest.CreateNMarkets(t, ctx, keeper, 1)

	// Before the duration passes, the market should be recently added.
	mockTimeProvider.On("Now").Return(constants.TimeT.Add(types.MarketIsRecentDuration - 1)).Once()
	require.True(t, keeper.IsRecentlyAdded(0))

	// After the duration passes, the market is no longer recently added.
	mockTimeProvider.On("Now").Return(constants.TimeT.Add(types.MarketIsRecentDuration)).Once()
	require.False(t, keeper.IsRecentlyAdded(0))
}

func TestCreateMarket_Errors(t *testing.T) {
	tests := map[string]struct {
		// Setup
		pair                                              string
		minExchanges                                      uint32
		minPriceChangePpm                                 uint32
		price                                             uint64
		marketPriceIdDoesntMatchMarketParamId             bool
		marketPriceExponentDoesntMatchMarketParamExponent bool
		// Expected
		expectedErr string
	}{
		"Empty pair": {
			pair:              "", // pair cannot be empty
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(50),
			price:             constants.FiveBillion,
			expectedErr:       sdkerrors.Wrap(types.ErrInvalidInput, constants.ErrorMsgMarketPairCannotBeEmpty).Error(),
		},
		"Invalid min price change: zero": {
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(0), // must be > 0
			price:             constants.FiveBillion,
			expectedErr:       sdkerrors.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Invalid min price change: ten thousand": {
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(10_000), // must be < 10,000
			price:             constants.FiveBillion,
			expectedErr:       sdkerrors.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Min exchanges cannot be zero": {
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(0), // cannot be zero
			minPriceChangePpm: uint32(50),
			price:             constants.FiveBillion,
			expectedErr:       types.ErrZeroMinExchanges.Error(),
		},
		"Market param and price ids don't match": {
			pair:                                  constants.BtcUsdPair,
			minExchanges:                          uint32(2),
			minPriceChangePpm:                     uint32(50),
			price:                                 constants.FiveBillion,
			marketPriceIdDoesntMatchMarketParamId: true,
			expectedErr: sdkerrors.Wrap(
				types.ErrInvalidInput,
				"market param id 0 does not match market price id 1",
			).Error(),
		},
		"Market param and price exponents don't match": {
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(50),
			price:             constants.FiveBillion,
			marketPriceExponentDoesntMatchMarketParamExponent: true,
			expectedErr: sdkerrors.Wrap(
				types.ErrInvalidInput,
				"market param 0 exponent -6 does not match market price 0 exponent -5",
			).Error(),
		},
		"Market price is 0": {
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(50),
			price:             uint64(0),
			expectedErr:       sdkerrors.Wrap(types.ErrInvalidInput, "market 0 price cannot be zero").Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)

			marketPriceIdOffset := uint32(0)
			if tc.marketPriceIdDoesntMatchMarketParamId {
				marketPriceIdOffset = uint32(1)
			}

			marketPriceExponentOffset := int32(0)
			if tc.marketPriceExponentDoesntMatchMarketParamExponent {
				marketPriceExponentOffset = int32(1)
			}

			_, err := keeper.CreateMarket(
				ctx,
				types.MarketParam{
					Id:                0,
					Pair:              tc.pair,
					Exponent:          int32(-6),
					MinExchanges:      tc.minExchanges,
					MinPriceChangePpm: tc.minPriceChangePpm,
				},
				types.MarketPrice{
					Id:       0 + marketPriceIdOffset,
					Exponent: int32(-6) + marketPriceExponentOffset,
					Price:    tc.price,
				},
			)
			require.EqualError(t, err, tc.expectedErr)

			// Verify no new MarketPrice created.
			_, err = keeper.GetMarketPrice(ctx, 0)
			require.EqualError(
				t,
				err,
				sdkerrors.Wrap(types.ErrMarketPriceDoesNotExist, lib.Uint32ToString(0)).Error(),
			)

			// Verify no new market event.
			keepertest.AssertMarketEventsNotInIndexerBlock(t, keeper, ctx)
		})
	}
}

func TestGetNumMarkets(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)

	require.Equal(t, uint32(0), keeper.GetNumMarkets(ctx))

	keepertest.CreateNMarkets(t, ctx, keeper, 10)
	require.Equal(t, uint32(10), keeper.GetNumMarkets(ctx))
}

func TestGetAllMarketParamPrices(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)

	allParamPrices, err := keeper.GetAllMarketParamPrices(ctx)
	require.NoError(t, err)
	require.ElementsMatch(
		t,
		items,
		allParamPrices,
	)
}
