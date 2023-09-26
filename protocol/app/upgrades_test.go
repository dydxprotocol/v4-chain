package app_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/upgrades"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
)

func TestSetupUpgradeHandlers(t *testing.T) {
	upgradeName := "test-upgrade"
	createHandlerWasCalled := 0
	handlerWasCalled := 0
	upgradeHandlerFn := func(manager *module.Manager, configurator module.Configurator) types.UpgradeHandler {
		createHandlerWasCalled++
		return func(ctx sdk.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			handlerWasCalled++
			return nil, errors.New("upgrade failed")
		}
	}

	tests := map[string]struct {
		upgrades      []upgrades.Upgrade
		expectedPanic string
	}{
		"Successfully setup handler": {
			upgrades: []upgrades.Upgrade{
				{
					UpgradeName:          upgradeName,
					CreateUpgradeHandler: upgradeHandlerFn,
				},
			},
		},
		"Panic due to duplicate handlers": {
			upgrades: []upgrades.Upgrade{
				{
					UpgradeName:          upgradeName,
					CreateUpgradeHandler: upgradeHandlerFn,
				},
				{
					UpgradeName:          upgradeName,
					CreateUpgradeHandler: upgradeHandlerFn,
				},
			},
			expectedPanic: "Cannot register duplicate upgrade handler 'test-upgrade'",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Restore the global variable back to original pre-test state.
			originalUpgrades := app.Upgrades
			defer func() {
				app.Upgrades = originalUpgrades
			}()

			app.Upgrades = tc.upgrades
			if tc.expectedPanic == "" {
				app := testapp.DefaultTestApp(nil)
				require.True(t, app.UpgradeKeeper.HasHandler(upgradeName))
				require.Equal(t, 1, createHandlerWasCalled)
				require.Equal(t, 0, handlerWasCalled)
			} else {
				require.PanicsWithValue(
					t,
					tc.expectedPanic,
					func() {
						testapp.DefaultTestApp(nil)
					},
				)
			}
		})
	}
}
