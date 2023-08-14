package keeper_test

import (
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgServerUpdateProposeParams(t *testing.T) {
	testMsg := types.MsgUpdateProposeParams{
		Authority: "authority",
		Params: types.ProposeParams{
			MaxBridgesPerBlock:           2,
			ProposeDelayDuration:         time.Second,
			SkipRatePpm:                  800_000,
			SkipIfBlockDelayedByDuration: time.Second,
		},
	}

	tests := map[string]struct {
		setupMocks   func(ctx sdk.Context, mck *mocks.BridgeKeeper)
		expectedResp *types.MsgUpdateProposeParamsResponse
		expectedErr  string
	}{
		"Success": {
			setupMocks: func(ctx sdk.Context, mck *mocks.BridgeKeeper) {
				mck.On("UpdateProposeParams", mock.Anything, testMsg.Params).Return(nil)
			},
			expectedResp: &types.MsgUpdateProposeParamsResponse{},
		},
		"Failure: keeper error is propagated": {
			setupMocks: func(ctx sdk.Context, mck *mocks.BridgeKeeper) {
				mck.On("UpdateProposeParams", mock.Anything, testMsg.Params).Return(
					errors.New("can't update event params"),
				)
			},
			expectedErr: "can't update event params",
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

			resp, err := msgServer.UpdateProposeParams(goCtx, &testMsg)

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
