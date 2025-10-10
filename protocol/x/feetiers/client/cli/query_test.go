//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
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

func TestQueryPerpetualFeeParams(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryPerpetualFeeParams(), []string{})

	require.NoError(t, err)
	var resp types.QueryPerpetualFeeParamsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

func TestQueryUserFeeTier(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdQueryUserFeeTier(),
		[]string{"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4"},
	)

	require.NoError(t, err)
	var resp types.QueryUserFeeTierResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}

func TestQueryFeeDiscountParams(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryMarketFeeDiscountParams(), []string{})
	require.NoError(t, err)
	var allResp types.QueryAllMarketFeeDiscountParamsResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &allResp))
}
