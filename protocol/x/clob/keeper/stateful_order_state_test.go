package keeper_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tracer"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func orderToStringId(
	o types.Order,
) string {
	return string(o.OrderId.ToStateKey())
}

func orderToStringSubaccountId(
	o types.Order,
) string {
	return string(o.OrderId.SubaccountId.ToStateKey())
}

// TODO(jonfung) make ticket and remove all conditional orderes
func createPartiallyFilledStatefulOrderInState(
	ctx sdk.Context,
	k keeper.Keeper,
	order types.Order,
	timestamp time.Time,
) {
	k.MustAddOrderToStatefulOrdersTimeSlice(
		ctx,
		timestamp,
		order.OrderId,
	)
	k.SetOrderFillAmount(ctx, order.OrderId, satypes.BaseQuantums(10), uint32(20))
	k.SetLongTermOrderPlacement(ctx, order, uint32(30))
}

func TestLongTermOrderInitMemStore_Success(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContextWithUninitializedMemStore(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	triggeredConditionalOrderStore := ks.ClobKeeper.GetTriggeredConditionalOrderPlacementStore(ks.Ctx)
	untriggeredConditionalOrderStore := ks.ClobKeeper.GetUntriggeredConditionalOrderPlacementStore(ks.Ctx)
	longTermOrderStore := ks.ClobKeeper.GetLongTermOrderPlacementStore(ks.Ctx)

	// Set orders only on the store, not the memstore.
	index := uint32(0)
	storeOrder := func(order types.Order, store prefix.Store) {
		longTermOrderPlacement := types.LongTermOrderPlacement{
			Order: order,
			PlacementIndex: types.TransactionOrdering{
				BlockHeight:      0,
				TransactionIndex: index,
			},
		}
		longTermOrderPlacementBytes := ks.Cdc.MustMarshal(&longTermOrderPlacement)
		orderKey := order.OrderId.ToStateKey()
		store.Set(orderKey, longTermOrderPlacementBytes)
		index++
	}

	// Set some long term orders.
	storeOrder(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20, longTermOrderStore)
	storeOrder(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25, longTermOrderStore)

	// Set a untriggered conditional order.
	storeOrder(
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
		untriggeredConditionalOrderStore,
	)

	// Set a triggered conditional order.
	storeOrder(
		constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
		triggeredConditionalOrderStore,
	)

	// Init the memstore.
	ks.ClobKeeper.InitMemStore(ks.Ctx)

	// Assert that the values can be read after memStore has been warmed.
	order, exists := ks.ClobKeeper.GetLongTermOrderPlacement(
		ks.Ctx, constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId)
	require.True(t, exists)
	require.Equal(t, constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20, order.Order)

	order, exists = ks.ClobKeeper.GetLongTermOrderPlacement(
		ks.Ctx, constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId)
	require.True(t, exists)
	require.Equal(t, constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25, order.Order)

	order, exists = ks.ClobKeeper.GetLongTermOrderPlacement(
		ks.Ctx, constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25.OrderId)
	require.False(t, exists)

	order, exists = ks.ClobKeeper.GetLongTermOrderPlacement(
		ks.Ctx, constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId)
	require.True(t, exists)
	require.Equal(t, constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20, order.Order)

	order, exists = ks.ClobKeeper.GetLongTermOrderPlacement(
		ks.Ctx, constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20.OrderId)
	require.True(t, exists)
	require.Equal(t, constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20, order.Order)

	order, exists = ks.ClobKeeper.GetLongTermOrderPlacement(
		ks.Ctx, constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId)
	require.False(t, exists)
}

func TestMustTriggerConditionalOrder(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	// Set the tracer on the multistore to verify the performed writes are correct.
	traceDecoder := &tracer.TraceDecoder{}
	ks.Ctx.MultiStore().SetTracer(traceDecoder)

	conditionalOrder := constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20

	// Verify you can create an untriggered conditional order that did not previously exist.
	ks.ClobKeeper.SetLongTermOrderPlacement(ks.Ctx, conditionalOrder, 0)

	// Verify you can trigger the untriggered conditional order.
	ks.ClobKeeper.MustTriggerConditionalOrder(ks.Ctx, conditionalOrder.OrderId)

	// Verify you can obtain the long term order placement.
	longTermOrderPlacement, found := ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, conditionalOrder.OrderId)
	require.True(t, found)
	require.Equal(
		t,
		conditionalOrder,
		longTermOrderPlacement.Order,
	)

	// Get order Id from a specified store.
	getOrderFromStore := func(orderId types.OrderId, store prefix.Store) (
		orderPlacement types.LongTermOrderPlacement,
		found bool,
	) {
		orderKey := orderId.ToStateKey()
		bytes := store.Get(orderKey)
		if bytes == nil {
			return orderPlacement, false
		}
		ks.Cdc.MustUnmarshal(bytes, &orderPlacement)
		return orderPlacement, true
	}

	triggeredConditionalOrderMemStore := ks.ClobKeeper.GetTriggeredConditionalOrderPlacementMemStore(ks.Ctx)
	triggeredConditionalOrderStore := ks.ClobKeeper.GetTriggeredConditionalOrderPlacementStore(ks.Ctx)
	untriggeredConditionalOrderMemStore := ks.ClobKeeper.GetUntriggeredConditionalOrderPlacementMemStore(ks.Ctx)
	untriggeredConditionalOrderStore := ks.ClobKeeper.GetUntriggeredConditionalOrderPlacementStore(ks.Ctx)

	// Verify that the triggered conditional order does not exist in untriggered memstore/store
	_, found = getOrderFromStore(conditionalOrder.OrderId, untriggeredConditionalOrderStore)
	require.False(t, found)
	_, found = getOrderFromStore(conditionalOrder.OrderId, untriggeredConditionalOrderMemStore)
	require.False(t, found)

	// Verify that the triggered conditional order does exist in triggered memstore/store
	longTermOrderPlacement, found = getOrderFromStore(conditionalOrder.OrderId, triggeredConditionalOrderStore)
	require.True(t, found)
	require.Equal(
		t,
		conditionalOrder,
		longTermOrderPlacement.Order,
	)
	require.Equal(
		t,
		uint32(0),
		longTermOrderPlacement.PlacementIndex.BlockHeight,
	)
	_, found = getOrderFromStore(conditionalOrder.OrderId, triggeredConditionalOrderMemStore)
	require.True(t, found)
	require.Equal(
		t,
		conditionalOrder,
		longTermOrderPlacement.Order,
	)
	require.Equal(
		t,
		uint32(0),
		longTermOrderPlacement.PlacementIndex.BlockHeight,
	)
	require.Equal(t,
		ks.ClobKeeper.GetStatefulOrderCount(ks.Ctx, conditionalOrder.OrderId.SubaccountId),
		uint32(1),
	)

	traceDecoder.RequireKeyPrefixWrittenInSequence(
		t,
		[]string{
			// Write the order to untriggered state and memStore and increment the stateful order
			// count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.UntriggeredConditionalOrderKeyPrefix +
				orderToStringId(conditionalOrder),
			types.UntriggeredConditionalOrderKeyPrefix +
				orderToStringId(conditionalOrder),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(conditionalOrder),
			types.NextStatefulOrderBlockTransactionIndexKey,
			// Write to triggered state and memstore
			types.TriggeredConditionalOrderKeyPrefix +
				orderToStringId(conditionalOrder),
			types.TriggeredConditionalOrderKeyPrefix +
				orderToStringId(conditionalOrder),
			// Delete from state and memstore
			types.UntriggeredConditionalOrderKeyPrefix +
				orderToStringId(conditionalOrder),
			types.UntriggeredConditionalOrderKeyPrefix +
				orderToStringId(conditionalOrder),
		},
	)

	// Assert triggering a conditional order that is not in state panics.
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustTriggerConditionalOrder: conditional order Id does not exist in Untriggered state: %+v",
			constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20.OrderId,
		),
		func() {
			ks.ClobKeeper.MustTriggerConditionalOrder(
				ks.Ctx,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20.OrderId,
			)
		},
	)
}

