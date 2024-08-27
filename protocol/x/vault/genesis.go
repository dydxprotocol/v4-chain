package vault

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	// Set default quoting params.
	if err := k.SetDefaultQuotingParams(ctx, genState.DefaultQuotingParams); err != nil {
		panic(err)
	}
	// Set total shares, owner shares, and locked shares.
	if err := k.SetTotalShares(ctx, genState.TotalShares); err != nil {
		panic(err)
	}
	for _, ownerShares := range genState.OwnerShares {
		if err := k.SetOwnerShares(ctx, ownerShares.Owner, ownerShares.Shares); err != nil {
			panic(err)
		}
	}
	for _, lockedShares := range genState.AllLockedShares {
		if err := k.SetLockedShares(ctx, lockedShares.OwnerAddress, lockedShares); err != nil {
			panic(err)
		}
	}

	// For each vault:
	// 1. Set vault params
	// 2. Set most recent client ids
	// 3. Add to address store
	for _, vault := range genState.Vaults {
		if err := k.SetVaultParams(ctx, vault.VaultId, vault.VaultParams); err != nil {
			panic(err)
		}
		k.SetMostRecentClientIds(ctx, vault.VaultId, vault.MostRecentClientIds)
		k.AddVaultToAddressStore(ctx, vault.VaultId)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export total shares, owner shares, and locked shares.
	genesis.TotalShares = k.GetTotalShares(ctx)
	genesis.OwnerShares = k.GetAllOwnerShares(ctx)
	genesis.AllLockedShares = k.GetAllLockedShares(ctx)

	// Export params.
	genesis.DefaultQuotingParams = k.GetDefaultQuotingParams(ctx)

	// Export vaults.
	genesis.Vaults = k.GetAllVaults(ctx)

	return genesis
}
