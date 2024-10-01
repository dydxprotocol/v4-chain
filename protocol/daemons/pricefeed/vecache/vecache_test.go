package vecache_test

import (
	"math/big"
	"sync"
	"testing"
	"time"

	vecache "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/vecache"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestVEPriceCaching(t *testing.T) {
	veCache := vecache.VeUpdatesCache{}
	ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

	t.Run("valid: set price updates for single round + height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		var updates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &veCache, 1, 1, updates)
	})

	t.Run("valid: set priced updates for multi round single height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(2)
		var updates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			updates = append(updates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, updates, 1)

		checkValidCacheState(t, &veCache, 2, 1, updates)
		var newUpdates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, newUpdates, 2)
		checkValidCacheState(t, &veCache, 2, 2, newUpdates)
	})

	t.Run("valid: set price updates for single rounds multi height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(3)
		var updates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &veCache, 3, 1, updates)

		ctx = ctx.WithBlockHeight(4)
		var newUpdates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, newUpdates, 1)
		checkValidCacheState(t, &veCache, 4, 1, newUpdates)
	})

	t.Run("valid: set diff update for same height + round", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(5)
		var updates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &veCache, 5, 1, updates)

		var newUpdates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, newUpdates, 1)
		checkValidCacheState(t, &veCache, 5, 1, newUpdates)
	})

	t.Run("invalid: No valid values, wrong block height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(6)
		var updates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, updates, 1)
		require.False(t, veCache.HasValidValues(7, 1))
	})

	t.Run("invalid: No valid values, wrong round", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(8)
		var updates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx, updates, 1)
		require.False(t, veCache.HasValidValues(8, 2))
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
					var updates vecache.PriceUpdates
					for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
						updates = append(updates, vecache.PriceUpdate{
							MarketId:  marketPrice.MarketId,
							SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
							PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
						})
					}
					veCache.SetPriceUpdates(ctx.WithBlockHeight(int64(i)), updates, int32(i))
				} else {
					// Odd goroutines read
					time.Sleep(time.Millisecond) // Slight delay to increase chances of interleaving
					_ = veCache.GetPriceUpdates()
					_ = veCache.GetHeight()
					_ = veCache.GetRound()
					_ = veCache.HasValidValues(int64(i), int32(i))
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
				var updates vecache.PriceUpdates
				for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
					updates = append(updates, vecache.PriceUpdate{
						MarketId:  marketPrice.MarketId,
						SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
						PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
					})
				}
				veCache.SetPriceUpdates(ctx.WithBlockHeight(int64(i)), updates, int32(i))
			}(i)
		}

		wg.Wait()

		// Verify the final state
		height := veCache.GetHeight()
		round := veCache.GetRound()
		require.True(t, height >= 0 && height < int64(numWrites))
		require.True(t, round >= 0 && round < int32(numWrites))
	})

	t.Run("concurrent reads", func(t *testing.T) {
		// Set initial state
		var updates vecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, vecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		veCache.SetPriceUpdates(ctx.WithBlockHeight(100), updates, 5)

		var wg sync.WaitGroup
		numReads := 1000

		for i := 0; i < numReads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				require.Equal(t, updates, veCache.GetPriceUpdates())
				require.Equal(t, int64(100), veCache.GetHeight())
				require.Equal(t, int32(5), veCache.GetRound())
				require.True(t, veCache.HasValidValues(100, 5))
			}()
		}

		wg.Wait()
	})
}

func checkValidCacheState(
	t *testing.T,
	veCache *vecache.VeUpdatesCache,
	shouldBeHight int64,
	shouldBeRound int32,
	shouldBeUpdates vecache.PriceUpdates,
) {
	require.True(t, veCache.HasValidValues(shouldBeHight, shouldBeRound))
	require.Equal(t, shouldBeHight, veCache.GetHeight())
	require.Equal(t, shouldBeRound, veCache.GetRound())
	require.Equal(t, shouldBeUpdates, veCache.GetPriceUpdates())
}
