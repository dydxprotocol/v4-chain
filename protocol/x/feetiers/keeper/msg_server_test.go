package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
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

// TestMsgSetFeeDiscountCampaignParams tests the SetFeeDiscountCampaignParams message handler
func TestMsgSetFeeDiscountCampaignParams(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	// Set a fixed current time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Create valid fee discount campaign params
	validParams := types.FeeDiscountCampaignParams{
		ClobPairId:    1,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
		ChargePpm:     500_000, // 50% discount
	}

	testCases := []struct {
		name      string
		input     *types.MsgSetFeeDiscountCampaignParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid single param",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    []types.FeeDiscountCampaignParams{validParams},
			},
			expErr: false,
		},
		{
			name: "valid multiple params",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeDiscountCampaignParams{
					validParams,
					{
						ClobPairId:    2,
						StartTimeUnix: 1100,
						EndTimeUnix:   1200,
						ChargePpm:     750_000, // 25% discount
					},
				},
			},
			expErr: false,
		},
		{
			name: "empty params",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    []types.FeeDiscountCampaignParams{},
			},
			expErr: false, // Empty list is valid (no-op)
		},
		{
			name: "invalid authority",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: "invalid",
				Params:    []types.FeeDiscountCampaignParams{validParams},
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "invalid param - end time before current time",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeDiscountCampaignParams{
					{
						ClobPairId:    1,
						StartTimeUnix: 900,
						EndTimeUnix:   950, // Before current time (1000)
						ChargePpm:     500_000,
					},
				},
			},
			expErr:    true,
			expErrMsg: "Invalid time range",
		},
		{
			name: "invalid param - start time after end time",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeDiscountCampaignParams{
					{
						ClobPairId:    1,
						StartTimeUnix: 1200,
						EndTimeUnix:   1100, // Before start time
						ChargePpm:     500_000,
					},
				},
			},
			expErr:    true,
			expErrMsg: "Invalid time range",
		},
		{
			name: "invalid param - too long duration",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeDiscountCampaignParams{
					{
						ClobPairId:    1,
						StartTimeUnix: 1100,
						EndTimeUnix:   1100 + 91*24*60*60, // 91 days
						ChargePpm:     500_000,
					},
				},
			},
			expErr:    true,
			expErrMsg: "Invalid time range",
		},
		{
			name: "invalid param - charge PPM exceeds maximum",
			input: &types.MsgSetFeeDiscountCampaignParams{
				Authority: lib.GovModuleAddress.String(),
				Params: []types.FeeDiscountCampaignParams{
					{
						ClobPairId:    1,
						StartTimeUnix: 1100,
						EndTimeUnix:   1200,
						ChargePpm:     1_500_000, // 150% - exceeds maximum
					},
				},
			},
			expErr:    true,
			expErrMsg: "Invalid charge PPM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.SetFeeDiscountCampaignParams(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestMsgSetFeeDiscountCampaignParamsUpdate tests updating existing fee discount campaign params
func TestMsgSetFeeDiscountCampaignParamsUpdate(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	// Set a fixed current time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Initial fee discount campaign params
	clobPairId := uint32(1)
	initialParams := types.FeeDiscountCampaignParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
		ChargePpm:     500_000, // 50% discount
	}

	// Set the initial fee discount campaign
	_, err := ms.SetFeeDiscountCampaignParams(ctx, &types.MsgSetFeeDiscountCampaignParams{
		Authority: lib.GovModuleAddress.String(),
		Params:    []types.FeeDiscountCampaignParams{initialParams},
	})
	require.NoError(t, err)

	// Verify it was set correctly
	getParams, err := k.GetFeeDiscountCampaignParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, initialParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, initialParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, initialParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, initialParams.ChargePpm, getParams.ChargePpm)

	// Update with new fee discount campaign params
	updatedParams := types.FeeDiscountCampaignParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1150,
		EndTimeUnix:   1250,
		ChargePpm:     250_000, // 75% discount
	}

	// Set the updated fee discount campaign
	_, err = ms.SetFeeDiscountCampaignParams(ctx, &types.MsgSetFeeDiscountCampaignParams{
		Authority: lib.GovModuleAddress.String(),
		Params:    []types.FeeDiscountCampaignParams{updatedParams},
	})
	require.NoError(t, err)

	// Verify it was updated correctly
	getParams, err = k.GetFeeDiscountCampaignParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, updatedParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, updatedParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, updatedParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, updatedParams.ChargePpm, getParams.ChargePpm)
}

// TestMsgSetMultipleFeeDiscountCampaignParams tests setting multiple fee discount campaign params
func TestMsgSetMultipleFeeDiscountCampaignParams(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	// Set a fixed current time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Multiple fee discount campaign params
	campaigns := []types.FeeDiscountCampaignParams{
		{
			ClobPairId:    1,
			StartTimeUnix: 1100,
			EndTimeUnix:   1200,
			ChargePpm:     0, // 100% discount (free)
		},
		{
			ClobPairId:    2,
			StartTimeUnix: 1150,
			EndTimeUnix:   1250,
			ChargePpm:     500_000, // 50% discount
		},
		{
			ClobPairId:    3,
			StartTimeUnix: 1200,
			EndTimeUnix:   1300,
			ChargePpm:     750_000, // 25% discount
		},
	}

	// Set all fee discount campaigns
	_, err := ms.SetFeeDiscountCampaignParams(ctx, &types.MsgSetFeeDiscountCampaignParams{
		Authority: lib.GovModuleAddress.String(),
		Params:    campaigns,
	})
	require.NoError(t, err)

	// Verify all campaigns were set correctly
	for _, campaign := range campaigns {
		getParams, err := k.GetFeeDiscountCampaignParams(ctx, campaign.ClobPairId)
		require.NoError(t, err)
		require.Equal(t, campaign.ClobPairId, getParams.ClobPairId)
		require.Equal(t, campaign.StartTimeUnix, getParams.StartTimeUnix)
		require.Equal(t, campaign.EndTimeUnix, getParams.EndTimeUnix)
		require.Equal(t, campaign.ChargePpm, getParams.ChargePpm)
	}

	// Check total count of fee discount campaigns
	allCampaigns := k.GetAllFeeDiscountCampaignParams(ctx)
	require.Len(t, allCampaigns, len(campaigns))
}
