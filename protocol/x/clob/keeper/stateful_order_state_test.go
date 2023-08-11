package keeper_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/proto"
	"github.com/dydxprotocol/v4/testutil/tracer"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
	k.SetStatefulOrderPlacement(ctx, order, uint32(30))
}

func TestStatefulOrderInitMemStore_Success(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper, _, _, _, _, _, _ := keepertest.ClobKeepersWithUninitializedMemStore(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	// Set some stateful orders.
	keeper.SetStatefulOrderPlacement(
		ctx,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
		0,
	)

	keeper.SetStatefulOrderPlacement(
		ctx,
		constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
		0,
	)

	// Init the memstore.
	keeper.InitMemStore(ctx)

	// Assert that the values can be read after memStore has been warmed.
	order, exists := keeper.GetStatefulOrderPlacement(
		ctx, constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId)
	require.True(t, exists)
	require.Equal(t, constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20, order.Order)

	order, exists = keeper.GetStatefulOrderPlacement(
		ctx, constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId)
	require.True(t, exists)
	require.Equal(t, constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25, order.Order)

	order, exists = keeper.GetStatefulOrderPlacement(
		ctx, constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25.OrderId)
	require.False(t, exists)
}

func TestGetSetDeleteStatefulOrderState(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx,
		clobKeeper,
		_,
		_,
		_,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

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
	ctx.MultiStore().SetTracer(traceDecoder)

	// Create each stateful order.
	for i, order := range orders {
		orderId := order.OrderId

		// Verify you cannot get an order that does not exist.
		_, found := clobKeeper.GetStatefulOrderPlacement(ctx, orderId)
		require.False(t, found)

		// Verify deleting a stateful order that does not exist succeeds.
		clobKeeper.DeleteStatefulOrderPlacement(ctx, orderId)

		// Verify you can create a stateful order that did not previously exist.
		clobKeeper.SetStatefulOrderPlacement(ctx, order, blockHeights[i])

		// Verify you can get each stateful order.
		foundOrderPlacement, found := clobKeeper.GetStatefulOrderPlacement(ctx, order.OrderId)
		require.True(t, found)
		require.Equal(
			t,
			types.StatefulOrderPlacement{
				Order:            order,
				BlockHeight:      blockHeights[i],
				TransactionIndex: nextExpectedTransactionIndex,
			},
			foundOrderPlacement,
		)

		// Increment the next expected transaction index, since it's incremented for each new stateful
		// order placement.
		nextExpectedTransactionIndex += 1
	}

	// Delete each stateful order and verify it cannot be found.
	for _, order := range orders {
		clobKeeper.DeleteStatefulOrderPlacement(ctx, order.OrderId)

		_, found := clobKeeper.GetStatefulOrderPlacement(ctx, order.OrderId)
		require.False(t, found)
	}

	// Re-create each stateful order with a different block height and transaction index, and
	// verify it can be found.
	for i, order := range orders {
		clobKeeper.SetStatefulOrderPlacement(ctx, order, blockHeights[i]+1)

		foundOrderPlacement, found := clobKeeper.GetStatefulOrderPlacement(ctx, order.OrderId)
		require.True(t, found)
		require.Equal(
			t,
			types.StatefulOrderPlacement{
				Order:            order,
				BlockHeight:      blockHeights[i] + 1,
				TransactionIndex: nextExpectedTransactionIndex,
			},
			foundOrderPlacement,
		)

		nextExpectedTransactionIndex += 1
	}

	traceDecoder.RequireKeyPrefixWrittenInSequence(
		t,
		[]string{
			// Delete the order from state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
			// Delete the order from state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
			// Delete the order from state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
			// Delete the order from state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			// Delete the order from state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			// Delete the order from state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal(),
				)),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal(),
				)),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId.Marshal(),
				)),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
		},
	)
}

