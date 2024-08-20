package keeper

import (
	"math/big"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetVaultEquity returns the equity of a vault (in quote quantums).
func (k Keeper) GetVaultEquity(
	ctx sdk.Context,
	vaultId types.VaultId,
) (*big.Int, error) {
	risk, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{
			SubaccountId: *vaultId.ToSubaccountId(),
		},
	)
	if err != nil {
		return nil, err
	}
	return risk.NC, nil
}

// GetVaultInventory returns the inventory of a vault in a given perpeutal (in base quantums).
func (k Keeper) GetVaultInventoryInPerpetual(
	ctx sdk.Context,
	vaultId types.VaultId,
	perpId uint32,
) *big.Int {
	// Get subaccount.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, *vaultId.ToSubaccountId())
	// Calculate inventory.
	inventory := big.NewInt(0)
	for _, p := range subaccount.PerpetualPositions {
		if p.GetPerpetualId() == perpId {
			inventory.Add(inventory, p.GetBigQuantums())
			break
		}
	}
	return inventory
}

// DecommissionVaults decommissions all deactivated vaults that have non-positive equities.
func (k Keeper) DecommissionNonPositiveEquityVaults(
	ctx sdk.Context,
) {
	// Iterate through all vaults.
	vaultParamsIterator := k.getVaultParamsIterator(ctx)
	defer vaultParamsIterator.Close()
	for ; vaultParamsIterator.Valid(); vaultParamsIterator.Next() {
		var vaultParams types.VaultParams
		k.cdc.MustUnmarshal(vaultParamsIterator.Value(), &vaultParams)

		// Skip if vault is not deactivated.
		if vaultParams.Status != types.VaultStatus_VAULT_STATUS_DEACTIVATED {
			continue
		}

		// Get vault equity.
		vaultId, err := types.GetVaultIdFromStateKey(vaultParamsIterator.Key())
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to get vault ID from state key", err)
			continue
		}
		equity, err := k.GetVaultEquity(ctx, *vaultId)
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to get vault equity", err)
			continue
		}

		// Decommission vault if equity is non-positive.
		if equity.Sign() <= 0 {
			k.DecommissionVault(ctx, *vaultId)
		}
	}
}

// DecommissionVault decommissions a vault by
// 1. deleting its total shares and owner shares
// 2. deleting its address from vault address store
// 3. deleting its quoting params if any
func (k Keeper) DecommissionVault(
	ctx sdk.Context,
	vaultId types.VaultId,
) {
	// Delete TotalShares of the vault.
	totalSharesStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.TotalSharesKeyPrefix))
	totalSharesStore.Delete(vaultId.ToStateKey())

	// Delete all OwnerShares of the vault.
	ownerSharesStore := k.getVaultOwnerSharesStore(ctx, vaultId)
	ownerSharesIterator := storetypes.KVStorePrefixIterator(ownerSharesStore, []byte{})
	defer ownerSharesIterator.Close()
	for ; ownerSharesIterator.Valid(); ownerSharesIterator.Next() {
		ownerSharesStore.Delete(ownerSharesIterator.Key())
	}

	// Delete from vault address store.
	vaultAddressStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultAddressKeyPrefix))
	vaultAddressStore.Delete([]byte(vaultId.ToModuleAccountAddress()))

	// Delete vault params.
	vaultParamsStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultParamsKeyPrefix))
	vaultParamsStore.Delete(vaultId.ToStateKey())
}

// AddVaultToAddressStore adds a vault's address to the vault address store.
func (k Keeper) AddVaultToAddressStore(
	ctx sdk.Context,
	vaultId types.VaultId,
) {
	vaultAddressStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultAddressKeyPrefix))
	vaultAddressStore.Set([]byte(vaultId.ToModuleAccountAddress()), []byte{})
}

// IsVault checks if a given address is the address of an existing vault.
func (k Keeper) IsVault(
	ctx sdk.Context,
	address string,
) bool {
	vaultAddressStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultAddressKeyPrefix))
	return vaultAddressStore.Has([]byte(address))
}

// GetAllVaults returns all vaults with their total shares, owner shares, and individual params.
// Note: This function is only used for exporting module state.
func (k Keeper) GetAllVaults(ctx sdk.Context) []*types.Vault {
	vaults := []*types.Vault{}
	vaultParamsIterator := k.getVaultParamsIterator(ctx)
	defer vaultParamsIterator.Close()
	for ; vaultParamsIterator.Valid(); vaultParamsIterator.Next() {
		vaultId, err := types.GetVaultIdFromStateKey(vaultParamsIterator.Key())
		if err != nil {
			panic(err)
		}

		var vaultParams types.VaultParams
		k.cdc.MustUnmarshal(vaultParamsIterator.Value(), &vaultParams)

		totalShares, exists := k.GetTotalShares(ctx, *vaultId)
		if !exists {
			panic("TotalShares not found for vault " + vaultId.ToString())
		}

		allOwnerShares := k.GetAllOwnerShares(ctx, *vaultId)

		mostRecentClientIds := k.GetMostRecentClientIds(ctx, *vaultId)

		vaults = append(vaults, &types.Vault{
			VaultId:             vaultId,
			TotalShares:         &totalShares,
			OwnerShares:         allOwnerShares,
			VaultParams:         vaultParams,
			MostRecentClientIds: mostRecentClientIds,
		})
	}
	return vaults
}
