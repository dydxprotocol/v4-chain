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

func TestVault(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Vault asset.
		asset *big.Int
		// Perp ID that corresponds to the vault.
		perpId uint32
		// Vault inventory.
		inventory *big.Int
		// Total shares.
		totalShares *big.Int
		// Query request.
		req *vaulttypes.QueryVaultRequest

		/* --- Expectations --- */
		expectedEquity *big.Int
		expectedErr    string
	}{
		"Success": {
			req: &vaulttypes.QueryVaultRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 0,
			},
			vaultId:        constants.Vault_Clob_0,
			asset:          big.NewInt(100),
			perpId:         0,
			inventory:      big.NewInt(200),
			totalShares:    big.NewInt(300),
			expectedEquity: big.NewInt(500),
		},
		"Success: negative inventory and equity": {
			req: &vaulttypes.QueryVaultRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 0,
			},
			vaultId:        constants.Vault_Clob_0,
			asset:          big.NewInt(100),
			perpId:         0,
			inventory:      big.NewInt(-200),
			totalShares:    big.NewInt(300),
			expectedEquity: big.NewInt(-300),
		},
		"Error: query non-existent vault": {
			req: &vaulttypes.QueryVaultRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 1, // Non-existent vault.
			},
			vaultId:     constants.Vault_Clob_0,
			asset:       big.NewInt(100),
			perpId:      0,
			inventory:   big.NewInt(200),
			totalShares: big.NewInt(300),
			expectedErr: "vault not found",
		},
		"Error: nil request": {
			req:         nil,
			vaultId:     constants.Vault_Clob_0,
			asset:       big.NewInt(100),
			perpId:      0,
			inventory:   big.NewInt(200),
			totalShares: big.NewInt(300),
			expectedErr: "invalid request",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: tc.vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assettypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromBigInt(tc.asset),
									},
								},
								PerpetualPositions: []*satypes.PerpetualPosition{
									{
										PerpetualId: tc.perpId,
										Quantums:    dtypes.NewIntFromBigInt(tc.inventory),
									},
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			// Set total shares.
			err := k.SetTotalShares(ctx, tc.vaultId, vaulttypes.BigIntToNumShares(tc.totalShares))
			require.NoError(t, err)

			// Check Vault query response is as expected.
			response, err := k.Vault(ctx, tc.req)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				expectedResponse := vaulttypes.QueryVaultResponse{
					VaultId:      tc.vaultId,
					SubaccountId: *tc.vaultId.ToSubaccountId(),
					Equity:       dtypes.NewIntFromBigInt(tc.expectedEquity),
					Inventory:    dtypes.NewIntFromBigInt(tc.inventory),
					TotalShares:  tc.totalShares.Uint64(),
				}
				require.Equal(t, expectedResponse, *response)
			}
		})
	}
}

func TestAllVaults(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Query request.
		req *vaulttypes.QueryAllVaultsRequest
		// Vault IDs.
		vaultIds []vaulttypes.VaultId
		// Total shares for each vault.
		totalShares map[vaulttypes.VaultId]*big.Int
		// Asset position of each vault.
		assets []*big.Int
		// Inventory of each vault.
		inventories []*big.Int
		// Perpetual ID of each vault.
		perpIds []uint32

		/* --- Expectations --- */
		expectedErr string
	}{
		"Success": {
			req: &vaulttypes.QueryAllVaultsRequest{},
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: map[vaulttypes.VaultId]*big.Int{
				constants.Vault_Clob_0: big.NewInt(100),
				constants.Vault_Clob_1: big.NewInt(200),
			},
			assets: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(2_000),
			},
			inventories: []*big.Int{
				big.NewInt(-555),
				big.NewInt(666),
			},
			perpIds: []uint32{0, 1},
		},
		"Error: nil request": {
			req: nil,
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: map[vaulttypes.VaultId]*big.Int{
				constants.Vault_Clob_0: big.NewInt(100),
				constants.Vault_Clob_1: big.NewInt(200),
			},
			assets: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(2_000),
			},
			inventories: []*big.Int{
				big.NewInt(-555),
				big.NewInt(666),
			},
			perpIds:     []uint32{0, 1},
			expectedErr: "invalid request",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						var subaccounts []satypes.Subaccount
						for i, vaultId := range tc.vaultIds {
							subaccounts = append(subaccounts, satypes.Subaccount{
								Id: vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assettypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromBigInt(tc.assets[i]),
									},
								},
								PerpetualPositions: []*satypes.PerpetualPosition{
									{
										PerpetualId: tc.perpIds[i],
										Quantums:    dtypes.NewIntFromBigInt(tc.inventories[i]),
									},
								},
							})
						}
						genesisState.Subaccounts = subaccounts
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			// Set total shares.
			for _, vaultId := range tc.vaultIds {
				err := k.SetTotalShares(ctx, vaultId, vaulttypes.BigIntToNumShares(tc.totalShares[vaultId]))
				require.NoError(t, err)
			}

			// Check AllVaults query response is as expected.
			response, err := k.AllVaults(ctx, tc.req)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				expectedVaults := make([]*vaulttypes.QueryVaultResponse, len(tc.vaultIds))
				for i, vaultId := range tc.vaultIds {
					vault, err := k.Vault(ctx, &vaulttypes.QueryVaultRequest{
						Type:   vaultId.Type,
						Number: vaultId.Number,
					})
					require.NoError(t, err)
					expectedVaults[i] = vault
				}
				require.ElementsMatch(
					t,
					expectedVaults,
					response.Vaults,
				)
			}
		})
	}
}
