//go:build all || container_test

package v_8_0_test

import (
	"testing"

	v_8_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v8.0"

	"github.com/cosmos/gogoproto/proto"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	aptypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

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

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_8_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Check that the listing module state has been initialized with the hard cap and default deposit params.
	postUpgradeSmartAccountActiveCheck(node, t)
	postUpgradeMarketIdsCheck(node, t)
}

func postUpgradeSmartAccountActiveCheck(node *containertest.Node, t *testing.T) {
	// query the smart account active
	resp, err := containertest.Query(
		node,
		aptypes.NewQueryClient,
		aptypes.QueryClient.Params,
		&aptypes.QueryParamsRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	queryResponse := aptypes.QueryParamsResponse{}
	err = proto.UnmarshalText(resp.String(), &queryResponse)
	require.NoError(t, err)
	require.Equal(t, true, queryResponse.Params.IsSmartAccountActive)
}

func postUpgradeMarketIdsCheck(node *containertest.Node, t *testing.T) {
	// query the next market id
	resp, err := containertest.Query(
		node,
		pricetypes.NewQueryClient,
		pricetypes.QueryClient.NextMarketId,
		&pricetypes.QueryNextMarketIdRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	nextMarketIdResp := pricetypes.QueryNextMarketIdResponse{}
	err = proto.UnmarshalText(resp.String(), &nextMarketIdResp)
	require.NoError(t, err)
	require.Equal(t, uint32(v_8_0.ID_NUM), nextMarketIdResp.NextMarketId)

	// query the next perpetual id
	resp, err = containertest.Query(
		node,
		perptypes.NewQueryClient,
		perptypes.QueryClient.NextPerpetualId,
		&perptypes.QueryNextPerpetualIdRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	nextPerpIdResp := perptypes.QueryNextPerpetualIdResponse{}
	err = proto.UnmarshalText(resp.String(), &nextPerpIdResp)
	require.NoError(t, err)
	require.Equal(t, uint32(v_8_0.ID_NUM), nextPerpIdResp.NextPerpetualId)

	// query the next clob pair id
	resp, err = containertest.Query(
		node,
		clobtypes.NewQueryClient,
		clobtypes.QueryClient.NextClobPairId,
		&clobtypes.QueryNextClobPairIdRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	nextClobPairIdResp := clobtypes.QueryNextClobPairIdResponse{}
	err = proto.UnmarshalText(resp.String(), &nextClobPairIdResp)
	require.NoError(t, err)
	require.Equal(t, uint32(v_8_0.ID_NUM), nextClobPairIdResp.NextClobPairId)
}
