package keeper_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUnlockShares(t *testing.T) {
	tests := map[string]struct {
		// Existing owner share unlocks.
		ownerShareUnlocks *types.OwnerShareUnlocks
		// Current block height.
		currentBlockHeight uint32
		// Msg.
		msg *types.MsgUnlockShares
		// Expected owner share unlocks after unlocking.
		expectedOwnerShareUnlocks *types.OwnerShareUnlocks
		// Expected error
		expectedErr string
	}{
		"Success - Authority is gov module, Unlocks all shares of Alice": {
			ownerShareUnlocks: &types.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(1)),
						UnlockBlockHeight: 3,
					},
				},
			},
			currentBlockHeight: 3,
			msg: &types.MsgUnlockShares{
				Authority:    lib.GovModuleAddress.String(),
				OwnerAddress: constants.AliceAccAddress.String(),
			},
			expectedOwnerShareUnlocks: nil,
		},
		"Success - Authority is delaymsg module, Unlocks some shares of Bob": {
			ownerShareUnlocks: &types.OwnerShareUnlocks{
				OwnerAddress: constants.BobAccAddress.String(),
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(55)),
						UnlockBlockHeight: 6,
					},
					{
						Shares:            types.BigIntToNumShares(big.NewInt(66)),
						UnlockBlockHeight: 7,
					},
				},
			},
			currentBlockHeight: 6,
			msg: &types.MsgUnlockShares{
				Authority:    delaymsgtypes.ModuleAddress.String(),
				OwnerAddress: constants.BobAccAddress.String(),
			},
			expectedOwnerShareUnlocks: &types.OwnerShareUnlocks{
				OwnerAddress: constants.BobAccAddress.String(),
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(66)),
						UnlockBlockHeight: 7,
					},
				},
			},
		},
		"Failure - Invalid authority": {
			msg: &types.MsgUnlockShares{
				Authority:    constants.AliceAccAddress.String(),
				OwnerAddress: constants.AliceAccAddress.String(),
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty owner": {
			msg: &types.MsgUnlockShares{
				Authority:    lib.GovModuleAddress.String(),
				OwnerAddress: "",
			},
			expectedErr: "owner address cannot be empty",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			k := tApp.App.VaultKeeper
			ms := keeper.NewMsgServerImpl(k)
			ctx := tApp.InitChain()

			if tc.currentBlockHeight > 1 {
				ctx = tApp.AdvanceToBlock(tc.currentBlockHeight, testapp.AdvanceToBlockOptions{})
			}

			if tc.ownerShareUnlocks != nil {
				err := k.SetOwnerShareUnlocks(ctx, tc.msg.OwnerAddress, *tc.ownerShareUnlocks)
				require.NoError(t, err)
			}

			res, err := ms.UnlockShares(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			newOwnerShareUnlocks, exists := k.GetOwnerShareUnlocks(ctx, tc.msg.OwnerAddress)
			if tc.expectedOwnerShareUnlocks != nil {
				// Verify that owner share unlocks are as expected.
				require.True(t, exists)
				require.Equal(t, *tc.expectedOwnerShareUnlocks, newOwnerShareUnlocks)
				// Verify that number of unlocked shares in response is as expected.
				expectedUnlockedShares := tc.ownerShareUnlocks.GetTotalLockedShares()
				expectedUnlockedShares.Sub(expectedUnlockedShares, newOwnerShareUnlocks.GetTotalLockedShares())
				require.Equal(
					t,
					types.BigIntToNumShares(expectedUnlockedShares),
					res.UnlockedShares,
				)
			} else {
				require.False(t, exists)
			}
		})
	}
}
