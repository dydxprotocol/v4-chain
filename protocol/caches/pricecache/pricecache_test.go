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
	dummyBytesEmpty := []byte{}
	dummyBytes1 := []byte{1}
	dummyBytes2 := []byte{1, 2, 3, 4}
	dummyBytes3 := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	t.Run("valid: set price updates for single round + height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, dummyBytes1)
		checkValidCacheState(t, &priceCache, dummyBytes1, updates)
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
		priceCache.SetPriceUpdates(ctx, updates, dummyBytes2)

		checkValidCacheState(t, &priceCache, dummyBytes2, updates)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, newUpdates, dummyBytes3)
		checkValidCacheState(t, &priceCache, dummyBytes3, newUpdates)
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
		priceCache.SetPriceUpdates(ctx, updates, dummyBytes2)
		checkValidCacheState(t, &priceCache, dummyBytes2, updates)

		ctx = ctx.WithBlockHeight(4)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, newUpdates, dummyBytes2)
		checkValidCacheState(t, &priceCache, dummyBytes2, newUpdates)
	})

	t.Run("invalid: No valid values, wrong txHash", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(6)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, dummyBytes2)
		require.False(t, priceCache.HasValidValues(dummyBytes3))
	})

	t.Run("invalid: No valid values, wrong tx hash (empty)", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(8)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId: marketPrice.MarketId,
				Price:    big.NewInt(int64(marketPrice.SpotPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, dummyBytesEmpty)
		require.False(t, priceCache.HasValidValues(dummyBytes1))
	})
}

func checkValidCacheState(
	t *testing.T,
	priceCache pricecache.PriceUpdatesCache,
	shouldBeHash []byte,
	shouldBeUpdates pricecache.PriceUpdates,
) {
	require.True(t, priceCache.HasValidValues(shouldBeHash))
	require.Equal(t, shouldBeUpdates, priceCache.GetPriceUpdates())
}