func TestGetSetDeleteLongTermOrderState(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	orders := []types.Order{
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
		constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
		constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10,
	}
	blockHeights := []uint32{
		2,
		3,
		1,
	}
	nextExpectedTransactionIndex := uint32(0)

	// Set the tracer on the multistore to verify the performed writes are correct.
	traceDecoder := &tracer.TraceDecoder{}
	ks.Ctx.MultiStore().SetTracer(traceDecoder)

	// Create each stateful order.
	for i, order := range orders {
		orderId := order.OrderId

		// Verify you cannot get an order that does not exist.
		_, found := ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, orderId)
		require.False(t, found)

		// Verify deleting a stateful order that does not exist succeeds.
		ks.ClobKeeper.DeleteLongTermOrderPlacement(ks.Ctx, orderId)

		// Verify you can create a stateful order that did not previously exist.
		ks.ClobKeeper.SetLongTermOrderPlacement(ks.Ctx, order, blockHeights[i])

		// Verify you can get each stateful order.
		foundOrderPlacement, found := ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, order.OrderId)
		require.True(t, found)
		require.Equal(
			t,
			types.LongTermOrderPlacement{
				Order: order,
				PlacementIndex: types.TransactionOrdering{
					BlockHeight:      blockHeights[i],
					TransactionIndex: nextExpectedTransactionIndex,
				},
			},
			foundOrderPlacement,
		)

		// Increment the next expected transaction index, since it's incremented for each new stateful
		// order placement.
		nextExpectedTransactionIndex += 1
	}

	// Delete each stateful order and verify it cannot be found.
	for _, order := range orders {
		ks.ClobKeeper.DeleteLongTermOrderPlacement(ks.Ctx, order.OrderId)

		_, found := ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, order.OrderId)
		require.False(t, found)
	}

	// Re-create each stateful order with a different block height and transaction index, and
	// verify it can be found.
	for i, order := range orders {
		ks.ClobKeeper.SetLongTermOrderPlacement(ks.Ctx, order, blockHeights[i]+1)

		foundOrderPlacement, found := ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, order.OrderId)
		require.True(t, found)
		require.Equal(
			t,
			types.LongTermOrderPlacement{
				Order: order,
				PlacementIndex: types.TransactionOrdering{
					BlockHeight:      blockHeights[i] + 1,
					TransactionIndex: nextExpectedTransactionIndex,
				},
			},
			foundOrderPlacement,
		)

		nextExpectedTransactionIndex += 1
	}

	traceDecoder.RequireKeyPrefixWrittenInSequence(
		t,
		[]string{
			// Delete the order from state and memStore and decrement the stateful order count.
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			// Write the order to state and memStore and increment the stateful order count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			// Delete the order from state and memStore and decrement the stateful order count.
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			// Write the order to state and memStore and increment the stateful order count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			// Delete the order from state and memStore and decrement the stateful order count.
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			// Write the order to state and memStore and increment the stateful order count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			// Delete the order from state and memStore and decrement the stateful order count.
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			// Delete the order from state and memStore and decrement the stateful order count.
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			// Delete the order from state and memStore and decrement the stateful order count.
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			// Write the order to state and memStore and increment the stateful order count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			// Write the order to state and memStore and increment the stateful order count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			// Write the order to state and memStore and increment the stateful order count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
		},
	)
}

