package v_4_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	ratelimitkeeper "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/keeper"
	ratelimittypes "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	rateLimitKeepr ratelimitkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
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
