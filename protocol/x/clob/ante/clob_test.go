package ante_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	txtest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk/tx"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/ante"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	msgs                      []sdk.Msg
	setupMocks                func(ctx sdk.Context, mck *mocks.ClobKeeper)
	useWithIsCheckTxContext   bool
	useWithIsRecheckTxContext bool
	isSimulate                bool
	expectedErr               error
}

func runTestCase(t *testing.T, tc TestCase) {
	// Setup Test Context.
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()

	if tc.useWithIsCheckTxContext {
		ctx = ctx.WithIsCheckTx(tc.useWithIsCheckTxContext)
	}

	if tc.useWithIsRecheckTxContext {
		ctx = ctx.WithIsReCheckTx(tc.useWithIsRecheckTxContext)
	}

	if tc.useWithIsRecheckTxContext && tc.useWithIsCheckTxContext {
		t.Error("Expected only one of useWithIsCheckTxContext or useWithIsCheckTxContext to be true")
	}

	// Setup AnteHandler.
	mockClobKeeper := &mocks.ClobKeeper{}
	mockClobKeeper.On("Logger", mock.Anything).Return(log.NewNopLogger()).Maybe()
	mockClobKeeper.On("IsInitialized").Return(true).Maybe()
	cd := ante.NewClobDecorator(mockClobKeeper)
	antehandler := sdk.ChainAnteDecorators(cd)
	if tc.setupMocks != nil {
		tc.setupMocks(ctx, mockClobKeeper)
	}

	// Create Test Transaction.
	priv1, _, _ := testdata.KeyTestPubAddr()
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := txtest.CreateTestTx(privs, accNums, accSeqs, "dydx", tc.msgs)
	require.NoError(t, err)

	// Call Antehandler.
	_, err = antehandler(ctx, tx, tc.isSimulate)

	// Assert error expectations.
	if tc.expectedErr != nil {
		require.ErrorIs(t, tc.expectedErr, err)
	} else {
		require.NoError(t, err)
	}

	// Assert mock expectations.
	result := mockClobKeeper.AssertExpectations(t)
	require.True(t, result)
}

func TestClobDecorator_MsgPlaceOrder(t *testing.T) {
	tests := map[string]TestCase{
		"Successfully places a short term order using a single message": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("PlaceShortTermOrder",
					ctx,
					constants.Msg_PlaceOrder,
				).Return(
					satypes.BaseQuantums(0),
					clobtypes.Success,
					nil,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"Successfully places a stateful order using a single message": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder_LongTerm},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("PlaceStatefulOrder",
					ctx,
					constants.Msg_PlaceOrder_LongTerm,
					false,
				).Return(
					nil,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"Successfully places a conditional order using a single message": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder_Conditional},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("PlaceStatefulOrder",
					ctx,
					constants.Msg_PlaceOrder_Conditional,
					false,
				).Return(
					nil,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"PlaceShortTermOrder is not called on keeper during deliver": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: false,
			expectedErr:             nil,
		},
		"PlaceShortTermOrder is not called on keeper during simulate": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: false,
			isSimulate:              true,
			expectedErr:             nil,
		},
		"PlaceShortTermOrder is not called on keeper during re-check": {
			msgs:                      []sdk.Msg{constants.Msg_PlaceOrder},
			useWithIsCheckTxContext:   false,
			useWithIsRecheckTxContext: true,
			isSimulate:                false,
			expectedErr:               nil,
		},
		"PlaceStatefulOrder is not called on keeper during deliver": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder_LongTerm},
			useWithIsCheckTxContext: false,
			expectedErr:             nil,
		},
		"PlaceStatefulOrder is not called on keeper during simulate": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder_LongTerm},
			useWithIsCheckTxContext: false,
			isSimulate:              true,
			expectedErr:             nil,
		},
		"PlaceStatefulOrder is called on keeper during re-check": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder_LongTerm},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("PlaceStatefulOrder",
					ctx,
					constants.Msg_PlaceOrder_LongTerm,
					false,
				).Return(
					nil,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"Fails if PlaceShortTermOrder returns an error": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("PlaceShortTermOrder",
					ctx,
					constants.Msg_PlaceOrder,
				).Return(
					satypes.BaseQuantums(0),
					clobtypes.OrderStatus(0),
					clobtypes.ErrHeightExceedsGoodTilBlock,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             clobtypes.ErrHeightExceedsGoodTilBlock,
		},
		"Fails if PlaceStatefulOrder returns an error": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder_LongTerm},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("PlaceStatefulOrder",
					ctx,
					constants.Msg_PlaceOrder_LongTerm,
					false,
				).Return(
					clobtypes.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             clobtypes.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
		},
		"Fails if there are multiple off-chain places": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: true,
			expectedErr:             sdkerrors.ErrInvalidRequest,
		},
		"Fails if there is a mix of long term and short term orders": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_PlaceOrder_LongTerm},
			useWithIsCheckTxContext: true,
			expectedErr:             sdkerrors.ErrInvalidRequest,
		},
		"Fails if there is a mix of conditional and short term orders": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_PlaceOrder_Conditional},
			useWithIsCheckTxContext: true,
			expectedErr:             sdkerrors.ErrInvalidRequest,
		},
		"Fails if there are a mix of off-chain and on-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_Send},
			useWithIsCheckTxContext: true,
			expectedErr:             sdkerrors.ErrInvalidRequest,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}

