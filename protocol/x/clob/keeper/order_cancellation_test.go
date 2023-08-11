package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/indexer/msgsender"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	clobtest "github.com/dydxprotocol/v4/testutil/clob"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/dydxprotocol/v4/x/perpetuals"
	"github.com/dydxprotocol/v4/x/prices"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShortTermCancelOrder_Success(t *testing.T) {
	memClob := &mocks.MemClob{}
	indexerMessageSender := &mocks.IndexerEventManager{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper, _, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexerMessageSender)
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx = ctx.WithBlockHeight(14)
	ctx = ctx.WithIsCheckTx(true)
	nextBlock := uint32(15)
	offchainUpdates := types.NewOffchainUpdates()

	// With a `GoodTilBlock` 0 blocks in the future.
	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderShortTerm(
		order.OrderId,
		uint32(nextBlock),
	)).Return(offchainUpdates, nil)
	err := keeper.CheckTxCancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock))
	require.NoError(t, err)
	indexerMessageSender.AssertExpectations(t)

	// With a `GoodTilBlock` exactly ShortBlockWindow blocks in the future
	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderShortTerm(
		order.OrderId,
		nextBlock+types.ShortBlockWindow,
	)).Return(offchainUpdates, nil)
	err = keeper.CheckTxCancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock+types.ShortBlockWindow))
	require.NoError(t, err)
	indexerMessageSender.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestShortTermCancelOrder_SuccessfullySendsOffchainData(t *testing.T) {
	memClob := &mocks.MemClob{}
	indexerMessageSender := &mocks.IndexerEventManager{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper, _, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexerMessageSender)
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx = ctx.WithBlockHeight(14).WithTxBytes(constants.TestTxBytes)
	ctx = ctx.WithIsCheckTx(true)
	nextBlock := uint32(15)
	offchainUpdates := types.NewOffchainUpdates()
	message := msgsender.Message{
		Key:   []byte("key"),
		Value: []byte("value"),
	}
	offchainUpdates.AddRemoveMessage(order.OrderId, message)

	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderShortTerm(
		order.OrderId,
		nextBlock,
	)).Return(offchainUpdates, nil)
	indexerMessageSender.On(
		"SendOffchainData",
		message.AddHeader(constants.TestTxHashHeader),
	).Return().Once()
	err := keeper.CheckTxCancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock))
	require.NoError(t, err)
	indexerMessageSender.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestShortTermCancelOrder_ErrGoodTilBlockExceedsHeight(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper,
		_, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx = ctx.WithIsCheckTx(true)

	// With a `GoodTilBlock` one less than the next height.
	ctx = ctx.WithBlockHeight(15)
	err := keeper.CheckTxCancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 15))
	require.ErrorIs(t, err, types.ErrHeightExceedsGoodTilBlock)

	// With a `GoodTilBlock` two less than the next height.
	ctx = ctx.WithBlockHeight(16)
	err = keeper.CheckTxCancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 15))
	require.ErrorIs(t, err, types.ErrHeightExceedsGoodTilBlock)

	memClob.AssertExpectations(t)
}

func TestShortTermCancelOrder_ErrGoodTilBlockExceedsShortBlockWindow(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper,
		_, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	ctx = ctx.WithIsCheckTx(true)

	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx = ctx.WithBlockHeight(14)
	nextBlock := uint32(15)

	// With a `GoodTilBlock` one more than ShortBlockWindow blocks in the future.
	err := keeper.CheckTxCancelOrder(
		ctx,
		types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock+types.ShortBlockWindow+1),
	)
	require.ErrorIs(t, err, types.ErrGoodTilBlockExceedsShortBlockWindow)

	memClob.AssertExpectations(t)
}

func TestCancelOrder_KeeperForwardsErrorsFromMemclob(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper,
		_, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	ctx = ctx.WithIsCheckTx(true)
	previousBlockTime := 50
	ctx = ctx.WithBlockTime(time.Unix(int64(previousBlockTime), 0))
	keeper.SetBlockTimeForLastCommittedBlock(ctx)

	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx = ctx.WithBlockHeight(14)

	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderShortTerm(
		shortTermOrder.OrderId,
		uint32(15),
	)).Return(nil, types.ErrMemClobCancelAlreadyExists)
	err := keeper.CheckTxCancelOrder(ctx, types.NewMsgCancelOrderShortTerm(shortTermOrder.OrderId, 15))
	require.ErrorIs(t, err, types.ErrMemClobCancelAlreadyExists)

	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	keeper.SetStatefulOrderPlacement(ctx, longTermOrder, 15)
	keeper.MustAddOrderToStatefulOrdersTimeSlice(
		ctx,
		longTermOrder.MustGetUnixGoodTilBlockTime(),
		longTermOrder.GetOrderId(),
	)
	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderStateful(
		longTermOrder.OrderId,
		uint32(100),
	)).Return(nil, types.ErrMemClobCancelAlreadyExists)
	err = keeper.CheckTxCancelOrder(ctx, types.NewMsgCancelOrderStateful(longTermOrder.OrderId, 100))
	require.ErrorIs(t, err, types.ErrMemClobCancelAlreadyExists)

	memClob.AssertExpectations(t)
}

