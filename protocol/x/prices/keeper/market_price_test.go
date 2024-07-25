package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func createNMarketPriceUpdates(
	n int,
) []*types.MsgUpdateMarketPrices_MarketPrice {
	items := make([]*types.MsgUpdateMarketPrices_MarketPrice, n)
	for i := range items {
		items[i] = &types.MsgUpdateMarketPrices_MarketPrice{
			MarketId: uint32(i),
			Price:    uint64(i),
		}
	}

	return items
}

func TestUpdateMarketPrices(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	require.Equal(t, uint32(10), keepertest.GetNumMarkets(t, ctx, keeper))

	// Create firstPriceUpdates which should be overwritten by secondPriceUpdates
	firstPriceUpdates := createNMarketPriceUpdates(10)
	secondPriceUpdates := createNMarketPriceUpdates(10)
	for _, pu := range secondPriceUpdates {
		pu.Price = 10 + (pu.Price * 10)
	}

	priceUpdates := append(firstPriceUpdates, secondPriceUpdates...)
	err := keeper.UpdateMarketPrices(
		ctx,
		priceUpdates,
	)
	require.NoError(t, err)

	marketPrices := make([]types.MarketPrice, 10)
	for i, item := range items {
		marketPrice, err := keeper.GetMarketPrice(ctx, item.Param.Id)
		require.NoError(t, err)
		require.Equal(t,
			secondPriceUpdates[i].Price,
			marketPrice.Price,
		)
		marketPrices = append(marketPrices, marketPrice)
	}

	keepertest.AssertPriceUpdateEventsInIndexerBlock(t, keeper, ctx, marketPrices)
}

func TestUpdateMarketPrices_NotFound(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	mockTimeProvider.On("Now").Return(constants.TimeT)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	priceUpdates := createNMarketPriceUpdates(10)
	err := keeper.UpdateMarketPrices(
		ctx,
		priceUpdates,
	)
	require.EqualError(t, err, "0: Market price does not exist")
	keepertest.AssertMarketEventsNotInIndexerBlock(t, keeper, ctx)

	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	err = keeper.UpdateMarketPrices(
		ctx,
		priceUpdates,
	)
	require.NoError(t, err)

	for i, item := range items {
		price, err := keeper.GetMarketPrice(ctx, item.Param.Id)
		require.NoError(t, err)
		require.Equal(t,
			priceUpdates[i].Price,
			price.Price,
		)
	}
}

func TestGetMarketPrice(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
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
	ctx, keeper, _, _, _, _, _ := keepertest.PricesKeepers(t)
	_, err := keeper.GetMarketPrice(ctx, uint32(0))
	require.EqualError(t, err, "0: Market price does not exist")
}

func TestGetAllMarketPrices(t *testing.T) {
	ctx, keeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
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

func TestGetMarketIdToValidIndexPrice(t *testing.T) {
	ctx, keeper, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
	// Now() is used by `GetMarketIdToValidIndexPrice` internally compare with the cutoff time
	// of each price.
	mockTimeProvider.On("Now").Return(constants.TimeT)
	keepertest.CreateTestPriceMarkets(t,
		ctx,
		keeper,
		[]types.MarketParamPrice{
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithPair("0-0"),
				pricestest.WithId(6),
				pricestest.WithMinExchanges(2),
			),
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithPair("1-1"),
				pricestest.WithId(7),
				pricestest.WithMinExchanges(2),
			),
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithPair("2-2"),
				pricestest.WithId(8),
				pricestest.WithMinExchanges(2),
				pricestest.WithExponent(-8),
			),
			*pricestest.GenerateMarketParamPrice(
				pricestest.WithPair("3-3"),
				pricestest.WithId(9),
				pricestest.WithMinExchanges(2),
				pricestest.WithExponent(-9),
			),
		},
	)

	// Set up index price cache values for market 7, 8, 9.
	indexPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)
	marketIdToIndexPrice := keeper.GetMarketIdToValidIndexPrice(ctx)
	// While there are 4 markets in state, only 7, 8, 9 have index prices,
	// and only 8, 9 have valid median index prices.
	// Market7 only has 1 valid price due to update time constraint.
	require.Len(t, marketIdToIndexPrice, 3)
	require.Equal(t,
		types.MarketPrice{
			Id:       constants.MarketId9,
			Price:    uint64(2002),
			Exponent: constants.Exponent9,
		},
		marketIdToIndexPrice[constants.MarketId9],
	) // Median of 1001, 2002, 3003
	require.Equal(t,
		types.MarketPrice{
			Id:       constants.MarketId8,
			Price:    uint64(2503),
			Exponent: constants.Exponent8,
		},
		marketIdToIndexPrice[constants.MarketId8],
	) // Median of 2002, 3003
}
