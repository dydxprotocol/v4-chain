package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetMegavaultEquity returns the equity of the megavault (in quote quantums), which consists of
// - equity of the megavault main subaccount
// - equity of all vaults (if positive)
func (k Keeper) GetMegavaultEquity(ctx sdk.Context) (*big.Int, error) {
	megavaultEquity, err := k.GetSubaccountEquity(ctx, types.MegavaultMainSubaccount)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to get megavault subaccount equity")
	}

	// Add equities of all vaults.
	vaultParamsIterator := k.getVaultParamsIterator(ctx)
	defer vaultParamsIterator.Close()
	for ; vaultParamsIterator.Valid(); vaultParamsIterator.Next() {
		var vaultParams types.VaultParams
		k.cdc.MustUnmarshal(vaultParamsIterator.Value(), &vaultParams)

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

		// Add equity if it is positive.
		if equity.Sign() > 0 {
			megavaultEquity.Add(megavaultEquity, equity)
		}
	}

	return megavaultEquity, nil
}

// GetVaultEquity returns the equity of a vault (in quote quantums).
func (k Keeper) GetVaultEquity(
	ctx sdk.Context,
	vaultId types.VaultId,
) (*big.Int, error) {
	return k.GetSubaccountEquity(ctx, *vaultId.ToSubaccountId())
}

// GetSubaccountEquity returns the equity of a subaccount (in quote quantums).
func (k Keeper) GetSubaccountEquity(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (*big.Int, error) {
	risk, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{
			SubaccountId: subaccountId,
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

// DecommissionVault decommissions a vault by deleting it from vault address store and vault params store.
func (k Keeper) DecommissionVault(
	ctx sdk.Context,
	vaultId types.VaultId,
) {
	// Delete from vault address store.
	vaultAddressStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultAddressKeyPrefix))
	vaultAddressStore.Delete([]byte(vaultId.ToModuleAccountAddress()))

	// Delete vault params.
	vaultParamsStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultParamsKeyPrefix))
	vaultParamsStore.Delete(vaultId.ToStateKey())

	// Delete most recent client IDs.
	clientIdsStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MostRecentClientIdsKeyPrefix))
	clientIdsStore.Delete(vaultId.ToStateKey())
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

// GetAllVaults returns all vaults with their vault params and most recent client IDs.
// Note: This function is only used for exporting module state.
func (k Keeper) GetAllVaults(ctx sdk.Context) []types.Vault {
	vaults := []types.Vault{}
	vaultParamsIterator := k.getVaultParamsIterator(ctx)
	defer vaultParamsIterator.Close()
	for ; vaultParamsIterator.Valid(); vaultParamsIterator.Next() {
		vaultId, err := types.GetVaultIdFromStateKey(vaultParamsIterator.Key())
		if err != nil {
			panic(err)
		}

		var vaultParams types.VaultParams
		k.cdc.MustUnmarshal(vaultParamsIterator.Value(), &vaultParams)

		mostRecentClientIds := k.GetMostRecentClientIds(ctx, *vaultId)

		vaults = append(vaults, types.Vault{
			VaultId:             *vaultId,
			VaultParams:         vaultParams,
			MostRecentClientIds: mostRecentClientIds,
		})
	}
	return vaults
}
