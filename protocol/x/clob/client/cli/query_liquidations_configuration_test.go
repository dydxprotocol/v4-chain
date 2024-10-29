//go:build all || integration_test

package cli_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/client/cli"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
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

// func TestCmdGetLiquidationsConfiguration(t *testing.T) {
// 	fmt.Println("TestCmdGetLiquidationsConfiguration")
// 	networkWithClobPairObjects(t, 2)

// 	cfg := network.DefaultConfig(nil)
// 	query := "docker exec interchain-security-instance interchain-security-cd query clob get-liquidations-config"
// 	data, _, err := network.QueryCustomNetwork(query)
// 	require.NoError(t, err)
// 	var resp types.QueryLiquidationsConfigurationResponse
// 	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
// 	require.NotNil(t, resp.LiquidationsConfig)
// 	require.Equal(
// 		t,
// 		types.LiquidationsConfig_Default,
// 		resp.LiquidationsConfig,
// 	)
// 	network.CleanupCustomNetwork()
// }
