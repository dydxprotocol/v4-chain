package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutildelaymsg "github.com/dydxprotocol/v4-chain/protocol/testutil/delaymsg"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetTotalShares(t *testing.T) {
	tests := map[string]struct {
		// Function to set total shares.
		setFunc func(k keeper.Keeper, ctx sdk.Context) error
		// Expected total shares.
		expectedTotalShares vaulttypes.NumShares
		// Expected error.
		expectedErr error
	}{
		"Success: default total shares is 0": {
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(0),
			},
		},
		"Success: set total shares to 0": {
			setFunc: func(k keeper.Keeper, ctx sdk.Context) error {
				return k.SetTotalShares(ctx, vaulttypes.NumShares{
					NumShares: dtypes.NewInt(0),
				})
			},
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(0),
			},
		},
		"Success: set total shares to 777": {
			setFunc: func(k keeper.Keeper, ctx sdk.Context) error {
				return k.SetTotalShares(ctx, vaulttypes.NumShares{
					NumShares: dtypes.NewInt(777),
				})
			},
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(777),
			},
		},
		"Failure: set total shares to -1": {
			setFunc: func(k keeper.Keeper, ctx sdk.Context) error {
				return k.SetTotalShares(ctx, vaulttypes.NumShares{
					NumShares: dtypes.NewInt(-1),
				})
			},
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(0),
			},
			expectedErr: vaulttypes.ErrNegativeShares,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.setFunc != nil {
				if tc.expectedErr != nil {
					require.ErrorIs(t, tc.setFunc(k, ctx), tc.expectedErr)
				} else {
					require.NoError(t, tc.setFunc(k, ctx))
				}
			}

			require.Equal(
				t,
				tc.expectedTotalShares,
				k.GetTotalShares(ctx),
			)
		})
	}
}

func TestGetSetOwnerShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	alice := constants.AliceAccAddress.String()
	bob := constants.BobAccAddress.String()

	// Get owners shares for Alice.
	_, exists := k.GetOwnerShares(ctx, alice)
	require.Equal(t, false, exists)

	// Set owner shares for Alice and get.
	numShares := vaulttypes.BigIntToNumShares(
		big.NewInt(7),
	)
	err := k.SetOwnerShares(ctx, alice, numShares)
	require.NoError(t, err)
	got, exists := k.GetOwnerShares(ctx, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Alice and then get.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(456),
	)
	err = k.SetOwnerShares(ctx, alice, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(0),
	)
	err = k.SetOwnerShares(ctx, bob, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob to a negative value.
	// Should get error and owner shares should remain unchanged.
	numSharesInvalid := vaulttypes.BigIntToNumShares(
		big.NewInt(-1),
	)
	err = k.SetOwnerShares(ctx, bob, numSharesInvalid)
	require.ErrorIs(t, err, vaulttypes.ErrNegativeShares)
	got, exists = k.GetOwnerShares(ctx, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)
}

func TestGetAllOwnerShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Get all owner shares when there's no owner.
	allOwnerShares := k.GetAllOwnerShares(ctx)
	require.Equal(t, []vaulttypes.OwnerShare{}, allOwnerShares)

	// Set alice and bob as owners and get all owner shares.
	alice := constants.AliceAccAddress.String()
	aliceShares := vaulttypes.BigIntToNumShares(big.NewInt(7))
	bob := constants.BobAccAddress.String()
	bobShares := vaulttypes.BigIntToNumShares(big.NewInt(123))

	err := k.SetOwnerShares(ctx, alice, aliceShares)
	require.NoError(t, err)
	err = k.SetOwnerShares(ctx, bob, bobShares)
	require.NoError(t, err)

	allOwnerShares = k.GetAllOwnerShares(ctx)
	require.ElementsMatch(
		t,
		[]vaulttypes.OwnerShare{
			{
				Owner:  alice,
				Shares: aliceShares,
			},
			{
				Owner:  bob,
				Shares: bobShares,
			},
		},
		allOwnerShares,
	)
}

