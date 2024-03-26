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
	numShares := vaulttypes.BigRatToNumShares(
		big.NewRat(7, 1),
	)
	err := k.SetTotalShares(ctx, constants.Vault_Clob_0, numShares)
	require.NoError(t, err)
	got, exists := k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for another vault and then get.
	numShares = vaulttypes.BigRatToNumShares(
		big.NewRat(456, 3),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_1, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_1)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for second vault to 0.
	numShares = vaulttypes.BigRatToNumShares(
		big.NewRat(0, 1),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_1, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_1)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for the first vault again and then get.
	numShares = vaulttypes.BigRatToNumShares(
		big.NewRat(73423, 59),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_0, numShares)
	require.NoError(t, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set total shares for the first vault to a negative value.
	// Should get error and total shares should remain unchanged.
	numShares = vaulttypes.BigRatToNumShares(
		big.NewRat(-1, 1),
	)
	err = k.SetTotalShares(ctx, constants.Vault_Clob_0, numShares)
	require.Equal(t, vaulttypes.ErrNegativeShares, err)
	got, exists = k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(
		t,
		vaulttypes.BigRatToNumShares(
			big.NewRat(73423, 59),
		),
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
	numShares := vaulttypes.BigRatToNumShares(
		big.NewRat(7, 1),
	)
	err := k.SetOwnerShares(ctx, constants.Vault_Clob_0, alice, numShares)
	require.NoError(t, err)
	got, exists := k.GetOwnerShares(ctx, constants.Vault_Clob_0, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Alice in vault clob 1 and then get.
	numShares = vaulttypes.BigRatToNumShares(
		big.NewRat(456, 3),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob_1, alice, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_1, alice)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob in vault clob 1.
	numShares = vaulttypes.BigRatToNumShares(
		big.NewRat(0, 1),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob_1, bob, numShares)
	require.NoError(t, err)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_1, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)

	// Set owner shares for Bob in vault clob 1 to a negative value.
	// Should get error and total shares should remain unchanged.
	numSharesInvalid := vaulttypes.BigRatToNumShares(
		big.NewRat(-1, 1),
	)
	err = k.SetOwnerShares(ctx, constants.Vault_Clob_1, bob, numSharesInvalid)
	require.ErrorIs(t, err, vaulttypes.ErrNegativeShares)
	got, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_1, bob)
	require.Equal(t, true, exists)
	require.Equal(t, numShares, got)
}

