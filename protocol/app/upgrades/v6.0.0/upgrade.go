package v_6_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
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

func initRevShareModuleState(
	ctx sdk.Context,
	revShareKeeper revsharetypes.RevShareKeeper,
	priceKeeper pricetypes.PricesKeeper,
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

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	clobKeeper clobtypes.ClobKeeper,
	revShareKeeper revsharetypes.RevShareKeeper,
	priceKeeper pricetypes.PricesKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		// Remove all stateful FOK orders from state.
		removeStatefulFOKOrders(sdkCtx, clobKeeper)

		// Initialize the rev share module state.
		initRevShareModuleState(sdkCtx, revShareKeeper, priceKeeper)

		sdkCtx.Logger().Info("Successfully removed stateful orders from state")

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
