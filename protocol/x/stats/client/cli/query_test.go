//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	statsQuery := "docker exec interchain-security-instance-setup interchain-security-cd query stats get-params"
	data, _, err := network.QueryCustomNetwork(statsQuery)

	require.NoError(t, err)
	var resp types.QueryParamsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

func TestQueryStatsMetadata(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	statsQuery := "docker exec interchain-security-instance-setup interchain-security-cd query stats get-stats-metadata"
	data, _, err := network.QueryCustomNetwork(statsQuery)

	require.NoError(t, err)
	var resp types.QueryStatsMetadataResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}

func TestQueryGlobalStats(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	statsQuery := "docker exec interchain-security-instance-setup interchain-security-cd query stats get-global-stats"
	data, _, err := network.QueryCustomNetwork(statsQuery)

	require.NoError(t, err)
	var resp types.QueryGlobalStatsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}

func TestQueryUserStats(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	statsQuery := "docker exec interchain-security-instance-setup interchain-security-cd query stats get-user-stats alice"
	data, _, err := network.QueryCustomNetwork(statsQuery)


	require.NoError(t, err)
	var resp types.QueryUserStatsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}
