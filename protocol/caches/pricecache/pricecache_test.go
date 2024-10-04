package pricecache_test

import (
	"math/big"
	"sync"
	"testing"
	"time"

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
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
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
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)

		checkValidCacheState(t, &priceCache, 2, 1, updates)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
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
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &priceCache, 3, 1, updates)

		ctx = ctx.WithBlockHeight(4)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
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
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &priceCache, 5, 1, updates)

		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
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
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
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
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx, updates, 1)
		require.False(t, priceCache.HasValidValues(8, 2))
	})

	t.Run("concurrent reads and writes", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if i%2 == 0 {
					// Even goroutines write
					var updates pricecache.PriceUpdates
					for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
						updates = append(updates, pricecache.PriceUpdate{
							MarketId:  marketPrice.MarketId,
							SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
							PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
						})
					}
					priceCache.SetPriceUpdates(ctx.WithBlockHeight(int64(i)), updates, int32(i))
				} else {
					// Odd goroutines read
					time.Sleep(time.Millisecond) // Slight delay to increase chances of interleaving
					_ = priceCache.GetPriceUpdates()
					_ = priceCache.GetHeight()
					_ = priceCache.GetRound()
					_ = priceCache.HasValidValues(int64(i), int32(i))
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("concurrent writes", func(t *testing.T) {
		var wg sync.WaitGroup
		numWrites := 1000

		for i := 0; i < numWrites; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				var updates pricecache.PriceUpdates
				for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
					updates = append(updates, pricecache.PriceUpdate{
						MarketId:  marketPrice.MarketId,
						SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
						PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
					})
				}
				priceCache.SetPriceUpdates(ctx.WithBlockHeight(int64(i)), updates, int32(i))
			}(i)
		}

		wg.Wait()

		// Verify the final state
		height := priceCache.GetHeight()
		round := priceCache.GetRound()
		require.True(t, height >= 0 && height < int64(numWrites))
		require.True(t, round >= 0 && round < int32(numWrites))
	})

	t.Run("concurrent reads", func(t *testing.T) {
		// Set initial state
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		priceCache.SetPriceUpdates(ctx.WithBlockHeight(100), updates, 5)

		var wg sync.WaitGroup
		numReads := 1000

		for i := 0; i < numReads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				require.Equal(t, updates, priceCache.GetPriceUpdates())
				require.Equal(t, int64(100), priceCache.GetHeight())
				require.Equal(t, int32(5), priceCache.GetRound())
				require.True(t, priceCache.HasValidValues(100, 5))
			}()
		}

		wg.Wait()
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
