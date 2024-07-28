package v_cosmwasm_0

import (
	"context"
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

var ModuleAccsToInitialize = []string{
	wasmtypes.ModuleName,
}

// copied from v3 upgrade handler
func InitializeModuleAccs(ctx sdk.Context, ak authkeeper.AccountKeeper) {
	for _, modAccName := range ModuleAccsToInitialize {
		// Get module account and relevant permissions from the accountKeeper.
		addr, perms := ak.GetModuleAddressAndPermissions(modAccName)
		if addr == nil {
			panic(fmt.Sprintf(
				"Did not find %v in `ak.GetModuleAddressAndPermissions`. This is not expected. Skipping.",
				modAccName,
			))
		}

		// Try to get the account in state.
		acc := ak.GetAccount(ctx, addr)
		if acc != nil {
			// Account has been initialized.
			macc, isModuleAccount := acc.(sdk.ModuleAccountI)
			if isModuleAccount {
				// Module account was correctly initialized. Skipping
				ctx.Logger().Info(fmt.Sprintf(
					"module account %+v was correctly initialized. No-op",
					macc,
				))
				continue
			}
			// Module account has been initialized as a BaseAccount. Change to module account.
			// Note: We need to get the base account to retrieve its account number, and convert it
			// in place into a module account.
			baseAccount, ok := acc.(*authtypes.BaseAccount)
			if !ok {
				panic(fmt.Sprintf(
					"cannot cast %v into a BaseAccount, acc = %+v",
					modAccName,
					acc,
				))
			}
			newModuleAccount := authtypes.NewModuleAccount(
				baseAccount,
				modAccName,
				perms...,
			)
			ak.SetModuleAccount(ctx, newModuleAccount)
			ctx.Logger().Info(fmt.Sprintf(
				"Successfully converted %v to module account in state: %+v",
				modAccName,
				newModuleAccount,
			))
			continue
		}

		// Account has not been initialized at all. Initialize it as module.
		// Implementation taken from
		// https://github.com/dydxprotocol/cosmos-sdk/blob/bdf96fdd/x/auth/keeper/keeper.go#L213
		newModuleAccount := authtypes.NewEmptyModuleAccount(modAccName, perms...)
		maccI := (ak.NewAccount(ctx, newModuleAccount)).(sdk.ModuleAccountI) // this set the account number
		ak.SetModuleAccount(ctx, maccI)
		ctx.Logger().Info(fmt.Sprintf(
			"Successfully initialized module account in state: %+v",
			newModuleAccount,
		))
	}
}

// TODO(OTE-535): remove duplicated code from v6 upgrade
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

// TODO(OTE-535): remove duplicated code from v6 upgrade
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
	ak authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		InitializeModuleAccs(sdkCtx, ak)
		// Remove all stateful FOK orders from state.
		removeStatefulFOKOrders(sdkCtx, clobKeeper)

		// Initialize the rev share module state.
		initRevShareModuleState(sdkCtx, revShareKeeper, priceKeeper)

		sdkCtx.Logger().Info("Successfully removed stateful orders from state")

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
