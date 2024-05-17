//go:build all || integration_test

package cli_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {

	genesisChanges := GetPerpetualGenesisShort()
	network.DeployCustomNetwork(genesisChanges)
	cfg := network.DefaultConfig(nil)

	perpQuery := "docker exec interchain-security-instance-setup interchain-security-cd query perpetuals get-params --node tcp://7.7.8.4:26658 -o json"
	data, _, err := network.QueryCustomNetwork(perpQuery)
	require.NoError(t, err)

	var resp types.QueryParamsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
	network.CleanupCustomNetwork()
}
