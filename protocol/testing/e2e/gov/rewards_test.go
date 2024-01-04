package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	rewardstypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

var (
	GENESIS_REWARD_PARAMS = rewardstypes.Params{
		TreasuryAccount:  "test_treasury",
		Denom:            "adv4tnt",
		DenomExponent:    -18,
		MarketId:         1234,
		FeeMultiplierPpm: 700_000,
	}
)

// This tests `MsgUpdateParams` in `x/rewards`.
func TestUpdateRewardsModuleParams(t *testing.T) {
	tests := map[string]struct {
		msg                       *rewardstypes.MsgUpdateParams
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &rewardstypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: rewardstypes.Params{
					TreasuryAccount:  "test_treasury",
					Denom:            "adv4tnt",
					DenomExponent:    -5,
					MarketId:         0,
					FeeMultiplierPpm: 700_001,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: treasury account is empty": {
			msg: &rewardstypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: rewardstypes.Params{
					TreasuryAccount:  "",
					Denom:            GENESIS_REWARD_PARAMS.Denom,
					DenomExponent:    GENESIS_REWARD_PARAMS.DenomExponent,
					MarketId:         GENESIS_REWARD_PARAMS.MarketId,
					FeeMultiplierPpm: GENESIS_REWARD_PARAMS.FeeMultiplierPpm,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: denom is invalid": {
			msg: &rewardstypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: rewardstypes.Params{
					TreasuryAccount:  GENESIS_REWARD_PARAMS.TreasuryAccount,
					Denom:            "7adv4tnt", // cannot start with number
					DenomExponent:    GENESIS_REWARD_PARAMS.DenomExponent,
					MarketId:         GENESIS_REWARD_PARAMS.MarketId,
					FeeMultiplierPpm: GENESIS_REWARD_PARAMS.FeeMultiplierPpm,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: fee multiplier ppm is greater than 1 million": {
			msg: &rewardstypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: rewardstypes.Params{
					TreasuryAccount:  GENESIS_REWARD_PARAMS.TreasuryAccount,
					Denom:            GENESIS_REWARD_PARAMS.Denom,
					DenomExponent:    GENESIS_REWARD_PARAMS.DenomExponent,
					MarketId:         GENESIS_REWARD_PARAMS.MarketId,
					FeeMultiplierPpm: 1_000_001,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &rewardstypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(rewardstypes.ModuleName).String(),
				Params:    GENESIS_REWARD_PARAMS,
			},
			expectSubmitProposalFails: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				// Initialize rewards module with genesis params.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *rewardstypes.GenesisState) {
						genesisState.Params = GENESIS_REWARD_PARAMS
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialParams := tApp.App.RewardsKeeper.GetParams(ctx)

			// Submit and tally governance proposal that includes `MsgUpdateParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFails,
				tc.expectedProposalStatus,
			)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that rewards params have been updated.
				require.Equal(t, tc.msg.Params, tApp.App.RewardsKeeper.GetParams(ctx))
			} else {
				// Otherwise, verify that rewards module params match the ones before proposal submission.
				require.Equal(t, initialParams, tApp.App.RewardsKeeper.GetParams(ctx))
			}
		})
	}
}