func TestMintShares(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Existing vault equity.
		equity *big.Int
		// Existing vault TotalShares.
		totalShares *big.Rat
		// Owner that deposits.
		owner string
		// Existing owner shares.
		ownerShares *big.Rat
		// Quote quantums to deposit.
		quantumsToDeposit *big.Int

		/* --- Expectations --- */
		// Expected TotalShares after minting.
		expectedTotalShares *big.Rat
		// Expected OwnerShares after minting.
		expectedOwnerShares *big.Rat
		// Expected error.
		expectedErr error
	}{
		"Equity 0, TotalShares 0, OwnerShares 0, Deposit 1000": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(0),
			totalShares:       big.NewRat(0, 1),
			owner:             constants.AliceAccAddress.String(),
			ownerShares:       big.NewRat(0, 1),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `1_000 / 1` shares.
			expectedTotalShares: big.NewRat(1_000, 1),
			expectedOwnerShares: big.NewRat(1_000, 1),
		},
		"Equity 0, TotalShares non-existent, OwnerShares non-existent, Deposit 12345654321": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(12_345_654_321),
			// Should mint `12_345_654_321 / 1` shares.
			expectedTotalShares: big.NewRat(12_345_654_321, 1),
			expectedOwnerShares: big.NewRat(12_345_654_321, 1),
		},
		"Equity 1000, TotalShares non-existent, OwnerShares non-existent, Deposit 500": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(1_000),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(500),
			// Should mint `500 / 1` shares.
			expectedTotalShares: big.NewRat(500, 1),
			expectedOwnerShares: big.NewRat(500, 1),
		},
		"Equity 4000, TotalShares 5000, OwnerShares 2500, Deposit 1000": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(4_000),
			totalShares:       big.NewRat(5_000, 1),
			owner:             constants.AliceAccAddress.String(),
			ownerShares:       big.NewRat(2_500, 1),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `1_250 / 1` shares.
			expectedTotalShares: big.NewRat(6_250, 1),
			expectedOwnerShares: big.NewRat(3_750, 1),
		},
		"Equity 1_000_000, TotalShares 1, OwnerShares 1/2, Deposit 1": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1_000_000),
			totalShares:       big.NewRat(1, 1),
			owner:             constants.BobAccAddress.String(),
			ownerShares:       big.NewRat(1, 2),
			quantumsToDeposit: big.NewInt(1),
			// Should mint `1 / 1_000_000` shares.
			expectedTotalShares: big.NewRat(1_000_001, 1_000_000),
			expectedOwnerShares: big.NewRat(500_001, 1_000_000),
		},
		"Equity 8000, TotalShares 4000, OwnerShares  Deposit 455": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(8_000),
			totalShares:       big.NewRat(4_000, 1),
			owner:             constants.CarlAccAddress.String(),
			ownerShares:       big.NewRat(101, 7),
			quantumsToDeposit: big.NewInt(455),
			// Should mint `455 / 2` shares.
			expectedTotalShares: big.NewRat(8_455, 2),
			expectedOwnerShares: big.NewRat(3_387, 14),
		},
		"Equity 123456, TotalShares 654321, OwnerShares 0, Deposit 123456789": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(123_456),
			totalShares:       big.NewRat(654_321, 1),
			owner:             constants.DaveAccAddress.String(),
			quantumsToDeposit: big.NewInt(123_456_789),
			// Should mint `26_926_789_878_423 / 41_152` shares.
			expectedTotalShares: big.NewRat(26_953_716_496_215, 41_152),
			expectedOwnerShares: big.NewRat(26_926_789_878_423, 41_152),
		},
		"Equity -1, TotalShares 10, Deposit 1": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(-1),
			totalShares:       big.NewRat(10, 1),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(1),
			expectedErr:       vaulttypes.ErrNonPositiveEquity,
		},
		"Equity 1, TotalShares 1, Deposit 0": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1),
			totalShares:       big.NewRat(1, 1),
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
					vaulttypes.BigRatToNumShares(tc.totalShares),
				)
				require.NoError(t, err)
			}
			if tc.ownerShares != nil {
				err := tApp.App.VaultKeeper.SetOwnerShares(
					ctx,
					tc.vaultId,
					tc.owner,
					vaulttypes.BigRatToNumShares(tc.ownerShares),
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
					vaulttypes.BigRatToNumShares(tc.totalShares),
					totalShares,
				)
				// Check that OwnerShares is unchanged.
				ownerShares, _ := tApp.App.VaultKeeper.GetOwnerShares(ctx, tc.vaultId, tc.owner)
				require.Equal(t, vaulttypes.BigRatToNumShares(tc.ownerShares), ownerShares)
			} else {
				require.NoError(t, err)
				// Check that TotalShares is as expected.
				totalShares, exists := tApp.App.VaultKeeper.GetTotalShares(ctx, tc.vaultId)
				require.True(t, exists)
				require.Equal(
					t,
					vaulttypes.BigRatToNumShares(tc.expectedTotalShares),
					totalShares,
				)
				// Check that OwnerShares is as expected.
				ownerShares, exists := tApp.App.VaultKeeper.GetOwnerShares(ctx, tc.vaultId, tc.owner)
				require.True(t, exists)
				require.Equal(
					t,
					vaulttypes.BigRatToNumShares(tc.expectedOwnerShares),
					ownerShares,
				)
			}
		})
	}
}
