//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestPendingSendPackets(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	cmd := exec.Command("docker", "exec", "interchain-security-instance", "interchain-security-cd", "query", "ratelimit", "pending-send-packets", param, "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.QueryAllPendingSendPacketsResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	assert.Equal(t, 0, len(resp.PendingSendPackets))
}

// func TestPendingSendPackets(t *testing.T) {
// 	net, ctx := setupNetwork(t)

// 	out, err := clitestutil.ExecTestCLICmd(ctx,
// 		cli.CmdPendingSendPackets(),
// 		[]string{
// 			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
// 		})

// 	require.NoError(t, err)
// 	var resp types.QueryAllPendingSendPacketsResponse
// 	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
// 	assert.Equal(t, 0, len(resp.PendingSendPackets))
// }
