package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
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
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	items := keepertest.CreateNMarkets(t, ctx, keeper, 10)
	require.Equal(t, uint32(10), keeper.GetNumMarkets(ctx))

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
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
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
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
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
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
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
