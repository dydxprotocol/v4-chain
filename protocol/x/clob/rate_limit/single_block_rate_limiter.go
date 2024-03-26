package rate_limit

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ RateLimiter[int] = (*singleBlockRateLimiter[int])(nil)

// A RateLimiter optimized when there is only a single block to record and capture information for.
type singleBlockRateLimiter[K comparable] struct {
	context      string
	config       types.MaxPerNBlocksRateLimit
	perKeyCounts map[K]uint32
}

func NewSingleBlockRateLimiter[K comparable](context string, config types.MaxPerNBlocksRateLimit) RateLimiter[K] {
	if config.NumBlocks != 1 {
		panic(fmt.Sprintf(
			"Expected NumBlocks == 1 but found %d, "+
				"did you mean to use NewMultiBlockRateLimiter?",
			config.NumBlocks,
		))
	}
	return &singleBlockRateLimiter[K]{
		context:      context,
		config:       config,
		perKeyCounts: make(map[K]uint32),
	}
}

func (r *singleBlockRateLimiter[K]) RateLimit(ctx sdk.Context, key K) error {
	return r.RateLimitIncrBy(ctx, key, 1)
}

func (r *singleBlockRateLimiter[K]) RateLimitIncrBy(ctx sdk.Context, key K, incrBy uint32) error {
	count := r.perKeyCounts[key] + incrBy
	r.perKeyCounts[key] = count
	if count > r.config.Limit {
		return errorsmod.Wrapf(
			types.ErrBlockRateLimitExceeded,
			"Rate of %d exceeds configured block rate limit of %+v for %s and key %+v",
			count,
			r.config,
			r.context,
			key,
		)
	}
	return nil
}

func (r *singleBlockRateLimiter[K]) PruneRateLimits(ctx sdk.Context) {
	// Initialize the new map to be about half as big as the last one so that we limit how
	// many times it is resized if the number of RateLimit invocations remains consistent but
	// still shrink if the number of RateLimit invocations decreases.
	r.perKeyCounts = make(map[K]uint32, len(r.perKeyCounts)/2)
}