func TestStatefulCancelOrder_Success(t *testing.T) {
	memClob := &mocks.MemClob{}
	indexerMessageSender := &mocks.IndexerEventManager{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper, _, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexerMessageSender)
	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	previousBlockTime := 15
	ctx = ctx.WithBlockHeight(14)
	ctx = ctx.WithIsCheckTx(true)
	ctx = ctx.WithBlockTime(time.Unix(int64(previousBlockTime), 0))
	keeper.SetBlockTimeForLastCommittedBlock(ctx)
	offchainUpdates := types.NewOffchainUpdates()

	// Cancel with a `GoodTilBlockTime` one more than order's GTBT in future
	statefulOrderCancel := types.NewMsgCancelOrderStateful(
		order.OrderId,
		uint32(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime()+1),
	)
	memClob.On("CancelOrder", ctx, statefulOrderCancel).Return(offchainUpdates, nil)
	// Simulate stateful order placement to state.
	keeper.SetStatefulOrderPlacement(ctx, order, 10)
	keeper.MustAddOrderToStatefulOrdersTimeSlice(
		ctx,
		order.MustGetUnixGoodTilBlockTime(),
		order.GetOrderId(),
	)
	err := keeper.CheckTxCancelOrder(ctx, statefulOrderCancel)
	require.NoError(t, err)
	indexerMessageSender.AssertExpectations(t)

	// Cancel wth a `GoodTilBlockTime` exactly StatefulOrderTimeWindow blocks in the future
	statefulCancellationGoodTilBlockTime := uint32(
		time.Unix(int64(previousBlockTime), 0).Add(types.StatefulOrderTimeWindow).Unix(),
	)
	// Simulate stateful order placement to state.
	keeper.SetStatefulOrderPlacement(ctx, order, 10)
	keeper.MustAddOrderToStatefulOrdersTimeSlice(
		ctx,
		order.MustGetUnixGoodTilBlockTime(),
		order.GetOrderId(),
	)
	// Simulate some order fills being added to state.
	keeper.SetOrderFillAmount(ctx, order.GetOrderId(), 1, 50)
	// Cancel the beforementioned stateful order.
	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderStateful(
		order.OrderId,
		uint32(time.Unix(int64(previousBlockTime), 0).Add(types.StatefulOrderTimeWindow).Unix()),
	)).Return(offchainUpdates, nil)
	err = keeper.CheckTxCancelOrder(
		ctx,
		types.NewMsgCancelOrderStateful(
			order.OrderId,
			statefulCancellationGoodTilBlockTime,
		),
	)

	require.NoError(t, err)
	// Verify that the stateful order placement was removed.
	_, found := keeper.GetStatefulOrderPlacement(
		ctx,
		order.OrderId,
	)
	require.False(t, found)
	// Verify that the stateful order placement was removed from time slice.
	orderIds := keeper.GetStatefulOrdersTimeSlice(
		ctx,
		time.Unix(int64(statefulCancellationGoodTilBlockTime), 0),
	)
	require.NotContains(t, orderIds, order.OrderId)
	// Verify that order fills were removed from fills.
	exists, _, _ := keeper.GetOrderFillAmount(
		ctx,
		order.OrderId,
	)
	require.False(t, exists)

	indexerMessageSender.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestStatefulCancelOrder_ErrStatefulOrderDoesNotExist(t *testing.T) {
	memClob := &mocks.MemClob{}
	indexerEventManager := &mocks.IndexerEventManager{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx, keeper, _, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexerEventManager)
	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	previousBlockTime := 15
	ctx = ctx.WithBlockHeight(14)
	ctx = ctx.WithIsCheckTx(true)
	ctx = ctx.WithBlockTime(time.Unix(int64(previousBlockTime), 0))
	keeper.SetBlockTimeForLastCommittedBlock(ctx)

	statefulOrderCancel := types.NewMsgCancelOrderStateful(
		order.OrderId,
		uint32(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime()+1),
	)
	err := keeper.CheckTxCancelOrder(ctx, statefulOrderCancel)
	require.ErrorContains(t, err, types.ErrStatefulOrderDoesNotExist.Error())
	indexerEventManager.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestPerformOrderCancellationStatefulValidation(t *testing.T) {
	blockHeight := uint32(5)
	blockTime := time.Unix(10, 0)

	tests := map[string]struct {
		setupState     func(ctx sdk.Context, k *keeper.Keeper)
		msgCancelOrder *types.MsgCancelOrder
		expectedErr    string
	}{
		"short-term cancellation succeeds with a GoodTilBlock of blockHeight": {
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_ShortTerm,
					ClobPairId:   uint32(0),
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: blockHeight},
			},
		},
		"short-term cancellation fails with a GoodTilBlock in the past": {
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_ShortTerm,
					ClobPairId:   uint32(0),
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 0},
			},
			expectedErr: types.ErrHeightExceedsGoodTilBlock.Error(),
		},
		"short-term cancellation fails with a GoodTilBlock exceeds ShortBlockWindow from block height": {
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_ShortTerm,
					ClobPairId:   uint32(0),
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: blockHeight + types.ShortBlockWindow + 1},
			},
			expectedErr: types.ErrGoodTilBlockExceedsShortBlockWindow.Error(),
		},
		"stateful cancellation fails when stateful cancel GoodTilBlockTime less than previous block blockTime": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
					ClobPairId:   uint32(0),
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: lib.MustConvertIntegerToUint32(blockTime.Unix() - 1),
				},
			},
			expectedErr: types.ErrTimeExceedsGoodTilBlockTime.Error(),
		},
		"stateful cancellation fails when stateful cancel GoodTilBlockTime equal to previous block blockTime": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
					ClobPairId:   uint32(0),
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: lib.MustConvertIntegerToUint32(blockTime.Unix()),
				},
			},
			expectedErr: types.ErrTimeExceedsGoodTilBlockTime.Error(),
		},
		`stateful cancellation fails when GoodTilBlockTime is over StatefulOrderTimeWindow greater
		than previous block blockTime`: {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
					ClobPairId:   uint32(0),
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					// One more than StatefulOrderTimeWindow in the future.
					GoodTilBlockTime: uint32(blockTime.Add(types.StatefulOrderTimeWindow).Unix() + 1),
				},
			},
			expectedErr: types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow.Error(),
		},
		`stateful cancellation succeeds when GoodTilBlockTime equal to StatefulOrderTimeWindow greater
		than previous block blockTime`: {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
				k.SetStatefulOrderPlacement(ctx, constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15, 15)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: uint32(blockTime.Add(types.StatefulOrderTimeWindow).Unix()),
				},
			},
		},
		"stateful cancellation fails if order to cancel does not exist": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime() + 1,
				},
			},
			expectedErr: "Order Id to cancel does not exist.",
		},
		"stateful cancellation succeeds if existing order GTBT is less than cancel order GTBT.": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					// Cancel GTBT is greater than the placed order GTBT.
					GoodTilBlockTime: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime() + 1,
				},
			},
		},
		"stateful cancellation succeeds if existing order GTBT is equal to cancel order GTBT.": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime(),
				},
			},
		},
		"stateful cancellation fails if existing order GTBT is greater than cancel order GTBT.": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
				k.SetStatefulOrderPlacement(ctx, constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15, 15)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime() - 1,
				},
			},
			expectedErr: "cancellation goodTilBlockTime less than stateful order goodTilBlockTime.",
		},
		"stateful cancellation success if existing cancellation with lower GTBT": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: 20, // New cancel is greater than existing cancel GTBT
				},
			},
		},
		"stateful cancellation success if existing cancellation with higher GTBT": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: 20, // New cancel is lower than existing cancel GTBT
				},
			},
		},
		"stateful cancellation success if existing cancellation with equal GTBT": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				// Sets to be 10 unix time.
				k.SetBlockTimeForLastCommittedBlock(
					ctx,
				)
				k.SetStatefulOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: 15, // New cancel is equal than existing cancel GTBT
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything, mock.Anything)
			memClob.On("SetClobKeeper", mock.Anything).Return()

			ctx, keeper, pricesKeeper, _, perpetualsKeeper, _, _, _ := keepertest.ClobKeepers(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)
			ctx = ctx.WithBlockTime(blockTime).WithBlockHeight(int64(blockHeight))
			prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			clobPair := types.ClobPair{
				Metadata: &types.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &types.PerpetualClobMetadata{
						PerpetualId: 0,
					},
				},
				Status:               types.ClobPair_STATUS_ACTIVE,
				StepBaseQuantums:     12,
				SubticksPerTick:      39,
				MinOrderBaseQuantums: 204,
			}

			_, err := keeper.CreatePerpetualClobPair(
				ctx,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				clobPair.Status,
				clobPair.MakerFeePpm,
				clobPair.TakerFeePpm,
			)
			require.NoError(t, err)
			keeper.SetBlockTimeForLastCommittedBlock(ctx)

			if tc.setupState != nil {
				tc.setupState(ctx, keeper)
			}

			err = keeper.PerformOrderCancellationStatefulValidation(ctx, tc.msgCancelOrder, blockHeight)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
