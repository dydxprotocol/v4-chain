package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateControlsDisabledProposal(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *govtypesv1.GenesisState) {
				genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	ctx = testapp.SubmitAndTallyProposal(
		t,
		ctx,
		tApp,
		[]sdk.Msg{&canonicalusdctypes.MsgUpdateControls{
			Authority: lib.GovModuleAddress.String(),
			Controls:  canonicalusdctypes.DefaultControls(),
		}},
		testapp.TestSubmitProposalTxHeight,
		false,
		false,
		govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
	)
	require.Equal(t, canonicalusdctypes.Mode_MODE_DISABLED, tApp.App.CanonicalUsdcKeeper.GetControls(ctx).Mode)
}
