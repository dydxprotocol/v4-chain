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

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	pricesKeeper pricestypes.PricesKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Initialize the currency pair ID cache for all existing market params.
		initCurrencyPairIDCache(sdkCtx, pricesKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
