//go:build all || integration_test

package cli_test

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func setupNetwork(
	t *testing.T,
) (
	*network.Network,
	client.Context,
) {
	t.Helper()
	cfg := network.DefaultConfig(nil)

	// Init state.
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	state = *types.DefaultGenesis()

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	net := network.New(t, cfg)
	ctx := net.Validators[0].ClientCtx

	return net, ctx
}

func TestQueryPerpetualFeeParams(t *testing.T) {

	cmd := exec.Command("docker", "exec", "interchain-security-instance", "interchain-security-cd", "query", "feetiers", "get-perpetual-fee-params", "--node", "tcp://7.7.8.4:26658")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryPerpetualFeeParamsResponse
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

// func TestQueryPerpetualFeeParams(t *testing.T) {
// 	net, ctx := setupNetwork(t)

// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryPerpetualFeeParams(), []string{})

// 	require.NoError(t, err)
// 	var resp types.QueryPerpetualFeeParamsResponse
// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
// }

func TestQueryUserFeeTier(t *testing.T) {

	cmd := exec.Command("docker", "exec", "interchain-security-instance", "interchain-security-cd", "query", "feetiers", "get-user-fee-tier", "alice", "--node", "tcp://7.7.8.4:26658")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryUserFeeTierResponse
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
}

// func TestQueryUserFeeTier(t *testing.T) {
// 	net, ctx := setupNetwork(t)

// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryUserFeeTier(), []string{"alice"})

// 	require.NoError(t, err)
// 	var resp types.QueryUserFeeTierResponse
// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// }