func TestGetSetOwnerShareUnlocks(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	alice := constants.AliceAccAddress.String()
	bob := constants.BobAccAddress.String()

	// Get share unlocks for Alice.
	_, exists := k.GetOwnerShareUnlocks(ctx, alice)
	require.Equal(t, false, exists)

	// Set share unlocks for Alice and get.
	aliceShareUnlocks := vaulttypes.OwnerShareUnlocks{
		OwnerAddress: alice,
		ShareUnlocks: []vaulttypes.ShareUnlock{
			{
				Shares:            vaulttypes.BigIntToNumShares(big.NewInt(7)),
				UnlockBlockHeight: 1,
			},
		},
	}
	err := k.SetOwnerShareUnlocks(ctx, alice, aliceShareUnlocks)
	require.NoError(t, err)
	got, exists := k.GetOwnerShareUnlocks(ctx, alice)
	require.Equal(t, true, exists)
	require.Equal(t, aliceShareUnlocks, got)

	// Set share unlocks for Bob and then get.
	bobLockedShares := vaulttypes.OwnerShareUnlocks{
		OwnerAddress: bob,
		ShareUnlocks: []vaulttypes.ShareUnlock{
			{
				Shares:            vaulttypes.BigIntToNumShares(big.NewInt(901)),
				UnlockBlockHeight: 76,
			},
			{
				Shares:            vaulttypes.BigIntToNumShares(big.NewInt(333)),
				UnlockBlockHeight: 965,
			},
		},
	}
	err = k.SetOwnerShareUnlocks(ctx, bob, bobLockedShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShareUnlocks(ctx, bob)
	require.Equal(t, true, exists)
	require.Equal(t, bobLockedShares, got)

	// Set invalid share unlocks for Bob.
	// Should get error and share unlocks should remain unchanged.
	bobLockedShares.OwnerAddress = ""
	err = k.SetOwnerShareUnlocks(ctx, bob, bobLockedShares)
	require.Error(t, err)
	bobLockedShares.OwnerAddress = constants.BobAccAddress.String()
	got, exists = k.GetOwnerShareUnlocks(ctx, bob)
	require.Equal(t, true, exists)
	require.Equal(t, bobLockedShares, got)
}

func TestLockShares(t *testing.T) {
	tests := map[string]struct {
		// Existing owner share unlocks.
		existingOwnerShareUnlocks *vaulttypes.OwnerShareUnlocks
		// Owner address.
		ownerAddress string
		// Owner shares.
		ownerShares *big.Int
		// Shares to lock.
		sharesToLock *big.Int
		// Block height to lock until.
		lockUntilBlock uint32
		// Current block height.
		currentBlockHeight uint32
		// Expected owner share unlocks.
		expectedOwnerShareUnlocks vaulttypes.OwnerShareUnlocks
		// Expected error.
		expectedErr string
	}{
		"Success - No existing locked shares and lock 7 shares until height 2": {
			existingOwnerShareUnlocks: nil,
			ownerAddress:              constants.AliceAccAddress.String(),
			ownerShares:               big.NewInt(7),
			sharesToLock:              big.NewInt(7),
			lockUntilBlock:            2,
			currentBlockHeight:        1,
			expectedOwnerShareUnlocks: vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(7)),
						UnlockBlockHeight: 2,
					},
				},
			},
		},
		"Success - 1234 existing locked shares and lock 789 shares until height 456": {
			existingOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.BobAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(1_234)),
						UnlockBlockHeight: 2,
					},
				},
			},
			ownerAddress:       constants.BobAccAddress.String(),
			ownerShares:        big.NewInt(2_078),
			sharesToLock:       big.NewInt(789),
			lockUntilBlock:     456,
			currentBlockHeight: 1,
			expectedOwnerShareUnlocks: vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.BobAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(1_234)),
						UnlockBlockHeight: 2,
					},
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(789)),
						UnlockBlockHeight: 456,
					},
				},
			},
		},
		"Error - Total locked shares would exceed total owner shares": {
			existingOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.CarlAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(17)),
						UnlockBlockHeight: 2,
					},
				},
			},
			ownerAddress:       constants.CarlAccAddress.String(),
			ownerShares:        big.NewInt(65),
			sharesToLock:       big.NewInt(49), // greater than 65-17=48 remaining unlocked shares.
			lockUntilBlock:     3,
			currentBlockHeight: 1,
			expectedErr:        vaulttypes.ErrLockedSharesExceedsOwnerShares.Error(),
		},
		"Error - Empty owner address": {
			existingOwnerShareUnlocks: nil,
			ownerAddress:              "",
			sharesToLock:              big.NewInt(7),
			lockUntilBlock:            2,
			currentBlockHeight:        1,
			expectedErr:               "invalid parameters",
		},
		"Error - 0 shares to lock": {
			existingOwnerShareUnlocks: nil,
			ownerAddress:              constants.AliceAccAddress.String(),
			sharesToLock:              big.NewInt(0),
			lockUntilBlock:            2,
			currentBlockHeight:        1,
			expectedErr:               "invalid parameters",
		},
		"Error - Negative shares to lock": {
			existingOwnerShareUnlocks: nil,
			ownerAddress:              constants.AliceAccAddress.String(),
			sharesToLock:              big.NewInt(-1),
			lockUntilBlock:            2,
			currentBlockHeight:        1,
			expectedErr:               "invalid parameters",
		},
		"Error - Lock until height same as current block height": {
			existingOwnerShareUnlocks: nil,
			ownerAddress:              constants.AliceAccAddress.String(),
			sharesToLock:              big.NewInt(7),
			lockUntilBlock:            14,
			currentBlockHeight:        14,
			expectedErr:               "invalid parameters",
		},
		"Error - Lock until height smaller than current block height": {
			existingOwnerShareUnlocks: nil,
			ownerAddress:              constants.AliceAccAddress.String(),
			sharesToLock:              big.NewInt(7),
			lockUntilBlock:            13,
			currentBlockHeight:        14,
			expectedErr:               "invalid parameters",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			if tc.currentBlockHeight > 1 {
				ctx = tApp.AdvanceToBlock(tc.currentBlockHeight, testapp.AdvanceToBlockOptions{})
			}
			k := tApp.App.VaultKeeper

			if tc.ownerAddress != "" {
				err := k.SetOwnerShares(ctx, tc.ownerAddress, vaulttypes.BigIntToNumShares(tc.ownerShares))
				require.NoError(t, err)
			}
			if tc.existingOwnerShareUnlocks != nil {
				err := k.SetOwnerShareUnlocks(ctx, tc.ownerAddress, *tc.existingOwnerShareUnlocks)
				require.NoError(t, err)
			}

			err := k.LockShares(
				ctx,
				tc.ownerAddress,
				vaulttypes.BigIntToNumShares(tc.sharesToLock),
				tc.lockUntilBlock,
			)
			allDelayedMsgUnlockShares := testutildelaymsg.FilterDelayedMsgsByType[*vaulttypes.MsgUnlockShares](
				t,
				tApp.App.DelayMsgKeeper.GetAllDelayedMessages(ctx),
			)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				l, exists := k.GetOwnerShareUnlocks(ctx, tc.ownerAddress)
				require.Equal(
					t,
					tc.existingOwnerShareUnlocks != nil,
					exists,
				)
				if exists {
					require.Equal(t, *tc.existingOwnerShareUnlocks, l)
				}
				require.Empty(t, allDelayedMsgUnlockShares)
			} else {
				require.NoError(t, err)
				o, exists := k.GetOwnerShareUnlocks(ctx, tc.ownerAddress)
				require.True(t, exists)
				require.Equal(t, tc.expectedOwnerShareUnlocks, o)
				require.Len(t, allDelayedMsgUnlockShares, 1)
				require.Equal(
					t,
					tc.lockUntilBlock,
					allDelayedMsgUnlockShares[0].GetBlockHeight(),
				)
				msg, err := allDelayedMsgUnlockShares[0].GetMessage()
				require.NoError(t, err)
				require.Equal(t, tc.ownerAddress, msg.(*vaulttypes.MsgUnlockShares).OwnerAddress)
			}
		})
	}
}

