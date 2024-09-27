//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize
var defaultAssetYieldIndex = "1/1"

func TestGetAssetYieldIndexQuery(t *testing.T) {
	cfg := network.DefaultConfig(nil)

	assetYieldIndexQuery := "docker exec interchain-security-instance-setup interchain-security-cd" +
		" query ratelimit get-asset-yield-index"
	data, _, err := network.QueryCustomNetwork(assetYieldIndexQuery)

	require.NoError(t, err)
	var resp types.GetAssetYieldIndexQueryResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, defaultAssetYieldIndex, resp.AssetYieldIndex)
}
