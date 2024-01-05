package keeper_test

import (
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
				mockLogger := &mocks.Logger{}
				mockLogger.On(
					"Error",
					[]interface{}{
						testError.Error(),
						mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything,
					}...,
				).Return()
			},
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

			resp, err := msgServer.ProposedOperations(ctx, msg)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

			// Assert mock expectations.
			result := mockKeeper.AssertExpectations(t)
			require.True(t, result)
		})
	}
}
