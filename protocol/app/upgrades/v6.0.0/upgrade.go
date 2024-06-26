package v_6_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	marketmapkeeper "github.com/skip-mev/slinky/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

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
		// todo fill out these fields
		MarketAuthorities: nil,
		Admin:             "",
	})
	if err != nil {
		panic(fmt.Sprintf("failed to set x/mm params %v", err))
	}
}

func migratePricesToMarketMap(ctx sdk.Context, pk pricestypes.PricesKeeper, mmk marketmapkeeper.Keeper) {
	h, err := dydx.NewAPIHandler(zap.NewNop(), config.APIConfig{
		Enabled:          true,
		Timeout:          1,
		Interval:         1,
		ReconnectTimeout: 1,
		MaxQueries:       1,
		Atomic:           false,
		Endpoints:        []config.Endpoint{{URL: "upgrade"}},
		BatchSize:        0,
		Name:             dydx.Name,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to construct dydx handler %v", err))
	}
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
	mm, err := h.ConvertMarketParamsToMarketMap(mpr)
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

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper clobtypes.ClobKeeper,
	pricesKeeper pricestypes.PricesKeeper,
	mmKeeper marketmapkeeper.Keeper,
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

		sdkCtx.Logger().Info("Successfully removed stateful orders from state")

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
