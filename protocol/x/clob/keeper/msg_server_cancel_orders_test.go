package keeper_test

import (
	"context"
	"testing"
	"time"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	keeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCancelOrder_PanicIfShortTermOrder(t *testing.T) {
	cancellation := constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20
	require.Panicsf(
		t,
		func() {
			msgServer := keeper.NewMsgServerImpl(nil)
			//nolint: errcheck
			msgServer.CancelOrder(context.Background(), &cancellation)
		},
		"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
		cancellation.OrderId,
	)
}

func TestCancelOrder_PanicIfValidationSucceedsButOrderNotFound(t *testing.T) {
	cancellation := constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20
	require.Panicsf(
		t,
		func() {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContext(
				t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			mockClobKeeper := &mocks.ClobKeeper{}
			mockClobKeeper.On("PerformOrderCancellationStatefulValidation", ks.Ctx, mock.Anything, mock.Anything).Return(nil)
			mockClobKeeper.On("GetLongTermOrderPlacement", ks.Ctx, mock.Anything, mock.Anything).Return(nil, false)
			msgServer := keeper.NewMsgServerImpl(mockClobKeeper)
			//nolint: errcheck
			msgServer.CancelOrder(ks.Ctx, &cancellation)
		},
		"CancelOrder: stateful cancelation passed validation, but order %v does not exist",
		cancellation.OrderId,
	)
}

func TestCancelOrder_InfoLogIfOrderNotFound(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)

	orderToCancel := constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15

	ctx := ks.Ctx.WithBlockHeight(2)
	ctx = ctx.WithIsCheckTx(false).WithIsReCheckTx(false)
	mockLogger := &mocks.Logger{}
	mockLogger.On("With",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(mockLogger)
	mockLogger.On("Info",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("*errors.wrappedError"),
	).Return()
	ctx = ctx.WithLogger(mockLogger)
	ctx = ctx.WithBlockTime(time.Unix(int64(2), 0))
	ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
		Height:    1,
		Timestamp: time.Unix(int64(2), 0),
	})

	// MsgCancelOrder will fail because the order is not in state, but it should not log an error because
	// because the order is present in ProcessProposerMatchesEvents.RemovedStatefulOrderIds. This will log
	// an info message instead.
	ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
		ctx,
		types.ProcessProposerMatchesEvents{BlockHeight: 2, RemovedStatefulOrderIds: []types.OrderId{orderToCancel.OrderId}},
	)

	_, err := msgServer.CancelOrder(ctx, &orderToCancel)
	require.ErrorIs(t, err, types.ErrStatefulOrderCancellationFailedForAlreadyRemovedOrder)
	mockLogger.AssertExpectations(t)
}

func TestCancelOrder_ErrorLogIfGTBTTooLow(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)

	orderToCancel := constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15

	ctx := ks.Ctx.WithBlockHeight(2)
	ctx = ctx.WithIsCheckTx(false).WithIsReCheckTx(false)
	mockLogger := &mocks.Logger{}
	mockLogger.On("With",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockLogger)
	mockLogger.On(
		"Error",
		"Error cancelling order",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return()
	ctx = ctx.WithLogger(mockLogger)
	ctx = ctx.WithBlockTime(time.Unix(int64(2), 0))
	ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
		Height:    1,
		Timestamp: time.Unix(int64(20), 0),
	})

	// MsgCancelOrder will fail because the GTBT of the cancellation message is lower than the current block time.
	// This should log an error.
	_, err := msgServer.CancelOrder(ctx, &orderToCancel)
	require.ErrorIs(t, err, types.ErrTimeExceedsGoodTilBlockTime)
	mockLogger.AssertExpectations(t)
}

func TestCancelOrder_Error(t *testing.T) {
	tests := map[string]struct {
		StatefulOrderCancellation types.MsgCancelOrder
		ExpectedError             error
	}{
		"Returns an error when validation fails": {
			StatefulOrderCancellation: constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
			ExpectedError:             types.ErrStatefulOrderDoesNotExist,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize mocks, context, msgServer.
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			indexerEventManager := &mocks.IndexerEventManager{}
			indexerEventManager.On(
				"AddTxnEvent",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return().Once()

			ks := keepertest.NewClobKeepersTestContext(
				t, memClob, &mocks.BankKeeper{}, indexerEventManager)
			msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)

			ctx := ks.Ctx.WithBlockHeight(2)
			ctx = ctx.WithBlockTime(time.Unix(int64(2), 0))
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    1,
				Timestamp: time.Unix(int64(2), 0),
			})

			// Run MsgHandler for cancellation.
			_, err := msgServer.CancelOrder(ctx, &tc.StatefulOrderCancellation)
			require.ErrorIs(t, err, tc.ExpectedError)
		})
	}
}

func TestCancelOrder_Success(t *testing.T) {
	tests := map[string]struct {
		StatefulOrderPlacement    types.Order
		StatefulOrderCancellation types.MsgCancelOrder
	}{
		"Succeeds when GTBT are not equal": {
			StatefulOrderPlacement:    constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.MustGetOrder(),
			StatefulOrderCancellation: constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
		},
		"Succeeds when GTBT are equal": {
			StatefulOrderPlacement:    constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.MustGetOrder(),
			StatefulOrderCancellation: constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize mocks, context, msgServer.
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			indexerEventManager := &mocks.IndexerEventManager{}

			ks := keepertest.NewClobKeepersTestContext(
				t, memClob, &mocks.BankKeeper{}, indexerEventManager)
			msgServer := keeper.NewMsgServerImpl(ks.ClobKeeper)

			ctx := ks.Ctx.WithBlockHeight(2)
			ctx = ctx.WithBlockTime(time.Unix(int64(2), 0))
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    1,
				Timestamp: time.Unix(int64(2), 0),
			})

			// Setup IndexerEventManager mock.
			indexerEventManager.On(
				"AddTxnEvent",
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexerevents.StatefulOrderEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewStatefulOrderRemovalEvent(
						tc.StatefulOrderPlacement.GetOrderId(),
						indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
					),
				),
			).Return().Once()

			// Add stateful order placement to state
			ks.ClobKeeper.SetLongTermOrderPlacement(ctx, tc.StatefulOrderPlacement, 1)
			ks.ClobKeeper.AddStatefulOrderIdExpiration(
				ctx,
				tc.StatefulOrderPlacement.MustGetUnixGoodTilBlockTime(),
				tc.StatefulOrderPlacement.GetOrderId(),
			)

			// Add BlockHeight to `ProcessProposerMatchesEvents`. This is normally done in `BeginBlock`.
			ks.ClobKeeper.MustSetProcessProposerMatchesEvents(
				ctx,
				types.ProcessProposerMatchesEvents{
					BlockHeight: lib.MustConvertIntegerToUint32(2),
				},
			)

			// Run MsgHandler for cancellation.
			_, err := msgServer.CancelOrder(ctx, &tc.StatefulOrderCancellation)
			require.NoError(t, err)

			// Ensure stateful order placement removed from state.
			_, found := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.StatefulOrderPlacement.GetOrderId())
			require.False(t, found)

			// Ensure cancellation exists in `ProcessProposerMatchesEvents`.
			cancellations := ks.ClobKeeper.GetDeliveredCancelledOrderIds(ctx)
			require.Len(t, cancellations, 1)
			require.Equal(t, cancellations[0], tc.StatefulOrderCancellation.OrderId)

			// Run mock assertions.
			indexerEventManager.AssertExpectations(t)
		})
	}
}
