package v_7_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	vaultkeeper "github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
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
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
