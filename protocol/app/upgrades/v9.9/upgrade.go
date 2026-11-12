package v_9_9

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	canonicalusdckeeper "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keeper *canonicalusdckeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}
		// A route that does not resolve leaves the module DISABLED for a
		// governance proposal rather than failing the upgrade.
		if err := enterGradual(sdkCtx, keeper); err != nil {
			sdkCtx.Logger().Error("canonicalusdc stays DISABLED after the upgrade", "err", err)
		}
		return vm, nil
	}
}

func enterGradual(ctx sdk.Context, keeper *canonicalusdckeeper.Keeper) error {
	if keeper.GetControls(ctx).Mode != canonicalusdctypes.Mode_MODE_DISABLED {
		return nil
	}
	routes, pinned := PinnedRoutes[ctx.ChainID()]
	if !pinned {
		var err error
		if routes.Noble, err = keeper.ResolveRoute(ctx, NobleChainID); err != nil {
			return err
		}
		if routes.Injective, err = keeper.ResolveRoute(ctx, InjectiveChainID); err != nil {
			return err
		}
	}
	noble, injective := routes.Noble, routes.Injective
	controls := canonicalusdctypes.DefaultControls()
	controls.Mode = canonicalusdctypes.Mode_MODE_GRADUAL
	controls.NobleChannel = noble.Channel
	controls.NobleConnection = noble.Connection
	controls.NobleClient = noble.Client
	controls.NoblePacketDenom = NoblePacketDenom
	controls.NobleWithdrawalsEnabled = true
	controls.InjectiveChannel = injective.Channel
	controls.InjectiveConnection = injective.Connection
	controls.InjectiveClient = injective.Client
	controls.InjectivePacketDenom = InjectivePacketDenom
	controls.InjectiveDenom = transfertypes.ParseDenomTrace(
		fmt.Sprintf("%s/%s/%s", transfertypes.PortID, injective.Channel, InjectivePacketDenom),
	).IBCDenom()
	controls.MaxPendingSettlements = MaxPendingSettlements
	return keeper.ApplyControls(ctx, controls)
}
