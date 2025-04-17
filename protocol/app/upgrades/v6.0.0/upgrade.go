package v_6_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/dydxprotocol/slinky/providers/apis/dydx"
	dydxtypes "github.com/dydxprotocol/slinky/providers/apis/dydx/types"
	marketmapkeeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	vaultkeeper "github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// GovAuthority is the module account address of x/gov.
var GovAuthority = authtypes.NewModuleAddress(govtypes.ModuleName).String()

func removeStatefulFOKOrders(ctx sdk.Context, k clobtypes.ClobKeeper) {
	allStatefulOrders := k.GetAllStatefulOrders(ctx)
	for _, order := range allStatefulOrders {
		if order.TimeInForce == clobtypes.Order_TIME_IN_FORCE_FILL_OR_KILL {
			// Remove the orders from state.
			k.MustRemoveStatefulOrder(ctx, order.OrderId)

			// Send indexer event for removal of stateful order.
			k.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexerevents.StatefulOrderEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewStatefulOrderRemovalEvent(
						order.OrderId,
						indexershared.ConvertOrderRemovalReasonToIndexerOrderRemovalReason(
							clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED,
						),
					),
				),
			)
		}
	}
}

func setMarketMapParams(ctx sdk.Context, mmk marketmapkeeper.Keeper) {
	err := mmk.SetParams(ctx, marketmaptypes.Params{
		// init so that gov is the admin and a market authority
		MarketAuthorities: []string{GovAuthority},
		Admin:             GovAuthority,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to set x/mm params %v", err))
	}
}

func migratePricesToMarketMap(ctx sdk.Context, pk pricestypes.PricesKeeper, mmk marketmapkeeper.Keeper) {
	allMarketParams := pk.GetAllMarketParams(ctx)
	var mpr dydxtypes.QueryAllMarketParamsResponse
	for _, mp := range allMarketParams {
		mpr.MarketParams = append(mpr.MarketParams, dydxtypes.MarketParam{
			Id:                 mp.Id,
			Pair:               mp.Pair,
			Exponent:           mp.Exponent,
			MinExchanges:       mp.MinExchanges,
			MinPriceChangePpm:  mp.MinPriceChangePpm,
			ExchangeConfigJson: mp.ExchangeConfigJson,
		})
	}
	mm, err := dydx.ConvertMarketParamsToMarketMap(mpr)
	if err != nil {
		panic(fmt.Sprintf("Couldn't convert markets %v", err))
	}
	for _, market := range mm.MarketMap.Markets {
		err = mmk.CreateMarket(ctx, market)
		if err != nil {
			panic(fmt.Sprintf("Failed to create market %s", market.Ticker.String()))
		}
	}
}

func initRevShareModuleState(
	ctx sdk.Context,
	revShareKeeper revsharetypes.RevShareKeeper,
	priceKeeper pricestypes.PricesKeeper,
) {
	// Initialize the rev share module state.
	params := revsharetypes.MarketMapperRevenueShareParams{
		Address:         authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
		RevenueSharePpm: 0,
		ValidDays:       0,
	}
	err := revShareKeeper.SetMarketMapperRevenueShareParams(ctx, params)
	if err != nil {
		panic(fmt.Sprintf("failed to set market mapper revenue share params: %s", err))
	}

	// Initialize the rev share details for all existing markets.
	markets := priceKeeper.GetAllMarketParams(ctx)
	for _, market := range markets {
		revShareDetails := revsharetypes.MarketMapperRevShareDetails{
			ExpirationTs: 0,
		}
		revShareKeeper.SetMarketMapperRevShareDetails(ctx, market.Id, revShareDetails)
	}
}

func initVaultDefaultQuotingParams(
	ctx sdk.Context,
	vaultKeeper vaultkeeper.Keeper,
) {
	// Initialize the default quoting params for the vault module.
	oldParams := vaultKeeper.UnsafeGetParams(ctx)
	if err := vaultKeeper.SetDefaultQuotingParams(
		ctx,
		vaulttypes.QuotingParams{
			Layers:                           oldParams.Layers,
			SpreadMinPpm:                     oldParams.SpreadMinPpm,
			SpreadBufferPpm:                  oldParams.SpreadBufferPpm,
			SkewFactorPpm:                    oldParams.SkewFactorPpm,
			OrderSizePctPpm:                  oldParams.OrderSizePctPpm,
			OrderExpirationSeconds:           oldParams.OrderExpirationSeconds,
			ActivationThresholdQuoteQuantums: oldParams.ActivationThresholdQuoteQuantums,
		},
	); err != nil {
		panic(fmt.Sprintf("failed to set vault default quoting params: %s", err))
	}

	// Delete deprecated `Params`.
	vaultKeeper.UnsafeDeleteParams(ctx)
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper clobtypes.ClobKeeper,
	pricesKeeper pricestypes.PricesKeeper,
	mmKeeper marketmapkeeper.Keeper,
	revShareKeeper revsharetypes.RevShareKeeper,
	vaultKeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Remove all stateful FOK orders from state.
		removeStatefulFOKOrders(sdkCtx, clobKeeper)

		// Migrate x/prices params to x/marketmap Markets
		migratePricesToMarketMap(sdkCtx, pricesKeeper, mmKeeper)
		// Set x/marketmap Params
		setMarketMapParams(sdkCtx, mmKeeper)

		// Initialize the rev share module state.
		initRevShareModuleState(sdkCtx, revShareKeeper, pricesKeeper)

		// Initialize x/vault default quoting params.
		initVaultDefaultQuotingParams(sdkCtx, vaultKeeper)

		sdkCtx.Logger().Info("Successfully removed stateful orders from state")

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
