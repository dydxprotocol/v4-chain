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
		// Msg.
		msg *types.MsgUnlockShares
		// Expected error
		expectedErr string
	}{
		"Success - Authority is gov module, Unlocks shares of Alice": {
			msg: &types.MsgUnlockShares{
				Authority:    lib.GovModuleAddress.String(),
				OwnerAddress: constants.AliceAccAddress.String(),
			},
		},
		"Success - Authority is delaymsg module, Unlocks shares of Bob": {
			msg: &types.MsgUnlockShares{
				Authority:    delaymsgtypes.ModuleAddress.String(),
				OwnerAddress: constants.BobAccAddress.String(),
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

			testBlockHeight := uint32(7)
			testShares := big.NewInt(123_456_789)
			ctx := tApp.AdvanceToBlock(testBlockHeight, testapp.AdvanceToBlockOptions{})

			ownerShareUnlocks := types.OwnerShareUnlocks{
				OwnerAddress: tc.msg.OwnerAddress,
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(testShares),
						UnlockBlockHeight: testBlockHeight,
					},
				},
			}
			if tc.msg.OwnerAddress != "" {
				err := k.SetOwnerShareUnlocks(ctx, tc.msg.OwnerAddress, ownerShareUnlocks)
				require.NoError(t, err)
			}

			res, err := ms.UnlockShares(ctx, tc.msg)
			newLockedShares, exists := k.GetOwnerShareUnlocks(ctx, tc.msg.OwnerAddress)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				if tc.msg.OwnerAddress != "" {
					require.True(t, exists)
					require.Equal(t, ownerShareUnlocks, newLockedShares)
				}
			} else {
				require.NoError(t, err)
				require.False(t, exists)
				require.Equal(t, types.BigIntToNumShares(testShares), res.UnlockedShares)
			}
		})
	}
}
