package cli_test

import (
	"fmt"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/dydxprotocol/v4-chain/protocol/app/stoppable"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
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

	t.Cleanup(func() {
		stoppable.StopServices(t, cfg.GRPCAddress)
	})

	return net, ctx
}

func TestQueryVestEntry(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryVestEntry(), []string{
		"rewards_vester",
		fmt.Sprintf("--%s=json", tmcli.OutputFlag), // specify output format as json
	})

	require.NoError(t, err)
	var resp types.QueryVestEntryResponse
	outBytes := out.Bytes()
	require.NoError(t, net.Config.Codec.UnmarshalJSON(outBytes, &resp))
	require.Equal(t, types.DefaultGenesis().VestEntries[0], resp.Entry)
}