func TestIsClobTransaction(t *testing.T) {
	tests := map[string]struct {
		msgs           []sdk.Msg
		expectedResult bool
		expectedErr    error
	}{
		"Returns false for MsgSend": {
			msgs:           []sdk.Msg{constants.Msg_Send},
			expectedResult: false,
			expectedErr:    nil,
		},
		"Returns false for MsgTransfer": {
			msgs:           []sdk.Msg{constants.Msg_Transfer},
			expectedResult: false,
			expectedErr:    nil,
		},
		"Returns false for no messages": {
			msgs:           []sdk.Msg{},
			expectedResult: false,
			expectedErr:    nil,
		},
		"Returns false and error for multiple `PlaceOrder` message": {
			msgs:           []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_PlaceOrder},
			expectedResult: false,
			expectedErr:    sdkerrors.ErrInvalidRequest,
		},
		"Returns false and error for multiple `CancelOrder` messages": {
			msgs:           []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_CancelOrder},
			expectedResult: false,
			expectedErr:    sdkerrors.ErrInvalidRequest,
		},
		"Returns false and error for mix of `PlaceOrder` and `CancelOrder` messages": {
			msgs:           []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_CancelOrder},
			expectedResult: false,
			expectedErr:    sdkerrors.ErrInvalidRequest,
		},
		"Returns false and error for mix of `MsgSend` and `PlaceOrder` messages": {
			msgs:           []sdk.Msg{constants.Msg_Send, constants.Msg_PlaceOrder},
			expectedResult: false,
			expectedErr:    sdkerrors.ErrInvalidRequest,
		},
		"Returns true for a `CancelOrder` message": {
			msgs:           []sdk.Msg{constants.Msg_CancelOrder},
			expectedResult: true,
			expectedErr:    nil,
		},
		"Returns true for a `PlaceOrder` message": {
			msgs:           []sdk.Msg{constants.Msg_PlaceOrder},
			expectedResult: true,
			expectedErr:    nil,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize some test setup which builds a test transaction from a slice of messages.
			var reg codectypes.InterfaceRegistry
			protoCfg := authtx.NewTxConfig(codec.NewProtoCodec(reg), authtx.DefaultSignModes)
			builder := protoCfg.NewTxBuilder()
			err := builder.SetMsgs(tc.msgs...)
			require.NoError(t, err)
			tx := builder.GetTx()

			// Invoke the function under test.
			result, err := ante.IsSingleClobMsgTx(tx)

			// Assert the results.
			require.Equal(t, tc.expectedResult, result)
			require.ErrorIs(t, tc.expectedErr, err)
		})
	}
}

