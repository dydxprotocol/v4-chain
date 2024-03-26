//go:build all || container_test

package v_5_0_0_test

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	v_5_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v5.0.0"
	containertest "github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	perpetuals "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	preUpgradeChecks(node, t)

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_5_0_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	preUpgradeCheckPerpetualMarketType(node, t)
	// Add test for your upgrade handler logic below
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	postUpgradecheckPerpetualMarketType(node, t)
	// Add test for your upgrade handler logic below
}

func preUpgradeCheckPerpetualMarketType(node *containertest.Node, t *testing.T) {
	perpetualsList := &perpetuals.QueryAllPerpetualsResponse{}
	resp, err := containertest.Query(
		node,
		perpetuals.NewQueryClient,
		perpetuals.QueryClient.AllPerpetuals,
		&perpetuals.QueryAllPerpetualsRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), perpetualsList)
	require.NoError(t, err)
	for _, perpetual := range perpetualsList.Perpetual {
		assert.Equal(t, perpetuals.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED, perpetual.Params.MarketType)
	}
}

func postUpgradecheckPerpetualMarketType(node *containertest.Node, t *testing.T) {
	perpetualsList := &perpetuals.QueryAllPerpetualsResponse{}
	resp, err := containertest.Query(
		node,
		perpetuals.NewQueryClient,
		perpetuals.QueryClient.AllPerpetuals,
		&perpetuals.QueryAllPerpetualsRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), perpetualsList)
	require.NoError(t, err)
	for _, perpetual := range perpetualsList.Perpetual {
		assert.Equal(t, perpetuals.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS, perpetual.Params.MarketType)
	}
}
