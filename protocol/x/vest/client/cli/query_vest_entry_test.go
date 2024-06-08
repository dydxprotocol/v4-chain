//go:build all || integration_test

package cli_test

import (
	"fmt"
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

	vestQuery := "docker exec interchain-security-instance-setup interchain-security-cd query vest vest-entry " + vester_account + " " + param
	data, _, err := network.QueryCustomNetwork(vestQuery)

	require.NoError(t, err)
	var resp types.QueryVestEntryResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().VestEntries[1], resp.Entry)
}
