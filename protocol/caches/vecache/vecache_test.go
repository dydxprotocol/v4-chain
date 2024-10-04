package vecache_test

import (
	"sync"
	"testing"

	vecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/vecache"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestVECache(t *testing.T) {
	veCache := vecache.NewVECache()
	ctx := sdk.Context{}.WithBlockHeight(0)

	t.Run("valid: set and get seen votes in cache", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		consAddresses := map[string]struct{}{
			"address1": {},
			"address2": {},
		}
		veCache.SetSeenVotesInCache(ctx, consAddresses)
		require.Equal(t, int64(1), veCache.GetHeight())
		require.Equal(t, consAddresses, veCache.GetSeenVotesInCache())
	})

	t.Run("valid: update seen votes in cache", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(2)
		consAddresses := map[string]struct{}{
			"address3": {},
		}
		veCache.SetSeenVotesInCache(ctx, consAddresses)
		require.Equal(t, int64(2), veCache.GetHeight())
		require.Equal(t, consAddresses, veCache.GetSeenVotesInCache())
	})

	t.Run("concurrent reads and writes", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if i%2 == 0 {
					// Even goroutines write
					consAddresses := map[string]struct{}{
						"address1": {},
						"address2": {},
					}
					veCache.SetSeenVotesInCache(ctx.WithBlockHeight(int64(i)), consAddresses)
				} else {
					// Odd goroutines read
					_ = veCache.GetSeenVotesInCache()
					_ = veCache.GetHeight()
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("concurrent writes", func(t *testing.T) {
		var wg sync.WaitGroup
		numWrites := 1000

		for i := 0; i < numWrites; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				consAddresses := map[string]struct{}{
					"address1": {},
					"address2": {},
				}
				veCache.SetSeenVotesInCache(ctx.WithBlockHeight(int64(i)), consAddresses)
			}(i)
		}

		wg.Wait()

		// Verify the final state
		height := veCache.GetHeight()
		require.True(t, height >= 0 && height < int64(numWrites))
	})

	t.Run("concurrent reads", func(t *testing.T) {
		// Set initial state
		consAddresses := map[string]struct{}{
			"address1": {},
			"address2": {},
		}
		veCache.SetSeenVotesInCache(ctx.WithBlockHeight(100), consAddresses)

		var wg sync.WaitGroup
		numReads := 1000

		for i := 0; i < numReads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				require.Equal(t, consAddresses, veCache.GetSeenVotesInCache())
				require.Equal(t, int64(100), veCache.GetHeight())
			}()
		}

		wg.Wait()
	})
}
