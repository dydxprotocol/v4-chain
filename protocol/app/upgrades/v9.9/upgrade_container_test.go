//go:build all || container_test

package v_9_9_test

import (
	"strings"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	v_9_9 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.9"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testing/version"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

const usdcDenom = "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"

func TestCanonicalUsdcUpgrade(t *testing.T) {
	if strings.TrimSpace(version.CurrentVersion) != v_9_9.UpgradeName {
		t.Skipf("container upgrade to %s requires VERSION_CURRENT=%s", v_9_9.UpgradeName, v_9_9.UpgradeName)
	}

	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	aliceBefore := queryUsdcBalance(node, t, nodeAddress)

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_9_9.UpgradeName)
	require.NoError(t, err)

	resp, err := containertest.Query(
		node,
		canonicalusdctypes.NewQueryClient,
		canonicalusdctypes.QueryClient.State,
		&canonicalusdctypes.QueryStateRequest{},
	)
	require.NoError(t, err)
	state, ok := resp.(*canonicalusdctypes.QueryStateResponse)
	require.True(t, ok)
	require.Equal(t, canonicalusdctypes.Mode_MODE_DISABLED, state.Controls.Mode)
	require.Equal(t, "0", state.Ledger.NobleBacking.String())
	require.Equal(t, "0", state.Ledger.InjectiveBacking.String())
	require.Equal(t, "0", state.Ledger.PendingInjective.String())
	require.Equal(t, "0", state.Ledger.PendingNoble.String())

	require.Equal(t, aliceBefore, queryUsdcBalance(node, t, nodeAddress))

	accountResp, err := containertest.Query(
		node,
		authtypes.NewQueryClient,
		authtypes.QueryClient.ModuleAccountByName,
		&authtypes.QueryModuleAccountByNameRequest{Name: canonicalusdctypes.ModuleName},
	)
	require.NoError(t, err)
	require.NotNil(t, accountResp)
}

func queryUsdcBalance(node *containertest.Node, t *testing.T, address string) string {
	t.Helper()
	resp, err := containertest.Query(
		node,
		banktypes.NewQueryClient,
		banktypes.QueryClient.Balance,
		&banktypes.QueryBalanceRequest{Address: address, Denom: usdcDenom},
	)
	require.NoError(t, err)
	balanceResp, ok := resp.(*banktypes.QueryBalanceResponse)
	require.True(t, ok)
	return balanceResp.Balance.Amount.String()
}
