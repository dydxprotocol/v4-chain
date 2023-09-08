package abci

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RunCached wraps a function with a cache context and writes the cache to the context if the
// function call succeeds. If the function call fails, the cache is discarded.
func RunCached(c sdk.Context, f func(sdk.Context) error) error {
	ctx, writeCache := c.CacheContext()

	if err := f(ctx); err != nil {
		return err
	}

	writeCache()
	return nil
}
