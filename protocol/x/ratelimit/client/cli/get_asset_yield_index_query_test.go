//go:build all || integration_test

package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize
var defaultAssetYieldIndex = "1/1"

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

func TestGetAssetYieldIndexQuery(t *testing.T) {
	net, ctx := setupNetwork(t)

	common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetAssetYieldIndexQuery(), common)
	require.NoError(t, err)

	var resp types.GetAssetYieldIndexQueryResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, defaultAssetYieldIndex, resp.AssetYieldIndex)
}
