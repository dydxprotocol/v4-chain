package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4/mocks"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx, keeper, _, _, _, _, _, _ := keepertest.ClobKeepers(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{})
	logger := keeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestInitMemStore_OnlyAllowedOnce(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx, k, _, _, _, _, _, _ := keepertest.ClobKeepersWithUninitializedMemStore(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{})

	k.InitMemStore(ctx)

	// Initializing a second time causes a panic
	require.Panics(t, func() {
		k.InitMemStore(ctx)
	})
}
