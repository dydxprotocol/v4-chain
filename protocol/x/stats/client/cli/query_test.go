//go:build all || integration_test

package cli_test

import (
	"bytes"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "stats", "get-params", "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryParamsResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

func TestQueryStatsMetadata(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "stats", "get-stats-metadata", "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryStatsMetadataResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}

func TestQueryGlobalStats(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "stats", "get-global-stats", "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryGlobalStatsResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}

func TestQueryUserStats(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "stats", "get-user-stats", "alice", "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryUserStatsResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}
