package rate_limit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Used to RateLimit for a key K.
type RateLimiter[K any] interface {
	// Returns an error if the RateLimiter exceeds any configured rate limits for the key K and context state.
	RateLimit(ctx sdk.Context, key K) error
	// Returns an error if the RateLimiter exceeds any configured rate limits for the key K and context state
	// Increments the rate limit counter by incrBy.
	RateLimitIncrBy(ctx sdk.Context, key K, incrBy uint32) error
	// Prunes rate limits for the provided context.
	PruneRateLimits(ctx sdk.Context)
}
