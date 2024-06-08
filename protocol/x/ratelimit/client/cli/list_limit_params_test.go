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

func TestListLimiterParams(t *testing.T) {

	cfg := network.DefaultConfig(nil)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	rateQuery := "docker exec interchain-security-instance-setup interchain-security-cd query ratelimit list-limit-params " + param
	data, _, err := network.QueryCustomNetwork(rateQuery)

	require.NoError(t, err)
	var resp types.ListLimitParamsResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, types.DefaultGenesis().LimitParamsList, resp.LimitParamsList)
}
