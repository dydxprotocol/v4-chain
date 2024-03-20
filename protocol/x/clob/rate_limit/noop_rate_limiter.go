package rate_limit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ RateLimiter[int] = (*noOpRateLimiter[int])(nil)

// A RateLimiter that does not enforce any limits.
type noOpRateLimiter[K any] struct {
}

// Returns a RateLimiter that does not enforce any limits.
func NewNoOpRateLimiter[K any]() RateLimiter[K] {
	return noOpRateLimiter[K]{}
}

func (n noOpRateLimiter[K]) RateLimit(ctx sdk.Context, key K) error {
	return nil
}

func (n noOpRateLimiter[K]) RateLimitIncrBy(ctx sdk.Context, key K, incrBy uint32) error {
	return nil
}

func (n noOpRateLimiter[K]) PruneRateLimits(ctx sdk.Context) {
}
