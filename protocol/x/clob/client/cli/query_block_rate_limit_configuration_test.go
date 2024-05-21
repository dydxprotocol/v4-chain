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
	emptyConfig = types.BlockRateLimitConfiguration{
		MaxShortTermOrdersPerNBlocks:             []types.MaxPerNBlocksRateLimit{},
		MaxStatefulOrdersPerNBlocks:              []types.MaxPerNBlocksRateLimit{},
		MaxShortTermOrderCancellationsPerNBlocks: []types.MaxPerNBlocksRateLimit{},
	}
)

func TestCmdGetBlockRateLimitConfiguration(t *testing.T) {
	fmt.Println("TestCmdGetBlockRateLimitConfiguration")
	networkWithClobPairObjects(t, 2)

	cfg := network.DefaultConfig(nil)
	query := "docker exec interchain-security-instance interchain-security-cd query clob get-block-rate-limit-config --node tcp://7.7.8.4:26658 -o json"
	data, _, err := network.QueryCustomNetwork(query)

	require.NoError(t, err)
	var resp types.QueryBlockRateLimitConfigurationResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.NotNil(t, resp.BlockRateLimitConfig)
	require.Equal(
		t,
		emptyConfig,
		resp.BlockRateLimitConfig,
	)
	network.CleanupCustomNetwork()
}