func TestClobDecorator_MsgCancelOrder(t *testing.T) {
	tests := map[string]TestCase{
		"Successfully cancels a short term order using a single message": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CancelShortTermOrder",
					ctx,
					clobtypes.NewMsgCancelOrderShortTerm(
						constants.Msg_CancelOrder.OrderId,
						constants.Msg_CancelOrder.GetGoodTilBlock(),
					),
				).Return(nil)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"Successfully cancels a long term order using a single message": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder_LongTerm},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CancelStatefulOrder",
					ctx,
					clobtypes.NewMsgCancelOrderStateful(
						constants.Msg_CancelOrder_LongTerm.OrderId,
						constants.Msg_CancelOrder_LongTerm.GetGoodTilBlockTime(),
					),
				).Return(nil)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"Successfully cancels a conditional order using a single message": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder_Conditional},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CancelStatefulOrder",
					ctx,
					clobtypes.NewMsgCancelOrderStateful(
						constants.Msg_CancelOrder_Conditional.OrderId,
						constants.Msg_CancelOrder_Conditional.GetGoodTilBlockTime(),
					),
				).Return(nil)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"Works with any number of off-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_Send, constants.Msg_Send},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"CancelShortTermOrder is not called on keeper during deliver": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder},
			useWithIsCheckTxContext: false,
			expectedErr:             nil,
		},
		"CancelShortTermOrder is not called on keeper during simulate": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder},
			useWithIsCheckTxContext: false,
			isSimulate:              true,
			expectedErr:             nil,
		},
		"CancelShortTermOrder is not called on keeper during re-check": {
			msgs:                      []sdk.Msg{constants.Msg_CancelOrder},
			useWithIsCheckTxContext:   false,
			useWithIsRecheckTxContext: true,
			isSimulate:                false,
			expectedErr:               nil,
		},
		"CancelStatefulOrder is not called on keeper during deliver": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder_LongTerm},
			useWithIsCheckTxContext: false,
			expectedErr:             nil,
		},
		"CancelStatefulOrder is not called on keeper during simulate": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder_LongTerm},
			useWithIsCheckTxContext: false,
			isSimulate:              true,
			expectedErr:             nil,
		},
		"CancelStatefulOrder is called on keeper during re-check": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder_LongTerm},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CancelStatefulOrder",
					ctx,
					clobtypes.NewMsgCancelOrderStateful(
						constants.Msg_CancelOrder_LongTerm.OrderId,
						constants.Msg_CancelOrder_LongTerm.GetGoodTilBlockTime(),
					),
				).Return(nil)
			},
			useWithIsCheckTxContext:   false,
			useWithIsRecheckTxContext: true,
			isSimulate:                false,
			expectedErr:               nil,
		},
		"Fails if CancelShortTermOrder returns an error": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CancelShortTermOrder",
					ctx,
					clobtypes.NewMsgCancelOrderShortTerm(
						constants.Msg_CancelOrder.OrderId,
						constants.Msg_CancelOrder.GetGoodTilBlock(),
					),
				).Return(clobtypes.ErrHeightExceedsGoodTilBlock)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             clobtypes.ErrHeightExceedsGoodTilBlock,
		},
		"Fails if CancelStatefulOrder returns an error": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder_LongTerm},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CancelStatefulOrder",
					ctx,
					clobtypes.NewMsgCancelOrderStateful(
						constants.Msg_CancelOrder_LongTerm.OrderId,
						constants.Msg_CancelOrder_LongTerm.GetGoodTilBlockTime(),
					),
				).Return(clobtypes.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             clobtypes.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
		},
		"Fails if there are multiple off-chain cancels": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_CancelOrder},
			useWithIsCheckTxContext: true,
			expectedErr:             sdkerrors.ErrInvalidRequest,
		},
		"Fails if there are multiple off-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: true,
			expectedErr:             sdkerrors.ErrInvalidRequest,
		},
		"Fails if there are a mix of off-chain and on-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_Send},
			useWithIsCheckTxContext: true,
			expectedErr:             sdkerrors.ErrInvalidRequest,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}
