package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// A node bootstrapped through state sync starts with an empty store, so the rate limiters the app
// installs at startup from state are still the placeholders when the first block arrives.
// Initialize must install the real limiters from the restored state.
func TestInitialize_InstallsRateLimitersFromState(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		indexer_manager.NewIndexerEventManagerNoop(),
	)
	// Rate limiting only runs in CheckTx.
	checkCtx := ks.Ctx.WithIsCheckTx(true)
	msg := &types.MsgUpdateLeverage{SubaccountId: &constants.Alice_Num0}

	// No config in state yet: Initialize leaves the placeholder limiter in place, which never limits.
	ks.ClobKeeper.Initialize(ks.Ctx)
	require.NoError(t, ks.ClobKeeper.RateLimitUpdateLeverage(checkCtx, msg))
	require.NoError(t, ks.ClobKeeper.RateLimitUpdateLeverage(checkCtx, msg))

	// The config lands in state without InitGenesis or governance, as a restored snapshot would.
	config := types.BlockRateLimitConfiguration{
		MaxShortTermOrdersAndCancelsPerNBlocks: []types.MaxPerNBlocksRateLimit{{NumBlocks: 1, Limit: 1}},
		MaxLeverageUpdatesPerNBlocks:           []types.MaxPerNBlocksRateLimit{{NumBlocks: 1, Limit: 1}},
	}
	ks.Ctx.KVStore(ks.StoreKey).Set([]byte(types.BlockRateLimitConfigKey), ks.Cdc.MustMarshal(&config))

	ks.ClobKeeper.Initialize(ks.Ctx)
	require.NoError(t, ks.ClobKeeper.RateLimitUpdateLeverage(checkCtx, msg))
	require.ErrorIs(t, ks.ClobKeeper.RateLimitUpdateLeverage(checkCtx, msg), types.ErrBlockRateLimitExceeded)

	// Later calls do not rebuild the limiters and lose the counts.
	ks.ClobKeeper.Initialize(ks.Ctx)
	require.ErrorIs(t, ks.ClobKeeper.RateLimitUpdateLeverage(checkCtx, msg), types.ErrBlockRateLimitExceeded)
}
