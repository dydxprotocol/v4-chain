//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
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

func TestQueryDowntimeParams(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryDowntimeParams(), []string{})

	require.NoError(t, err)
	var resp types.QueryDowntimeParamsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

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
