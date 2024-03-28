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
		// All owner shares (in descending number of shares).
		allOwnerShares []vaulttypes.OwnerShares
		// Query request.
		req *vaulttypes.QueryVaultRequest

		/* --- Expectations --- */
		expectedEquity uint64
		expectedErr    string
	}{
		"Success": {
			req: &vaulttypes.QueryVaultRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 0,
			},
			vaultId:     constants.Vault_Clob_0,
			asset:       big.NewInt(100),
			perpId:      0,
			inventory:   big.NewInt(200),
			totalShares: big.NewInt(300),
			allOwnerShares: []vaulttypes.OwnerShares{
				{
					Owner:  constants.Alice_Num0.Owner,
					Shares: 244,
				},
				{
					Owner:  constants.Bob_Num0.Owner,
					Shares: 56,
				},
			},
			expectedEquity: 500,
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
			allOwnerShares: []vaulttypes.OwnerShares{
				{
					Owner:  constants.Alice_Num0.Owner,
					Shares: 244,
				},
				{
					Owner:  constants.Bob_Num0.Owner,
					Shares: 56,
				},
			},
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

			// Set total shares and owner shares.
			err := k.SetTotalShares(ctx, tc.vaultId, vaulttypes.BigIntToNumShares(tc.totalShares))
			require.NoError(t, err)
			for _, ownerShares := range tc.allOwnerShares {
				err := k.SetOwnerShares(
					ctx,
					tc.vaultId,
					ownerShares.Owner,
					vaulttypes.BigIntToNumShares(big.NewInt(int64(ownerShares.Shares))),
				)
				require.NoError(t, err)
			}

			// Check Vault query response is as expected.
			response, err := k.Vault(ctx, tc.req)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				expectedResponse := vaulttypes.QueryVaultResponse{
					VaultId:        tc.vaultId,
					SubaccountId:   *tc.vaultId.ToSubaccountId(),
					Equity:         tc.expectedEquity,
					Inventory:      tc.inventory.Uint64(),
					TotalShares:    tc.totalShares.Uint64(),
					AllOwnerShares: tc.allOwnerShares,
				}
				require.Equal(t, expectedResponse, *response)
			}
		})
	}
}
