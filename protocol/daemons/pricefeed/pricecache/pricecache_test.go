package pricecache_test

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	pricecache "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/pricecache"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestVEPriceCaching(t *testing.T) {
	pc := pricecache.PriceCache{}
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
		pc.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &pc, 1, 1, updates)
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
		pc.SetPriceUpdates(ctx, updates, 1)

		checkValidCacheState(t, &pc, 2, 1, updates)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		pc.SetPriceUpdates(ctx, newUpdates, 2)
		checkValidCacheState(t, &pc, 2, 2, newUpdates)
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
		pc.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &pc, 3, 1, updates)

		ctx = ctx.WithBlockHeight(4)
		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		pc.SetPriceUpdates(ctx, newUpdates, 1)
		checkValidCacheState(t, &pc, 4, 1, newUpdates)
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
		pc.SetPriceUpdates(ctx, updates, 1)
		checkValidCacheState(t, &pc, 5, 1, updates)

		var newUpdates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidSingleMarketPriceUpdate {
			newUpdates = append(newUpdates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		pc.SetPriceUpdates(ctx, newUpdates, 1)
		checkValidCacheState(t, &pc, 5, 1, newUpdates)
	})

	t.Run("invalid: No valid prices, wrong block height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(6)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		pc.SetPriceUpdates(ctx, updates, 1)
		require.False(t, pc.HasValidPrices(7, 1))
	})

	t.Run("invalid: No valid prices, wrong round", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(8)
		var updates pricecache.PriceUpdates
		for _, marketPrice := range constants.ValidUpdateMarketPrices.MarketPriceUpdates {
			updates = append(updates, pricecache.PriceUpdate{
				MarketId:  marketPrice.MarketId,
				SpotPrice: big.NewInt(int64(marketPrice.SpotPrice)),
				PnlPrice:  big.NewInt(int64(marketPrice.PnlPrice)),
			})
		}
		pc.SetPriceUpdates(ctx, updates, 1)
		require.False(t, pc.HasValidPrices(8, 2))
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
					pc.SetPriceUpdates(ctx.WithBlockHeight(int64(i)), updates, int32(i))
				} else {
					// Odd goroutines read
					time.Sleep(time.Millisecond) // Slight delay to increase chances of interleaving
					_ = pc.GetPriceUpdates()
					_ = pc.GetHeight()
					_ = pc.GetRound()
					_ = pc.HasValidPrices(int64(i), int32(i))
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
				pc.SetPriceUpdates(ctx.WithBlockHeight(int64(i)), updates, int32(i))
			}(i)
		}

		wg.Wait()

		// Verify the final state
		height := pc.GetHeight()
		round := pc.GetRound()
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
		pc.SetPriceUpdates(ctx.WithBlockHeight(100), updates, 5)

		var wg sync.WaitGroup
		numReads := 1000

		for i := 0; i < numReads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				require.Equal(t, updates, pc.GetPriceUpdates())
				require.Equal(t, int64(100), pc.GetHeight())
				require.Equal(t, int32(5), pc.GetRound())
				require.True(t, pc.HasValidPrices(100, 5))
			}()
		}

		wg.Wait()
	})
}

func checkValidCacheState(
	t *testing.T,
	pc *pricecache.PriceCache,
	shouldBeHight int64,
	shouldBeRound int32,
	shouldBeUpdates pricecache.PriceUpdates,
) {
	fmt.Println()
	require.True(t, pc.HasValidPrices(shouldBeHight, shouldBeRound))
	require.Equal(t, shouldBeHight, pc.GetHeight())
	require.Equal(t, shouldBeRound, pc.GetRound())
	require.Equal(t, shouldBeUpdates, pc.GetPriceUpdates())
}