func TestGetSetDeleteLongTermOrderState_Replacements(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	orders := []types.Order{
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
	}
	blockHeights := []uint32{
		7,
		9,
	}

	require.Equal(t, orders[0].OrderId, orders[1].OrderId)

	// Set the tracer on the multistore to verify the performed writes are correct.
	traceDecoder := &tracer.TraceDecoder{}
	ks.Ctx.MultiStore().SetTracer(traceDecoder)

	// Create both stateful orders.
	for i, order := range orders {
		ks.ClobKeeper.SetLongTermOrderPlacement(ks.Ctx, order, blockHeights[i])
	}
	// Since the order is a replacement we expect the stateful order count to be 1.
	require.Equal(
		t,
		ks.ClobKeeper.GetStatefulOrderCount(
			ks.Ctx,
			constants.Alice_Num0,
		),
		uint32(1),
	)

	// Verify the last created order exists.
	foundOrderPlacement, found := ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, orders[1].OrderId)
	require.True(t, found)
	require.Equal(
		t,
		types.LongTermOrderPlacement{
			Order: orders[1],
			PlacementIndex: types.TransactionOrdering{
				BlockHeight:      blockHeights[1],
				TransactionIndex: 1,
			},
		},
		foundOrderPlacement,
	)

	// Verify the order can be deleted.
	ks.ClobKeeper.DeleteLongTermOrderPlacement(ks.Ctx, orders[1].OrderId)
	_, found = ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, orders[1].OrderId)
	require.False(t, found)

	// Verify the multistore writes are correct.
	traceDecoder.RequireKeyPrefixWrittenInSequence(
		t,
		[]string{
			// Write the order to state and memStore and increment the stateful order count.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			// Write the order to state and memStore. We should not expect the stateful order
			// count to change since this is a replacement.
			types.NextStatefulOrderBlockTransactionIndexKey,
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			// Delete the order from state and memStore and decrement the stateful order count.
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			types.LongTermOrderPlacementKeyPrefix +
				orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			types.StatefulOrderCountPrefix +
				orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
		},
	)
}

