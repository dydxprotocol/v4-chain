package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestDecommissionNonPositiveEquityVaults(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault IDs.
		vaultIds []vaulttypes.VaultId
		// Statuses of above vaults.
		statuses []vaulttypes.VaultStatus
		// Equities of above vaults.
		equities []*big.Int

		/* --- Expectations --- */
		// Whether the vaults are decommissioned.
		decommissioned []bool
	}{
		"Decommission no vault": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
			statuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
				vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
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
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
			statuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
				vaulttypes.VaultStatus_VAULT_STATUS_DEACTIVATED,
			},
			equities: []*big.Int{
				big.NewInt(1),
				big.NewInt(0),
			},
			decommissioned: []bool{
				false,
				true, // decommissioned as vault has 0 equity and is deactivated.
			},
		},
		"Decommission two vaults": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
				constants.Vault_Clob7,
			},
			statuses: []vaulttypes.VaultStatus{
				vaulttypes.VaultStatus_VAULT_STATUS_DEACTIVATED,
				vaulttypes.VaultStatus_VAULT_STATUS_DEACTIVATED,
				vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			},
			equities: []*big.Int{
				big.NewInt(0),
				big.NewInt(-1),
				big.NewInt(-1),
			},
			decommissioned: []bool{
				true,
				true,
				false, // not decommissioned (even though equity is negative) bc status is not deactivated.
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
											testutil.CreateSingleAssetPosition(
												assettypes.AssetUsdc.Id,
												tc.equities[i],
											),
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

			// Set statuses of vaults and add to vault address store.
			for i, vaultId := range tc.vaultIds {
				err := k.SetVaultParams(
					ctx,
					vaultId,
					vaulttypes.VaultParams{
						Status: tc.statuses[i],
					},
				)
				require.NoError(t, err)
				k.AddVaultToAddressStore(ctx, vaultId)
			}

			// Decommission all vaults.
			k.DecommissionNonPositiveEquityVaults(ctx)

			// Check that below are deleted for decommissioned vaults only:
			// - vault params
			// - vault address (from vault address store)
			// - most recent client IDs
			for i, decommissioned := range tc.decommissioned {
				_, exists := k.GetVaultParams(ctx, tc.vaultIds[i])
				require.Equal(t, !decommissioned, exists)
				require.Equal(
					t,
					!decommissioned,
					k.IsVault(ctx, tc.vaultIds[i].ToModuleAccountAddress()),
				)
				require.Empty(
					t,
					k.GetMostRecentClientIds(ctx, tc.vaultIds[i]),
				)
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
			vaultId: constants.Vault_Clob0,
		},
		"Total shares exists, no owners": {
			vaultId:           constants.Vault_Clob0,
			totalSharesExists: true,
		},
		"Total shares exists, one owner": {
			vaultId:           constants.Vault_Clob1,
			totalSharesExists: true,
			owners:            []string{constants.Alice_Num0.Owner},
		},
		"Total shares exists, two owners": {
			vaultId:           constants.Vault_Clob1,
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

			k.AddVaultToAddressStore(ctx, tc.vaultId)
			err := k.SetVaultParams(ctx, tc.vaultId, constants.VaultParams)
			require.NoError(t, err)

			// Decommission vault.
			k.DecommissionVault(ctx, tc.vaultId)

			// Check that vault address, vault params, and most recent client IDs are deleted.
			require.False(t, k.IsVault(ctx, tc.vaultId.ToModuleAccountAddress()))
			_, exists := k.GetVaultParams(ctx, tc.vaultId)
			require.False(t, exists)
			require.Empty(t, k.GetMostRecentClientIds(ctx, tc.vaultId))
		})
	}
}