func TestGetSetDeleteStatefulOrderState_Replacements(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx,
		clobKeeper,
		_,
		_,
		_,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

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
	ctx.MultiStore().SetTracer(traceDecoder)

	// Create both stateful orders.
	for i, order := range orders {
		clobKeeper.SetStatefulOrderPlacement(ctx, order, blockHeights[i])
	}

	// Verify the last created order exists.
	foundOrderPlacement, found := clobKeeper.GetStatefulOrderPlacement(ctx, orders[1].OrderId)
	require.True(t, found)
	require.Equal(
		t,
		types.StatefulOrderPlacement{
			Order:            orders[1],
			BlockHeight:      blockHeights[1],
			TransactionIndex: 1,
		},
		foundOrderPlacement,
	)

	// Verify the order can be deleted.
	clobKeeper.DeleteStatefulOrderPlacement(ctx, orders[1].OrderId)
	_, found = clobKeeper.GetStatefulOrderPlacement(ctx, orders[1].OrderId)
	require.False(t, found)

	// Verify the multistore writes are correct.
	traceDecoder.RequireKeyPrefixWrittenInSequence(
		t,
		[]string{
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
			),
			"NextStatefulOrderBlockTransactionIndex/value",
			// Write the order to state and memStore.
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
			),
			fmt.Sprintf(
				"StatefulOrderPlacement/value/%v",
				string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
			),
		},
	)
}

