//go:build all || integration_test

package cli_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestListLimiterParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	cmd := exec.Command("docker", "exec", "interchain-security-instance-setup", "interchain-security-cd", "query", "ratelimit", "list-limit-params", param, "--node", "tcp://7.7.8.4:26658", "-o json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	require.NoError(t, err)
	var resp types.ListLimitParamsResponse
	data := out.Bytes()
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().LimitParamsList, resp.LimitParamsList)
}
