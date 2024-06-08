//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	rewardQuery := "docker exec interchain-security-instance-setup interchain-security-cd query rewards params"
	data, _, err := network.QueryCustomNetwork(rewardQuery)

	require.NoError(t, err)
	var resp types.QueryParamsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}