func TestLongTermOrderState_ShortTermOrderPanics(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20
	errorString := fmt.Sprintf(
		"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
		shortTermOrder.OrderId,
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			ks.ClobKeeper.SetLongTermOrderPlacement(
				ks.Ctx,
				shortTermOrder,
				0,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			ks.ClobKeeper.GetLongTermOrderPlacement(
				ks.Ctx,
				shortTermOrder.OrderId,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			ks.ClobKeeper.DeleteLongTermOrderPlacement(
				ks.Ctx,
				shortTermOrder.OrderId,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
				ks.Ctx,
				constants.Time_21st_Feb_2021,
				shortTermOrder.OrderId,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			ks.ClobKeeper.MustRemoveStatefulOrder(
				ks.Ctx,
				shortTermOrder.OrderId,
			)
		},
	)
}

func TestGetAddAndRemoveStatefulOrderTimeSlice(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup func(ctx sdk.Context, k keeper.Keeper)

		// Expectations.
		expectedMultiStoreWrites []string
		expectedTimeSlices       map[time.Time][]types.OrderId
		expectedRemovedOrders    []types.OrderId
	}{
		"Can read an empty state": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {},

			expectedMultiStoreWrites: []string{},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
		},
		"Can read order IDs after they've been created": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.Time_21st_Feb_2021,
				)
			},

			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Place the first stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Add second order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set second stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Place the second stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Add third order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set third stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				// Place the third stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10.OrderId,
				},
			},
		},
		"Can read order IDs after they've been created and deleted and they're still sorted": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.MustGetUnixGoodTilBlockTime(),
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10.MustGetUnixGoodTilBlockTime(),
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
				)
				k.MustRemoveStatefulOrder(
					ctx,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10.OrderId,
				)
				k.MustRemoveStatefulOrder(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				)
			},

			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Place the first stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Add second order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:30.000000000",
				// Set second stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Place the second stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Add third order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				// Set third stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				// Place the third stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				// Remove first order from stateful order slice, which removes the fill amount, stateful
				// order placement from state and memStore, and decrement the stateful order count.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:30.000000000",
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Remove second order from stateful order slice, which removes the fill amount, stateful
				// order placement from state and memStore, and decrement the stateful order count.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.TimeFifteen: {
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
				},
			},
			expectedRemovedOrders: []types.OrderId{
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
		},
		"Can create order IDs in non-sorted order and they're sorted in state": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					constants.Time_21st_Feb_2021,
				)
			},

			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Place the first stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Add second order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.00000000",
				// Set second stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Place the second stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Add third order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				// Place the third stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				// Add fourth order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set fourth stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				// Place the fourth stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				// Add fifth order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set fifth stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				// Place the fifth stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				// Add sixth order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set sixth stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				// Place the sixth stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				// Add seventh order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				// Set seventh stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				// Place the seventh stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
				},
			},
		},
		"Can delete all order IDs that were created": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				orders := []types.Order{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				}
				for _, order := range orders {
					createPartiallyFilledStatefulOrderInState(
						ctx,
						k,
						order,
						order.MustGetUnixGoodTilBlockTime(),
					)
				}
				for _, order := range orders {
					k.MustRemoveStatefulOrder(
						ctx,
						order.OrderId,
					)
				}
			},
			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				// Place the first stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				// Add second order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:10.000000000",
				// Set second stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				// Place the second stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				// Add third order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				// Set third stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				// Place the third stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				// Remove first order from stateful order slice, which removes the fill amount, stateful
				// order placement from state and memStore, and decrement the stateful order count.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				// Remove second order from stateful order slice, which removes the fill amount, stateful
				// order placement from state and memStore, and decrement the stateful order count.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:10.000000000",
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				// Remove third order from stateful order slice, which removes the fill amount, stateful
				// order placement from state and memStore, and decrement the stateful order count.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.TimeTen:     {},
				constants.TimeFifteen: {},
			},
			expectedRemovedOrders: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
		},
		"Can add and remove order IDs from multiple time slices": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				orders := []types.Order{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				}
				for _, order := range orders {
					createPartiallyFilledStatefulOrderInState(
						ctx,
						k,
						order,
						order.MustGetUnixGoodTilBlockTime(),
					)
				}
				// Remove an order from two of the timestamps.
				k.MustRemoveStatefulOrder(
					ctx,
					orders[6].OrderId,
				)
				k.MustRemoveStatefulOrder(
					ctx,
					orders[2].OrderId,
				)
			},

			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				// Place the first stateful order in state and memStore and increment the stateful
				// order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20),
				// Add second order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:20.000000000",
				// Set second stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				// Place the second stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20),
				// Add third order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:25.000000000",
				// Set third stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				// Place the third stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				// Add fourth order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:15.000000000",
				// Set fourth stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Place the fourth stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15),
				// Add fifth order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:30.000000000",
				// Set fifth stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Place the fifth stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.UntriggeredConditionalOrderKeyPrefix +
					orderToStringId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10),
				// Add sixth order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:10.000000000",
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				// Place the sixth stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10),
				// Add seventh order to stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:10.000000000",
				// Set seventh stateful order fill amount to a non-zero value in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				// Place the seventh stateful order in state and memStore and increment the stateful order count.
				types.NextStatefulOrderBlockTransactionIndexKey,
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				// Remove seventh order from stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:10.000000000",
				// Remove seventh stateful order fill amount in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				// Remove the seventh stateful order placement from state and memStore and decrement the stateful
				// order count.
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
				// Remove third order from stateful order slice.
				types.StatefulOrdersTimeSlicePrefix + "1970-01-01T00:00:25.000000000",
				// Remove third stateful order fill amount in state.
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.OrderAmountFilledKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				// Remove the third stateful order placement from state and memStore and decrement the stateful
				// order count.
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.LongTermOrderPlacementKeyPrefix +
					orderToStringId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
				types.StatefulOrderCountPrefix +
					orderToStringSubaccountId(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.TimeFifteen: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
				},
				constants.TimeTwenty: {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				},
				constants.TimeThirty: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10.OrderId,
				},
				constants.TimeTen: {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ks.Ctx.MultiStore().SetTracer(traceDecoder)

			tc.setup(ks.Ctx, *ks.ClobKeeper)

			// Verify the writes were done in the expected order.
			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)

			// Verify the state is correct.
			for goodTilTime, expectedOrderIds := range tc.expectedTimeSlices {
				orderIds := ks.ClobKeeper.GetStatefulOrdersTimeSlice(ks.Ctx, goodTilTime)
				sort.Sort(types.SortedOrders(expectedOrderIds))
				require.Equal(
					t,
					expectedOrderIds,
					orderIds,
					"Mismatch of order IDs for timestamp",
					goodTilTime.String(),
				)
				for _, orderId := range orderIds {
					exists, _, _ := ks.ClobKeeper.GetOrderFillAmount(ks.Ctx, orderId)
					require.True(t, exists)
					_, exists = ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, orderId)
					require.True(t, exists)
				}
			}

			for _, orderId := range tc.expectedRemovedOrders {
				exists, _, _ := ks.ClobKeeper.GetOrderFillAmount(ks.Ctx, orderId)
				require.False(t, exists)
				_, exists = ks.ClobKeeper.GetLongTermOrderPlacement(ks.Ctx, orderId)
				require.False(t, exists)
			}
		})
	}
}

