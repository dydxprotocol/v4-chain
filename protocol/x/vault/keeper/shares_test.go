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