func TestStatefulOrderState_ShortTermOrderPanics(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx,
		clobKeeper,
		_,
		_,
		_,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20
	errorString := fmt.Sprintf(
		"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
		shortTermOrder.OrderId,
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			clobKeeper.SetStatefulOrderPlacement(
				ctx,
				shortTermOrder,
				0,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			clobKeeper.GetStatefulOrderPlacement(
				ctx,
				shortTermOrder.OrderId,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			clobKeeper.DeleteStatefulOrderPlacement(
				ctx,
				shortTermOrder.OrderId,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			clobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
				ctx,
				constants.Time_21st_Feb_2021,
				shortTermOrder.OrderId,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			clobKeeper.MustRemoveStatefulOrder(
				ctx,
				constants.Time_21st_Feb_2021,
				shortTermOrder.OrderId,
			)
		},
	)

	require.PanicsWithValue(
		t,
		errorString,
		func() {
			clobKeeper.DoesStatefulOrderExistInState(
				ctx,
				shortTermOrder,
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
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30,
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
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				// Place the first stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add second order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set second stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				// Place the second stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add third order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set third stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				// Place the third stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId,
				},
			},
		},
		"Can read order IDs after they've been created and deleted and they're still sorted": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.Time_21st_Feb_2021,
				)
				k.MustRemoveStatefulOrder(
					ctx,
					constants.Time_21st_Feb_2021,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId,
				)
				k.MustRemoveStatefulOrder(
					ctx,
					constants.Time_21st_Feb_2021,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				)
			},

			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				// Place the first stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add second order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set second stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				// Place the second stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add third order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set third stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				// Place the third stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Remove first order from stateful order slice, which removes the fill amount and stateful
				// order placement from state and memStore as well.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal(),
					)),
				),
				// Remove second order from stateful order slice, which removes the fill amount and stateful
				// order placement from state and memStore as well.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId,
				},
			},
			expectedRemovedOrders: []types.OrderId{
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
		},
		`Can create and delete an order ID that doesn't have a fill amount or stateful order
		placement in state`: {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					constants.Time_21st_Feb_2021,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId,
				)
				k.MustRemoveStatefulOrder(
					ctx,
					constants.Time_21st_Feb_2021,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId,
				)
			},

			expectedMultiStoreWrites: []string{
				// Add order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Remove order from stateful order slice, which removes the fill amount and stateful
				// order placement from state and memStore as well.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal(),
					)),
				),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
			expectedRemovedOrders: []types.OrderId{
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId,
			},
		},
		"Can create order IDs in non-sorted order and they're sorted in state": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30,
					constants.Time_21st_Feb_2021,
				)
				createPartiallyFilledStatefulOrderInState(
					ctx,
					k,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
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
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				// Place the first stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add second order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.00000000",
				// Set second stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				// Place the second stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add third order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				// Place the third stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add fourth order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set fourth stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Place the fourth stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add fifth order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set fifth stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Place the fifth stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add sixth order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set sixth stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				// Place the sixth stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add seventh order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set seventh stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				// Place the seventh stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
				},
			},
		},
		"Can delete all order IDs that were created": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				orders := []types.Order{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				}
				for _, order := range orders {
					createPartiallyFilledStatefulOrderInState(
						ctx,
						k,
						order,
						constants.Time_21st_Feb_2021,
					)
				}
				for _, order := range orders {
					k.MustRemoveStatefulOrder(
						ctx,
						constants.Time_21st_Feb_2021,
						order.OrderId,
					)
				}
			},
			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				// Place the first stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add second order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set second stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Place the second stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add third order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set third stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				// Place the third stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Remove first order from stateful order slice, which removes the fill amount and stateful
				// order placement from state and memStore as well.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				// Remove second order from stateful order slice, which removes the fill amount and stateful
				// order placement from state and memStore as well.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Remove third order from stateful order slice, which removes the fill amount and stateful
				// order placement from state and memStore as well.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
			expectedRemovedOrders: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
		},
		"Can add and remove order IDs from multiple time slices": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				timestamps := []time.Time{
					constants.Time_21st_Feb_2021,
					constants.Time_21st_Feb_2021.Add(1),
					constants.Time_21st_Feb_2021.Add(77),
				}
				orders := []types.Order{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				}

				for i, order := range orders {
					createPartiallyFilledStatefulOrderInState(
						ctx,
						k,
						order,
						timestamps[i%3],
					)
				}

				// Remove an order from two of the timestamps.
				k.MustRemoveStatefulOrder(
					ctx,
					timestamps[0],
					orders[6].OrderId,
				)
				k.MustRemoveStatefulOrder(
					ctx,
					timestamps[2],
					orders[2].OrderId,
				)
			},

			expectedMultiStoreWrites: []string{
				// Add first order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set first stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				// Place the first stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add second order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000001",
				// Set second stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				// Place the second stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add third order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000077",
				// Set third stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				// Place the third stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add fourth order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set fourth stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				// Place the fourth stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add fifth order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000001",
				// Set fifth stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				// Place the fifth stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add sixth order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000077",
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Place the sixth stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Add seventh order to stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Set seventh stateful order fill amount to a non-zero value in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Place the seventh stateful order in state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				// Remove seventh order from stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				// Remove seventh stateful order fill amount in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Remove the seventh stateful order placement from state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId.Marshal())),
				),
				// Remove third order from stateful order slice.
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000077",
				// Remove third stateful order fill amount in state.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				// Remove the third stateful order placement from state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId.Marshal())),
				),
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(1): {
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30.OrderId,
				},
				constants.Time_21st_Feb_2021.Add(77): {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ctx,
				clobKeeper,
				_,
				_,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			tc.setup(ctx, *clobKeeper)

			// Verify the writes were done in the expected order.
			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)

			// Verify the state is correct.
			for goodTilTime, expectedOrderIds := range tc.expectedTimeSlices {
				orderIds := clobKeeper.GetStatefulOrdersTimeSlice(ctx, goodTilTime)
				sort.Sort(types.SortedOrders(expectedOrderIds))
				require.Equal(
					t,
					expectedOrderIds,
					orderIds,
					"Mismatch of order IDs for timestamp",
					goodTilTime.String(),
				)
				for _, orderId := range orderIds {
					exists, _, _ := clobKeeper.GetOrderFillAmount(ctx, orderId)
					require.True(t, exists)
					_, exists = clobKeeper.GetStatefulOrderPlacement(ctx, orderId)
					require.True(t, exists)
				}
			}

			for _, orderId := range tc.expectedRemovedOrders {
				exists, _, _ := clobKeeper.GetOrderFillAmount(ctx, orderId)
				require.False(t, exists)
				_, exists = clobKeeper.GetStatefulOrderPlacement(ctx, orderId)
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
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
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
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000001",
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000077",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
			expectedExpiredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
			},
		},
		"Deletes all time slices before blockTime inclusive": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
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
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000001",
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000077",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {},
			},
			expectedExpiredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30.OrderId,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
			},
		},
		"Does not delete time slices after blockTime": {
			timeSlicesToOrderIds: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
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
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000000",
				"StatefulOrdersTimeSlice/value/2021-02-21T00:00:00.000000001",
			},
			expectedTimeSlices: map[time.Time][]types.OrderId{
				constants.Time_21st_Feb_2021.Add(77): {
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
			},
			expectedExpiredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
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
			ctx,
				clobKeeper,
				_,
				_,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Create all order IDs in state.
			for timestamp, orderIds := range tc.timeSlicesToOrderIds {
				for _, orderId := range orderIds {
					clobKeeper.MustAddOrderToStatefulOrdersTimeSlice(ctx, timestamp, orderId)
				}
			}

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			// Run the test.
			expiredOrderIds := clobKeeper.RemoveExpiredStatefulOrdersTimeSlices(ctx, tc.blockTime)

			// Verify the correct orders were expired.
			require.Equal(t, tc.expectedExpiredOrderIds, expiredOrderIds)

			// Verify the writes were done in the expected order.
			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)

			// Verify the state is correct.
			for goodTilTime, expectedOrderIds := range tc.expectedTimeSlices {
				orderIds := clobKeeper.GetStatefulOrdersTimeSlice(ctx, goodTilTime)
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
	ctx,
		clobKeeper,
		_,
		_,
		_,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	goodTilTime := constants.Time_21st_Feb_2021
	createPartiallyFilledStatefulOrderInState(ctx, *clobKeeper, order, goodTilTime)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddOrderToStatefulOrdersTimeSlice: order ID %v is already contained in state for time %v",
			order.OrderId,
			goodTilTime,
		),
		func() {
			clobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
				ctx,
				goodTilTime,
				order.OrderId,
			)
		},
	)
}

