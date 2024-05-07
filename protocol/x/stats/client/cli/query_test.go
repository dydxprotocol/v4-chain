//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/client/cli"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/types"
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

// func TestQueryParams(t *testing.T) {
// 	net, ctx := setupNetwork(t)

// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryParams(), []string{})

// 	require.NoError(t, err)
// 	var resp types.QueryParamsResponse
// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
// }

func TestQueryParams(t *testing.T) {
	cmd := exec.Command("/usr/local/bin/docker", "exec", "interchain-security-instance", "interchain-security-cd", "query", "blocktime", "get-downtime-params", "--node", "tcp://7.7.8.4:26658")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("error executing Docker command: %s", err)
	}

	// Parse the output
	// Assuming the output is in YAML format
	// You can use a YAML parser to parse it
	fmt.Println(out.String())
	// Example of parsing YAML
	// var result YourStruct
	// err = yaml.Unmarshal(out.Bytes(), &result)
	// if err != nil {
	// 	t.Fatalf("error parsing output: %s", err)
	// }

	// Your assertion logic
	// Example of assertions
	// require.Equal(t, expectedValue, result.SomeField)
}

func TestQueryParamsDocker(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryParams(), []string{})

	require.NoError(t, err)
	var resp types.QueryParamsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

func TestQueryStatsMetadata(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryStatsMetadata(), []string{})

	require.NoError(t, err)
	var resp types.QueryStatsMetadataResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}

func TestQueryGlobalStats(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryGlobalStats(), []string{})

	require.NoError(t, err)
	var resp types.QueryGlobalStatsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}

func TestQueryUserStats(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryUserStats(), []string{"alice"})

	require.NoError(t, err)
	var resp types.QueryUserStatsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}
