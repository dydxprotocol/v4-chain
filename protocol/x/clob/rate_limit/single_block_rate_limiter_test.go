package rate_limit_test

import (
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSingleBlockRateLimiter_InvalidNumBlocks(t *testing.T) {
	require.Panics(t, func() {
		rate_limit.NewSingleBlockRateLimiter[string]("test", types.MaxPerNBlocksRateLimit{
			NumBlocks: 2,
			Limit:     1,
		})
	}, "Expected NumBlocks == 1")
}

func TestSingleBlockRateLimiter(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	rl := rate_limit.NewSingleBlockRateLimiter[string]("test", types.MaxPerNBlocksRateLimit{
		NumBlocks: 1,
		Limit:     10,
	})

	for block := int64(1); block < 3; block += 1 {
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
		// First and subsequent requests should error.
		require.Error(t, rl.RateLimit(ctx, "B"), "Rate of %d exceeds configured block rate limit", 11)
		require.Error(t, rl.RateLimit(ctx, "B"), "Rate of %d exceeds configured block rate limit", 12)

		// The next iteration of the loop should allow 10 more and then fail again.
	}
}
