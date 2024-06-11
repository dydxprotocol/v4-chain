package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestWindDownMarketProposal(t *testing.T) {
	tests := map[string]struct {
		subaccounts               []satypes.Subaccount
		preexistingStatefulOrders []clobtypes.MsgPlaceOrder
		orders                    []clobtypes.MsgPlaceOrder

		expectedSubaccounts []satypes.Subaccount
	}{
		`Succeeds with final settlement deleveraging, non-negative TNC accounts deleveraged
			at oracle price`: {
			subaccounts: []satypes.Subaccount{
				// well-collateralized long and short positions
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_50_000,
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
				},
			},
		},
		`Succeeds with final settlement deleveraging, negative TNC accounts deleveraged at
			bankruptcy price`: {
			subaccounts: []satypes.Subaccount{
				// negative TNC position
				constants.Carl_Num0_1BTC_Short_49999USD,
				// offsetting position
				constants.Dave_Num0_1BTC_Long_50001USD,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
				},
			},
		},
		`Succeeds cancelling open stateful orders on both sides from different subaccounts`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			preexistingStatefulOrders: []clobtypes.MsgPlaceOrder{
				{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price5_GTBT5,
				},
				{
					Order: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO,
				},
				{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
				},
			},
		},
		`Succeeds blocking new orders from being placed`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
			},
			orders: []clobtypes.MsgPlaceOrder{
				// Short term orders
				{
					Order: constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
				},
				{
					Order: constants.Order_Alice_Num0_Id1_Clob0_Buy5_Price15_GTB20_IOC,
				},
				{
					Order: constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO,
				},
				// Stateful orders
				{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price5_GTBT5,
				},
				{
					Order: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO,
				},
				{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}

			// Initialize test app
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{
							constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						// Set oracle prices in the genesis.
						pricesGenesis := constants.TestPricesGenesisState
						*genesisState = pricesGenesis
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{constants.ClobPair_Btc}
						genesisState.LiquidationsConfig = constants.LiquidationsConfig_No_Limit
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			for _, order := range tc.preexistingStatefulOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					order,
				) {
					resp := tApp.CheckTx(checkTx)
					require.True(
						t,
						resp.IsOK(),
						"Expected CheckTx to succeed. Response: %+v",
						resp,
					)
				}
			}

			// Place stateful orders in state, verify they were placed
			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{})
			for _, order := range tc.preexistingStatefulOrders {
				_, exists := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.Order.OrderId)
				require.True(t, exists)
			}

			// Build MsgUpdateClobPair
			clobPairId := 0
			clobPair, exists := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(clobPairId))
			require.True(t, exists)
			clobPair.Status = clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT
			msgUpdateClobPairToFinalSettlement := &clobtypes.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
				ClobPair:  clobPair,
			}

			// Submit and Tally Proposal, proposal is executed in this step
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{
					msgUpdateClobPairToFinalSettlement,
				},
				uint32(ctx.BlockHeight())+1,
				false,
				false,
				govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
			)

			// Verify events emitted by indexer in the last block, the block in which the gov proposal was executed
			events := []*indexer_manager.IndexerTendermintEvent{
				{
					Subtype:             indexerevents.SubtypeUpdateClobPair,
					OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
					EventIndex:          0,
					Version:             indexerevents.UpdateClobPairEventVersion,
					DataBytes: indexer_manager.GetBytes(
						indexerevents.NewUpdateClobPairEvent(
							clobPair.GetClobPairId(),
							clobPair.Status,
							clobPair.QuantumConversionExponent,
							clobtypes.SubticksPerTick(clobPair.GetSubticksPerTick()),
							satypes.BaseQuantums(clobPair.GetStepBaseQuantums()),
						),
					),
				},
			}
			for i, order := range tc.preexistingStatefulOrders {
				events = append(
					events,
					&indexer_manager.IndexerTendermintEvent{
						Subtype:             indexerevents.SubtypeStatefulOrder,
						OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
						EventIndex:          uint32(i + 1),
						Version:             indexerevents.StatefulOrderEventVersion,
						DataBytes: indexer_manager.GetBytes(
							indexerevents.NewStatefulOrderRemovalEvent(
								order.Order.OrderId,
								indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_FINAL_SETTLEMENT,
							),
						),
					},
				)
			}
			expectedOnChainMessageAfterGovProposal := indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: uint32(ctx.BlockHeight()),
					Time:   ctx.BlockTime(),
					Events: events,
					TxHashes: []string{
						string(lib.GetTxHash(
							[]byte{},
						)),
					},
				},
			)
			onchainMessages := msgSender.GetOnchainMessages()
			require.Equal(
				t,
				expectedOnChainMessageAfterGovProposal,
				onchainMessages[len(onchainMessages)-1],
			)

			// Verify clob pair is transitioned to final settlement
			updatedClobPair, exists := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(clobPairId))
			require.True(t, exists)
			require.Equal(t, clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT, updatedClobPair.Status)

			// Verify that open stateful orders are removed from state
			for _, order := range tc.preexistingStatefulOrders {
				_, exists := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.Order.OrderId)
				require.False(t, exists)
			}

			// Set liquidation daemon info, to simulate liquidations daemon updating SubaccountOpenPositionInfo
			_, err := tApp.App.Server.LiquidateSubaccounts(
				ctx,
				&api.LiquidateSubaccountsRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
						},
					},
				},
			)
			require.NoError(t, err)

			// Advance block again to trigger final settlement deleveraging in PrepareCheckState
			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{})

			// Verify that final settlement deleveraging occurs
			for _, expectedSubaccount := range tc.expectedSubaccounts {
				subaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id)
				require.Equal(t, expectedSubaccount, subaccount)
			}

			// Attempt to place new orders, should fail validation
			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					order,
				) {
					resp := tApp.CheckTx(checkTx)
					require.Contains(t, resp.Log, "trading is disabled for clob pair")
					require.False(
						t,
						resp.IsOK(),
						"Expected CheckTx to fail. Response: %+v",
						resp,
					)
				}
			}
		})
	}
}
