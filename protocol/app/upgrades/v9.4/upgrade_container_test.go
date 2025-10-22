//go:build all || container_test

package v_9_4_test

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	v_9_4 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.4"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
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

func preUpgradeSetups(node *containertest.Node, t *testing.T) {}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Verify affiliate tiers are set to default values
	tiersResp := &affiliatetypes.AllAffiliateTiersResponse{}
	resp, err := containertest.Query(
		node,
		affiliatetypes.NewQueryClient,
		affiliatetypes.QueryClient.AllAffiliateTiers,
		&affiliatetypes.AllAffiliateTiersRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), tiersResp)
	require.NoError(t, err)
	require.Equal(t, v_9_4.PreviousAffilliateTiers, tiersResp.Tiers)
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Verify affiliate tiers are set to default values
	tiersResp := &affiliatetypes.AllAffiliateTiersResponse{}
	resp, err := containertest.Query(
		node,
		affiliatetypes.NewQueryClient,
		affiliatetypes.QueryClient.AllAffiliateTiers,
		&affiliatetypes.AllAffiliateTiersRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), tiersResp)
	require.NoError(t, err)
	require.Equal(t, v_9_4.DefaultAffiliateTiers, tiersResp.Tiers)

	// Verify affiliate parameters are set to default values
	paramsResp := &affiliatetypes.AffiliateParametersResponse{}
	resp, err = containertest.Query(
		node,
		affiliatetypes.NewQueryClient,
		affiliatetypes.QueryClient.AffiliateParameters,
		&affiliatetypes.AffiliateParametersRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), paramsResp)
	require.NoError(t, err)
	require.Equal(t, v_9_4.DefaultAffiliateParameters, paramsResp.Parameters)

	// Verify affiliate overrides were migrated from whitelist
	overridesResp := &affiliatetypes.AffiliateOverridesResponse{}
	resp, err = containertest.Query(
		node,
		affiliatetypes.NewQueryClient,
		affiliatetypes.QueryClient.AffiliateOverrides,
		&affiliatetypes.AffiliateOverridesRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), overridesResp)
	require.NoError(t, err)
	// Overrides should contain addresses from the pre-upgrade whitelist
	expectedOverrides := affiliatetypes.AffiliateOverrides{
		Addresses: []string{
			"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4", // Carl
			"dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs", // Dave
		},
	}
	require.Equal(t, expectedOverrides, overridesResp.Overrides)
}