func TestAddVaultToAddressStore(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultIds []vaulttypes.VaultId
	}{
		"Add 1 vault to vault address store": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
			},
		},
		"Add 2 vaults to vault address store": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob0,
				constants.Vault_Clob1,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			// Add vaults to vault address store.
			for _, vaultId := range tc.vaultIds {
				k.AddVaultToAddressStore(ctx, vaultId)
			}

			// Verify that vault is added to address store.
			for _, vaultId := range tc.vaultIds {
				require.True(t, k.IsVault(ctx, vaultId.ToModuleAccountAddress()))
			}
		})
	}
}

func TestVaultIsBestFeeTier(t *testing.T) {
	// Initialize genesis with a positive-equity vault and fee tiers.
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *vaulttypes.GenesisState) {
				genesisState.Vaults = []vaulttypes.Vault{
					{
						VaultId: constants.Vault_Clob0,
						VaultParams: vaulttypes.VaultParams{
							Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
						},
					},
				}
				genesisState.OperatorParams = vaulttypes.OperatorParams{
					Operator: constants.AliceAccAddress.String(),
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					{
						Id: constants.Vault_Clob0.ToSubaccountId(),
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  assettypes.AssetUsdc.Id,
								Quantums: dtypes.NewInt(1),
							},
						},
					},
					{
						Id: &vaulttypes.MegavaultMainSubaccount,
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  assettypes.AssetUsdc.Id,
								Quantums: dtypes.NewInt(1),
							},
						},
					},
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *feetierstypes.GenesisState) {
				genesisState.Params = feetierstypes.PerpetualFeeParams{
					Tiers: []*feetierstypes.PerpetualFeeTier{
						{
							Name:        "1",
							TakerFeePpm: 33,
							MakerFeePpm: 3,
						},
						{
							Name:                      "2",
							AbsoluteVolumeRequirement: 1_000,
							TakerFeePpm:               22,
							MakerFeePpm:               2,
						},
						{
							Name:                           "3",
							AbsoluteVolumeRequirement:      1_000_000_000,
							MakerVolumeShareRequirementPpm: 500_000,
							TakerFeePpm:                    11,
							MakerFeePpm:                    1,
						},
					},
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	vaultClob0Address := constants.Vault_Clob0.ToModuleAccountAddress()
	vaultClob1Address := constants.Vault_Clob1.ToModuleAccountAddress()
	aliceAddress := constants.AliceAccAddress.String()

	require.NoError(t, tApp.App.AffiliatesKeeper.UpdateAffiliateParameters(
		ctx,
		&affiliatetypes.MsgUpdateAffiliateParameters{
			AffiliateParameters: affiliatetypes.AffiliateParameters{
				RefereeMinimumFeeTierIdx: 2,
			},
		},
	))

	// Vault in genesis state should be in best fee tier.
	takerFee := tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob0Address, true, 2, uint32(1))
	require.Equal(t, int32(11), takerFee)
	makerFee := tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob0Address, false, 2, uint32(1))
	require.Equal(t, int32(1), makerFee)

	// A regular user Alice should be in worst fee tier.
	takerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, aliceAddress, true, 2, uint32(1))
	require.Equal(t, int32(33), takerFee)
	makerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, aliceAddress, false, 2, uint32(1))
	require.Equal(t, int32(3), makerFee)

	// A newly allocated-to vault should be in best fee tier.
	checkTx_AllocateToVault := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.AliceAccAddress.String(),
			Gas:                  constants.TestGasLimit,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		&vaulttypes.MsgAllocateToVault{
			Authority:     constants.AliceAccAddress.String(),
			VaultId:       constants.Vault_Clob1,
			QuoteQuantums: dtypes.NewInt(1),
		},
	)
	checkTxResp := tApp.CheckTx(checkTx_AllocateToVault)
	require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

	ctx = tApp.AdvanceToBlock(
		uint32(ctx.BlockHeight())+1,
		testapp.AdvanceToBlockOptions{},
	)
	takerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob1Address, true, 2, uint32(1))
	require.Equal(t, int32(11), takerFee)
	makerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob1Address, false, 2, uint32(1))
	require.Equal(t, int32(1), makerFee)
}

