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
	// For each vault:
	// 1. Set total shares
	// 2. Set owner shares
	// 3. Set vault params
	// 4. Set most recent client ids
	// 5. Add to address store
	for _, vault := range genState.Vaults {
		if err := k.SetTotalShares(ctx, *vault.VaultId, *vault.TotalShares); err != nil {
			panic(err)
		}
		for _, ownerShares := range vault.OwnerShares {
			if err := k.SetOwnerShares(ctx, *vault.VaultId, ownerShares.Owner, *ownerShares.Shares); err != nil {
				panic(err)
			}
		}
		if err := k.SetVaultParams(ctx, *vault.VaultId, vault.VaultParams); err != nil {
			panic(err)
		}
		k.SetMostRecentClientIds(ctx, *vault.VaultId, vault.MostRecentClientIds)
		k.AddVaultToAddressStore(ctx, *vault.VaultId)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export params.
	genesis.DefaultQuotingParams = k.GetDefaultQuotingParams(ctx)

	// Export vaults.
	genesis.Vaults = k.GetAllVaults(ctx)

	return genesis
}
