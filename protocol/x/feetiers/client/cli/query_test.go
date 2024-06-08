//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryPerpetualFeeParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	feeQuery := "docker exec interchain-security-instance-setup interchain-security-cd query feetiers get-perpetual-fee-params"
	data, _, err := network.QueryCustomNetwork(feeQuery)
	require.NoError(t, err)
	var resp types.QueryPerpetualFeeParamsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().Params, resp.Params)
}

func TestQueryUserFeeTier(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	feeQuery := "docker exec interchain-security-instance-setup interchain-security-cd query feetiers get-user-fee-tier alice"
	data, _, err := network.QueryCustomNetwork(feeQuery)
	require.NoError(t, err)
	var resp types.QueryUserFeeTierResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
}
