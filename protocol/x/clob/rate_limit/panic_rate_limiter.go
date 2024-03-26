package rate_limit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ RateLimiter[int] = (*noOpRateLimiter[int])(nil)

// A RateLimiter that panics on any invocation.
type panicRateLimiter[K any] struct {
}

// Returns a RateLimiter that will panic on invocation.
func NewPanicRateLimiter[K any]() RateLimiter[K] {
	return panicRateLimiter[K]{}
}

func (n panicRateLimiter[K]) RateLimit(ctx sdk.Context, key K) error {
	panic("Unexpected invocation of RateLimit")
}

func (n panicRateLimiter[K]) RateLimitIncrBy(ctx sdk.Context, key K, incrBy uint32) error {
	panic("Unexpected invocation of RateLimitIncrBy")
}

func (n panicRateLimiter[K]) PruneRateLimits(ctx sdk.Context) {
	panic("Unexpected invocation of PruneRateLimits")
}
