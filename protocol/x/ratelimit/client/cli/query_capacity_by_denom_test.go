//go:build all || integration_test

package cli_test

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	ratelimitutil "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/util"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryCapacityByDenom(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	rateQuery := "docker exec interchain-security-instance-setup interchain-security-cd query ratelimit capacity-by-denom " + assettypes.AssetUsdc.Denom
	data, _, err := network.QueryCustomNetwork(rateQuery)

	require.NoError(t, err)
	var resp types.QueryCapacityByDenomResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t,
		// LimiterCapacity resulting from default limiter params and 0 TVL.
		[]types.LimiterCapacity{
			{
				Limiter: types.DefaultUsdcHourlyLimter,
				Capacity: dtypes.NewIntFromBigInt(
					ratelimitutil.GetBaseline(
						big.NewInt(0),
						types.DefaultUsdcHourlyLimter,
					),
				),
			},
			{
				Limiter: types.DefaultUsdcDailyLimiter,
				Capacity: dtypes.NewIntFromBigInt(
					ratelimitutil.GetBaseline(
						big.NewInt(0),
						types.DefaultUsdcDailyLimiter,
					),
				),
			},
		},
		resp.LimiterCapacityList)
}
