package rate_limit_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestMultiBlockRateLimiter_SingleRateLimit(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	rl := rate_limit.NewMultiBlockRateLimiter[string]("test", []types.MaxPerNBlocksRateLimit{
		{
			NumBlocks: 2,
			Limit:     10,
		},
	})

	for block := int64(1); block < 10; block += 1 {
		ctx = ctx.WithBlockHeight(block)
		rl.PruneRateLimits(ctx)

		block += 1
		ctx = ctx.WithBlockHeight(block)
		rl.PruneRateLimits(ctx)

		for i := 0; i < 10; i += 1 {
			require.NoError(t, rl.RateLimit(ctx, "A"))
		}
		// First and subsequent requests should error.
		require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 11)
		require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 12)

		// The rate limit should be different for this other key.
		for i := 0; i < 10; i += 1 {
			require.NoError(t, rl.RateLimit(ctx, "B"))
		}
		// First and subsequent requests should error
		require.Error(t, rl.RateLimit(ctx, "B"), "Rate of %d exceeds configured block rate limit", 11)
		require.Error(t, rl.RateLimit(ctx, "B"), "Rate of %d exceeds configured block rate limit", 12)

		// The next iteration of the loop should allow 10 more and then fail again.
	}
}

func TestMultiBlockRateLimiter_MultipleRateLimits(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	rl := rate_limit.NewMultiBlockRateLimiter[string]("test", []types.MaxPerNBlocksRateLimit{
		{
			NumBlocks: 1,
			Limit:     10,
		},
		{
			NumBlocks: 3,
			Limit:     20,
		},
	})

	ctx = ctx.WithBlockHeight(1)
	rl.PruneRateLimits(ctx)

	for i := 0; i < 10; i += 1 {
		require.NoError(t, rl.RateLimit(ctx, "A"))
	}
	// First and subsequent requests should error for the first rate limit.
	require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 11)
	require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 12)

	ctx = ctx.WithBlockHeight(2)
	rl.PruneRateLimits(ctx)

	// This should bring up the count to 18 requests for second rate limit.
	for i := 0; i < 6; i += 1 {
		require.NoError(t, rl.RateLimit(ctx, "A"))
	}

	ctx = ctx.WithBlockHeight(3)
	rl.PruneRateLimits(ctx)

	// 2 more should bring up the second rate limit to 20.
	for i := 0; i < 2; i += 1 {
		require.NoError(t, rl.RateLimit(ctx, "A"))
	}
	// First and subsequent requests should error for the second rate limit
	require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 21)
	require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 22)
	require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 23)

	// Pruning will reduce the second rate limit to a count of 23 to 11 (since 12 were done in the first block).
	ctx = ctx.WithBlockHeight(4)
	rl.PruneRateLimits(ctx)

	// 9 more should bring up the second rate limit to 20.
	for i := 0; i < 9; i += 1 {
		require.NoError(t, rl.RateLimit(ctx, "A"))
	}
	// First request should error for the second rate limit and second request should error for the first rate limit.
	require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 21)
	require.Error(t, rl.RateLimit(ctx, "A"), "Rate of %d exceeds configured block rate limit", 11)
}
