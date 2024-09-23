package keeper

import (
	"fmt"
	"math/big"

	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"

	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetMegavaultEquity returns the equity of the megavault (in quote quantums), which consists of
// - equity of the megavault main subaccount
// - equity of all vaults (if not-deactivated and positive)
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

		if vaultParams.Status == types.VaultStatus_VAULT_STATUS_DEACTIVATED {
			continue
		}

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

// GetVaultLeverageAndEquity returns a vault's leverage and equity.
// - leverage = open notional / equity.
// Note that an error is returned if equity is non-positive.
func (k Keeper) GetVaultLeverageAndEquity(
	ctx sdk.Context,
	vaultId types.VaultId,
	perpetual *perptypes.Perpetual,
	marketPrice *pricestypes.MarketPrice,
) (
	leverage *big.Rat,
	equity *big.Int,
	err error,
) {
	equity, err = k.GetVaultEquity(ctx, vaultId)
	if err != nil {
		return nil, nil, err
	}
	if equity.Sign() <= 0 {
		return nil, equity, errorsmod.Wrap(
			types.ErrNonPositiveEquity,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}

	inventory := k.GetVaultInventoryInPerpetual(ctx, vaultId, perpetual.GetId())
	openNotional := lib.BaseToQuoteQuantums(
		inventory,
		perpetual.Params.AtomicResolution,
		marketPrice.GetPrice(),
		marketPrice.GetExponent(),
	)
	leverage = new(big.Rat).SetFrac(
		openNotional,
		equity,
	)

	return leverage, equity, nil
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

// GetVaultClobPerpAndMarket returns the clob pair, perpetual, market param, and market price
// that correspond to a vault.
func (k Keeper) GetVaultClobPerpAndMarket(
	ctx sdk.Context,
	vaultId types.VaultId,
) (
	clobPair clobtypes.ClobPair,
	perpetual perptypes.Perpetual,
	marketParam pricestypes.MarketParam,
	marketPrice pricestypes.MarketPrice,
	err error,
) {
	clobPair, exists := k.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(vaultId.Number))
	if !exists {
		return clobPair, perpetual, marketParam, marketPrice, errorsmod.Wrap(
			types.ErrClobPairNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	perpId := clobPair.Metadata.(*clobtypes.ClobPair_PerpetualClobMetadata).PerpetualClobMetadata.PerpetualId
	perpetual, err = k.perpetualsKeeper.GetPerpetual(ctx, perpId)
	if err != nil {
		return clobPair, perpetual, marketParam, marketPrice, errorsmod.Wrap(
			err,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	marketParam, exists = k.pricesKeeper.GetMarketParam(ctx, perpetual.Params.MarketId)
	if !exists {
		return clobPair, perpetual, marketParam, marketPrice, errorsmod.Wrap(
			types.ErrMarketParamNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	marketPrice, err = k.pricesKeeper.GetMarketPrice(ctx, perpetual.Params.MarketId)
	if err != nil {
		return clobPair, perpetual, marketParam, marketPrice, errorsmod.Wrap(
			err,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}

	return clobPair, perpetual, marketParam, marketPrice, nil
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

// AllocateToVault transfers funds from main vault to a specified vault.
func (k Keeper) AllocateToVault(
	ctx sdk.Context,
	vaultId types.VaultId,
	quantums *big.Int,
) error {
	// Check if vault has a corresponding clob pair.
	_, exists := k.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(vaultId.Number))
	if !exists {
		return types.ErrClobPairNotFound
	}

	// If vault doesn't exist:
	// 1. initialize params with `STAND_BY` status.
	// 2. add vault to address store.
	_, exists = k.GetVaultParams(ctx, vaultId)
	if !exists {
		err := k.SetVaultParams(
			ctx,
			vaultId,
			types.VaultParams{
				Status: types.VaultStatus_VAULT_STATUS_STAND_BY,
			},
		)
		if err != nil {
			return err
		}
		k.AddVaultToAddressStore(ctx, vaultId)
	}

	// Transfer from main vault to the specified vault.
	if err := k.sendingKeeper.ProcessTransfer(
		ctx,
		&sendingtypes.Transfer{
			Sender:    types.MegavaultMainSubaccount,
			Recipient: *vaultId.ToSubaccountId(),
			AssetId:   assetstypes.AssetUsdc.Id,
			Amount:    quantums.Uint64(), // validated to be positive above.
		},
	); err != nil {
		return err
	}
	return nil
}
