package keeper_test

import (
	"math/big"
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetTotalShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Get total shares for a non-existing vault.
	_, exists := k.GetTotalShares(ctx, constants.Vault_Clob0)
	require.Equal(t, false, exists)

	// Set total shares for a vault and then get.
	numShares := vaulttypes.BigIntToNumShares(
		big.NewInt(7),
	)
	err := k.SetTotalShares(ctx, constants.Vault_Clob0, numShares)
	require.NoError(t, err)
	got, exists := k.GetTotalShares(ctx, constants.Vault_Clob0)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for another vault and then get.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(456),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob1, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob1)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for second vault to 0.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(0),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob1, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob1)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for the first vault again and then get.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(7283133),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob0, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob0)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for the first vault to a negative value.
	// Should get error and total shares should remain unchanged.
	negativeShares := vaulttypes.BigIntToNumShares(
		big.NewInt(-1),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob0, negativeShares)
	require.Equal(t, vaulttypes.ErrNegativeShares, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob0)
	require.Equal(t, true, exists)
	require.Equal(
		t,
		numShares,
		got,
	)
}

func TestGetSetOwnerShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	alice := constants.AliceAccAddress.String()
	bob := constants.BobAccAddress.String()

	// Get owners shares for Alice in vault clob 0.
	_, exists := k.GetOwnerShares(ctx, constants.Vault_Clob0, alice)
	require.Equal(t, false, exists)

	// Set owner shares for Alice in vault clob 0 and get.
	numShares := vaulttypes.BigIntToNumShares(
		big.NewInt(7),
	)
	err := k.SetOwnerShares(ctx, constants.Vault_Clob0, alice, numShares)
	require.NoError(t, err)
	got, exists := k.GetOwnerShares(ctx, constants.Vault_Clob0, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Alice in vault clob 1 and then get.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(456),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob1, alice, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob1, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob in vault clob 1.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(0),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob1, bob, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob1, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob in vault clob 1 to a negative value.
	// Should get error and total shares should remain unchanged.
	numSharesInvalid := vaulttypes.BigIntToNumShares(
		big.NewInt(-1),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob1, bob, numSharesInvalid)
	require.ErrorIs(t, err, vaulttypes.ErrNegativeShares)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob1, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)
}

func TestGetAllOwnerShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Get all owner shares of a vault that has no owners.
	allOwnerShares := k.GetAllOwnerShares(ctx, constants.Vault_Clob0)
	require.Equal(t, []*vaulttypes.OwnerShare{}, allOwnerShares)

	// Set alice and bob as owners of a vault and get all owner shares.
	alice := constants.AliceAccAddress.String()
	aliceShares := vaulttypes.BigIntToNumShares(big.NewInt(7))
	bob := constants.BobAccAddress.String()
	bobShares := vaulttypes.BigIntToNumShares(big.NewInt(123))

	err := k.SetOwnerShares(ctx, constants.Vault_Clob0, alice, aliceShares)
	require.NoError(t, err)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob0, bob, bobShares)
	require.NoError(t, err)

	allOwnerShares = k.GetAllOwnerShares(ctx, constants.Vault_Clob0)
	require.ElementsMatch(
		t,
		[]*vaulttypes.OwnerShare{
			{
				Owner:  alice,
				Shares: &aliceShares,
			},
			{
				Owner:  bob,
				Shares: &bobShares,
			},
		},
		allOwnerShares,
	)
}
