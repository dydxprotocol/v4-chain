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

	upgrade "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

const govModuleAddress = "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky"

var nodeAddresses = []string{
	constants.AliceAccAddress.String(),
	constants.BobAccAddress.String(),
	constants.CarlAccAddress.String(),
	constants.DaveAccAddress.String(),
}

func TestUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()
	err = upgradeTestnet(nodeAddress, t, node, v_5_0_0.UpgradeName)
	require.NoError(t, err)
}

func TestStateUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	preUpgradeChecks(node, t)

	err = upgradeTestnet(nodeAddress, t, node, v_5_0_0.UpgradeName)
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

func upgradeTestnet(nodeAddress string, t *testing.T, node *containertest.Node, upgradeToVersion string) error {
	proposal, err := gov.NewMsgSubmitProposal(
		[]sdk.Msg{
			&upgrade.MsgSoftwareUpgrade{
				Authority: govModuleAddress,
				Plan: upgrade.Plan{
					Name:   upgradeToVersion,
					Height: 10,
				},
			},
		},
		testapp.TestDeposit,
		nodeAddress,
		testapp.TestMetadata,
		testapp.TestTitle,
		testapp.TestSummary,
		false,
	)
	require.NoError(t, err)

	require.NoError(t, containertest.BroadcastTx(
		node,
		proposal,
		nodeAddress,
	))
	err = node.Wait(2)
	require.NoError(t, err)

	for _, address := range nodeAddresses {
		require.NoError(t, containertest.BroadcastTx(
			node,
			&gov.MsgVote{
				ProposalId: 1,
				Voter:      address,
				Option:     gov.VoteOption_VOTE_OPTION_YES,
			},
			address,
		))
	}

	err = node.WaitUntilBlockHeight(12)
	return err
}