func TestGetMegavaultEquity(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Equity of subaccount 0 of `megavault` module account.
		megavaultSaEquity *big.Int
		// Vaults
		vaults []vaulttypes.Vault
		// Equity of each vault above.
		vaultEquities []*big.Int

		/* --- Expectations --- */
		// Expected megavault equity.
		expectedMegavaultEquity *big.Int
	}{
		"Megavault subaccount with 1 equity, One quoting vault with 1 equity": {
			megavaultSaEquity: big.NewInt(1),
			vaults: []vaulttypes.Vault{
				{
					VaultId: constants.Vault_Clob0,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
				},
			},
			vaultEquities: []*big.Int{
				big.NewInt(1),
			},
			expectedMegavaultEquity: big.NewInt(2),
		},
		"Megavault subaccount with 94 equity, One quoting vault with 7 equity, One quoting vault with 0 equity,": {
			megavaultSaEquity: big.NewInt(94),
			vaults: []vaulttypes.Vault{
				{
					VaultId: constants.Vault_Clob0,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
				},
				{
					VaultId: constants.Vault_Clob1,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
				},
			},
			vaultEquities: []*big.Int{
				big.NewInt(7),
				big.NewInt(0),
			},
			expectedMegavaultEquity: big.NewInt(101),
		},
		"Megavault subaccount with 0 equity, One quoting vault with 123 equity, One stand-by vault with 6789 equity,": {
			megavaultSaEquity: big.NewInt(0),
			vaults: []vaulttypes.Vault{
				{
					VaultId: constants.Vault_Clob0,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
				},
				{
					VaultId: constants.Vault_Clob1,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
					},
				},
			},
			vaultEquities: []*big.Int{
				big.NewInt(123),
				big.NewInt(6_789),
			},
			expectedMegavaultEquity: big.NewInt(6_912),
		},
		"Megavault subaccount with 1000 equity, One quoting vault with 345 equity, One close-only vault with -1 equity,": {
			megavaultSaEquity: big.NewInt(1_000),
			vaults: []vaulttypes.Vault{
				{
					VaultId: constants.Vault_Clob0,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
				},
				{
					VaultId: constants.Vault_Clob1,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
					},
				},
			},
			vaultEquities: []*big.Int{
				big.NewInt(345),
				big.NewInt(-1),
			},
			expectedMegavaultEquity: big.NewInt(1_345),
		},
		"Megavault subaccount with 1000 equity, One quoting vault with 345 equity, One deactivated vault with -5 equity,": {
			megavaultSaEquity: big.NewInt(1_000),
			vaults: []vaulttypes.Vault{
				{
					VaultId: constants.Vault_Clob0,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
				},
				{
					VaultId: constants.Vault_Clob1,
					VaultParams: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_DEACTIVATED,
					},
				},
			},
			vaultEquities: []*big.Int{
				big.NewInt(345),
				big.NewInt(-5),
			},
			expectedMegavaultEquity: big.NewInt(1_345), // deactivated vault is not counted.
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize equities of megavault main subaccount and vaults.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						subaccounts := []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										0,
										tc.megavaultSaEquity,
									),
								},
							},
						}
						for i, vault := range tc.vaults {
							subaccounts = append(
								subaccounts,
								satypes.Subaccount{
									Id: vault.VaultId.ToSubaccountId(),
									AssetPositions: []*satypes.AssetPosition{
										testutil.CreateSingleAssetPosition(
											0,
											tc.vaultEquities[i],
										),
									},
								},
							)
						}
						genesisState.Subaccounts = subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.Vaults = tc.vaults
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			megavaultEquity, err := k.GetMegavaultEquity(ctx)
			require.NoError(t, err)
			require.Equal(
				t,
				tc.expectedMegavaultEquity,
				megavaultEquity,
			)
		})
	}
}
