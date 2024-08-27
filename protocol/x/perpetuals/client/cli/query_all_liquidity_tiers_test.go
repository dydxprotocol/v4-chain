//go:build all || integration_test

package cli_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestAllLiquidityTiers(t *testing.T) {
	genesisChanges := GetPerpetualGenesisShort()
	network.DeployCustomNetwork(genesisChanges)
	cfg := network.DefaultConfig(nil)

	perpQuery := "docker exec interchain-security-instance interchain-security-cd query perpetuals get-all-liquidity-tiers"
	data, _, err := network.QueryCustomNetwork(perpQuery)
	require.NoError(t, err)

	var resp types.QueryAllLiquidityTiersResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, 2, len(resp.LiquidityTiers))
	network.CleanupCustomNetwork()
}
