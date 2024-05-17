//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize
func TestQueryParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "blocktime", "get-downtime-params", "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	fmt.Println("Output from Docker command:", out.String())

	require.NoError(t, err)
	var resp types.QueryDowntimeParamsResponse
	data := out.Bytes()
	data2 := cfg.GenesisState[types.ModuleName]

	fmt.Println("Data from Docker command:", data)
	fmt.Println("Data from Genesis state:", data2)

	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

func TestQueryAllDowntimeInfo(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "blocktime", "get-all-downtime-info", "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryAllDowntimeInfoResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}

func TestQueryPreviousBlockInfo(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "blocktime", "get-previous-block-info", "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryPreviousBlockInfoResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}
