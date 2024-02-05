package v_4_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	ratelimitkeeper "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/keeper"
	ratelimittypes "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	rateLimitKeepr ratelimitkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		if err := rateLimitKeepr.SetLimitParams(
			sdkCtx,
			ratelimittypes.DefaultUsdcRateLimitParams(),
		); err != nil {
			panic(fmt.Sprintf("failed to set default x/ratelimit params: %s", err))
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
