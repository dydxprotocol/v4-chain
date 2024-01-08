package bridge_test

import (
	"context"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"

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
	TEST_DENOM = "adv4tnt"
)

func TestBridge_Success(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// bridge events.
		bridgeEvents []bridgetypes.BridgeEvent
		// propose params.
		proposeParams bridgetypes.ProposeParams
		// safety params.
		safetyParams bridgetypes.SafetyParams
		// block time to advance to.
		blockTime time.Time

		/* --- Expectations --- */
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
		"Success: 2 bridge events, delay 0 blocks": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			proposeParams: bridgetypes.ProposeParams{
				MaxBridgesPerBlock:           2,
				ProposeDelayDuration:         0,
				SkipRatePpm:                  0, // do not skip proposing bridge events.
				SkipIfBlockDelayedByDuration: time.Second * 10,
			},
			safetyParams: bridgetypes.SafetyParams{
				IsDisabled:  false,
				DelayBlocks: 0,
			},
			blockTime:              time.Now(),
			expectNonEmptyBridgeTx: true,
		},
		"Success: 1 bridge event with 0 coin amount, delay 5 blocks": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id4_Height0_EmptyCoin,
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
		"Success: 4 bridge events, delay 27 blocks": {
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
			tApp := testapp.NewTestAppBuilder(t).
				// These tests only contact the tApp.App.Server causing non-determinism in the
				// other App instances in TestApp used for non-determinism checking.
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *bridgetypes.GenesisState) {
							genesisState.ProposeParams = tc.proposeParams
							genesisState.SafetyParams = tc.safetyParams
						},
					)
					genesis.GenesisTime = tc.blockTime
					return genesis
				}).Build()
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

			// 'MsgCompleteBridge's should be executed at block height `DelayBlocks+2` because
			// 1. Bridge events are recognized by server at block 1.
			// 2. Bridge events are proposed at block 2 and complete bridge messages are delayed for
			//    `DelayBlocks` number of blocks.
			// 3. Complete bridge messages are executed at block `DelayBlocks+2`.
			blockHeightOfBridgeCompletion := tc.safetyParams.DelayBlocks + 2

			// Advance to block right before bridge completion, if necessary.
			if blockHeightOfBridgeCompletion-1 > uint32(ctx.BlockHeight()) {
				ctx = tApp.AdvanceToBlock(blockHeightOfBridgeCompletion-1, testapp.AdvanceToBlockOptions{})
			}
			// Verify that balances have not changed yet.
			for _, event := range tc.bridgeEvents {
				balance := tApp.App.BankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(event.Address),
					TEST_DENOM,
				)
				require.Equal(t, initialBalances[event.Address], balance)
			}

			// Verify that balances are updated, if bridge events were proposed, at the block of
			// bridge completion.
			ctx = tApp.AdvanceToBlock(blockHeightOfBridgeCompletion, testapp.AdvanceToBlockOptions{})
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
		// bridge events to propose.
		bridgeEvents []bridgetypes.BridgeEvent
		// whether bridging is disabled.
		bridgingDisabled bool
	}{
		"Bad coin denom": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				{
					Id: e1.Id,
					Coin: sdk.NewCoin(
						e1.Coin.Denom+"a", // bad denom.
						e1.Coin.Amount,
					),
					Address:        e1.Address,
					EthBlockHeight: e1.EthBlockHeight,
				},
			},
		},
		"Bad coin amount": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				{
					Id: e1.Id,
					Coin: sdk.NewCoin(
						e1.Coin.Denom,
						e1.Coin.Amount.Add(sdkmath.NewInt(1)), // bad amount.
					),
					Address:        e1.Address,
					EthBlockHeight: e1.EthBlockHeight,
				},
			},
		},
		"Bad address": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				{
					Id: e1.Id,
					Coin: sdk.NewCoin(
						e1.Coin.Denom,
						e1.Coin.Amount,
					),
					Address:        e1.Address + "a", // bad address.
					EthBlockHeight: e1.EthBlockHeight,
				},
			},
		},
		"Bad eth block height": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				{
					Id: e1.Id,
					Coin: sdk.NewCoin(
						e1.Coin.Denom,
						e1.Coin.Amount,
					),
					Address:        e1.Address,
					EthBlockHeight: e1.EthBlockHeight + 1, // bad eth block height.
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
		"Bridging is disabled and non-empty bridge events are proposed": {
			bridgeEvents: []bridgetypes.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			bridgingDisabled: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.SafetyParams = bridgetypes.SafetyParams{
							IsDisabled:  tc.bridgingDisabled,
							DelayBlocks: 5,
						}
					},
				)
				return genesis
			}).Build()
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

			proposal, err := tApp.PrepareProposal()
			require.NoError(t, err)
			// Propose bridge events by overriding bridge tx, which is the third-to-last tx in the proposal.
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
			processProposalResp, err := tApp.App.ProcessProposal(&processRequest)
			require.NoError(t, err)
			require.Equal(t, abcitypes.ResponseProcessProposal_REJECT, processProposalResp.Status)
		})
	}
}

// This test case makes sure that the bridge server still accepts bridge events from bridge daemon when
// Acknowledged Event ID (on-chain) is greater than Recognized Event ID (off-chain), which can happen if
// - `AcknowledgedEventInfo` is initialized with Id > 0 in genesis.
// - A node falls behind the chain in terms of acknowledged events, a scenario that can be caused by:
//   - the node experiences issues with its Eth RPC endpoint and doesn't recognized events that are
//     accepted by the rest of the chain.
//   - the node restarts.
func TestBridge_AcknowledgedEventIdGreaterThanRecognizedEventId(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *bridgetypes.GenesisState) {
				genesisState.AcknowledgedEventInfo = bridgetypes.BridgeEventInfo{
					NextId:         2,
					EthBlockHeight: 123,
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	// Verify that AcknowledgedEventInfo.NextId is 2.
	aei := tApp.App.BridgeKeeper.GetAcknowledgedEventInfo(ctx)
	require.Equal(t, uint32(2), aei.NextId)

	// Verify that RecognizedEventInfo.NextId is still 0.
	rei := tApp.App.BridgeKeeper.GetRecognizedEventInfo(ctx)
	require.Equal(t, uint32(0), rei.NextId)

	// Verify that bridge query `RecognizedEventInfo` returns whichever of AcknowledgedEventInfo and
	// RecognizedEventInfo has a greater `NextId` (which is AcknowledgedEventInfo in this case).
	reiRequest := bridgetypes.QueryRecognizedEventInfoRequest{}
	abciResponse, err := tApp.App.Query(
		context.Background(),
		&abcitypes.RequestQuery{
			Path: "/dydxprotocol.bridge.Query/RecognizedEventInfo",
			Data: tApp.App.AppCodec().MustMarshal(&reiRequest),
		},
	)
	require.True(t, abciResponse.IsOK())
	require.NoError(t, err)
	var reiResponse bridgetypes.QueryRecognizedEventInfoResponse
	tApp.App.AppCodec().MustUnmarshal(abciResponse.Value, &reiResponse)
	require.Equal(t, aei, reiResponse.Info) // Verify that AcknowledgedEventInfo is returned.

	// Verify that it's ok to add events starting from `NextId` in above query response.
	_, err = tApp.App.Server.AddBridgeEvents(ctx, &api.AddBridgeEventsRequest{
		BridgeEvents: []bridgetypes.BridgeEvent{
			{
				Id:             reiResponse.Info.NextId,
				Coin:           sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
				Address:        constants.BobAccAddress.String(),
				EthBlockHeight: 234,
			},
		},
	})
	require.NoError(t, err)
}
