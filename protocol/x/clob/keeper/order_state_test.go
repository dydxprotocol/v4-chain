package keeper_test

import (
	"testing"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tracer"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAllOrderFillStates(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup func(ctx sdk.Context, k keeper.Keeper)

		// Expectations.
		expectedFillStates []keeper.OrderIdFillState
	}{
		"Reads an empty state": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
			},
			expectedFillStates: []keeper.OrderIdFillState{},
		},
		"Reads a single OrderFillState": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					100,
					10,
				)
			},
			expectedFillStates: []keeper.OrderIdFillState{
				{
					OrderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					OrderFillState: types.OrderFillState{
						FillAmount:          100,
						PrunableBlockHeight: 10,
					},
				},
			},
		},
		"Reads multiple OrderFillStates": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					100,
					10,
				)

				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
					101,
					11,
				)
			},
			expectedFillStates: []keeper.OrderIdFillState{
				{
					OrderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					OrderFillState: types.OrderFillState{
						FillAmount:          100,
						PrunableBlockHeight: 10,
					},
				},
				{
					OrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
					OrderFillState: types.OrderFillState{
						FillAmount:          101,
						PrunableBlockHeight: 11,
					},
				},
			},
		},
		"Writes same OrderFillState multiple times and the last update is reflected": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					100,
					10,
				)

				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					102,
					12,
				)
			},
			expectedFillStates: []keeper.OrderIdFillState{
				{
					OrderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					OrderFillState: types.OrderFillState{
						FillAmount:          102,
						PrunableBlockHeight: 12,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)

			tc.setup(ks.Ctx, *ks.ClobKeeper)

			fillStates := ks.ClobKeeper.GetAllOrderFillStates(ks.Ctx)
			require.ElementsMatch(t, fillStates, tc.expectedFillStates)
		})
	}
}

func TestSetGetOrderFillAmount(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId)

		// Invocation.
		orderId types.OrderId

		// Expectations.
		expectedExists              bool
		expectedFillAmount          satypes.BaseQuantums
		expectedPrunableBlockHeight uint32
	}{
		"SetOrderFillAmount then GetOrderFillAmount": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					10,
				)
			},

			expectedExists:              true,
			expectedFillAmount:          100,
			expectedPrunableBlockHeight: 10,
		},
		"SetOrderFillAmount twice, GetOrderFillAmount returns the most up-to-date values": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					10,
				)
				k.SetOrderFillAmount(
					ctx,
					orderId,
					200,
					20,
				)
			},

			expectedExists:              true,
			expectedFillAmount:          200,
			expectedPrunableBlockHeight: 20,
		},
		"GetOrderFillAmount with non-existent OrderFillState": {
			orderId:        constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup:          func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {},
			expectedExists: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)

			tc.setup(ks.Ctx, *ks.ClobKeeper, tc.orderId)

			exists, fillAmount, prunableBlockHeight := ks.ClobKeeper.GetOrderFillAmount(ks.Ctx, tc.orderId)

			require.Equal(t, exists, tc.expectedExists)
			if tc.expectedExists {
				require.Equal(t, fillAmount, tc.expectedFillAmount)
				require.Equal(t, prunableBlockHeight, tc.expectedPrunableBlockHeight)
			}
		})
	}
}

