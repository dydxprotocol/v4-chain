package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	sdktest "github.com/dydxprotocol/v4/testutil/sdk"
	txtest "github.com/dydxprotocol/v4/testutil/sdk/tx"
	"github.com/dydxprotocol/v4/x/clob/ante"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
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
	cd := ante.NewClobDecorator(mockClobKeeper)
	antehandler := sdk.ChainAnteDecorators(cd)
	if tc.setupMocks != nil {
		tc.setupMocks(ctx, mockClobKeeper)
	}

	// Create Test Transcation.
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
		"Successfully places an order using a single message": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CheckTxPlaceOrder",
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
		"Successfully places a long term order using a single message": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder_LongTerm},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CheckTxPlaceOrder",
					ctx,
					constants.Msg_PlaceOrder_LongTerm,
				).Return(
					satypes.BaseQuantums(0),
					clobtypes.Success,
					nil,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"Successfully places a conditional order using a single message": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder_Conditional},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CheckTxPlaceOrder",
					ctx,
					constants.Msg_PlaceOrder_Conditional,
				).Return(
					satypes.BaseQuantums(0),
					clobtypes.Success,
					nil,
				)
			},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"PlaceOrder is not called on keeper during deliver": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: false,
			expectedErr:             nil,
		},
		"PlaceOrder is not called on keeper during simulate": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: false,
			isSimulate:              true,
			expectedErr:             nil,
		},
		"PlaceOrder is not called on keeper during re-check": {
			msgs:                      []sdk.Msg{constants.Msg_PlaceOrder},
			useWithIsCheckTxContext:   false,
			useWithIsRecheckTxContext: true,
			isSimulate:                false,
			expectedErr:               nil,
		},
		"Fails if PlaceOrder returns an error": {
			msgs: []sdk.Msg{constants.Msg_PlaceOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CheckTxPlaceOrder",
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
		"Fails if there are multiple off-chain places": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
		"Fails if there are multiple long term and conditional orders": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder_LongTerm, constants.Msg_PlaceOrder_Conditional},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
		"Fails if there is a mix of long term and short term orders": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_PlaceOrder_LongTerm},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
		"Fails if there is a mix of conditional and short term orders": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_PlaceOrder_Conditional},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
		"Fails if there are a mix of off-chain and on-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_Send},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}

func TestIsClobOffChainTransaction(t *testing.T) {
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
			expectedErr:    errors.ErrInvalidRequest,
		},
		"Returns false and error for multiple `CancelOrder` messages": {
			msgs:           []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_CancelOrder},
			expectedResult: false,
			expectedErr:    errors.ErrInvalidRequest,
		},
		"Returns false and error for mix of `PlaceOrder` and `CancelOrder` messages": {
			msgs:           []sdk.Msg{constants.Msg_PlaceOrder, constants.Msg_CancelOrder},
			expectedResult: false,
			expectedErr:    errors.ErrInvalidRequest,
		},
		"Returns false and error for mix of `MsgSend` and `PlaceOrder` messages": {
			msgs:           []sdk.Msg{constants.Msg_Send, constants.Msg_PlaceOrder},
			expectedResult: false,
			expectedErr:    errors.ErrInvalidRequest,
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
			ctx, _, _ := sdktest.NewSdkContextWithMultistore()

			// Invoke the function under test.
			result, err := ante.IsOffChainSingleClobMsgTx(ctx, tx)

			// Assert the results.
			require.Equal(t, tc.expectedResult, result)
			require.ErrorIs(t, tc.expectedErr, err)
		})
	}
}

func TestClobDecorator_MsgCancelOrder(t *testing.T) {
	tests := map[string]TestCase{
		"Successfully cancels an order using a single message": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CheckTxCancelOrder",
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
		"Works with any number of off-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_Send, constants.Msg_Send},
			useWithIsCheckTxContext: true,
			expectedErr:             nil,
		},
		"CancelOrder is not called on keeper during deliver": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder},
			useWithIsCheckTxContext: false,
			expectedErr:             nil,
		},
		"CancelOrder is not called on keeper during simulate": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder},
			useWithIsCheckTxContext: false,
			isSimulate:              true,
			expectedErr:             nil,
		},
		"CancelOrder is not called on keeper during re-check": {
			msgs:                      []sdk.Msg{constants.Msg_CancelOrder},
			useWithIsCheckTxContext:   false,
			useWithIsRecheckTxContext: true,
			isSimulate:                false,
			expectedErr:               nil,
		},
		"Fails if CancelOrder returns an error": {
			msgs: []sdk.Msg{constants.Msg_CancelOrder},
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("CheckTxCancelOrder",
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
		"Fails if there are multiple off-chain cancels": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_CancelOrder},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
		"Fails if there are multiple off-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_PlaceOrder},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
		"Fails if there are a mix of off-chain and on-chain messages": {
			msgs:                    []sdk.Msg{constants.Msg_CancelOrder, constants.Msg_Send},
			useWithIsCheckTxContext: true,
			expectedErr:             errors.ErrInvalidRequest,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}
