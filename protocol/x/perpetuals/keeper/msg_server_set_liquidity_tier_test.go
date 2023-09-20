package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
		lttest.WithBasePositionNotional(3_000),
		lttest.WithImpactNotional(4_000),
	)

	tests := map[string]struct {
		msg         *types.MsgSetLiquidityTier
		expectedErr string
	}{
		"Success: update name and initial margin ppm": {
			msg: &types.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "large-cap",
					InitialMarginPpm:       123_432,
					MaintenanceFractionPpm: testLt.MaintenanceFractionPpm,
					BasePositionNotional:   testLt.BasePositionNotional,
					ImpactNotional:         testLt.ImpactNotional,
				},
			},
		},
		"Success: update all parameters": {
			msg: &types.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       567_123,
					MaintenanceFractionPpm: 500_001,
					BasePositionNotional:   400_202,
					ImpactNotional:         1_300_303,
				},
			},
		},
		"Success: create a new liquidity tier": {
			msg: &types.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id + 1,
					Name:                   "medium-cap",
					InitialMarginPpm:       567_123,
					MaintenanceFractionPpm: 500_001,
					BasePositionNotional:   400_202,
					ImpactNotional:         1_300_303,
				},
			},
		},
		"Failure: initial margin ppm exceeds max": {
			msg: &types.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       1_000_001,
					MaintenanceFractionPpm: 500_001,
					BasePositionNotional:   400_202,
					ImpactNotional:         1_300_303,
				},
			},
			expectedErr: "InitialMarginPpm exceeds maximum value",
		},
		"Failure: maintenance fraction ppm exceeds max": {
			msg: &types.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: types.LiquidityTier{
					Id:                     testLt.Id,
					Name:                   "medium-cap",
					InitialMarginPpm:       500_001,
					MaintenanceFractionPpm: 1_000_001,
					BasePositionNotional:   400_202,
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
					BasePositionNotional:   400_202,
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
					BasePositionNotional:   400_202,
					ImpactNotional:         1_300_303,
				},
			},
			expectedErr: "invalid authority",
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
				testLt.BasePositionNotional,
				testLt.ImpactNotional,
			)
			require.NoError(t, err)

			msgServer := perpkeeper.NewMsgServerImpl(pc.PerpetualsKeeper)
			wrappedCtx := sdk.WrapSDKContext(pc.Ctx)

			_, err = msgServer.SetLiquidityTier(wrappedCtx, tc.msg)
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
