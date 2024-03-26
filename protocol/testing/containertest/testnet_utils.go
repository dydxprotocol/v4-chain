package containertest

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"

	upgrade "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

const GovModuleAddress = "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky"

var NodeAddresses = []string{
	constants.AliceAccAddress.String(),
	constants.BobAccAddress.String(),
	constants.CarlAccAddress.String(),
	constants.DaveAccAddress.String(),
}

func UpgradeTestnet(nodeAddress string, t *testing.T, node *Node, upgradeToVersion string) error {
	proposal, err := gov.NewMsgSubmitProposal(
		[]sdk.Msg{
			&upgrade.MsgSoftwareUpgrade{
				Authority: GovModuleAddress,
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

	require.NoError(t, BroadcastTx(
		node,
		proposal,
		nodeAddress,
	))
	err = node.Wait(2)
	require.NoError(t, err)

	for _, address := range NodeAddresses {
		require.NoError(t, BroadcastTx(
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
