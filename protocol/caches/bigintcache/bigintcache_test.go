package bigintcache_test

import (
	"math/big"
	"testing"

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
