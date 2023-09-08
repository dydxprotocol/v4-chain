package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	keeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProposedOperations(t *testing.T) {
	testError := errors.New("error")
	operationsQueue := []types.OperationRaw{}

	tests := map[string]struct {
		setupMocks  func(ctx sdk.Context, mck *mocks.ClobKeeper)
		shouldPanic bool
		expectedErr error
	}{
		"Success": {
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("ProcessProposerOperations", ctx, operationsQueue).Return(nil)
			},
		},
		"Propagate Process Error": {
			setupMocks: func(ctx sdk.Context, mck *mocks.ClobKeeper) {
				mck.On("ProcessProposerOperations", ctx, operationsQueue).Return(testError)
			},
			shouldPanic: true,
			expectedErr: testError,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize Mocks and Context.
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			mockKeeper := &mocks.ClobKeeper{}
			ks := keepertest.NewClobKeepersTestContext(
				t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Setup mocks.
			blockHeight := int64(20)
			ctx := ks.Ctx.WithBlockHeight(blockHeight)
			tc.setupMocks(ctx, mockKeeper)

			// Define ProposedOperations receiver and arguments.
			msgServer := keeper.NewMsgServerImpl(mockKeeper)
			msg := &types.MsgProposedOperations{
				OperationsQueue: make([]types.OperationRaw, 0),
			}
			goCtx := sdk.WrapSDKContext(ctx)

			// Call ProposedOperations.
			if tc.shouldPanic {
				require.PanicsWithError(
					t,
					errorsmod.Wrapf(
						tc.expectedErr,
						"Block height: %d",
						blockHeight,
					).Error(),
					func() {
						msgServer.ProposedOperations(goCtx, msg) //nolint:errcheck
					},
				)
				return
			}

			resp, err := msgServer.ProposedOperations(goCtx, msg)
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Assert mock expectations.
			result := mockKeeper.AssertExpectations(t)
			require.True(t, result)
		})
	}
}
