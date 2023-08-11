package keeper_test

import (
	"context"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/sending/keeper"
	"github.com/dydxprotocol/v4/x/sending/types"
	"github.com/stretchr/testify/require"
)

type MsgServerTransferTestCase struct {
	setupMocks  func(ctx sdk.Context, mck *mocks.SendingKeeper)
	expectedErr error
	shouldPanic bool
}

func createMsgServerTransferTestCases[
	T *types.Transfer | *types.MsgDepositToSubaccount | *types.MsgWithdrawFromSubaccount,
](
	mockMethodName string,
	msg T,
) map[string]MsgServerTransferTestCase {
	testError := errors.New("error")

	return map[string]MsgServerTransferTestCase{
		"Success": {
			setupMocks: func(ctx sdk.Context, mck *mocks.SendingKeeper) {
				mck.On(mockMethodName, ctx, msg).Return(nil)
			},
			expectedErr: nil,
		},
		"Propagate Error": {
			setupMocks: func(ctx sdk.Context, mck *mocks.SendingKeeper) {
				mck.On(mockMethodName, ctx, msg).Return(testError)
			},
			expectedErr: testError,
		},
		"Propagate Panic": {
			setupMocks: func(ctx sdk.Context, mck *mocks.SendingKeeper) {
				mck.On(mockMethodName, ctx, msg).Panic(testError.Error())
			},
			shouldPanic: true,
			expectedErr: testError,
		},
	}
}

func setUpTestCase(
	t *testing.T,
	tc MsgServerTransferTestCase,
) (
	mockKeeper *mocks.SendingKeeper,
	msgServer types.MsgServer,
	goCtx context.Context,
) {
	// Initialize Mocks and Context.
	mockKeeper = &mocks.SendingKeeper{}
	ctx, _, _, _, _, _, _, _ := keepertest.SendingKeepers(t)
	ctx = ctx.WithBlockHeight(25)

	// Setup mocks.
	tc.setupMocks(ctx, mockKeeper)

	// Return message server and sdk context.
	return mockKeeper, keeper.NewMsgServerImpl(mockKeeper), sdk.WrapSDKContext(ctx)
}

func TestCreateTransfer(t *testing.T) {
	msg := constants.Msg_Transfer
	tests := createMsgServerTransferTestCases("ProcessTransfer", msg.Transfer)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockKeeper, msgServer, goCtx := setUpTestCase(t, tc)

			if tc.shouldPanic {
				// Call CreateTransfer.
				require.PanicsWithValue(t, tc.expectedErr.Error(), func() {
					//nolint:errcheck
					msgServer.CreateTransfer(goCtx, msg)
				})
			} else {
				// Call CreateTransfer.
				resp, err := msgServer.CreateTransfer(goCtx, msg)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
					require.NotNil(t, resp)
				}
			}

			// Assert mock expectations.
			result := mockKeeper.AssertExpectations(t)
			require.True(t, result)
		})
	}
}

func TestDepositToSubaccount(t *testing.T) {
	msg := constants.MsgDepositToSubaccount_Alice_To_Alice_Num0_500
	tests := createMsgServerTransferTestCases("ProcessDepositToSubaccount", &msg)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockKeeper, msgServer, goCtx := setUpTestCase(t, tc)

			if tc.shouldPanic {
				// Call DepositToSubaccount.
				require.PanicsWithValue(t, tc.expectedErr.Error(), func() {
					//nolint:errcheck
					msgServer.DepositToSubaccount(goCtx, &msg)
				})
			} else {
				// Call DepositToSubaccount.
				resp, err := msgServer.DepositToSubaccount(goCtx, &msg)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
					require.NotNil(t, resp)
				}
			}

			// Assert mock expectations.
			result := mockKeeper.AssertExpectations(t)
			require.True(t, result)
		})
	}
}

func TestWithdrawFromSubaccount(t *testing.T) {
	msg := constants.MsgWithdrawFromSubaccount_Alice_Num0_To_Alice_500
	tests := createMsgServerTransferTestCases("ProcessWithdrawFromSubaccount", &msg)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockKeeper, msgServer, goCtx := setUpTestCase(t, tc)

			if tc.shouldPanic {
				// Call DepositToSubaccount.
				require.PanicsWithValue(t, tc.expectedErr.Error(), func() {
					//nolint:errcheck
					msgServer.WithdrawFromSubaccount(goCtx, &msg)
				})
			} else {
				// Call DepositToSubaccount.
				resp, err := msgServer.WithdrawFromSubaccount(goCtx, &msg)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
					require.NotNil(t, resp)
				}
			}

			// Assert mock expectations.
			result := mockKeeper.AssertExpectations(t)
			require.True(t, result)
		})
	}
}
