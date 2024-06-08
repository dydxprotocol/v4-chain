//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	blockQuery := "docker exec interchain-security-instance-setup interchain-security-cd query blocktime get-downtime-params"
	data, _, err := network.QueryCustomNetwork(blockQuery)

	require.NoError(t, err)
	var resp types.QueryDowntimeParamsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

func TestQueryAllDowntimeInfo(t *testing.T) {
	cfg := network.DefaultConfig(nil)

	blockQuery := "docker exec interchain-security-instance-setup interchain-security-cd query blocktime get-all-downtime-info"
	data, _, err := network.QueryCustomNetwork(blockQuery)

	require.NoError(t, err)
	var resp types.QueryAllDowntimeInfoResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}

func TestQueryPreviousBlockInfo(t *testing.T) {
	cfg := network.DefaultConfig(nil)

	blockQuery := "docker exec interchain-security-instance-setup interchain-security-cd query blocktime get-previous-block-info"
	data, _, err := network.QueryCustomNetwork(blockQuery)
	require.NoError(t, err)
	var resp types.QueryPreviousBlockInfoResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}
