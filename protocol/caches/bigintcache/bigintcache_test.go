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
	dummyBytesEmpty := []byte{}
	dummyBytes1 := []byte{1}
	dummyBytes2 := []byte{1, 2, 3, 4}
	dummyBytes3 := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	t.Run("valid: set value for single round + height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		value := big.NewInt(1000)
		cache.SetValue(ctx, value, dummyBytes2)
		checkValidCacheState(t, &cache, dummyBytes2, value)
	})

	t.Run("valid: set value for multi round single height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(2)
		value1 := big.NewInt(2000)
		cache.SetValue(ctx, value1, dummyBytes1)
		checkValidCacheState(t, &cache, dummyBytes1, value1)

		value2 := big.NewInt(3000)
		cache.SetValue(ctx, value2, dummyBytes2)
		checkValidCacheState(t, &cache, dummyBytes2, value2)
	})

	t.Run("valid: set value for single rounds multi height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(3)
		value1 := big.NewInt(4000)
		cache.SetValue(ctx, value1, dummyBytesEmpty)
		checkValidCacheState(t, &cache, dummyBytesEmpty, value1)

		ctx = ctx.WithBlockHeight(4)
		value2 := big.NewInt(5000)
		cache.SetValue(ctx, value2, dummyBytes3)
		checkValidCacheState(t, &cache, dummyBytes3, value2)
	})

	t.Run("invalid: No valid value, wrong tx hash", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(6)
		value := big.NewInt(8000)
		cache.SetValue(ctx, value, dummyBytes1)
		require.False(t, cache.HasValidValue(dummyBytes2))
	})

	t.Run("invalid: No valid value, tx hash (empty)", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(8)
		value := big.NewInt(9000)
		cache.SetValue(ctx, value, dummyBytesEmpty)
		require.False(t, cache.HasValidValue(dummyBytes3))
	})
}

func checkValidCacheState(
	t *testing.T,
	cache bigintcache.BigIntCache,
	shouldBeTxHash []byte,
	shouldBeValue *big.Int,
) {
	require.True(t, cache.HasValidValue(shouldBeTxHash))
	require.Equal(t, shouldBeValue, cache.GetValue())
}