func TestOrderFillAmountInitMemStore_Success(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContextWithUninitializedMemStore(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	// Set some fill amounts.
	ks.ClobKeeper.SetOrderFillAmount(
		ks.Ctx,
		constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
		satypes.BaseQuantums(100),
		uint32(0),
	)

	ks.ClobKeeper.SetOrderFillAmount(
		ks.Ctx,
		constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
		satypes.BaseQuantums(100),
		uint32(0),
	)

	// This fill amount overwrites the first previous fill amount.
	ks.ClobKeeper.SetOrderFillAmount(
		ks.Ctx,
		constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
		satypes.BaseQuantums(200),
		uint32(0),
	)

	// Init the memstore.
	ks.ClobKeeper.InitMemStore(ks.Ctx)

	// Assert that the values can be read after memStore has been warmed.
	exists, amount, _ := ks.ClobKeeper.GetOrderFillAmount(
		ks.Ctx, constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId)
	require.True(t, exists)
	require.Equal(t, satypes.BaseQuantums(100), amount)

	exists, amount, _ = ks.ClobKeeper.GetOrderFillAmount(
		ks.Ctx, constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId)
	require.True(t, exists)
	require.Equal(t, satypes.BaseQuantums(200), amount)

	exists, _, _ = ks.ClobKeeper.GetOrderFillAmount(
		ks.Ctx, constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId)
	require.False(t, exists)
}

func TestPruning(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup func(t *testing.T, ctx sdk.Context, k keeper.Keeper, orderId types.OrderId)

		// Invocation.
		orderId types.OrderId

		// Expectations.
		expectedExists                                    bool
		expectedFillAmount                                satypes.BaseQuantums
		expectedPrunableBlockHeight                       uint32
		expectedEmptyPotentiallyPrunableOrderBlockHeights []uint32
	}{
		"Setting a fill amount, prune block, followed by pruning results in non-existent OrderFillState": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(t *testing.T, ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				blockHeight := uint32(10)

				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					blockHeight,
				)

				k.AddOrdersForPruning(
					ctx,
					[]types.OrderId{orderId},
					blockHeight,
				)

				result := k.PruneOrdersForBlockHeight(
					ctx,
					blockHeight,
				)

				require.Contains(t, result, orderId)
			},

			expectedExists: false,
			expectedEmptyPotentiallyPrunableOrderBlockHeights: []uint32{10},
		},
		`Updating the prunableBlockHeight on the OrderFillState results in order not getting pruned at previous
			prunableBlockHeight`: {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(t *testing.T, ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				blockHeight := uint32(10)
				nextBlockHeight := uint32(11)

				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					blockHeight,
				)

				k.AddOrdersForPruning(
					ctx,
					[]types.OrderId{orderId},
					blockHeight,
				)

				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					nextBlockHeight,
				)

				result := k.PruneOrdersForBlockHeight(
					ctx,
					blockHeight,
				)

				require.Empty(t, result)
			},

			expectedExists:              true,
			expectedFillAmount:          100,
			expectedPrunableBlockHeight: 11,
			expectedEmptyPotentiallyPrunableOrderBlockHeights: []uint32{10},
		},
		`Prunes orders for a block height that never had orders`: {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(t *testing.T, ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				k.PruneOrdersForBlockHeight(
					ctx,
					10,
				)
			},

			expectedExists: false,
			expectedEmptyPotentiallyPrunableOrderBlockHeights: []uint32{10},
		},
		"Prunes orders for a block before the current block": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(t *testing.T, ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				blockHeight := uint32(10)
				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
					100,
					blockHeight,
				)

				k.AddOrdersForPruning(
					ctx,
					[]types.OrderId{
						orderId,
					},
					blockHeight-1,
				)

				result := k.PruneOrdersForBlockHeight(
					ctx,
					blockHeight-1,
				)

				require.Empty(t, result)
			},

			expectedExists: false,
			expectedEmptyPotentiallyPrunableOrderBlockHeights: []uint32{9},
		},
		`Updates existing orders with new orders to be pruned`: {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(t *testing.T, ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				blockHeight := uint32(10)
				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
					100,
					blockHeight,
				)

				k.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
					100,
					blockHeight,
				)

				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					blockHeight,
				)

				k.AddOrdersForPruning(
					ctx,
					[]types.OrderId{
						orderId,
					},
					blockHeight,
				)

				k.AddOrdersForPruning(
					ctx,
					[]types.OrderId{
						constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
						constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
						constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
						constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
					},
					blockHeight,
				)

				result := k.PruneOrdersForBlockHeight(
					ctx,
					blockHeight,
				)

				require.ElementsMatch(
					t,
					result,
					[]types.OrderId{
						orderId,
						constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
						constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
					})

				// Ensure other two orders also don't exist.
				exists, _, _ := k.GetOrderFillAmount(ctx, constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId)
				require.False(t, exists)
				exists, _, _ = k.GetOrderFillAmount(ctx, constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId)
				require.False(t, exists)
			},

			expectedExists: false,
			expectedEmptyPotentiallyPrunableOrderBlockHeights: []uint32{10},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)

			tc.setup(t, ks.Ctx, *ks.ClobKeeper, tc.orderId)

			exists, fillAmount, prunableBlockHeight := ks.ClobKeeper.GetOrderFillAmount(ks.Ctx, tc.orderId)

			require.Equal(t, exists, tc.expectedExists)
			if tc.expectedExists {
				require.Equal(t, fillAmount, tc.expectedFillAmount)
				require.Equal(t, prunableBlockHeight, tc.expectedPrunableBlockHeight)
			}

			// Verify all prune order keys were deleted for specified heights
			for _, blockHeight := range tc.expectedEmptyPotentiallyPrunableOrderBlockHeights {
				potentiallyPrunableOrdersStore := ks.ClobKeeper.GetPruneableOrdersStore(ks.Ctx, blockHeight)
				it := potentiallyPrunableOrdersStore.Iterator(nil, nil)
				defer it.Close()
				require.False(t, it.Valid())
			}
		})
	}
}

