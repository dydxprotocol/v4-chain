//go:build all || container_test

package v_8_0_0_test

import (
	"testing"

	v_8_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v8.0.0"
	listingtypes "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"

	"github.com/cosmos/gogoproto/proto"

	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

const (
	AliceBobBTCQuantums = 1_000_000
	CarlDaveBTCQuantums = 2_000_000
	CarlDaveETHQuantums = 4_000_000
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

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_8_0_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
	postUpgradeListingModuleStateCheck(node, t)
}

func postUpgradeListingModuleStateCheck(node *containertest.Node, t *testing.T) {
	// Check that the listing module state has been initialized with the hard cap and default deposit params.
	resp, err := containertest.Query(
		node,
		listingtypes.NewQueryClient,
		listingtypes.QueryClient.ListingVaultDepositParams,
		&listingtypes.QueryListingVaultDepositParams{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	listingVaultDepositParamsResp := listingtypes.QueryListingVaultDepositParamsResponse{}
	err = proto.UnmarshalText(resp.String(), &listingVaultDepositParamsResp)
	require.NoError(t, err)
	require.Equal(t, listingtypes.DefaultParams(), listingVaultDepositParamsResp.Params)

	resp, err = containertest.Query(
		node,
		listingtypes.NewQueryClient,
		listingtypes.QueryClient.MarketsHardCap,
		&listingtypes.QueryMarketsHardCap{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	marketsHardCapResp := listingtypes.QueryMarketsHardCapResponse{}
	err = proto.UnmarshalText(resp.String(), &marketsHardCapResp)
	require.NoError(t, err)
	require.Equal(t, uint32(500), marketsHardCapResp.HardCap)
}
