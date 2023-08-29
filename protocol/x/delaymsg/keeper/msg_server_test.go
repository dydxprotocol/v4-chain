package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	AcceptedAuthority = authtypes.NewModuleAddress(bridgemoduletypes.ModuleName).String()
	InvalidAuthority  = "INVALID_AUTHORITY"
	TestError         = fmt.Errorf("test error")
	TestMsgId         = uint32(0)

	ValidAuthorities = map[string]struct{}{
		AcceptedAuthority: {},
	}
	IsValidAuthority = func(authority string) bool {
		_, ok := ValidAuthorities[authority]
		return ok
	}

	ValidDelayMsg = &types.MsgDelayMessage{
		Authority: AcceptedAuthority,
	}
	InvalidDelayMsg = &types.MsgDelayMessage{
		Authority: InvalidAuthority,
	}

	DelayMsgResponse = &types.MsgDelayMessageResponse{
		Id: uint64(TestMsgId),
	}
)

func setupMockWithValidReturnValues(ctx sdk.Context, mck *mocks.DelayMsgKeeper) {
	mck.On("DelayMessageByBlocks", ctx, mock.Anything, mock.Anything).Return(TestMsgId, nil)
	mck.On("DecodeMessage", mock.Anything, mock.Anything).Return(nil)
	mck.On("HasAuthority", mock.MatchedBy(IsValidAuthority)).Return(true)
	mck.On("HasAuthority", mock.Anything).Return(false)
}

func setupMockWithDecodeFailure(ctx sdk.Context, mck *mocks.DelayMsgKeeper) {
	mck.On("DelayMessageByBlocks", ctx, mock.Anything, mock.Anything).Return(TestMsgId, nil)
	mck.On("DecodeMessage", mock.Anything, mock.Anything).Return(TestError)
	mck.On("HasAuthority", mock.MatchedBy(IsValidAuthority)).Return(true)
	mck.On("HasAuthority", mock.Anything).Return(false)
}

func setupMockWithDelayMessageFailure(ctx sdk.Context, mck *mocks.DelayMsgKeeper) {
	mck.On("DelayMessageByBlocks", ctx, mock.Anything, mock.Anything).Return(TestMsgId, TestError)
	mck.On("DecodeMessage", mock.Anything, mock.Anything).Return(nil)
	mck.On("HasAuthority", mock.MatchedBy(IsValidAuthority)).Return(true)
	mck.On("HasAuthority", mock.Anything).Return(false)
}

func TestDelayMessage(t *testing.T) {
	tests := map[string]struct {
		msg         *types.MsgDelayMessage
		setupMocks  func(ctx sdk.Context, mck *mocks.DelayMsgKeeper)
		expectedErr error
	}{
		"Success": {
			setupMocks: setupMockWithValidReturnValues,
			msg:        ValidDelayMsg,
		},
		"Panics when signed by invalid authority": {
			setupMocks: setupMockWithValidReturnValues,
			msg: &types.MsgDelayMessage{
				Authority: InvalidAuthority,
			},
			expectedErr: fmt.Errorf(
				"%v is not recognized as a valid authority for sending delayed messages",
				InvalidAuthority,
			),
		},
		"Panics if message does not decode": {
			setupMocks:  setupMockWithDecodeFailure,
			msg:         ValidDelayMsg,
			expectedErr: fmt.Errorf("UnmarshalInterface for DelayedMessage failed, err = %w", TestError),
		},
		"Panics if DelayMessageByBlocks returns an error": {
			setupMocks:  setupMockWithDelayMessageFailure,
			msg:         ValidDelayMsg,
			expectedErr: fmt.Errorf("DelayMessageByBlocks failed, err  = %w", TestError),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockKeeper := &mocks.DelayMsgKeeper{}
			msgServer := keeper.NewMsgServerImpl(mockKeeper)
			ctx, _, _, _, _, _ := keepertest.DelayMsgKeepers(t)
			tc.setupMocks(ctx, mockKeeper)
			goCtx := sdk.WrapSDKContext(ctx)

			if tc.expectedErr != nil {
				require.PanicsWithError(
					t,
					tc.expectedErr.Error(),
					func() {
						_, _ = msgServer.DelayMessage(goCtx, tc.msg)
					},
				)
			} else {
				resp, err := msgServer.DelayMessage(goCtx, tc.msg)
				require.NoError(t, err)
				require.Equal(t, DelayMsgResponse, resp)
			}
		})
	}
}
