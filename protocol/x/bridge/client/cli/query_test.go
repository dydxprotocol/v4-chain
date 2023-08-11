//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4/testutil/network"
	"github.com/dydxprotocol/v4/x/bridge/client/cli"
	"github.com/dydxprotocol/v4/x/bridge/types"
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

func TestQueryEventParams(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryEventParams(), []string{})

	require.NoError(t, err)
	var resp types.QueryEventParamsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().EventParams, resp.Params)
}

func TestQueryProposeParams(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryProposeParams(), []string{})

	require.NoError(t, err)
	var resp types.QueryProposeParamsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().ProposeParams, resp.Params)
}

func TestQuerySafetyParams(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQuerySafetyParams(), []string{})

	require.NoError(t, err)
	var resp types.QuerySafetyParamsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().SafetyParams, resp.Params)
}

func TestQueryNextAcknowledgedEventId(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryNextAcknowledgedEventId(), []string{})

	require.NoError(t, err)
	var resp types.QueryNextAcknowledgedEventIdResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().NextAcknowledgedEventId, resp.Id)
}
