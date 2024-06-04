//go:build all || integration_test

package cli_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestPendingSendPackets(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	rateQuery := "docker exec interchain-security-instance-setup interchain-security-cd query ratelimit pending-send-packets"
	data, _, err := network.QueryCustomNetwork(rateQuery)

	require.NoError(t, err)
	var resp types.QueryAllPendingSendPacketsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	assert.Equal(t, 0, len(resp.PendingSendPackets))
}
