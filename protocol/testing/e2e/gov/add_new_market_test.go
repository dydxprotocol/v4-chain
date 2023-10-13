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
	testId := uint32(1001)
	testMarketParam := pricestest.GenerateMarketParamPrice(
		pricestest.WithId(testId),
	)
	testClobPair := clobtest.GenerateClobPair(
		clobtest.WithId(testId),
		clobtest.WithPerpetualId(testId),
		clobtest.WithStatus(clobtypes.ClobPair_STATUS_INITIALIZING),
	)
	testPerpetual := perptest.GeneratePerpetual(
		perptest.WithId(testId),
		perptest.WithMarketId(testId),
	)
	msgUpdateClobPairToActive := &clobtypes.MsgUpdateClobPair{
		Authority: authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String(),
		ClobPair: *clobtest.GenerateClobPair(
			clobtest.WithId(testId),
			clobtest.WithPerpetualId(testId),
			clobtest.WithStatus(clobtypes.ClobPair_STATUS_ACTIVE),
		),
	}
	msgUpdateClobPairToActive_WrongClobPairId := &clobtypes.MsgUpdateClobPair{
		Authority: authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String(),
		ClobPair: *clobtest.GenerateClobPair(
			clobtest.WithId(9999), // non existing clob pair
			clobtest.WithPerpetualId(testId),
			clobtest.WithStatus(clobtypes.ClobPair_STATUS_ACTIVE),
		),
	}
	msgUpdateClobPairToActive_WrongAuthority := &clobtypes.MsgUpdateClobPair{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ClobPair: *clobtest.GenerateClobPair(
			clobtest.WithId(testId),
			clobtest.WithPerpetualId(testId),
			clobtest.WithStatus(clobtypes.ClobPair_STATUS_ACTIVE),
		),
	}

	tests := map[string]struct {
		proposedMsgs                          []sdk.Msg
		updateClobDelayBlocks                 uint32
		expectSubmitProposalFail              bool
		expectDelayedUpdateClobPairMsgFailure bool
		expectedProposalStatus                govtypesv1.ProposalStatus
	}{
		"Success with 4 standard messages": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:  10,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Success with 4 standard messages, delay blocks = 1": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive),
					DelayBlocks: 1,
				},
			},
			updateClobDelayBlocks:  1,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Success with 4 standard messages, delay blocks = 0": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive),
					DelayBlocks: 0,
				},
			},
			updateClobDelayBlocks:  0,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Success with 4 standard messages, delayed `UpdateClobPair` msg failure": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive_WrongClobPairId),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:                 10,
			expectDelayedUpdateClobPairMsgFailure: true,
			expectedProposalStatus:                govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: proposal execution fails due to incorrectly ordered messages": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam.Param,
				},
				// Create clob pair before creating perpetual, which will fail.
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual.Params,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:  10,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Fail: proposal execution fails due to existing objects": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params: pricestest.GenerateMarketParamPrice(
						pricestest.WithId(5), // already exists
					).Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params: perptest.GeneratePerpetual(
						perptest.WithId(5), // already exists
					).Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair: *clobtest.GenerateClobPair(
						clobtest.WithId(5), // already exists
					),
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:  10,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Fail: proposal submission fails, due to invalid signer": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(clobtypes.ModuleName).String(), // should be gov
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive),
					DelayBlocks: 10,
				},
			},
			expectSubmitProposalFail: true,
			updateClobDelayBlocks:    10,
			expectedProposalStatus:   govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Fail: proposal execution fails - invalid signer on `MsgDelayMessage`": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Msg:         encoding.EncodeMessageToAny(t, msgUpdateClobPairToActive_WrongAuthority),
					DelayBlocks: 10,
				},
			},
			updateClobDelayBlocks:  10,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
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

			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				&tApp,
				tc.proposedMsgs,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			if tc.expectSubmitProposalFail {
				require.Equal(t, initMarketParams, tApp.App.PricesKeeper.GetAllMarketParams(ctx))
				require.Equal(t, initPerpetuals, tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx))
				require.Equal(t, initClobPairs, tApp.App.ClobKeeper.GetAllClobPairs(ctx))
				require.Len(t, tApp.App.GovKeeper.GetProposals(ctx), 0)
				return
			}

			switch tc.expectedProposalStatus {
			case govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED:
				require.Equal(t, initMarketParams, tApp.App.PricesKeeper.GetAllMarketParams(ctx))
				require.Equal(t, initPerpetuals, tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx))
				require.Equal(t, initClobPairs, tApp.App.ClobKeeper.GetAllClobPairs(ctx))
			case govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED:
				// Proposal passed and successfully executed, check states are updated.
				// Check market
				marketParam, exists := tApp.App.PricesKeeper.GetMarketParam(ctx, testMarketParam.Param.Id)
				require.True(t, exists)
				require.Equal(t, testMarketParam.Param, marketParam)

				marketPrice, err := tApp.App.PricesKeeper.GetMarketPrice(ctx, testMarketParam.Param.Id)
				require.NoError(t, err)
				require.Equal(t, pricestypes.MarketPrice{
					Id:       testMarketParam.Param.Id,
					Price:    0, // expect oracle price to be initialized as zero.
					Exponent: testMarketParam.Param.Exponent,
				}, marketPrice)
				// Check perpeutal
				perp, err := tApp.App.PerpetualsKeeper.GetPerpetual(ctx, testPerpetual.Params.Id)
				require.NoError(t, err)
				require.Equal(t, testPerpetual.Params, perp.Params)

				// If `DelayBlocks` is not 0, check that clob pair is created in initial state.
				if tc.updateClobDelayBlocks != 0 {
					clobPair, exists := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(testClobPair.Id))
					require.True(t, exists)
					require.Equal(t, *testClobPair, clobPair)
				}

				ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+tc.updateClobDelayBlocks+1, testapp.AdvanceToBlockOptions{})

				clobPair, exists := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(testClobPair.Id))
				require.True(t, exists)

				// Check that clob pair is not updated.
				if tc.expectDelayedUpdateClobPairMsgFailure {
					require.Equal(t, *testClobPair, clobPair)
					return
				}

				// Check that clob pair is updated.
				require.Equal(t, msgUpdateClobPairToActive.ClobPair, clobPair)
				// TODO(CORE-585): Check that orders cannot be placed if no valid oracle price update has occurred.
			default:
				t.Errorf("unexpected proposal status: %s", tc.expectedProposalStatus)
			}
		})
	}
}
