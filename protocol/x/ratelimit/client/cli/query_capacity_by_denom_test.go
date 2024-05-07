//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"math/big"
	"os/exec"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	ratelimitutil "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/util"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryCapacityByDenom(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	cmd := exec.Command("docker", "exec", "interchain-security-instance", "interchain-security-cd", "query", "ratelimit", "capacity-by-denom", param, assettypes.AssetUsdc.Denom, "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryCapacityByDenomResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.MarshalJSON(data, &resp))
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

// func TestQueryCapacityByDenom(t *testing.T) {
// 	net, ctx := setupNetwork(t)

// 	out, err := clitestutil.ExecTestCLICmd(ctx,
// 		cli.CmdQueryCapacityByDenom(),
// 		[]string{
// 			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
// 			assettypes.AssetUsdc.Denom,
// 		})

// 	require.NoError(t, err)
// 	var resp types.QueryCapacityByDenomResponse
// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 	require.Equal(t,
// 		// LimiterCapacity resulting from default limiter params and 0 TVL.
// 		[]types.LimiterCapacity{
// 			{
// 				Limiter: types.DefaultUsdcHourlyLimter,
// 				Capacity: dtypes.NewIntFromBigInt(
// 					ratelimitutil.GetBaseline(
// 						big.NewInt(0),
// 						types.DefaultUsdcHourlyLimter,
// 					),
// 				),
// 			},
// 			{
// 				Limiter: types.DefaultUsdcDailyLimiter,
// 				Capacity: dtypes.NewIntFromBigInt(
// 					ratelimitutil.GetBaseline(
// 						big.NewInt(0),
// 						types.DefaultUsdcDailyLimiter,
// 					),
// 				),
// 			},
// 		},
// 		resp.LimiterCapacityList)
// }
