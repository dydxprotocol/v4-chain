package keeper_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetPerMarketFeeDiscountParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set the fee discount params for a CLOB pair
	clobPairId := uint32(42)

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	setParams := types.PerMarketFeeDiscountParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
		ChargePpm:     500_000, // 50% discount
	}

	err := k.SetPerMarketFeeDiscountParams(ctx, setParams)
	require.NoError(t, err)

	// Get the fee discount params for the CLOB pair
	getParams, err := k.GetPerMarketFeeDiscountParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, setParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, setParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, setParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, setParams.ChargePpm, getParams.ChargePpm)
}

func TestGetPerMarketFeeDiscountParamsNotFound(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Get the fee discount params for a non-existent CLOB pair
	_, err := k.GetPerMarketFeeDiscountParams(ctx, 42)
	require.ErrorIs(t, err, types.ErrMarketFeeDiscountNotFound)
}

func TestGetAllMarketsFeeDiscountParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Set up multiple fee discounts
	discountParams := []types.PerMarketFeeDiscountParams{
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

	// Store the fee discount params
	for _, params := range discountParams {
		err := k.SetPerMarketFeeDiscountParams(ctx, params)
		require.NoError(t, err)
	}

	// Get all fee discount params
	allDiscountParams := k.GetAllMarketFeeDiscountParams(ctx)

	// Check that we got all the expected discount params
	require.Len(t, allDiscountParams, len(discountParams))

	// Create a map of CLOB pair IDs to discount params for easier checking
	discountParamsMap := make(map[uint32]types.PerMarketFeeDiscountParams)
	for _, params := range allDiscountParams {
		discountParamsMap[params.ClobPairId] = params
	}

	// Check each expected discount params is in the map
	for _, expectedParams := range discountParams {
		params, found := discountParamsMap[expectedParams.ClobPairId]
		require.True(t, found)
		require.Equal(t, expectedParams.ClobPairId, params.ClobPairId)
		require.Equal(t, expectedParams.StartTimeUnix, params.StartTimeUnix)
		require.Equal(t, expectedParams.EndTimeUnix, params.EndTimeUnix)
		require.Equal(t, expectedParams.ChargePpm, params.ChargePpm)
	}
}

func TestGetDiscountPpm(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	tests := []struct {
		name              string
		setupTime         int64
		setupParams       *types.PerMarketFeeDiscountParams
		checkTime         int64
		expectedChargePpm uint32
	}{
		{
			name:              "discount not found",
			setupTime:         1000,
			checkTime:         1100,
			expectedChargePpm: types.MaxChargePpm, // 100% charge (no discount)
		},
		{
			name:      "current time before start time",
			setupTime: 1000,
			setupParams: &types.PerMarketFeeDiscountParams{
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
			setupParams: &types.PerMarketFeeDiscountParams{
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
			setupParams: &types.PerMarketFeeDiscountParams{
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
			setupParams: &types.PerMarketFeeDiscountParams{
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
			setupParams: &types.PerMarketFeeDiscountParams{
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

			// If there's a discount params to set up, do it
			if tc.setupParams != nil {
				err := k.SetPerMarketFeeDiscountParams(setupCtx, *tc.setupParams)
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

func TestSetPerMarketFeeDiscountParamsUpdate(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set current block time to a fixed time
	baseTime := time.Unix(1000, 0)
	ctx = ctx.WithBlockTime(baseTime)

	// Initial fee discount params
	clobPairId := uint32(1)
	initialParams := types.PerMarketFeeDiscountParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
		ChargePpm:     500_000, // 50% discount
	}

	// Set the initial fee discount params
	err := k.SetPerMarketFeeDiscountParams(ctx, initialParams)
	require.NoError(t, err)

	// Verify it was set correctly
	getParams, err := k.GetPerMarketFeeDiscountParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, initialParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, initialParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, initialParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, initialParams.ChargePpm, getParams.ChargePpm)

	// Update with new fee discount params
	updatedParams := types.PerMarketFeeDiscountParams{
		ClobPairId:    clobPairId,
		StartTimeUnix: 1150,
		EndTimeUnix:   1250,
		ChargePpm:     250_000, // 75% discount
	}

	// Set the updated fee discount params
	err = k.SetPerMarketFeeDiscountParams(ctx, updatedParams)
	require.NoError(t, err)

	// Verify it was updated correctly
	getParams, err = k.GetPerMarketFeeDiscountParams(ctx, clobPairId)
	require.NoError(t, err)
	require.Equal(t, updatedParams.ClobPairId, getParams.ClobPairId)
	require.Equal(t, updatedParams.StartTimeUnix, getParams.StartTimeUnix)
	require.Equal(t, updatedParams.EndTimeUnix, getParams.EndTimeUnix)
	require.Equal(t, updatedParams.ChargePpm, getParams.ChargePpm)
}

func TestEmptyGetAllMarketsFeeDiscountParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Get all fee discounts when none exist
	allDiscountParams := k.GetAllMarketFeeDiscountParams(ctx)

	// Check that we got an empty slice, not nil
	require.NotNil(t, allDiscountParams)
	require.Len(t, allDiscountParams, 0)
}
