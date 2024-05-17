//go:build all || integration_test

package cli_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestCmdGetLiquidationsConfiguration(t *testing.T) {
	networkWithClobPairObjects(t, 2)

	cfg := network.DefaultConfig(nil)
	query := "docker exec interchain-security-instance interchain-security-cd query clob get-liquidations-config --node tcp://7.7.8.4:26658 -o json"
	data, _, err := network.QueryCustomNetwork(query)
	require.NoError(t, err)
	var resp types.QueryLiquidationsConfigurationResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.NotNil(t, resp.LiquidationsConfig)
	require.Equal(
		t,
		types.LiquidationsConfig_Default,
		resp.LiquidationsConfig,
	)
	network.CleanupCustomNetwork()
}