func TestMigratePruneableOrders(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()

	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	ordersA := []types.OrderId{
		constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
		constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.OrderId,
		constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
	}
	ordersB := []types.OrderId{
		constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
		constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,
		constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.OrderId,
	}

	ks.ClobKeeper.LegacyAddOrdersForPruning(
		ks.Ctx,
		ordersA,
		10,
	)
	ks.ClobKeeper.LegacyAddOrdersForPruning(
		ks.Ctx,
		ordersB,
		100,
	)

	ks.ClobKeeper.MigratePruneableOrders(ks.Ctx)

	getPostMigrationOrdersAtHeight := func(height uint32) []types.OrderId {
		postMigrationOrders := []types.OrderId{}
		store := ks.ClobKeeper.GetPruneableOrdersStore(ks.Ctx, height)
		it := store.Iterator(nil, nil)
		defer it.Close()
		for ; it.Valid(); it.Next() {
			var orderId types.OrderId
			err := orderId.Unmarshal(it.Value())
			require.NoError(t, err)
			postMigrationOrders = append(postMigrationOrders, orderId)
		}
		return postMigrationOrders
	}

	require.ElementsMatch(t, ordersA, getPostMigrationOrdersAtHeight(10))
	require.ElementsMatch(t, ordersB, getPostMigrationOrdersAtHeight(100))

	oldStore := prefix.NewStore(
		ks.Ctx.KVStore(ks.StoreKey),
		[]byte(types.LegacyBlockHeightToPotentiallyPrunableOrdersPrefix), // nolint:staticcheck
	)
	it := oldStore.Iterator(nil, nil)
	defer it.Close()
	require.False(t, it.Valid())
}

func TestRemoveOrderFillAmount(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId)

		// Invocation.
		orderId types.OrderId

		// Expectations.
		expectedExists           bool
		expectedFillAmount       satypes.BaseQuantums
		expectedMultiStoreWrites []string
	}{
		"SetOrderFillAmount then RemoveOrderFillAmount removes the fill amount": {
			orderId: constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId,
			setup: func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					0,
				)
				k.RemoveOrderFillAmount(
					ctx,
					orderId,
				)
			},

			expectedExists: false,
			expectedMultiStoreWrites: []string{
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId.ToStateKey()),
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId.ToStateKey()),
			},
		},
		"SetOrderFillAmount twice and then RemoveOrderFillAmount removes the fill amount": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					10,
				)
				k.SetOrderFillAmount(
					ctx,
					orderId,
					200,
					20,
				)
				k.RemoveOrderFillAmount(
					ctx,
					orderId,
				)
			},

			expectedExists: false,
			expectedMultiStoreWrites: []string{
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
			},
		},
		"RemoveOrderFillAmount with non-existent order": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				k.RemoveOrderFillAmount(ctx, constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId)
			},
			expectedExists: false,
			expectedMultiStoreWrites: []string{
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
			},
		},
		"SetOrderFillAmount, RemoveOrderFillAmount, SetOrderFillAmount re-creates the fill amount": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					0,
				)
				k.RemoveOrderFillAmount(
					ctx,
					orderId,
				)
				k.SetOrderFillAmount(
					ctx,
					orderId,
					50,
					0,
				)
			},

			expectedExists:     true,
			expectedFillAmount: 50,
			expectedMultiStoreWrites: []string{
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
			},
		},
		"RemoveOrderFillAmount does not delete fill amounts for other orders": {
			orderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			setup: func(ctx sdk.Context, k keeper.Keeper, orderId types.OrderId) {
				removedOrderId := constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId
				require.NotEqual(t, removedOrderId, orderId)
				k.SetOrderFillAmount(
					ctx,
					orderId,
					100,
					0,
				)
				k.SetOrderFillAmount(
					ctx,
					removedOrderId,
					10,
					0,
				)
				k.RemoveOrderFillAmount(
					ctx,
					removedOrderId,
				)
			},

			expectedExists:     true,
			expectedFillAmount: 100,
			expectedMultiStoreWrites: []string{
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.ToStateKey()),
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId.ToStateKey()),
				types.OrderAmountFilledKeyPrefix +
					string(constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId.ToStateKey()),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ks.Ctx.MultiStore().SetTracer(traceDecoder)

			tc.setup(ks.Ctx, *ks.ClobKeeper, tc.orderId)

			exists, fillAmount, prunableBlockHeight := ks.ClobKeeper.GetOrderFillAmount(ks.Ctx, tc.orderId)

			require.Equal(t, exists, tc.expectedExists)
			if tc.expectedExists {
				require.Equal(t, fillAmount, tc.expectedFillAmount)
				require.Equal(t, prunableBlockHeight, uint32(0))
			}

			traceDecoder.RequireKeyPrefixWrittenInSequence(
				t,
				tc.expectedMultiStoreWrites,
			)
		})
	}
}
