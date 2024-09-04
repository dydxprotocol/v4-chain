package types_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultSDaiRateLimitParams(t *testing.T) {
	defaultParams := types.DefaultSDaiRateLimitParams()
	bigLimitOneHour, worked := big.NewInt(0).SetString("1000000000000000000000000", 10)
	require.Equal(t, true, worked)
	limitOneHour := dtypes.NewIntFromBigInt(bigLimitOneHour)
	bigLimitOneDay, worked := big.NewInt(0).SetString("10000000000000000000000000", 10)
	require.Equal(t, true, worked)
	limitOneDay := dtypes.NewIntFromBigInt(bigLimitOneDay)
	require.Equal(t,
		types.LimitParams{
			Denom: types.SDaiDenom,
			Limiters: []types.Limiter{
				{
					Period:          3600 * time.Second,
					BaselineMinimum: limitOneHour,
					BaselineTvlPpm:  10_000,
				},
				{
					Period:          24 * time.Hour,
					BaselineMinimum: limitOneDay,
					BaselineTvlPpm:  100_000,
				},
			},
		},
		defaultParams,
	)
}

func TestLimitParamsValidate(t *testing.T) {
	validDenom := "uatom"
	validPeriod := time.Hour
	validBaselineMinimum := dtypes.NewIntFromBigInt(big.NewInt(1000000))
	validBaselineTvlPpm := uint32(500000)

	tests := []struct {
		name      string
		params    types.LimitParams
		expectErr bool
	}{
		{
			name: "valid params",
			params: types.LimitParams{
				Denom: validDenom,
				Limiters: []types.Limiter{
					{
						Period:          validPeriod,
						BaselineMinimum: validBaselineMinimum,
						BaselineTvlPpm:  validBaselineTvlPpm,
					},
				},
			},
			expectErr: false,
		},
		{
			name: "invalid denom",
			params: types.LimitParams{
				Denom: "invalid denom",
				Limiters: []types.Limiter{
					{
						Period:          validPeriod,
						BaselineMinimum: validBaselineMinimum,
						BaselineTvlPpm:  validBaselineTvlPpm,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "zero period",
			params: types.LimitParams{
				Denom: validDenom,
				Limiters: []types.Limiter{
					{
						Period:          0,
						BaselineMinimum: validBaselineMinimum,
						BaselineTvlPpm:  validBaselineTvlPpm,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "negative baseline minimum",
			params: types.LimitParams{
				Denom: validDenom,
				Limiters: []types.Limiter{
					{
						Period:          validPeriod,
						BaselineMinimum: dtypes.NewIntFromBigInt(big.NewInt(-1)),
						BaselineTvlPpm:  validBaselineTvlPpm,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "zero baseline minimum",
			params: types.LimitParams{
				Denom: validDenom,
				Limiters: []types.Limiter{
					{
						Period:          validPeriod,
						BaselineMinimum: dtypes.NewInt(0),
						BaselineTvlPpm:  validBaselineTvlPpm,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "zero baseline tvl ppm",
			params: types.LimitParams{
				Denom: validDenom,
				Limiters: []types.Limiter{
					{
						Period:          validPeriod,
						BaselineMinimum: validBaselineMinimum,
						BaselineTvlPpm:  0,
					},
				},
			},
			expectErr: true,
		},
		{
			name: "baseline tvl ppm exceeds one million",
			params: types.LimitParams{
				Denom: validDenom,
				Limiters: []types.Limiter{
					{
						Period:          validPeriod,
						BaselineMinimum: validBaselineMinimum,
						BaselineTvlPpm:  1000001,
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
