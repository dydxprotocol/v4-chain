package v_7_0_0

import (
	"context"
	"fmt"
	"math/big"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	vaultkeeper "github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

const (
	// Each megavault share is worth 1 USDC.
	QUOTE_QUANTUMS_PER_MEGAVAULT_SHARE = 1_000_000
)

func initCurrencyPairIDCache(ctx sdk.Context, k pricestypes.PricesKeeper) {
	marketParams := k.GetAllMarketParams(ctx)
	for _, mp := range marketParams {
		currencyPair, err := slinky.MarketPairToCurrencyPair(mp.Pair)
		if err != nil {
			panic(fmt.Sprintf("failed to convert market param pair to currency pair: %s", err))
		}
		k.AddCurrencyPairIDToStore(ctx, mp.Id, currencyPair)
	}
}

func migrateVaultQuotingParamsToVaultParams(ctx sdk.Context, k vaultkeeper.Keeper) {
	vaultIds := k.UnsafeGetAllVaultIds(ctx)
	ctx.Logger().Info(fmt.Sprintf("Migrating quoting parameters of %d vaults", len(vaultIds)))
	for _, vaultId := range vaultIds {
		quotingParams, exists := k.UnsafeGetQuotingParams(ctx, vaultId)
		vaultParams := vaulttypes.VaultParams{
			Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
		}
		if exists {
			vaultParams.QuotingParams = &quotingParams
		}
		err := k.SetVaultParams(ctx, vaultId, vaultParams)
		if err != nil {
			panic(
				fmt.Sprintf(
					"failed to set vault params for vault %+v with params %+v: %s",
					vaultId,
					vaultParams,
					err,
				),
			)
		}
		k.UnsafeDeleteQuotingParams(ctx, vaultId)
		ctx.Logger().Info(fmt.Sprintf(
			"Successfully migrated vault %+v",
			vaultId,
		))
	}
}

// In 6.x,
// Total shares store (key prefix `TotalShares:`) is `vaultId -> shares`
// Owner shares store (key prefix `OwnerShares:`) is `vaultId -> owner -> shares`
// In 7.x,
// Total shares store is just `"TotalShares" -> shares`
// Owner shares store (key prefix `OwnerShares:`) is `owner -> shares`
// Thus, this function
// 1. Calculate how much equity each owner owns
// 2. Delete all keys in deprecated total shares and owner shares stores
// 3. Grant each owner 1 megavault share per usdc of equity owned
// 4. Set total megavault shares to sum of all owner shares granted
func migrateVaultSharesToMegavaultShares(ctx sdk.Context, k vaultkeeper.Keeper) {
	ctx.Logger().Info("Migrating vault shares to megavault shares")
	quoteQuantumsPerShare := big.NewInt(QUOTE_QUANTUMS_PER_MEGAVAULT_SHARE)

	ownerEquities := k.UnsafeGetAllOwnerEquities(ctx)
	ctx.Logger().Info(fmt.Sprintf("Calculated owner equities %s", ownerEquities))
	k.UnsafeDeleteAllVaultTotalShares(ctx)
	ctx.Logger().Info("Deleted all keys in deprecated vault total shares store")
	k.UnsafeDeleteAllVaultOwnerShares(ctx)
	ctx.Logger().Info("Deleted all keys in deprecated vault owner shares store")

	totalShares := big.NewInt(0)
	for owner, equity := range ownerEquities {
		ownerShares := new(big.Int).Quo(
			equity.Num(),
			equity.Denom(),
		)
		ownerShares.Quo(ownerShares, quoteQuantumsPerShare)

		if ownerShares.Sign() <= 0 {
			ctx.Logger().Warn(fmt.Sprintf(
				"Owner %s has non-positive shares %s from %s quote quantums",
				owner,
				ownerShares,
				equity,
			))
			continue
		}

		err := k.SetOwnerShares(ctx, owner, vaulttypes.BigIntToNumShares(ownerShares))
		if err != nil {
			panic(err)
		}
		ctx.Logger().Info(fmt.Sprintf(
			"Set megavault owner shares of %s: shares=%s, equity=%s",
			owner,
			ownerShares,
			equity,
		))

		totalShares.Add(totalShares, ownerShares)
	}

	err := k.SetTotalShares(ctx, vaulttypes.BigIntToNumShares(totalShares))
	if err != nil {
		panic(err)
	}
	ctx.Logger().Info(fmt.Sprintf("Set megavault total shares to: %s", totalShares))
	ctx.Logger().Info("Successfully migrated vault shares to megavault shares")
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	pricesKeeper pricestypes.PricesKeeper,
	vaultKeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Initialize the currency pair ID cache for all existing market params.
		initCurrencyPairIDCache(sdkCtx, pricesKeeper)

		// Migrate vault quoting params to vault params.
		migrateVaultQuotingParamsToVaultParams(sdkCtx, vaultKeeper)

		// Migrate vault shares to megavault shares.
		migrateVaultSharesToMegavaultShares(sdkCtx, vaultKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
