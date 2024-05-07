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
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/client/cli"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
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

// func TestQueryDowntimeParams(t *testing.T) {

// 	fmt.Println(cli.CmdQueryDowntimeParams())
// 	fmt.Println(types.DefaultGenesis().Params)

// 	net, ctx := setupNetwork(t)

// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryDowntimeParams(), []string{})

// 	require.NoError(t, err)
// 	var resp types.QueryDowntimeParamsResponse
// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 	require.Equal(t, types.DefaultGenesis().Params, resp.Params)

// 	fmt.Println(ctx)
// 	fmt.Println(cli.CmdQueryDowntimeParams())
// 	fmt.Println(resp)
// 	fmt.Println(resp.Params)
// }

func TestQueryAllDowntimeInfo(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryAllDowntimeInfo(), []string{})

	require.NoError(t, err)
	var resp types.QueryAllDowntimeInfoResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}

func TestQueryPreviousBlockInfo(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryPreviousBlockInfo(), []string{})

	require.NoError(t, err)
	var resp types.QueryPreviousBlockInfoResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}
