package keeper_test

import (
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
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
			ExchangeConfigJson: `{"test_config_placeholder":{}}`,
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
	require.Equal(t, `{"test_config_placeholder":{}}`, marketParam.ExchangeConfigJson)
	require.Equal(t, uint32(2), marketParam.MinExchanges)
	require.Equal(t, uint32(9999), marketParam.MinPriceChangePpm)

	// Verify expected price of 0 created.
	require.Equal(t, uint32(0), marketPrice.Id)
	require.Equal(t, int32(-6), marketPrice.Exponent)
	require.Equal(t, constants.FiveBillion, marketPrice.Price)

	require.Equal(t, marketParam.Pair, metrics.GetMarketPairForTelemetry(marketParam.Id))

	// Verify expected market event.
	keepertest.AssertMarketCreateEventInIndexerBlock(t, keeper, ctx, marketParam)
}

func TestMarketIsRecentlyAvailable(t *testing.T) {
	tests := map[string]struct {
		blockHeight      int64
		now              time.Time
		expectedIsRecent bool
	}{
		"Recent: << block height, << elapsed since market creation time": {
			blockHeight:      0,
			now:              constants.TimeT.Add(types.MarketIsRecentDuration - 1),
			expectedIsRecent: true,
		},
		"Recent: >> block height, << elapsed since market creation time": {
			blockHeight:      types.PriceDaemonInitializationBlocks + 1,
			now:              constants.TimeT.Add(types.MarketIsRecentDuration - 1),
			expectedIsRecent: true,
		},
		"Recent: << block height, >> elapsed since market creation time": {
			blockHeight:      0,
			now:              constants.TimeT.Add(types.MarketIsRecentDuration + 1),
			expectedIsRecent: true,
		},
		"Not recent: >> block height, >> elapsed since market creation time": {
			blockHeight:      types.PriceDaemonInitializationBlocks + 1,
			now:              constants.TimeT.Add(types.MarketIsRecentDuration + 1),
			expectedIsRecent: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)

			// Create market with TimeT creation timestamp.
			mockTimeProvider.On("Now").Return(constants.TimeT).Once()
			require.False(t, keeper.IsRecentlyAvailable(ctx, 0))

			keepertest.CreateNMarkets(t, ctx, keeper, 1)

			ctx = ctx.WithBlockHeight(tc.blockHeight)
			mockTimeProvider.On("Now").Return(tc.now).Once()

			require.Equal(t, tc.expectedIsRecent, keeper.IsRecentlyAvailable(ctx, 0))
		})
	}
}

func TestCreateMarket_Errors(t *testing.T) {
	validExchangeConfigJson := `{"exchanges":[{"exchangeName":"Binance","ticker":"BTCUSDT"}]}`
	tests := map[string]struct {
		// Setup
		pair                                              string
		minExchanges                                      uint32
		minPriceChangePpm                                 uint32
		price                                             uint64
		marketPriceIdDoesntMatchMarketParamId             bool
		marketPriceExponentDoesntMatchMarketParamExponent bool
		exchangeConfigJson                                string
		// Expected
		expectedErr string
	}{
		"Empty pair": {
			pair:               "", // pair cannot be empty
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(50),
			price:              constants.FiveBillion,
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgMarketPairCannotBeEmpty).Error(),
		},
		"Invalid min price change: zero": {
			pair:               constants.BtcUsdPair,
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(0), // must be > 0
			price:              constants.FiveBillion,
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Invalid min price change: ten thousand": {
			pair:              constants.BtcUsdPair,
			minExchanges:      uint32(2),
			minPriceChangePpm: uint32(10_000), // must be < 10,000
			price:             constants.FiveBillion,
			expectedErr:       errorsmod.Wrap(types.ErrInvalidInput, constants.ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Min exchanges cannot be zero": {
			pair:               constants.BtcUsdPair,
			minExchanges:       uint32(0), // cannot be zero
			minPriceChangePpm:  uint32(50),
			price:              constants.FiveBillion,
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr:        types.ErrZeroMinExchanges.Error(),
		},
		"Market param and price ids don't match": {
			pair:                                  constants.BtcUsdPair,
			minExchanges:                          uint32(2),
			minPriceChangePpm:                     uint32(50),
			price:                                 constants.FiveBillion,
			marketPriceIdDoesntMatchMarketParamId: true,
			exchangeConfigJson:                    validExchangeConfigJson,
			expectedErr: errorsmod.Wrap(
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
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr: errorsmod.Wrap(
				types.ErrInvalidInput,
				"market param 0 exponent -6 does not match market price 0 exponent -5",
			).Error(),
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
					Id:                 0,
					Pair:               tc.pair,
					Exponent:           int32(-6),
					MinExchanges:       tc.minExchanges,
					MinPriceChangePpm:  tc.minPriceChangePpm,
					ExchangeConfigJson: tc.exchangeConfigJson,
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
				errorsmod.Wrap(types.ErrMarketPriceDoesNotExist, lib.UintToString(uint32(0))).Error(),
			)

			// Verify no new market event.
			keepertest.AssertMarketEventsNotInIndexerBlock(t, keeper, ctx)
		})
	}
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
