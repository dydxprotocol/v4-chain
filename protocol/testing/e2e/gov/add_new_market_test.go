package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestAddNewMarketProposal(t *testing.T) {
	testMarketParam1001 := pricestest.GenerateMarketParamPrice(
		pricestest.WithId(1001),
	)
	testClobPair1001 := clobtest.GenerateClobPair(
		clobtest.WithId(1001),
		clobtest.WithPerpetualId(1001),
		clobtest.WithStatus(clobtypes.ClobPair_STATUS_INITIALIZING),
	)
	testPerpetual1001 := perptest.GeneratePerpetual(
		perptest.WithId(1001),
		perptest.WithMarketId(1001),
	)
	msgUpdateClobPair1001ToActive := &clobtypes.MsgUpdateClobPair{
		Authority: authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String(),
		ClobPair: *clobtest.GenerateClobPair(
			clobtest.WithId(1001),
			clobtest.WithPerpetualId(1001),
			clobtest.WithStatus(clobtypes.ClobPair_STATUS_ACTIVE),
		),
	}

	tests := map[string]struct {
		proposedMsgs             []sdk.Msg
		updateClobDelayBlocks    uint32
		expectSubmitProposalFail bool
		expectedProposalStatus   govtypesv1.ProposalStatus
	}{
		"Success with 4 standard messages": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam1001.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual1001.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair1001,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPair1001ToActive),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:  10,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Proposals execution fails due to incorrectly ordered messages": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam1001.Param,
				},
				// Create clob pair before creating perpetual, which will fail.
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair1001,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual1001.Params,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPair1001ToActive),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:  10,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Proposal execution fails due to existing objects": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params: pricestest.GenerateMarketParamPrice(
						pricestest.WithId(5),
					).Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params: perptest.GeneratePerpetual(
						perptest.WithId(5),
					).Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair: *clobtest.GenerateClobPair(
						clobtest.WithId(5),
					),
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPair1001ToActive),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:  10,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Proposal submission fails, due to invalid signer": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(clobtypes.ModuleName).String(),
					Params: pricestest.GenerateMarketParamPrice(
						pricestest.WithId(5),
					).Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params: perptest.GeneratePerpetual(
						perptest.WithId(5),
					).Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair: *clobtest.GenerateClobPair(
						clobtest.WithId(5),
					),
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPair1001ToActive),
					DelayBlocks: 10,
				},
			},
			expectSubmitProposalFail: true,
			updateClobDelayBlocks:    10,
			expectedProposalStatus:   govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		// TODO(): Uncomment this case. Currently gov EndBlocker doesn't recover from panic in
		// message handler, and invalid message signer results in panic. We should fix this behavior,
		// and verifies that the proposal just fails instead.
		// "Invalid signer on `MsgDelayMessage`": {
		// 	proposedMsgs: []sdk.Msg{
		// 		&pricestypes.MsgCreateOracleMarket{
		// 			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		// 			Params:    testMarketParam1001.Param,
		// 		},
		// 		&perptypes.MsgCreatePerpetual{
		// 			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		// 			Params:    testPerpetual1001.Params,
		// 		},
		// 		&clobtypes.MsgCreateClobPair{
		// 			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		// 			ClobPair:  *testClobPair1001,
		// 		},
		// 		&delaymsgtypes.MsgDelayMessage{
		// 			Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		// 			Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPair1001ToActive_WrongAuthority),
		// 			DelayBlocks: 10,
		// 		},
		// 	},
		// 	updateClobDelayBlocks: 10,
		//  expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			initMarketParams := tApp.App.PricesKeeper.GetAllMarketParams(ctx)
			initPerpetuals := tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx)
			initClobPairs := tApp.App.ClobKeeper.GetAllClobPairs(ctx)

			ctx = testapp.SubmitAndPassProposal(
				t,
				ctx,
				&tApp,
				tc.proposedMsgs,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			if tc.expectSubmitProposalFail {
				return
			}

			switch tc.expectedProposalStatus {
			case govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED:
				require.Equal(t, initMarketParams, tApp.App.PricesKeeper.GetAllMarketParams(ctx))
				require.Equal(t, initPerpetuals, tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx))
				require.Equal(t, initClobPairs, tApp.App.ClobKeeper.GetAllClobPairs(ctx))
			case govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED:
				// Proposal passed and and succesfully executed, check states are updated.
				// Check market
				marketParam, exists := tApp.App.PricesKeeper.GetMarketParam(ctx, testMarketParam1001.Param.Id)
				require.True(t, exists)
				require.Equal(t, testMarketParam1001.Param, marketParam)

				marketPrice, err := tApp.App.PricesKeeper.GetMarketPrice(ctx, testMarketParam1001.Param.Id)
				require.NoError(t, err)
				require.Equal(t, pricestypes.MarketPrice{
					Id:       testMarketParam1001.Param.Id,
					Price:    0,
					Exponent: testMarketParam1001.Param.Exponent,
				}, marketPrice)
				// Check perpeutal
				perp, err := tApp.App.PerpetualsKeeper.GetPerpetual(ctx, testPerpetual1001.Params.Id)
				require.NoError(t, err)
				require.Equal(t, testPerpetual1001.Params, perp.Params)

				// Check clob
				clobPair, exists := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(testClobPair1001.Id))
				require.True(t, exists)
				require.Equal(t, clobPair, *testClobPair1001)

				ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+tc.updateClobDelayBlocks+1, testapp.AdvanceToBlockOptions{})
				clobPair, exists = tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(testClobPair1001.Id))
				require.True(t, exists)
				// Check that clob pair is updated.
				require.Equal(t, msgUpdateClobPair1001ToActive.ClobPair, clobPair)
			default:
				t.Errorf("unexpected proposal status: %s", tc.expectedProposalStatus)
			}
		})
	}
}
