//go:build all || container_test

package v_7_0_3_test

import (
	"testing"

	v_7_0_3 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v7.0.3"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	v_7_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v7.0"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
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

	preUpgradeSetups(node, t)
	preUpgradeChecks(node, t)

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_7_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Check that the listing module state has been initialized with the hard cap and default deposit params.
	postUpgradeMarketIdsCheck(node, t)
}

func postUpgradeMarketIdsCheck(node *containertest.Node, t *testing.T) {
	// query the next market id
	resp, err := containertest.Query(
		node,
		pricetypes.NewQueryClient,
		pricetypes.NextMarketId,
		pricetypes.QueryNextMarketIdRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, uint32(v_7_0_3.ID_NUM), resp.NextMarketId)

	// query the next perpetual id
	resp, err = containertest.Query(
		node,
		perptypes.NewQueryClient,
		perptypes.NextPerpetualId,
		perptypes.QueryNextPerpetualIdRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, uint32(v_7_0_3.ID_NUM), resp.NextPerpetualId)

	// query the next clob pair id
	resp, err = containertest.Query(
		node,
		clobtypes.NewQueryClient,
		clobtypes.NextClobPairId,
		clobtypes.QueryNextClobPairIdRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, uint32(v_7_0_3.ID_NUM), resp.NextClobPairId)
}
