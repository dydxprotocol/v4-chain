package cli_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func setupNetwork(t *testing.T) (*network.Network, client.Context) {
	t.Helper()
	cfg := network.DefaultConfig(nil)

	// Init state.
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	// Modify default genesis state
	state = *types.DefaultGenesis()

	// Add test affiliate tiers
	state.AffiliateTiers = types.DefaultAffiliateTiers

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	net := network.New(t, cfg)
	ctx := net.Validators[0].ClientCtx

	return net, ctx
}

func TestQueryAffiliateTiers(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryAffiliateTiers(), []string{})
	require.NoError(t, err)

	var resp types.AllAffiliateTiersResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.DefaultAffiliateTiers, resp.Tiers)
}

func TestQueryAffiliateInfo(t *testing.T) {
	net, ctx := setupNetwork(t)

	testAddress := constants.AliceAccAddress.String()
	out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryAffiliateInfo(), []string{testAddress})
	require.NoError(t, err)

	var resp types.AffiliateInfoResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}

func TestQueryReferredBy(t *testing.T) {
	net, ctx := setupNetwork(t)

	testAddress := constants.AliceAccAddress.String()
	out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryReferredBy(), []string{testAddress})
	require.NoError(t, err)

	var resp types.ReferredByResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}

func TestQueryAffiliateWhitelist(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryAffiliateWhitelist(), []string{})
	require.NoError(t, err)

	var resp types.AffiliateWhitelistResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
}
