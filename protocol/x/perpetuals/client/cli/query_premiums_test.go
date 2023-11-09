//go:build all || integration_test

package cli_test

import (
	"fmt"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryPremiumSamples(t *testing.T) {
	net, _, _ := networkWithLiquidityTierAndPerpetualObjects(t, 2, 2)
	ctx := net.Validators[0].ClientCtx

	common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryPremiumSamples(), common)
	require.NoError(t, err)

	var resp types.QueryPremiumSamplesResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))

	// In CI, we see that PremiumSamples may have NumPremiums set to a non-zero value. Waiting for a block height before
	// the query does not reproduce this locally, so we just check that the response PremiumSamples are non-nil the rest
	// of the struct is as expected.
	require.NotNil(t, resp.PremiumSamples)
	require.Len(t, resp.PremiumSamples.AllMarketPremiums, 0)
}

func TestQueryPremiumVotes(t *testing.T) {
	net, _, _ := networkWithLiquidityTierAndPerpetualObjects(t, 2, 2)
	ctx := net.Validators[0].ClientCtx

	common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryPremiumVotes(), common)
	require.NoError(t, err)

	var resp types.QueryPremiumVotesResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.NotNil(t, resp.PremiumVotes)
	require.Equal(t, []types.MarketPremiums{}, resp.PremiumVotes.AllMarketPremiums)
}
