package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestDecommissionNonPositiveEquityVaults(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault IDs.
		vaultIds []vaulttypes.VaultId
		// Total shares of above vaults.
		totalShares []*big.Int
		// Equities of above vaults.
		equities []*big.Int

		/* --- Expectations --- */
		// Whether the vaults are decommissioned.
		decommissioned []bool
	}{
		"Decommission no vault": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(7),
				big.NewInt(7),
			},
			equities: []*big.Int{
				big.NewInt(1),
				big.NewInt(1),
			},
			decommissioned: []bool{
				false,
				false,
			},
		},
		"Decommission one vault": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(7),
				big.NewInt(7),
			},
			equities: []*big.Int{
				big.NewInt(1),
				big.NewInt(0),
			},
			decommissioned: []bool{
				false,
				true, // this vault should be decommissioned.
			},
		},
		"Decommission two vaults": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(7),
				big.NewInt(7),
			},
			equities: []*big.Int{
				big.NewInt(0),
				big.NewInt(-1),
			},
			decommissioned: []bool{
				true,
				true,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize vaults with their equities.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						subaccounts := []satypes.Subaccount{}
						for i, vaultId := range tc.vaultIds {
							if tc.equities[i].Sign() != 0 {
								subaccounts = append(
									subaccounts,
									satypes.Subaccount{
										Id: vaultId.ToSubaccountId(),
										AssetPositions: []*satypes.AssetPosition{
											{
												AssetId:  assettypes.AssetUsdc.Id,
												Quantums: dtypes.NewIntFromBigInt(tc.equities[i]),
											},
										},
									},
								)
							}
						}
						genesisState.Subaccounts = subaccounts
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			// Set total shares and owner shares for all vaults.
			testOwner := constants.Alice_Num0.Owner
			for i, vaultId := range tc.vaultIds {
				err := k.SetTotalShares(
					ctx,
					vaultId,
					vaulttypes.BigIntToNumShares(tc.totalShares[i]),
				)
				require.NoError(t, err)
				err = k.SetOwnerShares(
					ctx,
					vaultId,
					testOwner,
					vaulttypes.BigIntToNumShares(big.NewInt(7)),
				)
				require.NoError(t, err)
			}

			// Decommission all vaults.
			k.DecommissionNonPositiveEquityVaults(ctx)

			// Check that total shares and owner shares are deleted for decommissioned
			// vaults and not deleted for non-decommissioned vaults.
			for i, decommissioned := range tc.decommissioned {
				_, exists := k.GetTotalShares(ctx, tc.vaultIds[i])
				require.Equal(t, !decommissioned, exists)
				_, exists = k.GetOwnerShares(ctx, tc.vaultIds[i], testOwner)
				require.Equal(t, !decommissioned, exists)
			}
		})
	}
}

func TestDecommissionVault(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Whether total shares exists.
		totalSharesExists bool
		// Owners.
		owners []string
	}{
		"Total shares doesn't exist, no owners": {
			vaultId: constants.Vault_Clob_0,
		},
		"Total shares exists, no owners": {
			vaultId:           constants.Vault_Clob_0,
			totalSharesExists: true,
		},
		"Total shares exists, one owner": {
			vaultId:           constants.Vault_Clob_1,
			totalSharesExists: true,
			owners:            []string{constants.Alice_Num0.Owner},
		},
		"Total shares exists, two owners": {
			vaultId:           constants.Vault_Clob_1,
			totalSharesExists: true,
			owners: []string{
				constants.Alice_Num0.Owner,
				constants.Bob_Num0.Owner,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			shares := vaulttypes.BigIntToNumShares(
				big.NewInt(7),
			)

			if tc.totalSharesExists {
				err := k.SetTotalShares(
					ctx,
					tc.vaultId,
					shares,
				)
				require.NoError(t, err)
			}
			for _, owner := range tc.owners {
				err := k.SetOwnerShares(
					ctx,
					tc.vaultId,
					owner,
					shares,
				)
				require.NoError(t, err)
			}

			// Decommission vault.
			k.DecommissionVault(ctx, tc.vaultId)

			// Check that total shares and owner shares are deleted.
			_, exists := k.GetTotalShares(ctx, tc.vaultId)
			require.Equal(t, false, exists)
			for _, owner := range tc.owners {
				_, exists = k.GetOwnerShares(ctx, tc.vaultId, owner)
				require.Equal(t, false, exists)
			}
		})
	}
}