func TestRemoveExpiredStatefulOrdersTimeSlices(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		timeSlicesToOrderIds map[time.Time][]types.OrderId

		// Parameters.
		blockTime time.Time

		// Expectations.
		expectedMultiStoreWrites []string
		expectedTimeSlices       map[time.Time][]types.OrderId
		expectedExpiredOrderIds  []types.OrderId
	}{
		"Can delete an empty state": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},

			blockTime: constants.Time_21st_Feb_2021,

			expectedMultiStoreWrites: []string{},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
			expectedExpiredOrderIds: []types.OrderId{},
		},
		"Deletes all time slices before blockTime": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(1): {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(77): {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
			},

			blockTime: constants.Time_21st_Feb_2021.Add(1_000_000_000_000),

			expectedMultiStoreWrites: []string{
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000001",
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000077",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
			expectedExpiredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
			},
		},
		"Deletes all time slices before blockTime inclusive": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(1): {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(77): {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
			},

			blockTime: constants.Time_21st_Feb_2021.Add(77),

			expectedMultiStoreWrites: []string{
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000001",
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000077",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
			expectedExpiredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
			},
		},
		"Does not delete time slices after blockTime": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(1): {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(77): {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
			},

			blockTime: constants.Time_21st_Feb_2021.Add(76),

			expectedMultiStoreWrites: []string{
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000000",
				types.StatefulOrdersTimeSlicePrefix + "2021-02-21T00:00:00.000000001",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021.Add(77): {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
			},
			expectedExpiredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Create all order IDs in state.
			for timestamp, orderIds := range tc.timeSlicesToOrderIds {
				for _, orderId := range orderIds {
					ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(ks.Ctx, timestamp, orderId)
				}
			}

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ks.Ctx.MultiStore().SetTracer(traceDecoder)

			// Run the test.
			expiredOrderIds := ks.ClobKeeper.RemoveExpiredStatefulOrdersTimeSlices(ks.Ctx, tc.blockTime)

			// Verify the correct orders were expired.
			require.Equal(t, tc.expectedExpiredOrderIds, expiredOrderIds)

			// Verify the writes were done in the expected order.
			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)

			// Verify the state is correct.
			for goodTilTime, expectedOrderIds := range tc.expectedTimeSlices {
				orderIds := ks.ClobKeeper.GetStatefulOrdersTimeSlice(ks.Ctx, goodTilTime)
				require.Equal(
					t,
					expectedOrderIds,
					orderIds,
					"Mismatch of order IDs for timestamp",
					goodTilTime.String(),
				)
			}
		})
	}
}

