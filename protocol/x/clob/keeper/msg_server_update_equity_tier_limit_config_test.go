package keeper_test

import (
	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateEquityTierLimitConfig(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		genesis := testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
			state.Subaccounts = []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_100_000USD,
			}
		})
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
			state.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_001_000_000), // $5,001
						Limit:          1,
					},
				},
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_002_000_000), // $5,002
						Limit:          2,
					},
				},
			}
		})
		return genesis
	}).Build()

	expectedConfig := clobtypes.EquityTierLimitConfiguration{
		ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
			{
				UsdTncRequired: dtypes.NewInt(0),
				Limit:          0,
			},
			{
				UsdTncRequired: dtypes.NewInt(5_003_000_000), // $5,003
				Limit:          3,
			},
		},
		StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
			{
				UsdTncRequired: dtypes.NewInt(0),
				Limit:          0,
			},
			{
				UsdTncRequired: dtypes.NewInt(5_004_000_000), // $5,004
				Limit:          4,
			},
		},
	}

	ctx := tApp.InitChain()

	originalConfig := tApp.App.ClobKeeper.GetEquityTierLimitConfiguration(ctx)
	require.NotEqual(t, expectedConfig, originalConfig)
	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateEquityTierLimitConfiguration{})

	requestWithoutAuthority := clobtypes.MsgUpdateEquityTierLimitConfiguration{
		Authority:             "fake authority",
		EquityTierLimitConfig: expectedConfig,
	}
	_, err := handler(ctx, &requestWithoutAuthority)
	require.Error(t, err, "invalid authority")
	require.Equal(t, originalConfig, tApp.App.ClobKeeper.GetEquityTierLimitConfiguration(ctx))

	requestWithAuthority := clobtypes.MsgUpdateEquityTierLimitConfiguration{
		Authority:             lib.GovModuleAddress.String(),
		EquityTierLimitConfig: expectedConfig,
	}
	_, err = handler(ctx, &requestWithAuthority)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, tApp.App.ClobKeeper.GetEquityTierLimitConfiguration(ctx))
}
