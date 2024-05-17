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
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryVestEntry(t *testing.T) {

	queryAndCheckVestEntry(t, "rewards_vester", types.DefaultGenesis().VestEntries[0])
	queryAndCheckVestEntry(t, "rewards_vester", types.DefaultGenesis().VestEntries[1])
}

func queryAndCheckVestEntry(
	t *testing.T,
	vester_account string,
	expectedEntry types.VestEntry,
) {

	cfg := network.DefaultConfig(nil)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "vest", "vest-entry", vester_account, param, "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryVestEntryResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().VestEntries[1], resp.Entry)
}
