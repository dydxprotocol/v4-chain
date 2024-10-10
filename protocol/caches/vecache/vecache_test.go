package vecache_test

import (
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
}
