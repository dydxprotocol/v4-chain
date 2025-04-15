package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
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

	require.True(t, ks.ClobKeeper.GetMemstoreInitialized(ks.Ctx))
}

func TestInitMemStore_StatefulOrderCount(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContextWithUninitializedMemStore(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	// Long term order.
	ks.ClobKeeper.SetLongTermOrderPlacement(
		ks.Ctx,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
		1,
	)

	// Triggered conditional order.
	ks.ClobKeeper.SetLongTermOrderPlacement(
		ks.Ctx,
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
		1,
	)
	ks.ClobKeeper.MustTriggerConditionalOrder(
		ks.Ctx,
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId,
	)

	// Untriggered conditional order.
	ks.ClobKeeper.SetLongTermOrderPlacement(
		ks.Ctx,
		constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_StopLoss20, // Clob 1
		1,
	)

	// Reset the stateful order count to zero.
	ks.ClobKeeper.SetStatefulOrderCount(
		ks.Ctx,
		constants.Alice_Num0,
		0,
	)

	// InitMemStore should repopulate the count.
	ks.ClobKeeper.InitMemStore(ks.Ctx)
	require.Equal(t, uint32(3), ks.ClobKeeper.GetStatefulOrderCount(ks.Ctx, constants.Alice_Num0))
}
