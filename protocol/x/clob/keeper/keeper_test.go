package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{})
	logger := ks.ClobKeeper.Logger(ks.Ctx)
	require.NotNil(t, logger)
}

func TestInitMemStore_OnlyAllowedOnce(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContextWithUninitializedMemStore(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{})

	ks.ClobKeeper.InitMemStore(ks.Ctx)

	// Initializing a second time causes a panic
	require.Panics(t, func() {
		ks.ClobKeeper.InitMemStore(ks.Ctx)
	})
}
