package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetTotalShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Get total shares for a non-existing vault.
	_, exists := k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, false, exists)

	// Set total shares for a vault and then get.
	numShares := vaulttypes.BigIntToNumShares(
		big.NewInt(7),
	)
	err := k.SetTotalShares(ctx, constants.Vault_Clob_0, numShares)
	require.NoError(t, err)
	got, exists := k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for another vault and then get.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(456),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_1, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_1)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for second vault to 0.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(0),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_1, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_1)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for the first vault again and then get.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(7283133),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_0, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for the first vault to a negative value.
	// Should get error and total shares should remain unchanged.
	negativeShares := vaulttypes.BigIntToNumShares(
		big.NewInt(-1),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_0, negativeShares)
	require.Equal(t, vaulttypes.ErrNegativeShares, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_0)
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
	_, exists := k.GetOwnerShares(ctx, constants.Vault_Clob_0, alice)
	require.Equal(t, false, exists)

	// Set owner shares for Alice in vault clob 0 and get.
	numShares := vaulttypes.BigIntToNumShares(
		big.NewInt(7),
	)
	err := k.SetOwnerShares(ctx, constants.Vault_Clob_0, alice, numShares)
	require.NoError(t, err)
	got, exists := k.GetOwnerShares(ctx, constants.Vault_Clob_0, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Alice in vault clob 1 and then get.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(456),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob_1, alice, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_1, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob in vault clob 1.
	numShares = vaulttypes.BigIntToNumShares(
		big.NewInt(0),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob_1, bob, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_1, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob in vault clob 1 to a negative value.
	// Should get error and total shares should remain unchanged.
	numSharesInvalid := vaulttypes.BigIntToNumShares(
		big.NewInt(-1),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob_1, bob, numSharesInvalid)
	require.ErrorIs(t, err, vaulttypes.ErrNegativeShares)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_1, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)
}

func TestGetAllOwnerShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Get all owner shares of a vault that has no owners.
	allOwnerShares := k.GetAllOwnerShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, []*vaulttypes.OwnerShare{}, allOwnerShares)

	// Set alice and bob as owners of a vault and get all owner shares.
	alice := constants.AliceAccAddress.String()
	aliceShares := vaulttypes.BigIntToNumShares(big.NewInt(7))
	bob := constants.BobAccAddress.String()
	bobShares := vaulttypes.BigIntToNumShares(big.NewInt(123))

	k.SetOwnerShares(ctx, constants.Vault_Clob_0, alice, aliceShares)
	k.SetOwnerShares(ctx, constants.Vault_Clob_0, bob, bobShares)

	allOwnerShares = k.GetAllOwnerShares(ctx, constants.Vault_Clob_0)
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

