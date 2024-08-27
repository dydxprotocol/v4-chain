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
