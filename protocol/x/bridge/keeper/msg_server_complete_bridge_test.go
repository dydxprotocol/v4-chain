package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgServerCompleteBridge(t *testing.T) {
	testMsg := types.MsgCompleteBridge{
		Authority: "authority",
		Event:     constants.BridgeEvent_Id0_Height0,
	}

	tests := map[string]struct {
		setupMocks   func(ctx sdk.Context, mck *mocks.BridgeKeeper)
		expectedResp *types.MsgCompleteBridgeResponse
		expectedErr  string
	}{
		"Success": {
			setupMocks: func(ctx sdk.Context, mck *mocks.BridgeKeeper) {
				mck.On("CompleteBridge", mock.Anything, testMsg.Event).Return(nil)
			},
			expectedResp: &types.MsgCompleteBridgeResponse{},
		},
		"Failure: keeper error is propagated": {
			setupMocks: func(ctx sdk.Context, mck *mocks.BridgeKeeper) {
				mck.On("CompleteBridge", mock.Anything, testMsg.Event).Return(
					errors.New("can't complete bridge"),
				)
			},
			expectedErr: "can't complete bridge",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize Mocks and Context.
			mockKeeper := &mocks.BridgeKeeper{}
			msgServer := keeper.NewMsgServerImpl(mockKeeper)
			ctx, _, _, _, _, _ := keepertest.BridgeKeepers(t)
			tc.setupMocks(ctx, mockKeeper)
			goCtx := sdk.WrapSDKContext(ctx)

			resp, err := msgServer.CompleteBridge(goCtx, &testMsg)

			// Assert msg server response.
			require.Equal(t, tc.expectedResp, resp)
			if tc.expectedErr != "" {
				require.Equal(t, tc.expectedErr, err.Error())
			}

			// Assert mock expectations.
			result := mockKeeper.AssertExpectations(t)
			require.True(t, result)
		})
	}
}
