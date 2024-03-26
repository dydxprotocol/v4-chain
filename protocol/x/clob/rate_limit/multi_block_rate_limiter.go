package rate_limit

import (
	"fmt"
	"sort"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ RateLimiter[int] = (*multiBlockRateLimiter[int])(nil)

// sortedMaxPerNBlocksRateLimit is type alias for a []types.MaxPerNBlocksRateLimit
// sorted in ascending order by types.MaxPerNBlocksRateLimit#NumBlocks.
type sortedMaxPerNBlocksRateLimit []types.MaxPerNBlocksRateLimit

// A RateLimiter over multiple blocks.
//
// `ctx.BlockHeight()` is used in RateLimit to track which block the rate limit should apply to.
//
// `ctx.BlockHeight()` is used in PruneRateLimits to remove rate limits for a block `maxNumBlocks` ago.
// If invoked during `EndBlocker` then you can pass in the `ctx` as is but if invoked during `PrepareCheckState`
// one must supply a `ctx` with the previous block height via `ctx.WithBlockHeight(ctx.BlockHeight()-1)`.
type multiBlockRateLimiter[K comparable] struct {
	// Context to use for rate limiting errors returned to the client.
	context string
	// A list of rate limits, sorted in ascending NumBlocks order.
	config sortedMaxPerNBlocksRateLimit
	// The count per rate limit.
	perKeyRateLimitCounts map[K][]uint32
	// The count per block.
	perKeyBlockCounts map[K]map[uint32]uint32
	// The rate limit with the most number of blocks.
	maxNumBlocks uint32
	// Which keys need to be pruned for a specific block, implemented as a circular array where
	// offset = height % maxNumBlocks.
	dirtyPerBlock []map[K]bool
}

// Returns a RateLimiter over multiple blocks.
//
// `ctx.BlockHeight()` is used in RateLimit to track which block the rate limit should apply to.
//
// `ctx.BlockHeight()` is used in PruneRateLimits to remove rate limits for a block `maxNumBlocks` ago.
// If invoked during `EndBlocker` then you can pass in the `ctx` as is but if invoked during `PrepareCheckState`
// one must supply a `ctx` with the previous block height via `ctx.WithBlockHeight(ctx.BlockHeight()-1)`.
func NewMultiBlockRateLimiter[K comparable](context string, config []types.MaxPerNBlocksRateLimit) RateLimiter[K] {
	// Ensure that we sort the number of blocks so that we are checking the lowest block number rate limits first.
	sort.Slice(config, func(i, j int) bool {
		return config[i].NumBlocks < config[j].NumBlocks
	})

	// The last rate limit will have the maximum number of blocks.
	maxNumBlocks := config[len(config)-1].NumBlocks

	return &multiBlockRateLimiter[K]{
		context:               context,
		config:                config,
		perKeyRateLimitCounts: make(map[K][]uint32),
		perKeyBlockCounts:     make(map[K]map[uint32]uint32),
		maxNumBlocks:          maxNumBlocks,
		dirtyPerBlock:         make([]map[K]bool, maxNumBlocks),
	}
}

func (r *multiBlockRateLimiter[K]) RateLimit(ctx sdk.Context, key K) error {
	return r.RateLimitIncrBy(ctx, key, 1)
}

func (r *multiBlockRateLimiter[K]) RateLimitIncrBy(ctx sdk.Context, key K, incrBy uint32) error {
	blockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())
	offset := blockHeight % r.maxNumBlocks

	// Mark this key as the one that needs to be cleaned up at the next block.
	dirty := r.dirtyPerBlock[offset]
	if dirty == nil {
		dirty = make(map[K]bool)
		r.dirtyPerBlock[offset] = dirty
	}
	dirty[key] = true

	// Update the per block count.
	perBlockCounts, found := r.perKeyBlockCounts[key]
	if !found {
		perBlockCounts = make(map[uint32]uint32)
		r.perKeyBlockCounts[key] = perBlockCounts
	}
	count := perBlockCounts[blockHeight] + incrBy
	perBlockCounts[blockHeight] = count

	// Update the per rate limit count.
	perRateLimitCounts, found := r.perKeyRateLimitCounts[key]
	if !found {
		perRateLimitCounts = make([]uint32, len(r.config))
		r.perKeyRateLimitCounts[key] = perRateLimitCounts
	}
	for i := range perRateLimitCounts {
		perRateLimitCounts[i] += incrBy
	}

	// Check the accumulated rate limit count to see if any rate limit has been exceeded.
	for i, rl := range r.config {
		if perRateLimitCounts[i] > rl.Limit {
			return errorsmod.Wrapf(
				types.ErrBlockRateLimitExceeded,
				"Rate of %d exceeds configured block rate limit of %+v for %s and %+v",
				perRateLimitCounts[i],
				rl,
				r.context,
				key,
			)
		}
	}
	return nil
}

func (r *multiBlockRateLimiter[K]) PruneRateLimits(ctx sdk.Context) {
	blockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())

	// For each rate limit we need to update the per rate limit counts subtracting away how many orders were placed
	// rl.NumBlock ago.
	for i, rl := range r.config {
		// If the block height we need to update for would be negative then stop updating
		// We can do this since r.config is sorted by NumBlocks in ascending order.
		if blockHeight < rl.NumBlocks {
			break
		}
		offset := (blockHeight - rl.NumBlocks) % r.maxNumBlocks
		// For each key we need to update the perKeyRateLimitCounts based upon how many orders
		// were for the key by using perKeyBlockCounts. Afterwards we can remove
		// the perKeyBlockCounts for blockHeight - r.maxNumBlocks.
		for key := range r.dirtyPerBlock[offset] {
			perBlockCounts, found := r.perKeyBlockCounts[key]
			if !found {
				panic(fmt.Sprintf(
					"Expected to have found perBlockCounts for %+v because it was marked dirty",
					key,
				))
			}

			perRateLimitCounts, found := r.perKeyRateLimitCounts[key]
			if !found {
				panic(fmt.Sprintf(
					"Expected to have found perRateLimitCounts for %+v because it was marked dirty",
					key,
				))
			}
			// Subtract away how many orders were placed RateLimit.NumBlocks ago for the perRateLimitCounts.
			perRateLimitCounts[i] -= perBlockCounts[blockHeight-rl.NumBlocks]
		}
	}

	// After cleaning up the per rate limit counts we can clear any per block data we stored for the old
	// block making way for the next block's per block data.
	offset := blockHeight % r.maxNumBlocks
	for key := range r.dirtyPerBlock[offset] {
		perBlockCounts, found := r.perKeyBlockCounts[key]
		if !found {
			panic(fmt.Sprintf(
				"Expected to have found perBlockCounts for %+v because it was marked dirty",
				key,
			))
		}

		delete(perBlockCounts, blockHeight-r.maxNumBlocks)
		if len(perBlockCounts) == 0 {
			delete(r.perKeyBlockCounts, key)
			// If there are no block perKeyCounts then the per rate limit counts must also all be zero.
			delete(r.perKeyRateLimitCounts, key)
		}
	}
	r.dirtyPerBlock[offset] = make(map[K]bool)
}
