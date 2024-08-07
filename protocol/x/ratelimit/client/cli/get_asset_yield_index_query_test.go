//go:build all || integration_test

package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize
var defaultAssetYieldIndex = "0/1"

func TestGetAssetYieldIndexQuery(t *testing.T) {
	cfg := network.DefaultConfig(nil)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	assetYieldIndexQuery := "docker exec interchain-security-instance-setup interchain-security-cd" +
		" query get-asset-yield-index " + param
	data, _, err := network.QueryCustomNetwork(assetYieldIndexQuery)

	require.NoError(t, err)
	var resp types.GetAssetYieldIndexQueryResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, defaultAssetYieldIndex, resp.AssetYieldIndex)
}
