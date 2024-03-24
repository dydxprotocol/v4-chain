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
	err := k.SetTotalShares(ctx, constants.Vault_Clob_0, vaulttypes.NumShares{
		NumShares: dtypes.NewInt(7),
	})
	require.NoError(t, err)
	numShares, exists := k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(7), numShares.NumShares)

	// Set total shares for another vault and then get.
	err = k.SetTotalShares(ctx, constants.Vault_Clob_1, vaulttypes.NumShares{
		NumShares: dtypes.NewInt(456),
	})
	require.NoError(t, err)
	numShares, exists = k.GetTotalShares(ctx, constants.Vault_Clob_1)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(456), numShares.NumShares)

	// Set total shares for second vault to 0.
	err = k.SetTotalShares(ctx, constants.Vault_Clob_1, vaulttypes.NumShares{
		NumShares: dtypes.NewInt(0),
	})
	require.NoError(t, err)
	numShares, exists = k.GetTotalShares(ctx, constants.Vault_Clob_1)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(0), numShares.NumShares)

	// Set total shares for the first vault again and then get.
	err = k.SetTotalShares(ctx, constants.Vault_Clob_0, vaulttypes.NumShares{
		NumShares: dtypes.NewInt(123),
	})
	require.NoError(t, err)
	numShares, exists = k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(123), numShares.NumShares)

	// Set total shares for the first vault to a negative value.
	// Should get error and total shares should remain unchanged.
	err = k.SetTotalShares(ctx, constants.Vault_Clob_0, vaulttypes.NumShares{
		NumShares: dtypes.NewInt(-1),
	})
	require.Equal(t, vaulttypes.ErrNegativeShares, err)
	numShares, exists = k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(123), numShares.NumShares)
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
		// Quote quantums to deposit.
		quantumsToDeposit *big.Int

		/* --- Expectations --- */
		// Expected TotalShares after minting.
		expectedTotalShares vaulttypes.NumShares
		// Expected error.
		expectedErr error
	}{
		"Equity 0, TotalShares 0, Deposit 1000": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(0),
			totalShares:       big.NewInt(0),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint 1_000 shares.
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(1_000),
			},
		},
		"Equity 0, TotalShares non-existent, Deposit 12345654321": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(0),
			quantumsToDeposit: big.NewInt(12_345_654_321),
			// Should mint 12_345_654_321 shares.
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(12_345_654_321),
			},
		},
		"Equity 1000, TotalShares non-existent, Deposit 500": {
			vaultId:           constants.Vault_Clob_0,
			equity:            big.NewInt(1_000),
			quantumsToDeposit: big.NewInt(500),
			// Should mint 500 shares.
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(500),
			},
		},
		"Equity 4000, TotalShares 5000, Deposit 1000": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(4_000),
			totalShares:       big.NewInt(5_000),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint 1_250 shares.
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(6_250),
			},
		},
		"Equity 1_000_000, TotalShares 1, Deposit 1": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1_000_000),
			totalShares:       big.NewInt(1),
			quantumsToDeposit: big.NewInt(1),
			// 1 * 1 / 1_000_000 = 1 / 1_000_000
			// Should thus mint 1 share and scale existing shares by 1_000_000.
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(1_000_001),
			},
		},
		"Equity 8000, TotalShares 4000, Deposit 455": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(8_000),
			totalShares:       big.NewInt(4_000),
			quantumsToDeposit: big.NewInt(455),
			// 455 * 4_000 / 8_000 = 455 / 2
			// Should thus mint 455 shares and scale existing shares by 2.
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(8_455),
			},
		},
		"Equity 123456, TotalShares 654321, Deposit 123456789": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(123_456),
			totalShares:       big.NewInt(654_321),
			quantumsToDeposit: big.NewInt(123_456_789),
			// 123_456_789 * 654_321 / 123_456 = 26_926_789_878_423 / 41_152
			// Should thus mint 26_926_789_878_423 shares and scale existing shares by 41_152.
			expectedTotalShares: vaulttypes.NumShares{
				NumShares: dtypes.NewInt(26_953_716_496_215),
			},
		},
		"Equity -1, TotalShares 10, Deposit 1": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(-1),
			totalShares:       big.NewInt(10),
			quantumsToDeposit: big.NewInt(1),
			expectedErr:       vaulttypes.ErrNonPositiveEquity,
		},
		"Equity 1, TotalShares 1, Deposit 0": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(1),
			totalShares:       big.NewInt(1),
			quantumsToDeposit: big.NewInt(0),
			expectedErr:       vaulttypes.ErrInvalidDepositAmount,
		},
		"Equity 0, TotalShares non-existent, Deposit -1": {
			vaultId:           constants.Vault_Clob_1,
			equity:            big.NewInt(0),
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
					vaulttypes.NumShares{
						NumShares: dtypes.NewIntFromBigInt(tc.totalShares),
					},
				)
				require.NoError(t, err)
			}

			// Mint shares.
			err := tApp.App.VaultKeeper.MintShares(
				ctx,
				tc.vaultId,
				"", // TODO (TRA-170): Increase owner shares.
				tc.quantumsToDeposit,
			)
			if tc.expectedErr != nil {
				// Check that error is as expected.
				require.ErrorContains(t, err, tc.expectedErr.Error())
				// Check that TotalShares is unchanged.
				totalShares, _ := tApp.App.VaultKeeper.GetTotalShares(ctx, tc.vaultId)
				require.Equal(t, tc.totalShares, totalShares.NumShares.BigInt())
			} else {
				require.NoError(t, err)
				// Check that TotalShares is as expected.
				totalShares, exists := tApp.App.VaultKeeper.GetTotalShares(ctx, tc.vaultId)
				require.True(t, exists)
				require.Equal(
					t,
					tc.expectedTotalShares,
					totalShares,
				)
			}
		})
	}
}