func TestMintShares(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Existing vault equity.
		equity *big.Int
		// Existing vault TotalShares.
		totalShares *big.Int
		// Owner that deposits.
		owner string
		// Existing owner shares.
		ownerShares *big.Int
		// Quote quantums to deposit.
		quantumsToDeposit *big.Int

		/* --- Expectations --- */
		// Expected TotalShares after minting.
		expectedTotalShares *big.Int
		// Expected OwnerShares after minting.
		expectedOwnerShares *big.Int
		// Expected error.
		expectedErr error
	}{
		"Equity 0, TotalShares 0, OwnerShares 0, Deposit 1000": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(0),
			totalShares:       big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			ownerShares:       big.NewInt(0),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `1_000` shares.
			expectedTotalShares: big.NewInt(1_000),
			expectedOwnerShares: big.NewInt(1_000),
		},
		"Equity 0, TotalShares non-existent, OwnerShares non-existent, Deposit 12345654321": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(12_345_654_321),
			// Should mint `12_345_654_321` shares.
			expectedTotalShares: big.NewInt(12_345_654_321),
			expectedOwnerShares: big.NewInt(12_345_654_321),
		},
		"Equity 1000, TotalShares non-existent, OwnerShares non-existent, Deposit 500": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(1_000),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(500),
			// Should mint `500` shares.
			expectedTotalShares: big.NewInt(500),
			expectedOwnerShares: big.NewInt(500),
		},
		"Equity 4000, TotalShares 5000, OwnerShares 2500, Deposit 1000": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(4_000),
			totalShares:       big.NewInt(5_000),
			owner:             constants.AliceAccAddress.String(),
			ownerShares:       big.NewInt(2_500),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `1_250` shares.
			expectedTotalShares: big.NewInt(6_250),
			expectedOwnerShares: big.NewInt(3_750),
		},
		"Equity 1_000_000, TotalShares 2_000, OwnerShares 1, Deposit 1_000": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1_000_000),
			totalShares:       big.NewInt(2_000),
			owner:             constants.BobAccAddress.String(),
			ownerShares:       big.NewInt(1),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `2` shares.
			expectedTotalShares: big.NewInt(2_002),
			expectedOwnerShares: big.NewInt(3),
		},
		"Equity 8000, TotalShares 4000, OwnerShares 101, Deposit 455": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(8_000),
			totalShares:       big.NewInt(4_000),
			owner:             constants.CarlAccAddress.String(),
			ownerShares:       big.NewInt(101),
			quantumsToDeposit: big.NewInt(455),
			// Should mint `227.5` shares, round down to 227.
			expectedTotalShares: big.NewInt(4_227),
			expectedOwnerShares: big.NewInt(328),
		},
		"Equity 123456, TotalShares 654321, OwnerShares 0, Deposit 123456789": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(123_456),
			totalShares:       big.NewInt(654_321),
			owner:             constants.DaveAccAddress.String(),
			quantumsToDeposit: big.NewInt(123_456_789),
			// Should mint `654_325_181.727` shares, round down to 654_325_181.
			expectedTotalShares: big.NewInt(654_979_502),
			expectedOwnerShares: big.NewInt(654_325_181),
		},
		"Equity 1000000, TotalShares 1000, OwnerShares 0, Deposit 9_900": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1_000_000),
			totalShares:       big.NewInt(1_000),
			owner:             constants.DaveAccAddress.String(),
			quantumsToDeposit: big.NewInt(9_900),
			// Should mint `9_900 * 1_000 / 1_000_000` shares, round down to 9.
			expectedTotalShares: big.NewInt(1_009),
			expectedOwnerShares: big.NewInt(9),
		},
		"Equity -1, TotalShares 10, Deposit 1": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(-1),
			totalShares:       big.NewInt(10),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(1),
			expectedErr:       vaulttypes.ErrNonPositiveEquity,
		},
		"Equity 1, TotalShares 1, Deposit 0": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1),
			totalShares:       big.NewInt(1),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(0),
			expectedErr:       vaulttypes.ErrInvalidDepositAmount,
		},
		"Equity 0, TotalShares non-existent, Deposit -1": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(-1),
			expectedErr:       vaulttypes.ErrInvalidDepositAmount,
		},
		"Equity 1000, TotalShares 1, Deposit 100": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1_000),
			totalShares:       big.NewInt(1),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(100),
			expectedErr:       vaulttypes.ErrZeroSharesToMint,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize vault with its existing equity.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: tc.vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  0,
										Quantums: dtypes.NewIntFromBigInt(tc.equity),
									},
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Set vault's existing total shares if specified.
			if tc.totalShares != nil {
				err := tApp.App.VaultKeeper.SetTotalShares(
					ctx,
					tc.vaultId,
					vaulttypes.BigIntToNumShares(tc.totalShares),
				)
				require.NoError(t, err)
			}
			// Set vault's existing owner shares if specified.
			if tc.ownerShares != nil {
				err := tApp.App.VaultKeeper.SetOwnerShares(
					ctx,
					tc.vaultId,
					tc.owner,
					vaulttypes.BigIntToNumShares(tc.ownerShares),
				)
				require.NoError(t, err)
			}

			// Mint shares.
			err := tApp.App.VaultKeeper.MintShares(
				ctx,
				tc.vaultId,
				tc.owner,
				tc.quantumsToDeposit,
			)
			if tc.expectedErr != nil {
				// Check that error is as expected.
				require.ErrorContains(t, err, tc.expectedErr.Error())
				// Check that TotalShares is unchanged.
				totalShares, _ := tApp.App.VaultKeeper.GetTotalShares(ctx, tc.vaultId)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(tc.totalShares),
					totalShares,
				)
				// Check that OwnerShares is unchanged.
				ownerShares, _ := tApp.App.VaultKeeper.GetOwnerShares(ctx, tc.vaultId, tc.owner)
				require.Equal(t, vaulttypes.BigIntToNumShares(tc.ownerShares), ownerShares)
			} else {
				require.NoError(t, err)
				// Check that TotalShares is as expected.
				totalShares, exists := tApp.App.VaultKeeper.GetTotalShares(ctx, tc.vaultId)
				require.True(t, exists)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(tc.expectedTotalShares),
					totalShares,
				)
				// Check that OwnerShares is as expected.
				ownerShares, exists := tApp.App.VaultKeeper.GetOwnerShares(ctx, tc.vaultId, tc.owner)
				require.True(t, exists)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(tc.expectedOwnerShares),
					ownerShares,
				)
			}
		})
	}
}
