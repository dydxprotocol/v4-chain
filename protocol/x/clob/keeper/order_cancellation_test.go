package keeper_test

import (
	"testing"
	"time"

	cmt "github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShortTermCancelOrder_Success(t *testing.T) {
	memClob := &mocks.MemClob{}
	indexerMessageSender := &mocks.IndexerEventManager{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerMessageSender)
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx := ks.Ctx.WithBlockHeight(14)
	ctx = ctx.WithIsCheckTx(true)
	nextBlock := uint32(15)
	offchainUpdates := types.NewOffchainUpdates()

	// With a `GoodTilBlock` 0 blocks in the future.
	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderShortTerm(
		order.OrderId,
		uint32(nextBlock),
	)).Return(offchainUpdates, nil)
	err := ks.ClobKeeper.CancelShortTermOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock))
	require.NoError(t, err)
	indexerMessageSender.AssertExpectations(t)

	// With a `GoodTilBlock` exactly ShortBlockWindow blocks in the future
	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderShortTerm(
		order.OrderId,
		nextBlock+types.ShortBlockWindow,
	)).Return(offchainUpdates, nil)
	err = ks.ClobKeeper.CancelShortTermOrder(
		ctx,
		types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock+types.ShortBlockWindow),
	)
	require.NoError(t, err)
	indexerMessageSender.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestShortTermCancelOrder_SuccessfullySendsOffchainData(t *testing.T) {
	memClob := &mocks.MemClob{}
	indexerMessageSender := &mocks.IndexerEventManager{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerMessageSender)
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx := ks.Ctx.WithBlockHeight(14).WithTxBytes(constants.TestTxBytes)
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
	err := ks.ClobKeeper.CancelShortTermOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock))
	require.NoError(t, err)
	indexerMessageSender.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestShortTermCancelOrder_ErrGoodTilBlockExceedsHeight(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx := ks.Ctx.WithIsCheckTx(true)

	// With a `GoodTilBlock` one less than the next height.
	ctx = ctx.WithBlockHeight(15)
	err := ks.ClobKeeper.CancelShortTermOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 15))
	require.ErrorIs(t, err, types.ErrHeightExceedsGoodTilBlock)

	// With a `GoodTilBlock` two less than the next height.
	ctx = ctx.WithBlockHeight(16)
	err = ks.ClobKeeper.CancelShortTermOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 15))
	require.ErrorIs(t, err, types.ErrHeightExceedsGoodTilBlock)

	memClob.AssertExpectations(t)
}

func TestShortTermCancelOrder_ErrGoodTilBlockExceedsShortBlockWindow(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	ctx := ks.Ctx.WithIsCheckTx(true)

	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx = ctx.WithBlockHeight(14)
	nextBlock := uint32(15)

	// With a `GoodTilBlock` one more than ShortBlockWindow blocks in the future.
	err := ks.ClobKeeper.CancelShortTermOrder(
		ctx,
		types.NewMsgCancelOrderShortTerm(order.OrderId, nextBlock+types.ShortBlockWindow+1),
	)
	require.ErrorIs(t, err, types.ErrGoodTilBlockExceedsShortBlockWindow)

	memClob.AssertExpectations(t)
}

func TestCancelOrder_KeeperForwardsErrorsFromMemclob(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	ctx := ks.Ctx.WithIsCheckTx(true)
	ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
		Height:    14,
		Timestamp: time.Unix(int64(50), 0),
	})

	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	ctx = ctx.WithBlockHeight(14)

	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderShortTerm(
		shortTermOrder.OrderId,
		uint32(15),
	)).Return(nil, types.ErrMemClobCancelAlreadyExists)
	err := ks.ClobKeeper.CancelShortTermOrder(ctx, types.NewMsgCancelOrderShortTerm(shortTermOrder.OrderId, 15))
	require.ErrorIs(t, err, types.ErrMemClobCancelAlreadyExists)

	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	ks.ClobKeeper.SetLongTermOrderPlacement(ctx.WithIsCheckTx(false), longTermOrder, 15)
	ks.ClobKeeper.AddStatefulOrderIdExpiration(
		ctx,
		longTermOrder.MustGetUnixGoodTilBlockTime(),
		longTermOrder.GetOrderId(),
	)
	memClob.On("CancelOrder", ctx, types.NewMsgCancelOrderStateful(
		longTermOrder.OrderId,
		uint32(100),
	)).Return(nil, types.ErrMemClobCancelAlreadyExists)
	err = ks.ClobKeeper.CancelShortTermOrder(ctx, types.NewMsgCancelOrderStateful(longTermOrder.OrderId, 100))
	require.ErrorIs(t, err, types.ErrMemClobCancelAlreadyExists)

	memClob.AssertExpectations(t)
}

func TestPerformOrderCancellationStatefulValidation(t *testing.T) {
	blockHeight := uint32(5)
	blockTime := time.Unix(10, 0)

	tests := map[string]struct {
		setupDeliverTxState func(ctx sdk.Context, k *keeper.Keeper)
		msgCancelOrder      *types.MsgCancelOrder
		expectedErr         string
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
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: blockHeight - 1},
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
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
				k.AddStatefulOrderIdExpiration(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				)
			},
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: uint32(blockTime.Add(types.StatefulOrderTimeWindow).Unix()),
				},
			},
		},
		"stateful cancellation fails if order to cancel does not exist": {
			msgCancelOrder: &types.MsgCancelOrder{
				OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
					GoodTilBlockTime: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime() + 1,
				},
			},
			expectedErr: "Order Id to cancel does not exist.",
		},
		"stateful cancellation succeeds if existing order GTBT is less than cancel order GTBT.": {
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
				k.AddStatefulOrderIdExpiration(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
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
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
				k.AddStatefulOrderIdExpiration(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
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
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
				k.AddStatefulOrderIdExpiration(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				)
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
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
				k.AddStatefulOrderIdExpiration(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
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
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
				k.AddStatefulOrderIdExpiration(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
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
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					5,
				)
				k.AddStatefulOrderIdExpiration(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
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
			tApp := testapp.NewTestAppBuilder(t).
				// Disable non-determinism checks since we mutate keeper state directly.
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() cmt.GenesisDoc {
					genesis := testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *types.GenesisState) {
						state.ClobPairs = []types.ClobPair{
							{
								Metadata: &types.ClobPair_PerpetualClobMetadata{
									PerpetualClobMetadata: &types.PerpetualClobMetadata{
										PerpetualId: 0,
									},
								},
								Status:           types.ClobPair_STATUS_ACTIVE,
								StepBaseQuantums: 12,
								SubticksPerTick:  39,
							},
						}
					})
					return genesis
				}).Build()

			ctx := tApp.AdvanceToBlock(
				// Stateful validation happens at blockHeight+1 for short term order cancelations.
				blockHeight-1,
				testapp.AdvanceToBlockOptions{BlockTime: blockTime},
			)

			if tc.setupDeliverTxState != nil {
				tc.setupDeliverTxState(ctx.WithIsCheckTx(false), tApp.App.ClobKeeper)
			}

			resp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				*tc.msgCancelOrder)[0],
			)

			if tc.expectedErr != "" {
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Contains(t, resp.Log, tc.expectedErr)
			} else {
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}
		})
	}
}
