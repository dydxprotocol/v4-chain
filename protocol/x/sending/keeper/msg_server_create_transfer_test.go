package keeper_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
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

					ctx := sdk.UnwrapSDKContext(goCtx)
					require.Len(t, ctx.EventManager().Events(), 1)
					event := ctx.EventManager().Events()[0]
					require.Equal(t, event.Type, types.EventTypeCreateTransfer)
					require.Equal(t, event.Attributes, []abci.EventAttribute{
						{
							Key:   types.AttributeKeySender,
							Value: msg.Transfer.Sender.Owner,
						},
						{
							Key:   types.AttributeKeySenderNumber,
							Value: fmt.Sprintf("%d", msg.Transfer.Sender.Number),
						},
						{
							Key:   types.AttributeKeyRecipient,
							Value: msg.Transfer.Recipient.Owner,
						},
						{
							Key:   types.AttributeKeyRecipientNumber,
							Value: fmt.Sprintf("%d", msg.Transfer.Recipient.Number),
						},
						{
							Key:   types.AttributeKeyAssetId,
							Value: fmt.Sprintf("%d", msg.Transfer.AssetId),
						},
						{
							Key:   types.AttributeKeyQuantums,
							Value: fmt.Sprintf("%d", msg.Transfer.Amount),
						},
					})
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

					ctx := sdk.UnwrapSDKContext(goCtx)
					require.Len(t, ctx.EventManager().Events(), 1)
					event := ctx.EventManager().Events()[0]
					require.Equal(t, event.Type, types.EventTypeDepositToSubaccount)
					require.Equal(t, event.Attributes, []abci.EventAttribute{
						{
							Key:   types.AttributeKeySender,
							Value: msg.Sender,
						},
						{
							Key:   types.AttributeKeyRecipient,
							Value: msg.Recipient.Owner,
						},
						{
							Key:   types.AttributeKeyRecipientNumber,
							Value: fmt.Sprintf("%d", msg.Recipient.Number),
						},
						{
							Key:   types.AttributeKeyAssetId,
							Value: fmt.Sprintf("%d", msg.AssetId),
						},
						{
							Key:   types.AttributeKeyQuantums,
							Value: fmt.Sprintf("%d", msg.Quantums),
						},
					})
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

					ctx := sdk.UnwrapSDKContext(goCtx)
					require.Len(t, ctx.EventManager().Events(), 1)
					event := ctx.EventManager().Events()[0]
					require.Equal(t, event.Type, types.EventTypeWithdrawFromSubaccount)
					require.Equal(t, event.Attributes, []abci.EventAttribute{
						{
							Key:   types.AttributeKeySender,
							Value: msg.Sender.Owner,
						},
						{
							Key:   types.AttributeKeySenderNumber,
							Value: fmt.Sprintf("%d", msg.Sender.Number),
						},
						{
							Key:   types.AttributeKeyRecipient,
							Value: msg.Recipient,
						},
						{
							Key:   types.AttributeKeyAssetId,
							Value: fmt.Sprintf("%d", msg.AssetId),
						},
						{
							Key:   types.AttributeKeyQuantums,
							Value: fmt.Sprintf("%d", msg.Quantums),
						},
					})
				}
			}

			// Assert mock expectations.
			result := mockKeeper.AssertExpectations(t)
			require.True(t, result)
		})
	}
}
