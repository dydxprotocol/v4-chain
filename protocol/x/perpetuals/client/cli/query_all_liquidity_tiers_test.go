//go:build all || integration_test

package cli_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/client/cli"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
)

func TestAllLiquidityTiers(t *testing.T) {
	net, _, _ := networkWithLiquidityTierAndPerpetualObjects(t, 2, 2)
	ctx := net.Validators[0].ClientCtx

	common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryAllLiquidityTiers(), common)
	require.NoError(t, err)

	var resp types.QueryAllLiquidityTiersResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, 2, len(resp.LiquidityTiers))
	require.Equal(t, types.DefaultGenesis().LiquidityTiers, resp.LiquidityTiers)
}
