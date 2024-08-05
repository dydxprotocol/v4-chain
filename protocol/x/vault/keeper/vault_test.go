package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
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
				constants.Vault_Clob0,
				constants.Vault_Clob1,
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
				constants.Vault_Clob0,
				constants.Vault_Clob1,
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
				constants.Vault_Clob0,
				constants.Vault_Clob1,
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
		// Vault quoting params.
		quotingParams *vaulttypes.QuotingParams
	}{
		"Total shares doesn't exist, no owners, default quoting params": {
			vaultId: constants.Vault_Clob0,
		},
		"Total shares exists, no owners, default quoting params": {
			vaultId:           constants.Vault_Clob0,
			totalSharesExists: true,
		},
		"Total shares exists, one owner, non-default quoting params": {
			vaultId:           constants.Vault_Clob1,
			totalSharesExists: true,
			owners:            []string{constants.Alice_Num0.Owner},
			quotingParams:     &constants.QuotingParams,
		},
		"Total shares exists, two owners, non-default quoting params": {
			vaultId:           constants.Vault_Clob1,
			totalSharesExists: true,
			owners: []string{
				constants.Alice_Num0.Owner,
				constants.Bob_Num0.Owner,
			},
			quotingParams: &constants.QuotingParams,
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

			k.AddVaultToAddressStore(ctx, tc.vaultId)
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
			if tc.quotingParams != nil {
				err := k.SetVaultQuotingParams(ctx, tc.vaultId, *tc.quotingParams)
				require.NoError(t, err)
			}

			// Decommission vault.
			k.DecommissionVault(ctx, tc.vaultId)

			// Check that total shares, owner shares, and vault address are deleted.
			_, exists := k.GetTotalShares(ctx, tc.vaultId)
			require.Equal(t, false, exists)
			for _, owner := range tc.owners {
				_, exists = k.GetOwnerShares(ctx, tc.vaultId, owner)
				require.Equal(t, false, exists)
			}
			require.False(t, k.IsVault(ctx, tc.vaultId.ToModuleAccountAddress()))
			// Check that vault quoting params are back to default.
			require.Equal(
				t,
				k.GetDefaultQuotingParams(ctx),
				k.GetVaultQuotingParams(ctx, tc.vaultId),
			)
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
				genesisState.Vaults = []*vaulttypes.Vault{
					{
						VaultId: &constants.Vault_Clob0,
						TotalShares: &vaulttypes.NumShares{
							NumShares: dtypes.NewInt(10),
						},
						OwnerShares: []*vaulttypes.OwnerShare{
							{
								Owner: constants.AliceAccAddress.String(),
								Shares: &vaulttypes.NumShares{
									NumShares: dtypes.NewInt(10),
								},
							},
						},
					},
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
						Id: &constants.Alice_Num0,
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

	// Vault in genesis state should be in best fee tier.
	takerFee := tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob0Address, true)
	require.Equal(t, int32(11), takerFee)
	makerFee := tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob0Address, false)
	require.Equal(t, int32(1), makerFee)

	// A regular user Alice should be in worst fee tier.
	takerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, aliceAddress, true)
	require.Equal(t, int32(33), takerFee)
	makerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, aliceAddress, false)
	require.Equal(t, int32(3), makerFee)

	// A newly created vault should be in best fee tier.
	checkTx_DepositToVault := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
			Gas:                  constants.TestGasLimit,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		&vaulttypes.MsgDepositToVault{
			VaultId:       &constants.Vault_Clob1,
			SubaccountId:  &constants.Alice_Num0,
			QuoteQuantums: dtypes.NewInt(1),
		},
	)
	checkTxResp := tApp.CheckTx(checkTx_DepositToVault)
	require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

	ctx = tApp.AdvanceToBlock(
		uint32(ctx.BlockHeight())+1,
		testapp.AdvanceToBlockOptions{},
	)
	takerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob1Address, true)
	require.Equal(t, int32(11), takerFee)
	makerFee = tApp.App.FeeTiersKeeper.GetPerpetualFeePpm(ctx, vaultClob1Address, false)
	require.Equal(t, int32(1), makerFee)
}
