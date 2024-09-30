package keeper_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	"github.com/stretchr/testify/require"
)

func TestCreateMarket(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, revShareKeeper, marketMapKeeper := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)

	testMarketParams := types.MarketParam{
		Id:                 0,
		Pair:               constants.BtcUsdPair,
		Exponent:           int32(-6),
		ExchangeConfigJson: `{"test_config_placeholder":{}}`,
		MinExchanges:       2,
		MinPriceChangePpm:  uint32(9_999),
	}

	// Test that creating market fails if it does not exist in marketmap
	_, err := keeper.CreateMarket(
		ctx,
		testMarketParams,
		types.MarketPrice{
			Id:       0,
			Exponent: int32(-6),
			Price:    constants.FiveBillion,
		},
	)
	require.Error(t, err, types.ErrTickerNotFoundInMarketMap)

	// Create the test market in the market map and verify it is not enabled
	keepertest.CreateMarketsInMarketMapFromParams(
		t,
		ctx,
		keeper.MarketMapKeeper.(*marketmapkeeper.Keeper),
		[]types.MarketParam{testMarketParams},
	)
	currencyPair, _ := slinky.MarketPairToCurrencyPair(constants.BtcUsdPair)
	mmMarket, _ := marketMapKeeper.GetMarket(ctx, currencyPair.String())
	require.False(t, mmMarket.Ticker.Enabled)

	// Verify that currency pair is not in the CurrencyPairToID cache
	_, found := keeper.GetIDForCurrencyPair(ctx, currencyPair)
	require.False(t, found)

	marketParam, err := keeper.CreateMarket(
		ctx,
		testMarketParams,
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

	// Verify that currency pair is in the CurrencyPairToID cache
	cpID, found := keeper.GetIDForCurrencyPair(ctx, currencyPair)
	require.True(t, found)
	require.Equal(t, uint64(marketParam.Id), cpID)

	// Verify expected price of 0 created.
	require.Equal(t, uint32(0), marketPrice.Id)
	require.Equal(t, int32(-6), marketPrice.Exponent)
	require.Equal(t, constants.FiveBillion, marketPrice.Price)

	require.Equal(t, marketParam.Pair, metrics.GetMarketPairForTelemetry(marketParam.Id))

	// Verify expected market event.
	keepertest.AssertMarketCreateEventInIndexerBlock(t, keeper, ctx, marketParam)

	// Verify market revenue share creation
	revShareParams := revShareKeeper.GetMarketMapperRevenueShareParams(ctx)
	revShareDetails, err := revShareKeeper.GetMarketMapperRevShareDetails(ctx, marketParam.Id)
	require.NoError(t, err)

	// Verify market is enabled in market map
	mmMarket, _ = marketMapKeeper.GetMarket(ctx, currencyPair.String())
	require.True(t, mmMarket.Ticker.Enabled)

	expirationTs := uint64(ctx.BlockTime().Unix() + int64(revShareParams.ValidDays*24*3600))
	require.Equal(t, revShareDetails.ExpirationTs, expirationTs)
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
				"market param id 1 does not match market price id 2",
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
				"market param 1 exponent -6 does not match market price 1 exponent -5",
			).Error(),
		},
		"Pair already exists": {
			pair:               "0-0",
			minExchanges:       uint32(2),
			minPriceChangePpm:  uint32(50),
			price:              constants.FiveBillion,
			exchangeConfigJson: validExchangeConfigJson,
			expectedErr: errorsmod.Wrap(
				types.ErrMarketParamPairAlreadyExists,
				"0-0",
			).Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, mockTimeKeeper, _, _ := keepertest.PricesKeepers(t)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)

			mockTimeKeeper.On("Now").Return(constants.TimeT)
			keepertest.CreateNMarkets(t, ctx, keeper, 1)

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
					Id:                 1,
					Pair:               tc.pair,
					Exponent:           int32(-6),
					MinExchanges:       tc.minExchanges,
					MinPriceChangePpm:  tc.minPriceChangePpm,
					ExchangeConfigJson: tc.exchangeConfigJson,
				},
				types.MarketPrice{
					Id:       1 + marketPriceIdOffset,
					Exponent: int32(-6) + marketPriceExponentOffset,
					Price:    tc.price,
				},
			)
			require.EqualError(t, err, tc.expectedErr)

			// Verify no new MarketPrice created.
			_, err = keeper.GetMarketPrice(ctx, 1)
			require.EqualError(
				t,
				err,
				errorsmod.Wrap(types.ErrMarketPriceDoesNotExist, lib.UintToString(uint32(1))).Error(),
			)

			// Verify no new market event.
			keepertest.AssertNMarketEventsNotInIndexerBlock(t, keeper, ctx, 1)
		})
	}
}

func TestGetAllMarketParamPrices(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
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

func TestAcquireNextMarketID(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)

	keepertest.CreateNMarkets(t, ctx, keeper, 10)

	// Get the highest market ID from the existing markets.
	allParams := keeper.GetAllMarketParams(ctx)
	highestMarketID := uint32(0)
	for _, param := range allParams {
		if param.Id > highestMarketID {
			highestMarketID = param.Id
		}
	}

	// Acquire the next market ID.
	nextMarketID := keeper.AcquireNextMarketID(ctx)
	require.Equal(t, highestMarketID+1, nextMarketID)

	// Verify the next market ID is stored in the module store.
	nextMarketIDFromStore := keeper.GetNextMarketID(ctx)
	require.Equal(t, nextMarketID+1, nextMarketIDFromStore)

	// Create a market with the next market ID outside of acquire flow
	_, err := keepertest.CreateTestMarket(
		t,
		ctx,
		keeper,
		types.MarketParam{
			Id:                 nextMarketIDFromStore,
			Pair:               "TEST-USD",
			Exponent:           int32(-6),
			ExchangeConfigJson: `{"test_config_placeholder":{}}`,
			MinExchanges:       2,
			MinPriceChangePpm:  uint32(9_999),
		},
		types.MarketPrice{
			Id:       nextMarketIDFromStore,
			Exponent: int32(-6),
			Price:    constants.FiveBillion,
		},
	)
	require.NoError(t, err)

	// Verify the next market ID is incremented.
	nextMarketID = keeper.AcquireNextMarketID(ctx)
	require.Equal(t, nextMarketIDFromStore+1, nextMarketID)
}
