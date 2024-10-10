package keeper_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	pricestest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/prices"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func createNMarketPriceUpdates(
	n int,
) []*types.MarketPriceUpdate {
	items := make([]*types.MarketPriceUpdate, n)
	for i := range items {
		items[i] = &types.MarketPriceUpdate{
			MarketId:  uint32(i),
			SpotPrice: uint64(i),
			PnlPrice:  uint64(i),
		}
	}
	return items
}

func TestUpdateMarketPrices(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	require.Equal(t, uint32(10), keepertest.GetNumMarkets(t, ctx, keeper))

	// Create firstPriceUpdates which should be overwritten by secondPriceUpdates
	firstPriceUpdates := createNMarketPriceUpdates(10)
	secondPriceUpdates := createNMarketPriceUpdates(10)
	for _, pu := range secondPriceUpdates {
		pu.SpotPrice = 10 + (pu.SpotPrice * 10)
		pu.PnlPrice = 10 + (pu.PnlPrice * 10)
	}

	priceUpdates := append(firstPriceUpdates, secondPriceUpdates...)
	for _, update := range priceUpdates {
		err := keeper.UpdateSpotAndPnlMarketPrices(
			ctx,
			update,
		)
		require.NoError(t, err)
	}

	marketPrices := make([]types.MarketPrice, 10)
	for i, item := range items {
		marketPrice, err := keeper.GetMarketPrice(ctx, item.Param.Id)
		require.NoError(t, err)
		require.Equal(t,
			secondPriceUpdates[i].SpotPrice,
			marketPrice.SpotPrice,
		)
		require.Equal(t,
			secondPriceUpdates[i].PnlPrice,
			marketPrice.PnlPrice,
		)
		marketPrices = append(marketPrices, marketPrice)
	}

	keepertest.AssertPriceUpdateEventsInIndexerBlock(t, keeper, ctx, marketPrices)
}

func TestUpdateMarketPrices_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	priceUpdates := createNMarketPriceUpdates(10)
	for i, update := range priceUpdates {
		err := keeper.UpdateSpotAndPnlMarketPrices(
			ctx,
			update,
		)
		require.EqualError(t, err, fmt.Sprintf(`%d: Market price does not exist`, i))
	}
	keepertest.AssertMarketEventsNotInIndexerBlock(t, keeper, ctx)

	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	for _, update := range priceUpdates {
		err := keeper.UpdateSpotAndPnlMarketPrices(
			ctx,
			update,
		)
		require.NoError(t, err)
	}

	for i, item := range items {
		price, err := keeper.GetMarketPrice(ctx, item.Param.Id)
		require.NoError(t, err)
		require.Equal(t,
			priceUpdates[i].SpotPrice,
			price.SpotPrice,
		)
		require.Equal(t,
			priceUpdates[i].PnlPrice,
			price.PnlPrice,
		)
	}
}

func TestGetMarketPrice(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	for _, item := range items {
		rst, err := keeper.GetMarketPrice(ctx, item.Param.Id)
		require.NoError(t, err)
		require.Equal(
			t,
			&item.Price,
			&rst,
		)
	}
}

func TestGetMarketPrice_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	_, err := keeper.GetMarketPrice(ctx, uint32(0))
	require.EqualError(t, err, "0: Market price does not exist")
}

func TestGetAllMarketPrices(t *testing.T) {
	ctx, keeper, _, _, _, mockTimeProvider := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	prices := make([]types.MarketPrice, len(items))
	for i, item := range items {
		prices[i] = item.Price
	}

	require.ElementsMatch(
		t,
		prices,
		keeper.GetAllMarketPrices(ctx),
	)
}

func TestGetMarketIdToValidDaemonPrice(t *testing.T) {
	ctx, keeper, _, daemonPriceCache, _, mockTimeProvider := keepertest.PricesKeepers(t)
	// Now() is used by `GetMarketIdToValidDaemonPrice` internally compare with the cutoff time
	// of each price.
	mockTimeProvider.On("Now").Return(constants.TimeT)
	keepertest.CreateTestPriceMarkets(t,
		ctx,
		keeper,
		[]types.MarketParamPrice{
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithId(6),
				pricestest.WithMinExchanges(2),
			),
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithId(7),
				pricestest.WithMinExchanges(2),
			),
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithId(8),
				pricestest.WithMinExchanges(2),
				pricestest.WithExponent(-8),
			),
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithId(9),
				pricestest.WithMinExchanges(2),
				pricestest.WithExponent(-9),
			),
		},
	)

	// Set up daemon price cache values for market 7, 8, 9.
	daemonPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)
	marketIdToDaemonPrice := keeper.GetMarketIdToValidDaemonPrice(ctx)
	// While there are 4 markets in state, only 7, 8, 9 have daemon prices,
	// and only 8, 9 have valid median daemon prices.
	// Market7 only has 1 valid price due to update time constraint,
	// but the min exchanges required is 2. Therefore, no median price.
	require.Len(t, marketIdToDaemonPrice, 2)
	require.Equal(t,
		types.MarketSpotPrice{
			Id:        constants.MarketId9,
			SpotPrice: uint64(2002),
			Exponent:  constants.Exponent9,
		},
		marketIdToDaemonPrice[constants.MarketId9],
	) // Median of 1001, 2002, 3003
	require.Equal(t,
		types.MarketSpotPrice{
			Id:        constants.MarketId8,
			SpotPrice: uint64(2503),
			Exponent:  constants.Exponent8,
		},
		marketIdToDaemonPrice[constants.MarketId8],
	) // Median of 2002, 3003
}
