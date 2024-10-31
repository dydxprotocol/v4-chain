//go:build all || integration_test

package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/client/cli"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize
var defaultAssetYieldIndex = "1/1"

func TestGetAssetYieldIndexQuery(t *testing.T) {
	net, ctx := setupNetwork(t)

	common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetAssetYieldIndexQuery(), common)
	require.NoError(t, err)

	var resp types.GetAssetYieldIndexQueryResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, defaultAssetYieldIndex, resp.AssetYieldIndex)
}
