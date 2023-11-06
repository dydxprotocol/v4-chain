//go:build all || integration_test

package cli_test

import (
	"fmt"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCmdGetLiquidationsConfiguration(t *testing.T) {
	net, _ := networkWithClobPairObjects(t, 2)
	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdGetLiquidationsConfiguration(), common)
	require.NoError(t, err)
	var resp types.QueryLiquidationsConfigurationResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.NotNil(t, resp.LiquidationsConfig)
	require.Equal(
		t,
		types.LiquidationsConfig_Default,
		resp.LiquidationsConfig,
	)
}
