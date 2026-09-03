//go:build all || container_test

package v_9_7_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	v_9_7 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.7"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

// The preupgrade genesis contains isolated perpetual 301 (BOME-USD) with its clob pair already in
// STATUS_FINAL_SETTLEMENT and its isolated insurance fund funded — the stranded-funds state that
// the v9.7 upgrade sweeps into the cross insurance fund.
const (
	finalSettlementIsolatedPerpetualId  = "301"
	usdcDenom                           = "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
	isolatedInsuranceFundGenesisBalance = "5000000000"
)

func TestStateUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	preUpgradeSetups(node, t)
	preUpgradeChecks(node, t)

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_9_7.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	require.Equal(t, isolatedInsuranceFundGenesisBalance, queryUsdcBalance(node, t, isolatedInsuranceFundAddress()))
	require.Equal(t, "0", queryUsdcBalance(node, t, perptypes.InsuranceFundModuleAddress.String()))
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// The final settlement market's isolated insurance fund was swept into the cross fund.
	require.Equal(t, "0", queryUsdcBalance(node, t, isolatedInsuranceFundAddress()))
	require.Equal(
		t,
		isolatedInsuranceFundGenesisBalance,
		queryUsdcBalance(node, t, perptypes.InsuranceFundModuleAddress.String()),
	)
}

func isolatedInsuranceFundAddress() string {
	return authtypes.NewModuleAddress(
		perptypes.InsuranceFundName + ":" + finalSettlementIsolatedPerpetualId,
	).String()
}

func queryUsdcBalance(node *containertest.Node, t *testing.T, address string) string {
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
