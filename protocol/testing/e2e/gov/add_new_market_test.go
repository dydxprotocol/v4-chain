package gov_test

import (
	"testing"
	"time"

	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/lib/marketmap"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

const (
	NumBlocksAfterTradingEnabled = 50
	TestMarketId                 = 1001
)

var (
	GenesisTime                                     = time.Unix(1690000000, 0)
	OrderTemplate_Alice_Num0_Id0_Clob0_Buy_LongTerm = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   TestMarketId,
		},
		Quantums: 1_000_000_000_000,
		Subticks: 1_000_000_000,
		Side:     clobtypes.Order_SIDE_BUY,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(GenesisTime.Add(1 * time.Hour).Unix()),
		},
	}
)

func TestAddNewMarketProposal(t *testing.T) {
	testMarketParam := pricestest.GenerateMarketParamPrice(
		pricestest.WithId(TestMarketId),
	)

	testClobPair := clobtest.GenerateClobPair(
		clobtest.WithId(TestMarketId),
		clobtest.WithPerpetualId(TestMarketId),
		clobtest.WithStatus(clobtypes.ClobPair_STATUS_INITIALIZING),
	)
	testPerpetual := perptest.GeneratePerpetual(
		perptest.WithId(TestMarketId),
		perptest.WithMarketId(TestMarketId),
	)
	msgUpdateClobPairToActive := &clobtypes.MsgUpdateClobPair{
		Authority: delaymsgtypes.ModuleAddress.String(),
		ClobPair: *clobtest.GenerateClobPair(
			clobtest.WithId(TestMarketId),
			clobtest.WithPerpetualId(TestMarketId),
			clobtest.WithStatus(clobtypes.ClobPair_STATUS_ACTIVE),
		),
	}
	msgUpdateClobPairToActive_WrongClobPairId := &clobtypes.MsgUpdateClobPair{
		Authority: delaymsgtypes.ModuleAddress.String(),
		ClobPair: *clobtest.GenerateClobPair(
			clobtest.WithId(9999), // non existing clob pair
			clobtest.WithPerpetualId(TestMarketId),
			clobtest.WithStatus(clobtypes.ClobPair_STATUS_ACTIVE),
		),
	}
	msgUpdateClobPairToActive_WrongAuthority := &clobtypes.MsgUpdateClobPair{
		Authority: lib.GovModuleAddress.String(),
		ClobPair: *clobtest.GenerateClobPair(
			clobtest.WithId(TestMarketId),
			clobtest.WithPerpetualId(TestMarketId),
			clobtest.WithStatus(clobtypes.ClobPair_STATUS_ACTIVE),
		),
	}

	tests := map[string]struct {
		proposedMsgs                          []sdk.Msg
		updateClobDelayBlocks                 uint32
		expectCheckTxFails                    bool
		expectSubmitProposalFail              bool
		expectDelayedUpdateClobPairMsgFailure bool
		expectedProposalStatus                govtypesv1.ProposalStatus
	}{
		"Success with 4 standard messages": {
			proposedMsgs: []sdk.Msg{
				&pricestypes.MsgCreateOracleMarket{
					Authority: lib.GovModuleAddress.String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: lib.GovModuleAddress.String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
					Authority: lib.GovModuleAddress.String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: lib.GovModuleAddress.String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
					Authority: lib.GovModuleAddress.String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: lib.GovModuleAddress.String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
					Authority: lib.GovModuleAddress.String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: lib.GovModuleAddress.String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
					Authority: lib.GovModuleAddress.String(),
					Params:    testMarketParam.Param,
				},
				// Create clob pair before creating perpetual, which will fail.
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair:  *testClobPair,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: lib.GovModuleAddress.String(),
					Params:    testPerpetual.Params,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
					Authority: lib.GovModuleAddress.String(),
					Params: pricestest.GenerateMarketParamPrice(
						pricestest.WithId(5), // already exists
					).Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: lib.GovModuleAddress.String(),
					Params: perptest.GeneratePerpetual(
						perptest.WithId(5), // already exists
					).Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair: *clobtest.GenerateClobPair(
						clobtest.WithId(5), // already exists
					),
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
					Authority: lib.GovModuleAddress.String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
					Authority: lib.GovModuleAddress.String(),
					Params:    testMarketParam.Param,
				},
				&perptypes.MsgCreatePerpetual{
					Authority: lib.GovModuleAddress.String(),
					Params:    testPerpetual.Params,
				},
				&clobtypes.MsgCreateClobPair{
					Authority: lib.GovModuleAddress.String(),
					ClobPair:  *testClobPair,
				},
				&delaymsgtypes.MsgDelayMessage{
					Authority:   lib.GovModuleAddress.String(),
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
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *marketmaptypes.GenesisState) {
						// Add test market to market map genesis
						marketMap, err := marketmap.ConstructMarketMapFromParams([]pricestypes.MarketParam{testMarketParam.Param})
						require.NoError(t, err)
						for ticker, market := range marketMap.Markets {
							market.Ticker.Enabled = false
							genesisState.MarketMap.Markets[ticker] = market
						}
					},
				)
				genesis.GenesisTime = GenesisTime
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			currencyPair, err := slinky.MarketPairToCurrencyPair(testMarketParam.Param.Pair)
			require.NoError(t, err)
			market, err := tApp.App.MarketMapKeeper.GetMarket(ctx, currencyPair.String())
			require.NoError(t, err)
			require.False(t, market.Ticker.Enabled)

			initMarketParams := tApp.App.PricesKeeper.GetAllMarketParams(ctx)
			initPerpetuals := tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx)
			initClobPairs := tApp.App.ClobKeeper.GetAllClobPairs(ctx)

			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				tc.proposedMsgs,
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			if tc.expectSubmitProposalFail {
				proposalsIter, err := tApp.App.GovKeeper.Proposals.Iterate(ctx, nil)
				require.NoError(t, err)
				proposals, err := proposalsIter.Values()
				require.NoError(t, err)
				require.Equal(t, initMarketParams, tApp.App.PricesKeeper.GetAllMarketParams(ctx))
				require.Equal(t, initPerpetuals, tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx))
				require.Equal(t, initClobPairs, tApp.App.ClobKeeper.GetAllClobPairs(ctx))
				require.Len(t, proposals, 0)
				return
			}

			switch tc.expectedProposalStatus {
			case govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED:
				require.Equal(t, initMarketParams, tApp.App.PricesKeeper.GetAllMarketParams(ctx))
				require.Equal(t, initPerpetuals, tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx))
				require.Equal(t, initClobPairs, tApp.App.ClobKeeper.GetAllClobPairs(ctx))
				// Check that market is still disabled in market map.
				market, _ := tApp.App.MarketMapKeeper.GetMarket(ctx, currencyPair.String())
				require.False(t, market.Ticker.Enabled)
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

				// Check that market is enabled in market map.
				market, _ := tApp.App.MarketMapKeeper.GetMarket(ctx, currencyPair.String())
				require.True(t, market.Ticker.Enabled)

				// Advance to some blocks after, and place an order on the market.
				ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+NumBlocksAfterTradingEnabled, testapp.AdvanceToBlockOptions{})
				price, err := tApp.App.PricesKeeper.GetMarketPrice(ctx, testMarketParam.Param.Id)
				require.NoError(t, err)
				// No oracle price updates were made.
				require.Equal(t, uint64(0), price.Price)

				// Place an order on the market which is now ACTIVE with 0 oracle price.
				checkTx := testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(
					OrderTemplate_Alice_Num0_Id0_Clob0_Buy_LongTerm,
				))
				resp := tApp.CheckTx(checkTx[0])
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Contains(t,
					resp.Log,
					satypes.ErrProductPositionNotUpdatable.Error(),
					"expected CheckTx response log to contain: %s, got: %s",
					satypes.ErrProductPositionNotUpdatable,
					resp.Log,
				)

				// Advance to the next block and check chain is not halted.
				tApp.AdvanceToBlock(
					uint32(ctx.BlockHeight())+1,
					testapp.AdvanceToBlockOptions{},
				)
			default:
				t.Errorf("unexpected proposal status: %s", tc.expectedProposalStatus)
			}
		})
	}
}
