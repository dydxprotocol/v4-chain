//go:build all || integration_test

package cli_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyEquityTierLimitConfig = types.EquityTierLimitConfiguration{
		ShortTermOrderEquityTiers: []types.EquityTierLimit{},
		StatefulOrderEquityTiers:  []types.EquityTierLimit{},
	}
)

func TestCmdGetEquityTierLimitConfig(t *testing.T) {
	networkWithClobPairObjects(t, 2)

	cfg := network.DefaultConfig(nil)
	query := "docker exec interchain-security-instance interchain-security-cd query clob get-equity-tier-limit-config  --node tcp://7.7.8.4:26658 -o json"
	data, _, err := network.QueryCustomNetwork(query)

	require.NoError(t, err)
	var resp types.QueryEquityTierLimitConfigurationResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.NotNil(t, resp.EquityTierLimitConfig)
	require.Equal(
		t,
		emptyEquityTierLimitConfig,
		resp.EquityTierLimitConfig,
	)
	network.CleanupCustomNetwork()
}
