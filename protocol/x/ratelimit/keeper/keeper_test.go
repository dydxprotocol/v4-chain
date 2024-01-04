package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

const (
	testDenom = "ibc/xxx"
)

func TestSetGetDenomCapacity(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	capacityList := []dtypes.SerializableInt{
		dtypes.NewInt(123_456_789),
		dtypes.NewInt(500_000_000),
	}
	denomCapacity := types.DenomCapacity{
		Denom:        testDenom,
		CapacityList: capacityList,
	}

	// Test SetDenomCapacity
	k.SetDenomCapacity(ctx, denomCapacity)

	// Test GetDenomCapacity
	gotDenomCapacity := k.GetDenomCapacity(ctx, testDenom)
	require.Equal(t, denomCapacity, gotDenomCapacity, "retrieved DenomCapacity does not match the set value")

	k.SetDenomCapacity(ctx, types.DenomCapacity{
		Denom:        testDenom,
		CapacityList: []dtypes.SerializableInt{}, // Empty list, results in deletion of the key.
	})

	// Check that the key is deleted under `DenomCapacity` storage.
	require.Equal(t,
		types.DenomCapacity{
			Denom:        testDenom,
			CapacityList: nil,
		},
		k.GetDenomCapacity(ctx, testDenom),
		"retrieved LimitParams do not match the set value",
	)
}

func TestSetGetLimitParams_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	limiters := []types.Limiter{
		{
			PeriodSec:       3_600,
			BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
			BaselineTvlPpm:  10_000,                         // 1%
		},
		{
			PeriodSec:       86_400,
			BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m tokens (assuming 6 decimals)
			BaselineTvlPpm:  100_000,                          // 10%
		},
	}
	limitParams := types.LimitParams{
		Denom:    testDenom,
		Limiters: limiters,
	}

	// Test SetLimitParams
	k.SetLimitParams(ctx, limitParams)

	// Test GetLimitParams
	gotLimitParams := k.GetLimitParams(ctx, testDenom)
	require.Equal(t, limitParams, gotLimitParams, "retrieved LimitParams do not match the set value")

	// Query for `DenomCapacity` of `testDenom`.
	gotDenomCapacity := k.GetDenomCapacity(ctx, testDenom)
	// Expected `DenomCapacity` is initialized such that each capacity is equal to the baseline.
	expectedDenomCapacity := types.DenomCapacity{
		Denom: testDenom,
		CapacityList: []dtypes.SerializableInt{
			// TODO(CORE-834): Update expected value after `GetBaseline` depends on current TVL.
			dtypes.NewInt(200_000_000_000),   // 200k tokens (assuming 6 decimals)
			dtypes.NewInt(2_000_000_000_000), // 1m tokens (assuming 6 decimals)
		},
	}
	require.Equal(t, expectedDenomCapacity, gotDenomCapacity, "retrieved DenomCapacity does not match the set value")

	// Set empty `LimitParams` for `testDenom`.
	k.SetLimitParams(ctx, types.LimitParams{
		Denom:    testDenom,
		Limiters: []types.Limiter{}, // Empty list, results in deletion of the key.
	})

	// Check that the key is deleted under `LimitParams` storage.
	require.Equal(t,
		types.LimitParams{
			Denom:    testDenom,
			Limiters: nil,
		},
		k.GetLimitParams(ctx, testDenom),
		"retrieved LimitParams do not match the set value")

	// Check that the key is deleted under `DenomCapacity` storage.
	require.Equal(t,
		types.DenomCapacity{
			Denom:        testDenom,
			CapacityList: nil,
		},
		k.GetDenomCapacity(ctx, testDenom),
		"retrieved LimitParams do not match the set value")
}

func TestGetBaseline(t *testing.T) {
	// TODO(CORE-836): Add test for GetBaseline.
}
