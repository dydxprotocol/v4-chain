//go:build all || integration_test

package cli_test

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestAllLiquidityTiers(t *testing.T) {
	genesisChanges := GetPerpetualGenesisShort()
	network.DeployCustomNetwork(genesisChanges)
	cfg := network.DefaultConfig(nil)

	common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}

	out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryAllLiquidityTiers(), common)
	require.NoError(t, err)

	perpQuery := "docker exec interchain-security-instance-setup interchain-security-cd query perpetuals get-all-liquidity-tiers"
	data, _, err := network.QueryCustomNetwork(perpQuery)
	require.NoError(t, err)

	var resp types.QueryAllLiquidityTiersResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, 2, len(resp.LiquidityTiers))
	network.CleanupCustomNetwork()
}
