//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/vest/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"
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

func TestQueryVestEntry(t *testing.T) {

	// net, ctx := setupNetwork(t)

	queryAndCheckVestEntry(t, "treausry_vester", types.DefaultGenesis().VestEntries[0])
	queryAndCheckVestEntry(t, "rewards_vester", types.DefaultGenesis().VestEntries[1])
}

func queryAndCheckVestEntry(
	t *testing.T,
	vester_account string,
	expectedEntry types.VestEntry,
) {

	cfg := network.DefaultConfig(nil)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	cmd := exec.Command("docker", "exec", "interchain-security-instance", "interchain-security-cd", "query", "vest", "vest-entry", vester_account, param, "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryVestEntryResponse
	outBytes := out.Bytes()
	require.NoError(t, cfg.Codec.MarshalJSON(outBytes, &resp))
	require.Equal(t, types.DefaultGenesis().VestEntries[1], resp.Entry)
}

// func queryAndCheckVestEntry(
// 	t *testing.T,
// 	ctx client.Context,
// 	net *network.Network,
// 	vester_account string,
// 	expectedEntry types.VestEntry,
// ) {
// 	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryVestEntry(), []string{
// 		"rewards_vester",
// 		fmt.Sprintf("--%s=json", tmcli.OutputFlag), // specify output format as json
// 	})

// 	require.NoError(t, err)
// 	var resp types.QueryVestEntryResponse
// 	outBytes := out.Bytes()
// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(outBytes, &resp))
// 	require.Equal(t, types.DefaultGenesis().VestEntries[1], resp.Entry)
// }
