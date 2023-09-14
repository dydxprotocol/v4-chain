package bridge_test

import (
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

const (
	TEST_DENOM = "dv4tnt"
)

func TestBridge_Success(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		// bridge events.
		bridgeEvents []bridgetypes.BridgeEvent
		// propose params.
		proposeParams bridgetypes.ProposeParams
		// safety params.
		safetyParams bridgetypes.SafetyParams
		// block time to advance to.
		blockTime time.Time

		// Expectations.
		// whether bridge tx should have non-empty bridge events.
		expectNonEmptyBridgeTx bool
	}{
		"Success: 1 bridge event, delay 5 blocks": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
			},
			proposeParams: bridgetypes.ProposeParams{
				MaxBridgesPerBlock:           2,
				ProposeDelayDuration:         0,
				SkipRatePpm:                  0, // do not skip proposing bridge events.
				SkipIfBlockDelayedByDuration: time.Second * 10,
			},
			safetyParams: bridgetypes.SafetyParams{
				IsDisabled:  false,
				DelayBlocks: 5,
			},
			blockTime:              time.Now(),
			expectNonEmptyBridgeTx: true,
		},
		"Success: 4 bridge event, delay 27 blocks": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1,
				constants.BridgeEvent_Id3_Height3,
			},
			proposeParams: bridgetypes.ProposeParams{
				MaxBridgesPerBlock:           4,
				ProposeDelayDuration:         0,
				SkipRatePpm:                  0,
				SkipIfBlockDelayedByDuration: time.Second * 10,
			},
			safetyParams: bridgetypes.SafetyParams{
				IsDisabled:  false,
				DelayBlocks: 27,
			},
			blockTime:              time.Now(),
			expectNonEmptyBridgeTx: true,
		},
		"Skipped: wait for other validators to recognize bridge events": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
			},
			proposeParams: bridgetypes.ProposeParams{
				MaxBridgesPerBlock: 2,
				// wait for 10 seconds before proposing bridge events.
				ProposeDelayDuration:         time.Second * 10,
				SkipRatePpm:                  0,
				SkipIfBlockDelayedByDuration: time.Second * 10,
			},
			safetyParams: bridgetypes.SafetyParams{
				IsDisabled:  false,
				DelayBlocks: 5,
			},
			blockTime:              time.Now(),
			expectNonEmptyBridgeTx: false,
		},
		"Skipped: block delayed by too much": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
			},
			proposeParams: bridgetypes.ProposeParams{
				MaxBridgesPerBlock:           2,
				ProposeDelayDuration:         0,
				SkipRatePpm:                  0,
				SkipIfBlockDelayedByDuration: time.Second * 10,
			},
			safetyParams: bridgetypes.SafetyParams{
				IsDisabled:  false,
				DelayBlocks: 5,
			},
			// should skip proposing bridge events as block time is 11 seconds ago,
			// which is more than 10 seconds of `SkipIfBlockDelayedByDuration`.
			blockTime:              time.Now().Add(-time.Second * 11),
			expectNonEmptyBridgeTx: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.ProposeParams = tc.proposeParams
						genesisState.SafetyParams = tc.safetyParams
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			// Get initial balances of addresses and their expected balances after bridging.
			initialBalances := make(map[string]sdk.Coin)
			expectedBalances := make(map[string]sdk.Coin)
			for _, event := range tc.bridgeEvents {
				if _, exists := expectedBalances[event.Address]; exists {
					expectedBalances[event.Address] = expectedBalances[event.Address].Add(event.Coin)
				} else {
					initialBalance := tApp.App.BankKeeper.GetBalance(
						ctx,
						sdk.MustAccAddressFromBech32(event.Address),
						TEST_DENOM,
					)
					initialBalances[event.Address] = initialBalance
					expectedBalances[event.Address] = initialBalance.Add(event.Coin)
				}
			}

			res, error := tApp.App.Server.AddBridgeEvents(ctx, &api.AddBridgeEventsRequest{
				BridgeEvents: tc.bridgeEvents,
			})
			require.NoError(t, error)
			require.Equal(t, &api.AddBridgeEventsResponse{}, res)

			// Verify that balances have not changed at the block right before the one where complete
			// bridge messages should be executed, which is `DelayBlocks+2` because
			// 1. Bridge events are recognized by server at block 1.
			// 2. Bridge events are proposed at block 2 and complete bridge messages are delayed for
			//    `DelayBlocks` number of blocks.
			// 3. Complete bridge messages are executed at block `DelayBlocks+2`.
			ctx = tApp.AdvanceToBlock(tc.safetyParams.DelayBlocks+1, testapp.AdvanceToBlockOptions{
				BlockTime: tc.blockTime.Add(-time.Second * 1),
			})
			for _, event := range tc.bridgeEvents {
				balance := tApp.App.BankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(event.Address),
					TEST_DENOM,
				)
				require.Equal(t, initialBalances[event.Address], balance)
			}

			// Verify that balances are updated, if bridge events were proposed, at the block where
			// complete bridge messages are executed.
			ctx = tApp.AdvanceToBlock(tc.safetyParams.DelayBlocks+2, testapp.AdvanceToBlockOptions{
				BlockTime: tc.blockTime,
			})
			for _, event := range tc.bridgeEvents {
				balance := tApp.App.BankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(event.Address),
					TEST_DENOM,
				)
				if tc.expectNonEmptyBridgeTx { // bridge events were proposed.
					require.Equal(t, expectedBalances[event.Address], balance)
				} else {
					require.Equal(t, initialBalances[event.Address], balance)
				}
			}
		})
	}
}

func TestBridge_REJECT(t *testing.T) {
	e0 := constants.BridgeEvent_Id0_Height0
	e1 := constants.BridgeEvent_Id1_Height0

	tests := map[string]struct {
		// bridge events.
		bridgeEvents []bridgetypes.BridgeEvent
	}{
		"Bad coin": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				{
					Id: e1.Id,
					Coin: sdk.NewCoin(
						e1.Coin.Denom,
						e1.Coin.Amount.Add(sdk.NewInt(1)), // bad amount.
					),
					Address:        e1.Address,
					EthBlockHeight: e1.EthBlockHeight,
				},
			},
		},
		"Event not recognized": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1, // event not recognized.
			},
		},
		"First event not next to be acknowledged": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Add good bridge events to app server.
			res, error := tApp.App.Server.AddBridgeEvents(ctx, &api.AddBridgeEventsRequest{
				BridgeEvents: []bridgetypes.BridgeEvent{
					e0,
					e1,
				},
			})
			require.Equal(t, &api.AddBridgeEventsResponse{}, res)
			require.NoError(t, error)

			proposal := tApp.PrepareProposal()
			// Propose bad bridge events by overriding bridge tx, which is the third-to-last tx in the proposal.
			proposal.Txs[len(proposal.Txs)-3] = testtx.MustGetTxBytes(
				&bridgetypes.MsgAcknowledgeBridges{
					Events: tc.bridgeEvents,
				},
			)
			processRequest := abcitypes.RequestProcessProposal{
				Txs:                proposal.Txs,
				Hash:               tApp.GetHeader().AppHash,
				Height:             tApp.GetHeader().Height,
				Time:               tApp.GetHeader().Time,
				NextValidatorsHash: tApp.GetHeader().NextValidatorsHash,
				ProposerAddress:    tApp.GetHeader().ProposerAddress,
			}
			// Verify that the bad proposal is rejected.
			processProposalResp := tApp.App.ProcessProposal(processRequest)
			require.Equal(t, abcitypes.ResponseProcessProposal_REJECT, processProposalResp.Status)
		})
	}
}
