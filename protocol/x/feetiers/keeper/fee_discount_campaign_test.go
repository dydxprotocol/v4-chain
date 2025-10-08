package keeper_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetFeeDiscountCampaignParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set the fee discount campaign params for a CLOB pair
	clobPairId := uint32(42)

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	setParams := types.FeeDiscountCampaignParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
		ChargePpm:     500_000, // 50% discount
	}

	err := k.SetFeeDiscountCampaignParams(ctx, setParams)
	require.NoError(t, err)

	// Get the fee discount campaign params for the CLOB pair
	getParams, err := k.GetFeeDiscountCampaignParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, setParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, setParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, setParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, setParams.ChargePpm, getParams.ChargePpm)
}

func TestGetFeeDiscountCampaignParamsNotFound(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Get the fee discount campaign params for a non-existent CLOB pair
	_, err := k.GetFeeDiscountCampaignParams(ctx, 42)
	require.ErrorIs(t, err, types.ErrFeeDiscountCampaignNotFound)
}

func TestGetAllFeeDiscountCampaignParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Set up multiple fee discount campaigns
	campaigns := []types.FeeDiscountCampaignParams{
		{
			ClobPairId:    1,
			StartTimeUnix: 1100,
			EndTimeUnix:   1200,
			ChargePpm:     0, // 100% discount (free)
		},
		{
			ClobPairId:    2,
			StartTimeUnix: 1100,
			EndTimeUnix:   1300,
			ChargePpm:     500_000, // 50% discount
		},
		{
			ClobPairId:    3,
			StartTimeUnix: 1200,
			EndTimeUnix:   1400,
			ChargePpm:     750_000, // 25% discount
		},
	}

	// Store the fee discount campaigns
	for _, campaign := range campaigns {
		err := k.SetFeeDiscountCampaignParams(ctx, campaign)
		require.NoError(t, err)
	}

	// Get all fee discount campaigns
	allCampaigns := k.GetAllFeeDiscountCampaignParams(ctx)

	// Check that we got all the expected campaigns
	require.Len(t, allCampaigns, len(campaigns))

	// Create a map of CLOB pair IDs to campaigns for easier checking
	campaignMap := make(map[uint32]types.FeeDiscountCampaignParams)
	for _, campaign := range allCampaigns {
		campaignMap[campaign.ClobPairId] = campaign
	}

	// Check each expected campaign is in the map
	for _, expectedCampaign := range campaigns {
		campaign, found := campaignMap[expectedCampaign.ClobPairId]
		require.True(t, found)
		require.Equal(t, expectedCampaign.ClobPairId, campaign.ClobPairId)
		require.Equal(t, expectedCampaign.StartTimeUnix, campaign.StartTimeUnix)
		require.Equal(t, expectedCampaign.EndTimeUnix, campaign.EndTimeUnix)
		require.Equal(t, expectedCampaign.ChargePpm, campaign.ChargePpm)
	}
}

func TestGetDiscountPpm(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	tests := []struct {
		name              string
		setupTime         int64
		setupCampaign     *types.FeeDiscountCampaignParams
		checkTime         int64
		expectedChargePpm uint32
	}{
		{
			name:              "campaign not found",
			setupTime:         1000,
			checkTime:         1100,
			expectedChargePpm: types.MaxChargePpm, // 100% charge (no discount)
		},
		{
			name:      "current time before start time",
			setupTime: 1000,
			setupCampaign: &types.FeeDiscountCampaignParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
				ChargePpm:     500_000, // 50% discount
			},
			checkTime:         1050,
			expectedChargePpm: types.MaxChargePpm, // 100% charge (no discount)
		},
		{
			name:      "current time at start time",
			setupTime: 1000,
			setupCampaign: &types.FeeDiscountCampaignParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
				ChargePpm:     500_000, // 50% discount
			},
			checkTime:         1100,
			expectedChargePpm: 500_000, // 50% discount
		},
		{
			name:      "current time between start and end time",
			setupTime: 1000,
			setupCampaign: &types.FeeDiscountCampaignParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
				ChargePpm:     0, // 100% discount (free)
			},
			checkTime:         1150,
			expectedChargePpm: 0, // 100% discount (free)
		},
		{
			name:      "current time at end time",
			setupTime: 1000,
			setupCampaign: &types.FeeDiscountCampaignParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
				ChargePpm:     500_000, // 50% discount
			},
			checkTime:         1200,
			expectedChargePpm: types.MaxChargePpm, // 100% charge (no discount)
		},
		{
			name:      "current time after end time",
			setupTime: 1000,
			setupCampaign: &types.FeeDiscountCampaignParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
				ChargePpm:     500_000, // 50% discount
			},
			checkTime:         1250,
			expectedChargePpm: types.MaxChargePpm, // 100% charge (no discount)
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the test
			setupCtx := ctx.WithBlockTime(time.Unix(tc.setupTime, 0))
			clobPairId := uint32(1)

			// If there's a campaign to set up, do it
			if tc.setupCampaign != nil {
				err := k.SetFeeDiscountCampaignParams(setupCtx, *tc.setupCampaign)
				require.NoError(t, err)
			}

			// Create a context with the check time
			checkCtx := ctx.WithBlockTime(time.Unix(tc.checkTime, 0))

			// Get the discount PPM
			chargePpm := k.GetDiscountPpm(checkCtx, clobPairId)
			require.Equal(t, tc.expectedChargePpm, chargePpm)
		})
	}
}

func TestSetFeeDiscountCampaignParamsUpdate(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Initial fee discount campaign
	clobPairId := uint32(1)
	initialParams := types.FeeDiscountCampaignParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
		ChargePpm:     500_000, // 50% discount
	}

	// Set the initial fee discount campaign
	err := k.SetFeeDiscountCampaignParams(ctx, initialParams)
	require.NoError(t, err)

	// Verify it was set correctly
	getParams, err := k.GetFeeDiscountCampaignParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, initialParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, initialParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, initialParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, initialParams.ChargePpm, getParams.ChargePpm)

	// Update with new fee discount campaign
	updatedParams := types.FeeDiscountCampaignParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1150,
		EndTimeUnix:   1250,
		ChargePpm:     250_000, // 75% discount
	}

	// Set the updated fee discount campaign
	err = k.SetFeeDiscountCampaignParams(ctx, updatedParams)
	require.NoError(t, err)

	// Verify it was updated correctly
	getParams, err = k.GetFeeDiscountCampaignParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, updatedParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, updatedParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, updatedParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, updatedParams.ChargePpm, getParams.ChargePpm)
}

func TestEmptyGetAllFeeDiscountCampaignParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Get all fee discount campaigns when none exist
	allCampaigns := k.GetAllFeeDiscountCampaignParams(ctx)

	// Check that we got an empty slice, not nil
	require.NotNil(t, allCampaigns)
	require.Len(t, allCampaigns, 0)
}