func TestUnlockShares(t *testing.T) {
	tests := map[string]struct {
		// Existing owner share unlocks.
		existingOwnerShareUnlocks *vaulttypes.OwnerShareUnlocks
		// Owner address.
		ownerAddress string
		// Current block height
		currentBlockHeight uint32
		// Expected unlocked shares.
		expectedUnlockedShares *big.Int
		// Expected owner share unlocks after unlocking.
		expectedOwnerShareUnlocks *vaulttypes.OwnerShareUnlocks
		// Expected error.
		expectedErr string
	}{
		"Success - Unlock all shares of Alice": {
			existingOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(7)),
						UnlockBlockHeight: 11,
					},
				},
			},
			ownerAddress:              constants.AliceAccAddress.String(),
			currentBlockHeight:        11,
			expectedUnlockedShares:    big.NewInt(7),
			expectedOwnerShareUnlocks: nil,
		},
		"Success - Unlock zero shares of Alice": {
			existingOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(7)),
						UnlockBlockHeight: 11,
					},
				},
			},
			ownerAddress:           constants.AliceAccAddress.String(),
			currentBlockHeight:     10,
			expectedUnlockedShares: big.NewInt(0),
			expectedOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(7)),
						UnlockBlockHeight: 11,
					},
				},
			},
		},
		"Success - Unlock some shares of Alice": {
			existingOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(888)),
						UnlockBlockHeight: 14,
					},
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(1_457)),
						UnlockBlockHeight: 18,
					},
				},
			},
			ownerAddress:           constants.AliceAccAddress.String(),
			currentBlockHeight:     17,
			expectedUnlockedShares: big.NewInt(888),
			expectedOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(1_457)),
						UnlockBlockHeight: 18,
					},
				},
			},
		},
		"Success - Unlock all but one share of Bob": {
			existingOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.BobAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(987_000_000)),
						UnlockBlockHeight: 11,
					},
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(654_320)),
						UnlockBlockHeight: 12,
					},
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(1)),
						UnlockBlockHeight: 13,
					},
				},
			},
			ownerAddress:           constants.BobAccAddress.String(),
			currentBlockHeight:     12,
			expectedUnlockedShares: big.NewInt(987_654_320),
			expectedOwnerShareUnlocks: &vaulttypes.OwnerShareUnlocks{
				OwnerAddress: constants.BobAccAddress.String(),
				ShareUnlocks: []vaulttypes.ShareUnlock{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(1)),
						UnlockBlockHeight: 13,
					},
				},
			},
		},
		"Error - Unlock shares of non-existent owner": {
			existingOwnerShareUnlocks: nil,
			ownerAddress:              constants.AliceAccAddress.String(),
			currentBlockHeight:        11,
			expectedUnlockedShares:    big.NewInt(0),
			expectedOwnerShareUnlocks: nil,
			expectedErr:               vaulttypes.ErrOwnerNotFound.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.existingOwnerShareUnlocks != nil {
				err := k.SetOwnerShareUnlocks(ctx, tc.ownerAddress, *tc.existingOwnerShareUnlocks)
				require.NoError(t, err)
			}

			ctx = ctx.WithBlockHeight(int64(tc.currentBlockHeight))

			unlockedShares, err := k.UnlockShares(ctx, tc.ownerAddress)

			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				l, exists := k.GetOwnerShareUnlocks(ctx, tc.ownerAddress)
				require.Equal(
					t,
					tc.existingOwnerShareUnlocks != nil,
					exists,
				)
				if exists {
					require.Equal(t, *tc.existingOwnerShareUnlocks, l)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUnlockedShares, unlockedShares.NumShares.BigInt())
				o, exists := k.GetOwnerShareUnlocks(ctx, tc.ownerAddress)
				require.Equal(
					t,
					tc.expectedOwnerShareUnlocks != nil,
					exists,
				)
				if exists {
					require.Equal(t, *tc.expectedOwnerShareUnlocks, o)
				}
			}
		})
	}
}
