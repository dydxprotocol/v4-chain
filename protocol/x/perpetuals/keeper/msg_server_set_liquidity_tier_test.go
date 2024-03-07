package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	lttest "github.com/dydxprotocol/v4-chain/protocol/testutil/liquidity_tier"
	perpkeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestSetLiquidityTier(t *testing.T) {
	testLt := *lttest.GenerateLiquidityTier(
		lttest.WithId(1),
		lttest.WithName("test"),
		lttest.WithInitialMarginPpm(1_000),
		lttest.WithMaintenanceFractionPpm(2_000),
		lttest.WithImpactNotional(4_000),
	)

	tests := map[string]struct {
		msg         *types.MsgSetLiquidityTier
		expectedErr string
	}{
		"Success: update name and initial margin ppm": {
			msg: &types.MsgSetLiquidityTier{
				Authority: lib.GovModuleAddress.String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "large-cap",
					InitialMarginPpm:       123_432,
					MaintenanceFractionPpm: testLt.MaintenanceFractionPpm,
					ImpactNotional:         testLt.ImpactNotional,
				},
			},
		},
		"Success: update all parameters": {
			msg: &types.MsgSetLiquidityTier{
				Authority: lib.GovModuleAddress.String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       567_123,
					MaintenanceFractionPpm: 500_001,
					ImpactNotional:         1_300_303,
				},
			},
		},
		"Success: create a new liquidity tier": {
			msg: &types.MsgSetLiquidityTier{
				Authority: lib.GovModuleAddress.String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id + 1,
					Name:                   "medium-cap",
					InitialMarginPpm:       567_123,
					MaintenanceFractionPpm: 500_001,
					ImpactNotional:         1_300_303,
				},
			},
		},
		"Failure: initial margin ppm exceeds max": {
			msg: &types.MsgSetLiquidityTier{
				Authority: lib.GovModuleAddress.String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       1_000_001,
					MaintenanceFractionPpm: 500_001,
					ImpactNotional:         1_300_303,
				},
			},
			expectedErr: "InitialMarginPpm exceeds maximum value",
		},
		"Failure: maintenance fraction ppm exceeds max": {
			msg: &types.MsgSetLiquidityTier{
				Authority: lib.GovModuleAddress.String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       500_001,
					MaintenanceFractionPpm: 1_000_001,
					ImpactNotional:         1_300_303,
				},
			},
			expectedErr: "MaintenanceFractionPpm exceeds maximum value",
		},
		"Failure: invalid authority": {
			msg: &types.MsgSetLiquidityTier{
				Authority: constants.BobAccAddress.String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       567_123,
					MaintenanceFractionPpm: 500_001,
					ImpactNotional:         1_300_303,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure: empty authority": {
			msg: &types.MsgSetLiquidityTier{
				Authority: "",
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       567_123,
					MaintenanceFractionPpm: 500_001,
					ImpactNotional:         1_300_303,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure: invalid open interest caps": {
			msg: &types.MsgSetLiquidityTier{
				Authority: lib.GovModuleAddress.String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       567_123,
					MaintenanceFractionPpm: 500_001,
					ImpactNotional:         1_300_303,
					OpenInterestLowerCap:   100,
					OpenInterestUpperCap:   50,
				},
			},
			expectedErr: "open interest lower cap is larger than upper cap",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pc := keepertest.PerpetualsKeepers(t)
			initialLt, err := pc.PerpetualsKeeper.SetLiquidityTier(
				pc.Ctx,
				testLt.Id,
				testLt.Name,
				testLt.InitialMarginPpm,
				testLt.MaintenanceFractionPpm,
				testLt.ImpactNotional,
				testLt.OpenInterestLowerCap,
				testLt.OpenInterestUpperCap,
			)
			require.NoError(t, err)

			msgServer := perpkeeper.NewMsgServerImpl(pc.PerpetualsKeeper)

			_, err = msgServer.SetLiquidityTier(pc.Ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				// Verify that liquidity tier is same as before.
				lt, err := pc.PerpetualsKeeper.GetLiquidityTier(pc.Ctx, tc.msg.LiquidityTier.Id)
				require.NoError(t, err)
				require.Equal(t, initialLt, lt)
			} else {
				require.NoError(t, err)

				// Verify that liquidity tier is updated.
				lt, err := pc.PerpetualsKeeper.GetLiquidityTier(pc.Ctx, tc.msg.LiquidityTier.Id)
				require.NoError(t, err)
				require.Equal(t, tc.msg.LiquidityTier, lt)
			}
		})
	}
}