func TestRemoveStatefulOrder_PanicsIfNotFound(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx,
		clobKeeper,
		_,
		_,
		_,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	orderId := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId
	goodTilTime := constants.Time_21st_Feb_2021
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustRemoveStatefulOrder: order ID %v is not in state for time %v",
			orderId,
			goodTilTime,
		),
		func() {
			clobKeeper.MustRemoveStatefulOrder(
				ctx,
				goodTilTime,
				orderId,
			)
		},
	)
}

func TestGetSetBlockTimeForLastCommittedBlock(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx,
		clobKeeper,
		_,
		_,
		_,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	// Set the tracer on the multistore to verify the performed writes are correct.
	traceDecoder := &tracer.TraceDecoder{}
	ctx.MultiStore().SetTracer(traceDecoder)

	ctx = ctx.WithBlockTime(constants.TimeT)
	clobKeeper.SetBlockTimeForLastCommittedBlock(ctx)
	require.True(
		t,
		constants.TimeT.Equal(
			clobKeeper.MustGetBlockTimeForLastCommittedBlock(ctx),
		),
	)

	// Test overwrite.
	ctx = ctx.WithBlockTime(constants.TimeTPlus1)
	clobKeeper.SetBlockTimeForLastCommittedBlock(ctx)
	require.True(
		t,
		constants.TimeTPlus1.Equal(
			clobKeeper.MustGetBlockTimeForLastCommittedBlock(ctx),
		),
	)

	traceDecoder.RequireReadWriteInSequence(
		t,
		[]tracer.TraceOperation{
			{
				Operation: tracer.WriteOperation,
				Key:       types.LastCommittedBlockTimeKey,
				Value:     sdk.FormatTimeString(constants.TimeT),
			},
			{
				Operation: tracer.ReadOperation,
				Key:       types.LastCommittedBlockTimeKey,
				Value:     sdk.FormatTimeString(constants.TimeT),
			},
			{
				Operation: tracer.WriteOperation,
				Key:       types.LastCommittedBlockTimeKey,
				Value:     sdk.FormatTimeString(constants.TimeTPlus1),
			},
			{
				Operation: tracer.ReadOperation,
				Key:       types.LastCommittedBlockTimeKey,
				Value:     sdk.FormatTimeString(constants.TimeTPlus1),
			},
		},
	)
}

func TestMustGetBlockTimeForLastCommittedBlock_Panics(t *testing.T) {
	// Setup keeper state and test parameters.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx,
		clobKeeper,
		_,
		_,
		_,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	require.PanicsWithValue(
		t,
		"Failed to get the block time of the previously committed block",
		func() {
			clobKeeper.MustGetBlockTimeForLastCommittedBlock(ctx)
		},
	)

	ctx = ctx.WithBlockTime(constants.TimeZero)
	require.PanicsWithValue(
		t,
		"Block-time is zero",
		func() {
			clobKeeper.SetBlockTimeForLastCommittedBlock(ctx)
		},
	)
}

