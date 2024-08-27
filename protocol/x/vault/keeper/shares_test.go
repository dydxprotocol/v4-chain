package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
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

func TestGetSetLockedShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	alice := constants.AliceAccAddress.String()
	bob := constants.BobAccAddress.String()

	// Get locked shares for Alice.
	_, exists := k.GetLockedShares(ctx, alice)
	require.Equal(t, false, exists)

	// Set locked shares for Alice and get.
	aliceLockedShares := vaulttypes.LockedShares{
		OwnerAddress:      alice,
		TotalLockedShares: vaulttypes.BigIntToNumShares(big.NewInt(7)),
		UnlockDetails: []vaulttypes.UnlockDetail{
			{
				Shares:            vaulttypes.BigIntToNumShares(big.NewInt(7)),
				UnlockBlockHeight: 1,
			},
		},
	}
	err := k.SetLockedShares(ctx, alice, aliceLockedShares)
	require.NoError(t, err)
	got, exists := k.GetLockedShares(ctx, alice)
	require.Equal(t, true, exists)
	require.Equal(t, aliceLockedShares, got)

	// Set locked shares for Bob and then get.
	bobLockedShares := vaulttypes.LockedShares{
		OwnerAddress:      bob,
		TotalLockedShares: vaulttypes.BigIntToNumShares(big.NewInt(1_234)),
		UnlockDetails: []vaulttypes.UnlockDetail{
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
	err = k.SetLockedShares(ctx, bob, bobLockedShares)
	require.NoError(t, err)
	got, exists = k.GetLockedShares(ctx, bob)
	require.Equal(t, true, exists)
	require.Equal(t, bobLockedShares, got)

	// Set invalid locked shares for Bob.
	// Should get error and locked shares should remain unchanged.
	bobLockedShares.TotalLockedShares = vaulttypes.BigIntToNumShares(big.NewInt(1_235))
	err = k.SetLockedShares(ctx, bob, bobLockedShares)
	require.Error(t, err)
	bobLockedShares.TotalLockedShares = vaulttypes.BigIntToNumShares(big.NewInt(1_234))
	got, exists = k.GetLockedShares(ctx, bob)
	require.Equal(t, true, exists)
	require.Equal(t, bobLockedShares, got)
}

func TestLockShares(t *testing.T) {
	tests := map[string]struct {
		// Existing locked shares.
		existingLockedShares *vaulttypes.LockedShares
		// Owner address.
		ownerAddress string
		// Shares to lock.
		sharesToLock *big.Int
		// Block height to lock until.
		lockUntilBlock uint32
		// Expected locked shares.
		expectedLockedShares vaulttypes.LockedShares
		// Expected error.
		expectedErr string
	}{
		"Success - No existing locked shares and lock 7 shares until height 2": {
			existingLockedShares: nil,
			ownerAddress:         constants.AliceAccAddress.String(),
			sharesToLock:         big.NewInt(7),
			lockUntilBlock:       2,
			expectedLockedShares: vaulttypes.LockedShares{
				OwnerAddress:      constants.AliceAccAddress.String(),
				TotalLockedShares: vaulttypes.BigIntToNumShares(big.NewInt(7)),
				UnlockDetails: []vaulttypes.UnlockDetail{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(7)),
						UnlockBlockHeight: 2,
					},
				},
			},
		},
		"Success - 1234 existing locked shares and lock 789 shares until height 456": {
			existingLockedShares: &vaulttypes.LockedShares{
				OwnerAddress:      constants.BobAccAddress.String(),
				TotalLockedShares: vaulttypes.BigIntToNumShares(big.NewInt(1_234)),
				UnlockDetails: []vaulttypes.UnlockDetail{
					{
						Shares:            vaulttypes.BigIntToNumShares(big.NewInt(1_234)),
						UnlockBlockHeight: 2,
					},
				},
			},
			ownerAddress:   constants.BobAccAddress.String(),
			sharesToLock:   big.NewInt(789),
			lockUntilBlock: 456,
			expectedLockedShares: vaulttypes.LockedShares{
				OwnerAddress:      constants.BobAccAddress.String(),
				TotalLockedShares: vaulttypes.BigIntToNumShares(big.NewInt(2_023)),
				UnlockDetails: []vaulttypes.UnlockDetail{
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
		"Error - Empty owner address": {
			existingLockedShares: nil,
			ownerAddress:         "",
			sharesToLock:         big.NewInt(7),
			lockUntilBlock:       2,
			expectedErr:          "invalid parameters",
		},
		"Error - 0 shares to lock": {
			existingLockedShares: nil,
			ownerAddress:         constants.AliceAccAddress.String(),
			sharesToLock:         big.NewInt(0),
			lockUntilBlock:       2,
			expectedErr:          "invalid parameters",
		},
		"Error - negative shares to lock": {
			existingLockedShares: nil,
			ownerAddress:         constants.AliceAccAddress.String(),
			sharesToLock:         big.NewInt(-1),
			lockUntilBlock:       2,
			expectedErr:          "invalid parameters",
		},
		"Error - lock until height same as current block height": {
			existingLockedShares: nil,
			ownerAddress:         constants.AliceAccAddress.String(),
			sharesToLock:         big.NewInt(7),
			lockUntilBlock:       1,
			expectedErr:          "invalid parameters",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.existingLockedShares != nil {
				err := k.SetLockedShares(ctx, tc.ownerAddress, *tc.existingLockedShares)
				require.NoError(t, err)
			}

			err := k.LockShares(
				ctx,
				tc.ownerAddress,
				vaulttypes.BigIntToNumShares(tc.sharesToLock),
				tc.lockUntilBlock,
			)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				l, exists := k.GetLockedShares(ctx, tc.ownerAddress)
				require.Equal(
					t,
					tc.existingLockedShares != nil,
					exists,
				)
				if exists {
					require.Equal(t, tc.existingLockedShares, l)
				}
			} else {
				require.NoError(t, err)
				l, exists := k.GetLockedShares(ctx, tc.ownerAddress)
				require.True(t, exists)
				require.Equal(t, tc.expectedLockedShares, l)
			}
		})
	}
}
