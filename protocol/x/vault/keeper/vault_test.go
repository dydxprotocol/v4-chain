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

func TestDecomissionVaults(t *testing.T) {
	vault0 := constants.Vault_Clob_0
	vault1 := constants.Vault_Clob_1
	// Initialize vault 0 with positive equity.
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					{
						Id: vault0.ToSubaccountId(),
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
		return genesis
	}).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Set total shares and owner shares for both vaults.
	shares := vaulttypes.BigRatToNumShares(
		big.NewRat(7, 1),
	)
	err := k.SetTotalShares(
		ctx,
		vault0,
		shares,
	)
	require.NoError(t, err)
	err = k.SetOwnerShares(
		ctx,
		vault0,
		constants.Alice_Num0.Owner,
		shares,
	)
	require.NoError(t, err)
	err = k.SetTotalShares(
		ctx,
		vault1,
		shares,
	)
	require.NoError(t, err)
	err = k.SetOwnerShares(
		ctx,
		vault1,
		constants.Bob_Num0.Owner,
		shares,
	)
	require.NoError(t, err)

	// Decomission all vaults.
	k.DecommissionVaults(ctx)

	// Check that total shares and owner shares are not deleted for vault 0.
	got, exists := k.GetTotalShares(ctx, vault0)
	require.Equal(t, true, exists)
	require.Equal(t, shares, got)
	got, exists = k.GetOwnerShares(ctx, vault0, constants.Alice_Num0.Owner)
	require.Equal(t, true, exists)
	require.Equal(t, shares, got)
	// Check that total shares and owner shares are deleted for vault 1.
	_, exists = k.GetTotalShares(ctx, vault1)
	require.Equal(t, false, exists)
	_, exists = k.GetOwnerShares(ctx, vault1, constants.Bob_Num0.Owner)
	require.Equal(t, false, exists)
}

func TestDecomissionVault(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Decomission a non-existent vault.
	k.DecommissionVault(ctx, constants.Vault_Clob_0)

	// Set total shares and owner shares for two owners of a vault.
	shares := vaulttypes.BigRatToNumShares(
		big.NewRat(7, 1),
	)
	err := k.SetTotalShares(
		ctx,
		constants.Vault_Clob_0,
		shares,
	)
	require.NoError(t, err)
	err = k.SetOwnerShares(
		ctx,
		constants.Vault_Clob_0,
		constants.Alice_Num0.Owner,
		shares,
	)
	require.NoError(t, err)
	err = k.SetOwnerShares(
		ctx,
		constants.Vault_Clob_0,
		constants.Bob_Num0.Owner,
		shares,
	)
	require.NoError(t, err)

	// Decomission above vault.
	k.DecommissionVault(ctx, constants.Vault_Clob_0)

	// Check that total shares and owner shares are deleted.
	_, exists := k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, false, exists)
	_, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_0, constants.Alice_Num0.Owner)
	require.Equal(t, false, exists)
	_, exists = k.GetOwnerShares(ctx, constants.Vault_Clob_0, constants.Bob_Num0.Owner)
	require.Equal(t, false, exists)
}