func TestGetAllStatefulOrders(t *testing.T) {
	tests := map[string]struct {
		// State.
		statefulOrderPlacements []types.StatefulOrderPlacement

		// Setup.
		setup func(ctx sdk.Context, k keeper.Keeper, statefulOrderPlacements []types.StatefulOrderPlacement)

		// Expectations.
		expectedStatefulOrders []types.Order
	}{
		"Can read an empty state": {
			setup: func(ctx sdk.Context, k keeper.Keeper, statefulOrderPlacements []types.StatefulOrderPlacement) {
			},

			expectedStatefulOrders: []types.Order{},
		},
		"Can read stateful orders from state with same block height sorted in ascending order": {
			statefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:       constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					BlockHeight: 4,
				},
				{
					Order:       constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight: 8,
				},
				{
					Order:       constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight: 4,
				},
				{
					Order:       constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight: 8,
				},
			},

			setup: func(ctx sdk.Context, k keeper.Keeper, statefulOrderPlacements []types.StatefulOrderPlacement) {
				for _, statefulOrderPlacement := range statefulOrderPlacements {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrderPlacement.Order,
						statefulOrderPlacement.BlockHeight,
					)
				}
			},

			expectedStatefulOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
			},
		},
		"Can read stateful orders from state with same transaction index sorted in ascending order": {
			statefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:       constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					BlockHeight: 3,
				},
				{
					Order:       constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight: 2,
				},
				{
					Order:       constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight: 7,
				},
				{
					Order:       constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight: 8,
				},
			},

			setup: func(ctx sdk.Context, k keeper.Keeper, statefulOrderPlacements []types.StatefulOrderPlacement) {
				for _, statefulOrderPlacement := range statefulOrderPlacements {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrderPlacement.Order,
						statefulOrderPlacement.BlockHeight,
					)
				}
			},

			expectedStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ctx,
				clobKeeper,
				_,
				_,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			tc.setup(ctx, *clobKeeper, tc.statefulOrderPlacements)

			// Verify the stateful order placements are correct.
			statefulOrders := clobKeeper.GetAllStatefulOrders(ctx)
			require.Equal(t, tc.expectedStatefulOrders, statefulOrders)
		})
	}
}

func TestDoesStatefulOrderExistInState(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup func(ctx sdk.Context, k keeper.Keeper)

		// Parameters.
		order types.Order

		// Expectations.
		expectedExists bool
	}{
		"Returns false if no stateful orders were created": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
			},

			order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,

			expectedExists: false,
		},
		"Returns false if other stateful orders were created": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				for _, statefulOrder := range []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				} {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrder,
						0,
					)
				}
			},

			order: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10,

			expectedExists: false,
		},
		"Returns true if stateful order was created": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				for _, statefulOrder := range []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				} {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrder,
						0,
					)
				}
			},

			order: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,

			expectedExists: true,
		},
		"Returns false if stateful order was replaced": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				for _, statefulOrder := range []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				} {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrder,
						0,
					)
				}
			},

			order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,

			expectedExists: false,
		},
		"Returns false if stateful order was a replacement order": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				for _, statefulOrder := range []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				} {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrder,
						0,
					)
				}
			},

			order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,

			expectedExists: true,
		},
		"Returns false if stateful order was deleted": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				for _, statefulOrder := range []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				} {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrder,
						0,
					)
				}

				k.DeleteStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				)
			},

			order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,

			expectedExists: false,
		},
		"Returns true if stateful order was replaced, deleted, then re-created": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				for _, statefulOrder := range []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				} {
					k.SetStatefulOrderPlacement(
						ctx,
						statefulOrder,
						0,
					)
				}

				k.DeleteStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				)

				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					0,
				)
			},

			order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,

			expectedExists: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ctx,
				clobKeeper,
				_,
				_,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			tc.setup(ctx, *clobKeeper)

			// Run the test and verify expectations.
			exists := clobKeeper.DoesStatefulOrderExistInState(ctx, tc.order)
			require.Equal(t, tc.expectedExists, exists)
		})
	}
}
