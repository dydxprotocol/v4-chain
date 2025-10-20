//go:build all || container_test

package v_9_4_test

import (
	"testing"

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

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_9_4.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {
	// Set default affiliate tiers and parameters.
	node.SetAffiliateTiers(v_9_4.PreviousAffiliateTiers)
	node.SetAffiliateParameters(v_9_4.PreviousAffiliateParameters)
	node.SetAffiliateWhitelist(v_9_4.PreviousAffiliateWhitelist)
}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Check that the affiliate tiers and parameters are set to the default values.
	require.Equal(t, v_9_4.DefaultAffiliateTiers, node.GetAffiliateTiers())
	require.Equal(t, v_9_4.DefaultAffiliateParameters, node.GetAffiliateParameters())
	require.Equal(t, v_9_4.DefaultAffiliateOverrides, node.GetAffiliateOverrides())
}
