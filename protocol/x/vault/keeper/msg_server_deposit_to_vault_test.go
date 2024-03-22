package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgServerDepositToVault(t *testing.T) {
	testMsg := constants.MsgDepositToVault_Clob0_Alice0_100

	tests := map[string]struct {
		setupMocks   func(ctx sdk.Context, mck *mocks.VaultKeeper)
		expectedResp *types.MsgDepositToVaultResponse
		expectedErr  string
	}{
		"Success": {
			setupMocks: func(ctx sdk.Context, mck *mocks.VaultKeeper) {
				mck.On(
					"HandleMsgDepositToVault",
					mock.Anything,
					testMsg,
				).Return(nil)
			},
			expectedResp: &types.MsgDepositToVaultResponse{},
		},
		"Failure: keeper error is propagated": {
			setupMocks: func(ctx sdk.Context, mck *mocks.VaultKeeper) {
				mck.On(
					"HandleMsgDepositToVault",
					mock.Anything,
					testMsg,
				).Return(
					errors.New("deposit failed"),
				)
			},
			expectedErr: "deposit failed",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize mocks and context.
			mockKeeper := &mocks.VaultKeeper{}
			msgServer := keeper.NewMsgServerImpl(mockKeeper)
			ctx, _, _ := keepertest.VaultKeepers(t)
			tc.setupMocks(ctx, mockKeeper)

			resp, err := msgServer.DepositToVault(ctx, testMsg)

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
