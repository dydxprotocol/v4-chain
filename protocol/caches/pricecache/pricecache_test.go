package pricecache_test

import (
	"math/big"
	"testing"

	pricecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/pricecache"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestVEPriceCaching(t *testing.T) {
	priceCache := pricecache.PriceUpdatesCacheImpl{}
	ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

	t.Run("valid: set price updates for single round + height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &priceCache, 1, 1, updates)
	})

	t.Run("valid: set priced updates for multi round single height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(2)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)

		checkValidCacheState(t, &priceCache, 2, 1, updates)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, newUpdates, 2)
		checkValidCacheState(t, &priceCache, 2, 2, newUpdates)
	})

	t.Run("valid: set price updates for single rounds multi height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(3)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &priceCache, 3, 1, updates)

		ctx = ctx.WithBlockHeight(4)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, newUpdates, 1)
		checkValidCacheState(t, &priceCache, 4, 1, newUpdates)
	})

	t.Run("valid: set diff update for same height + round", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(5)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &priceCache, 5, 1, updates)

		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, newUpdates, 1)
		checkValidCacheState(t, &priceCache, 5, 1, newUpdates)
	})

	t.Run("invalid: No valid values, wrong block height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(6)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		require.False(t, priceCache.HasValidValues(7, 1))
	})

	t.Run("invalid: No valid values, wrong round", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(8)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		require.False(t, priceCache.HasValidValues(8, 2))
	})
}

func checkValidCacheState(
	t *testing.T,
	priceCache pricecache.PriceUpdatesCache,
	shouldBeHight int64,
	shouldBeRound int32,
	shouldBeUpdates pricecache.PriceUpdates,
) {
	require.True(t, priceCache.HasValidValues(shouldBeHight, shouldBeRound))
	require.Equal(t, shouldBeHight, priceCache.GetHeight())
	require.Equal(t, shouldBeRound, priceCache.GetRound())
	require.Equal(t, shouldBeUpdates, priceCache.GetPriceUpdates())
}
