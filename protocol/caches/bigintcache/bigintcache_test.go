package bigintcache_test

import (
	"math/big"
	"sync"
	"testing"
	"time"

	bigintcache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/bigintcache"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestBigIntCaching(t *testing.T) {
	cache := bigintcache.BigIntCacheImpl{}
	ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

	t.Run("valid: set value for single round + height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		value := big.NewInt(1000)
		cache.SetValue(ctx, value, 1)
		checkValidCacheState(t, &cache, 1, 1, value)
	})

	t.Run("valid: set value for multi round single height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(2)
		value1 := big.NewInt(2000)
		cache.SetValue(ctx, value1, 1)
		checkValidCacheState(t, &cache, 2, 1, value1)

		value2 := big.NewInt(3000)
		cache.SetValue(ctx, value2, 2)
		checkValidCacheState(t, &cache, 2, 2, value2)
	})

	t.Run("valid: set value for single rounds multi height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(3)
		value1 := big.NewInt(4000)
		cache.SetValue(ctx, value1, 1)
		checkValidCacheState(t, &cache, 3, 1, value1)

		ctx = ctx.WithBlockHeight(4)
		value2 := big.NewInt(5000)
		cache.SetValue(ctx, value2, 1)
		checkValidCacheState(t, &cache, 4, 1, value2)
	})

	t.Run("valid: set diff value for same height + round", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(5)
		value1 := big.NewInt(6000)
		cache.SetValue(ctx, value1, 1)
		checkValidCacheState(t, &cache, 5, 1, value1)

		value2 := big.NewInt(7000)
		cache.SetValue(ctx, value2, 1)
		checkValidCacheState(t, &cache, 5, 1, value2)
	})

	t.Run("invalid: No valid value, wrong block height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(6)
		value := big.NewInt(8000)
		cache.SetValue(ctx, value, 1)
		require.False(t, cache.HasValidValue(7, 1))
	})

	t.Run("invalid: No valid value, wrong round", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(8)
		value := big.NewInt(9000)
		cache.SetValue(ctx, value, 1)
		require.False(t, cache.HasValidValue(8, 2))
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
					value := big.NewInt(int64(i * 1000))
					cache.SetValue(ctx.WithBlockHeight(int64(i)), value, int32(i))
				} else {
					// Odd goroutines read
					time.Sleep(time.Millisecond) // Slight delay to increase chances of interleaving
					_ = cache.GetValue()
					_ = cache.GetHeight()
					_ = cache.GetRound()
					_ = cache.HasValidValue(int64(i), int32(i))
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
				value := big.NewInt(int64(i * 1000))
				cache.SetValue(ctx.WithBlockHeight(int64(i)), value, int32(i))
			}(i)
		}

		wg.Wait()

		// Verify the final state
		height := cache.GetHeight()
		round := cache.GetRound()
		require.True(t, height >= 0 && height < int64(numWrites))
		require.True(t, round >= 0 && round < int32(numWrites))
	})

	t.Run("concurrent reads", func(t *testing.T) {
		// Set initial state
		initialValue := big.NewInt(1000000)
		cache.SetValue(ctx.WithBlockHeight(100), initialValue, 5)

		var wg sync.WaitGroup
		numReads := 1000

		for i := 0; i < numReads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				require.Equal(t, initialValue, cache.GetValue())
				require.Equal(t, int64(100), cache.GetHeight())
				require.Equal(t, int32(5), cache.GetRound())
				require.True(t, cache.HasValidValue(100, 5))
			}()
		}

		wg.Wait()
	})
}

func checkValidCacheState(
	t *testing.T,
	cache bigintcache.BigIntCache,
	shouldBeHeight int64,
	shouldBeRound int32,
	shouldBeValue *big.Int,
) {
	require.True(t, cache.HasValidValue(shouldBeHeight, shouldBeRound))
	require.Equal(t, shouldBeHeight, cache.GetHeight())
	require.Equal(t, shouldBeRound, cache.GetRound())
	require.Equal(t, shouldBeValue, cache.GetValue())
}
