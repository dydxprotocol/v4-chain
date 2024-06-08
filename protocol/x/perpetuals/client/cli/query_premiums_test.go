//go:build all || integration_test

package cli_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestQueryPremiumSamples(t *testing.T) {
	genesisChanges := GetPerpetualGenesisShort()
	network.DeployCustomNetwork(genesisChanges)
	cfg := network.DefaultConfig(nil)

	perpQuery := "docker exec interchain-security-instance-setup interchain-security-cd query perpetuals get-premium-samples"
	data, _, err := network.QueryCustomNetwork(perpQuery)
	require.NoError(t, err)

	var resp types.QueryPremiumSamplesResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))

	// In CI, we see that PremiumSamples may have NumPremiums set to a non-zero value. Waiting for a block height before
	// the query does not reproduce this locally, so we just check that the response PremiumSamples are non-nil the rest
	// of the struct is as expected.
	require.NotNil(t, resp.PremiumSamples)
	require.Len(t, resp.PremiumSamples.AllMarketPremiums, 0)
	network.CleanupCustomNetwork()
}

func TestQueryPremiumVotes(t *testing.T) {
	genesisChanges := GetPerpetualGenesisShort()
	network.DeployCustomNetwork(genesisChanges)
	cfg := network.DefaultConfig(nil)
	perpQuery := "docker exec interchain-security-instance-setup interchain-security-cd query perpetuals get-premium-votes"
	data, _, err := network.QueryCustomNetwork(perpQuery)

	require.NoError(t, err)

	var resp types.QueryPremiumVotesResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.NotNil(t, resp.PremiumVotes)
	require.Equal(t, []types.MarketPremiums{}, resp.PremiumVotes.AllMarketPremiums)
	network.CleanupCustomNetwork()

}
