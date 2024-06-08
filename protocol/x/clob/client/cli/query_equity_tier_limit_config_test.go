//go:build all || integration_test

package cli_test

import (
	"fmt"
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
	fmt.Println("TestCmdGetEquityTierLimitConfig")
	networkWithClobPairObjects(t, 2)

	cfg := network.DefaultConfig(nil)
	query := "docker exec interchain-security-instance interchain-security-cd query clob get-equity-tier-limit-config"
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
