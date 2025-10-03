package keeper_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, sdk.Context) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	return *k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, k)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestMsgUpdateParams(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	testCases := []struct {
		name      string
		input     *types.MsgUpdatePerpetualFeeParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid params",
			input: &types.MsgUpdatePerpetualFeeParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    types.DefaultGenesis().Params,
			},
			expErr: false,
		},
		{
			name: "invalid authority",
			input: &types.MsgUpdatePerpetualFeeParams{
				Authority: "invalid",
				Params:    types.DefaultGenesis().Params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "invalid params: negative duration",
			input: &types.MsgUpdatePerpetualFeeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.PerpetualFeeParams{
					Tiers: []*types.PerpetualFeeTier{
						{TotalVolumeShareRequirementPpm: 1},
					},
				},
			},
			expErr:    true,
			expErrMsg: "First fee tier must not have volume requirements",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdatePerpetualFeeParams(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgSetFeeHolidayParams(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	// Set a fixed current time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Create valid fee holiday params
	validParams := types.FeeHolidayParams{
		ClobPairId:    1,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
	}

	testCases := []struct {
		name      string
		input     *types.MsgSetFeeHolidayParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid single param",
			input: &types.MsgSetFeeHolidayParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    []types.FeeHolidayParams{validParams},
			},
			expErr: false,
		},
		{
			name: "valid multiple params",
			input: &types.MsgSetFeeHolidayParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeHolidayParams{
					validParams,
					{
						ClobPairId:    2,
						StartTimeUnix: 1100,
						EndTimeUnix:   1200,
					},
				},
			},
			expErr: false,
		},
		{
			name: "empty params",
			input: &types.MsgSetFeeHolidayParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    []types.FeeHolidayParams{},
			},
			expErr: false, // Empty list is valid (no-op)
		},
		{
			name: "invalid authority",
			input: &types.MsgSetFeeHolidayParams{
				Authority: "invalid",
				Params:    []types.FeeHolidayParams{validParams},
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "invalid param - end time before current time",
			input: &types.MsgSetFeeHolidayParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeHolidayParams{
					{
						ClobPairId:    1,
						StartTimeUnix: 900,
						EndTimeUnix:   950, // Before current time (1000)
					},
				},
			},
			expErr:    true,
			expErrMsg: "Invalid time range",
		},
		{
			name: "invalid param - start time after end time",
			input: &types.MsgSetFeeHolidayParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeHolidayParams{
					{
						ClobPairId:    1,
						StartTimeUnix: 1200,
						EndTimeUnix:   1100, // Before start time
					},
				},
			},
			expErr:    true,
			expErrMsg: "Invalid time range",
		},
		{
			name: "invalid param - too long duration",
			input: &types.MsgSetFeeHolidayParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeHolidayParams{
					{
						ClobPairId:    1,
						StartTimeUnix: 1100,
						EndTimeUnix:   1100 + 31*24*60*60, // 31 days
					},
				},
			},
			expErr:    true,
			expErrMsg: "Invalid time range",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.SetFeeHolidayParams(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgSetFeeHolidayParamsUpdate(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	// Set a fixed current time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Initial fee holiday params
	clobPairId := uint32(1)
	initialParams := types.FeeHolidayParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
	}

	// Set the initial fee holiday
	_, err := ms.SetFeeHolidayParams(ctx, &types.MsgSetFeeHolidayParams{
		Authority: lib.GovModuleAddress.String(),
		Params:    []types.FeeHolidayParams{initialParams},
	})
	require.NoError(t, err)

	// Verify it was set correctly
	getParams, err := k.GetFeeHolidayParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, initialParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, initialParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, initialParams.EndTimeUnix, getParams.EndTimeUnix)

	// Update with new fee holiday params
	updatedParams := types.FeeHolidayParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1150,
		EndTimeUnix:   1250,
	}

	// Set the updated fee holiday
	_, err = ms.SetFeeHolidayParams(ctx, &types.MsgSetFeeHolidayParams{
		Authority: lib.GovModuleAddress.String(),
		Params:    []types.FeeHolidayParams{updatedParams},
	})
	require.NoError(t, err)

	// Verify it was updated correctly
	getParams, err = k.GetFeeHolidayParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, updatedParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, updatedParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, updatedParams.EndTimeUnix, getParams.EndTimeUnix)
}

func TestMsgSetMultipleFeeHolidayParams(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	// Set a fixed current time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Multiple fee holiday params
	holidays := []types.FeeHolidayParams{
		{
			ClobPairId:    1,
			StartTimeUnix: 1100,
			EndTimeUnix:   1200,
		},
		{
			ClobPairId:    2,
			StartTimeUnix: 1150,
			EndTimeUnix:   1250,
		},
		{
			ClobPairId:    3,
			StartTimeUnix: 1200,
			EndTimeUnix:   1300,
		},
	}

	// Set all fee holidays
	_, err := ms.SetFeeHolidayParams(ctx, &types.MsgSetFeeHolidayParams{
		Authority: lib.GovModuleAddress.String(),
		Params:    holidays,
	})
	require.NoError(t, err)

	// Verify all holidays were set correctly
	for _, holiday := range holidays {
		getParams, err := k.GetFeeHolidayParams(ctx, holiday.ClobPairId)
		require.NoError(t, err)
		require.Equal(t, holiday.ClobPairId, getParams.ClobPairId)
		require.Equal(t, holiday.StartTimeUnix, getParams.StartTimeUnix)
		require.Equal(t, holiday.EndTimeUnix, getParams.EndTimeUnix)
	}

	// Check total count of fee holidays
	allHolidays := k.GetAllFeeHolidayParams(ctx)
	require.Len(t, allHolidays, len(holidays))
}
