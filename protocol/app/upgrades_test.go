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
	tests := map[string]struct {
		upgradeNames  []string
		expectedPanic string
	}{
		"Successfully setup upgrade handler": {
			upgradeNames: []string{"test1"},
		},
		"Successfully setup multiple upgrade handlers": {
			upgradeNames: []string{"test1", "test2"},
		},
		"Panic due to duplicate uppgrade names": {
			upgradeNames:  []string{"test1", "test1"},
			expectedPanic: "Cannot register duplicate upgrade handler 'test1'",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Restore the global variable back to original pre-test state.
			originalUpgrades := app.Upgrades
			t.Cleanup(func() {
				app.Upgrades = originalUpgrades
			})

			createHandlerWasCalled := make([]int, len(tc.upgradeNames))
			handlerWasCalled := make([]int, len(tc.upgradeNames))
			for i, upgradeName := range tc.upgradeNames {
				ii := i
				app.Upgrades = append(app.Upgrades, upgrades.Upgrade{
					UpgradeName: upgradeName,
					CreateUpgradeHandler: func(manager *module.Manager, configurator module.Configurator) types.UpgradeHandler {
						createHandlerWasCalled[ii] = createHandlerWasCalled[ii] + 1
						return func(ctx sdk.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
							handlerWasCalled[ii]++
							return nil, errors.New("upgrade failed")
						}
					},
				})
			}

			if tc.expectedPanic == "" {
				app := testapp.DefaultTestApp(nil)
				for i, upgradeName := range tc.upgradeNames {
					require.True(t, app.UpgradeKeeper.HasHandler(upgradeName))
					require.Equal(t, 1, createHandlerWasCalled[i])
					require.Equal(t, 0, handlerWasCalled[i])
				}
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

func TestDefaultUpgradesAndForks(t *testing.T) {
	require.Empty(t, app.Upgrades, "Expected empty upgrades list")
	require.Empty(t, app.Forks, "Expected empty forks list")
}
