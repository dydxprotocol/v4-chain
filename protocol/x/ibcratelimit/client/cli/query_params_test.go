//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/x/ibcratelimit/types"
	"github.com/stretchr/testify/require"
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

func TestQueryParams(t *testing.T) {
	// TODO(CORE-823): implement query for `x/ibcratelimit`
}