func TestAddOrderToStatefulOrdersTimeSlice_PanicsIfAlreadyExists(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	goodTilTime := constants.Time_21st_Feb_2021
	createPartiallyFilledStatefulOrderInState(ks.Ctx, *ks.ClobKeeper, order, goodTilTime)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddOrderToStatefulOrdersTimeSlice: order ID %v is already contained in state for time %v",
			order.OrderId,
			goodTilTime,
		),
		func() {
			ks.ClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
				ks.Ctx,
				goodTilTime,
				order.OrderId,
			)
		},
	)
}

func TestRemoveLongTermOrder_PanicsIfNotFound(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	orderId := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustRemoveStatefulOrder: order %v does not exist",
			orderId,
		),
		func() {
			ks.ClobKeeper.MustRemoveStatefulOrder(
				ks.Ctx,
				orderId,
			)
		},
	)
}

// TODO(CLOB-786): Fix this test to verify sorting by transaction index works.
func TestGetAllStatefulOrders(t *testing.T) {
	tests := map[string]struct {
		// State.
		statefulOrderPlacements     []types.LongTermOrderPlacement
		isTriggeredConditionalOrder map[types.OrderId]bool

		// Expectations.
		expectedPlacedStatefulOrders         []types.Order
		expectedUntriggeredConditionalOrders []types.Order
	}{
		"Can read an empty state": {
			statefulOrderPlacements:     []types.LongTermOrderPlacement{},
			isTriggeredConditionalOrder: map[types.OrderId]bool{},

			expectedPlacedStatefulOrders:         []types.Order{},
			expectedUntriggeredConditionalOrders: []types.Order{},
		},
		`Can read stateful orders from state and untriggered conditional orders are returned separately
			from other orders`: {
			statefulOrderPlacements: []types.LongTermOrderPlacement{
				{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 4,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 8,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 4,
					},
				},
				{
					Order: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 8,
					},
				},
			},
			isTriggeredConditionalOrder: map[types.OrderId]bool{},

			expectedPlacedStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
			},
			expectedUntriggeredConditionalOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
			},
		},
		"Can read stateful orders from state with same block height sorted in ascending order": {
			statefulOrderPlacements: []types.LongTermOrderPlacement{
				{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 4,
					},
				},
				{
					Order: constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 4,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 8,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 4,
					},
				},
				{
					Order: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 8,
					},
				},
				{
					Order: constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 8,
					},
				},
			},
			isTriggeredConditionalOrder: map[types.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId: true,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId:            true,
			},

			expectedPlacedStatefulOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
			},
			expectedUntriggeredConditionalOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
		},
		"Can read stateful orders from state with same transaction index sorted in ascending order": {
			statefulOrderPlacements: []types.LongTermOrderPlacement{
				{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 3,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 2,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 7,
					},
				},
				{
					Order: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight: 8,
					},
				},
			},
			isTriggeredConditionalOrder: map[types.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId: true,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId:            true,
			},

			expectedPlacedStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
			},
			expectedUntriggeredConditionalOrders: []types.Order{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			for _, statefulOrderPlacement := range tc.statefulOrderPlacements {
				ks.ClobKeeper.SetLongTermOrderPlacement(
					ks.Ctx,
					statefulOrderPlacement.Order,
					statefulOrderPlacement.PlacementIndex.BlockHeight,
				)
				if tc.isTriggeredConditionalOrder[statefulOrderPlacement.Order.OrderId] {
					require.True(t, statefulOrderPlacement.Order.IsConditionalOrder())
					ctx := ks.Ctx.WithBlockHeight(int64(statefulOrderPlacement.PlacementIndex.BlockHeight))
					ks.ClobKeeper.MustTriggerConditionalOrder(ctx, statefulOrderPlacement.Order.OrderId)
				}
			}

			// Verify the stateful order placements are correct.
			placedStatefulOrders := ks.ClobKeeper.GetAllPlacedStatefulOrders(ks.Ctx)
			untriggeredConditionalOrders := ks.ClobKeeper.GetAllUntriggeredConditionalOrders(ks.Ctx)
			require.Equal(t, tc.expectedPlacedStatefulOrders, placedStatefulOrders)
			require.Equal(t, tc.expectedUntriggeredConditionalOrders, untriggeredConditionalOrders)
		})
	}
}
