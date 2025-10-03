package keeper_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetFeeHolidayParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set the fee holiday params for a CLOB pair
	clobPairId := uint32(42)

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	setParams := types.FeeHolidayParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
	}

	err := k.SetFeeHolidayParams(ctx, setParams)
	require.NoError(t, err)

	// Get the fee holiday params for the CLOB pair
	getParams, err := k.GetFeeHolidayParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, setParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, setParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, setParams.EndTimeUnix, getParams.EndTimeUnix)
}

func TestGetFeeHolidayParamsNotFound(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Get the fee holiday params for a non-existent CLOB pair
	_, err := k.GetFeeHolidayParams(ctx, 42)
	require.ErrorIs(t, err, types.ErrFeeHolidayNotFound)
}

func TestGetAllFeeHolidayParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Set up multiple fee holidays
	holidays := []types.FeeHolidayParams{
		{
			ClobPairId:    1,
			StartTimeUnix: 1100,
			EndTimeUnix:   1200,
		},
		{
			ClobPairId:    2,
			StartTimeUnix: 1100,
			EndTimeUnix:   1300,
		},
		{
			ClobPairId:    3,
			StartTimeUnix: 1200,
			EndTimeUnix:   1400,
		},
	}

	// Store the fee holidays
	for _, holiday := range holidays {
		err := k.SetFeeHolidayParams(ctx, holiday)
		require.NoError(t, err)
	}

	// Get all fee holidays
	allHolidays := k.GetAllFeeHolidayParams(ctx)

	// Check that we got all the expected holidays
	require.Len(t, allHolidays, len(holidays))

	// Create a map of CLOB pair IDs to holidays for easier checking
	holidayMap := make(map[uint32]types.FeeHolidayParams)
	for _, holiday := range allHolidays {
		holidayMap[holiday.ClobPairId] = holiday
	}

	// Check each expected holiday is in the map
	for _, expectedHoliday := range holidays {
		holiday, found := holidayMap[expectedHoliday.ClobPairId]
		require.True(t, found)
		require.Equal(t, expectedHoliday.ClobPairId, holiday.ClobPairId)
		require.Equal(t, expectedHoliday.StartTimeUnix, holiday.StartTimeUnix)
		require.Equal(t, expectedHoliday.EndTimeUnix, holiday.EndTimeUnix)
	}
}

func TestIsFeeHolidayActive(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	tests := []struct {
		name           string
		setupTime      int64
		setupHoliday   *types.FeeHolidayParams
		checkTime      int64
		expectedActive bool
	}{
		{
			name:           "holiday not found",
			setupTime:      1000,
			checkTime:      1100,
			expectedActive: false,
		},
		{
			name:      "current time before start time",
			setupTime: 1000,
			setupHoliday: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
			},
			checkTime:      1050,
			expectedActive: false,
		},
		{
			name:      "current time at start time",
			setupTime: 1000,
			setupHoliday: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
			},
			checkTime:      1100,
			expectedActive: true,
		},
		{
			name:      "current time between start and end time",
			setupTime: 1000,
			setupHoliday: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
			},
			checkTime:      1150,
			expectedActive: true,
		},
		{
			name:      "current time at end time",
			setupTime: 1000,
			setupHoliday: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
			},
			checkTime:      1200,
			expectedActive: false,
		},
		{
			name:      "current time after end time",
			setupTime: 1000,
			setupHoliday: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
			},
			checkTime:      1250,
			expectedActive: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the test
			setupCtx := ctx.WithBlockTime(time.Unix(tc.setupTime, 0))
			clobPairId := uint32(1)

			// If there's a holiday to set up, do it
			if tc.setupHoliday != nil {
				err := k.SetFeeHolidayParams(setupCtx, *tc.setupHoliday)
				require.NoError(t, err)
			}

			// Create a context with the check time
			checkCtx := ctx.WithBlockTime(time.Unix(tc.checkTime, 0))

			// Check if the holiday is active
			isActive := k.IsFeeHolidayActive(checkCtx, clobPairId)
			require.Equal(t, tc.expectedActive, isActive)
		})
	}
}

func TestSetFeeHolidayParamsUpdate(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Initial fee holiday
	clobPairId := uint32(1)
	initialParams := types.FeeHolidayParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
	}

	// Set the initial fee holiday
	err := k.SetFeeHolidayParams(ctx, initialParams)
	require.NoError(t, err)

	// Verify it was set correctly
	getParams, err := k.GetFeeHolidayParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, initialParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, initialParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, initialParams.EndTimeUnix, getParams.EndTimeUnix)

	// Update with new fee holiday
	updatedParams := types.FeeHolidayParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1150,
		EndTimeUnix:   1250,
	}

	// Set the updated fee holiday
	err = k.SetFeeHolidayParams(ctx, updatedParams)
	require.NoError(t, err)

	// Verify it was updated correctly
	getParams, err = k.GetFeeHolidayParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, updatedParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, updatedParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, updatedParams.EndTimeUnix, getParams.EndTimeUnix)
}

func TestEmptyGetAllFeeHolidayParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Get all fee holidays when none exist
	allHolidays := k.GetAllFeeHolidayParams(ctx)

	// Check that we got an empty slice, not nil
	require.NotNil(t, allHolidays)
	require.Len(t, allHolidays, 0)
}
