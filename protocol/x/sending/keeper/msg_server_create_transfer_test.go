package keeper_test

import (
	"context"
	sdkmath "cosmossdk.io/math"
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
	// Initialize mocks, keepers, and context.
	mockKeeper = &mocks.SendingKeeper{}
	ks := keepertest.SendingKeepers(t)
	ctx := ks.Ctx.WithBlockHeight(25)

	// Setup mocks.
	tc.setupMocks(ctx, mockKeeper)

	// Return message server and sdk context.
	return mockKeeper, keeper.NewMsgServerImpl(mockKeeper), ctx
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

func TestMsgServerSendFromModuleToAccount(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		testMsg    types.MsgSendFromModuleToAccount
		keeperResp error // mock keeper response
		// Expectations.
		expectedResp *types.MsgSendFromModuleToAccountResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgSendFromModuleToAccount{
				Authority:        lib.GovModuleAddress.String(),
				SenderModuleName: "community_treasury",
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
			},
			expectedResp: &types.MsgSendFromModuleToAccountResponse{},
		},
		"Failure: invalid authority": {
			testMsg: types.MsgSendFromModuleToAccount{
				Authority:        "12345",
				SenderModuleName: "community_treasury",
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
			},
			expectedErr: fmt.Sprintf(
				"invalid authority %s",
				"12345",
			),
		},
		"Failure: keeper method returns error": {
			testMsg: types.MsgSendFromModuleToAccount{
				Authority:        lib.GovModuleAddress.String(),
				SenderModuleName: "community_treasury",
				Recipient:        constants.CarlAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
			},
			keeperResp:  fmt.Errorf("keeper error"),
			expectedErr: "keeper error",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize Mocks and Context.
			mockKeeper := &mocks.SendingKeeper{}
			msgServer := keeper.NewMsgServerImpl(mockKeeper)
			ks := keepertest.SendingKeepers(t)
			mockKeeper.On("HasAuthority", tc.testMsg.Authority).Return(
				ks.SendingKeeper.HasAuthority(tc.testMsg.Authority),
			)
			mockKeeper.On("SendFromModuleToAccount", ks.Ctx, &tc.testMsg).Return(tc.keeperResp)

			resp, err := msgServer.SendFromModuleToAccount(ks.Ctx, &tc.testMsg)

			// Assert msg server response.
			require.Equal(t, tc.expectedResp, resp)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
